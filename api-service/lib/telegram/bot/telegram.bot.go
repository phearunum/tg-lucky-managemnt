package telegram

import (
	redis "api-service/config"
	luckyModel "api-service/lib/lucky/models"
	"api-service/lib/telegram/dto"
	"api-service/lib/telegram/model"
	repository "api-service/lib/telegram/repository"
	"api-service/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type TelegramBotService struct {
	repo *repository.TelegramAccountRepository
}

func NewTelegramBotService(repo *repository.TelegramAccountRepository) *TelegramBotService {

	return &TelegramBotService{repo: repo}
}

type RequestSettingDTO struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	TimeRest   float64 `json:"time_rest"`
	Status     string  `json:"status"`
	BotName    string  `json:"bot_name"`
	ButtonType string  `json:"button_type"`
}
type LocationSettingDTO struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Lat      string  `json:"lat"`
	Long     string  `json:"long"`
	BotToken string  `json:"token_bot"`
	Allow    float64 `json:"allow"`
}
type Bot struct {
	Name              string
	Token             string
	Debug             bool
	GroupID           int64
	BackBtn           bool
	locationRequested map[int64]time.Time
	ReplyCounts       map[int64]int
	Replies           map[int64][]string
	Repo              *repository.TelegramAccountRepository
	mu                sync.Mutex
}

func InitializeRedis(token string) {

}
func NewBot(name, token string, debug bool, groupID int64, backBtn bool, repo *repository.TelegramAccountRepository) *Bot {
	// Init
	keybords, errf := repo.FilterRequestSettingAll()
	if errf != nil {
		log.Println("Failed to acknowledge callback query:", errf)
	}
	utils.InfoLog(keybords, "Init Redis Bot keybords")
	redisKey := fmt.Sprintf("request:bot:%s", "setting")
	if err := redis.SetWithoutExpired(redisKey, keybords); err != nil {
		utils.WarnLog(err.Error(), "Error Redis Init")
	}
	// Init Location
	locations, errL := repo.GetLocationSettingAll()
	if errL != nil {
		log.Println("Failed to acknowledge callback query:", errL)
	}
	utils.InfoLog(locations, "Init Redis Bot Location")
	// Marshal locations into a JSON byte slice
	locationsBytes, err := json.Marshal(locations)
	if err != nil {
		utils.ErrorLog(fmt.Errorf("failed to marshal locations: %v", err), err.Error())
	}
	// Now unmarshal the JSON byte slice into locationDTO
	var locationDTO []LocationSettingDTO
	if err := json.Unmarshal(locationsBytes, &locationDTO); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
	}
	for _, loca := range locationDTO {
		redisKeyL := fmt.Sprintf("location:bot:%s", loca.Name)
		if err := redis.SetWithoutExpired(redisKeyL, loca); err != nil {
			utils.WarnLog(err.Error(), "Error Redis Init")
		}
	}
	return &Bot{
		Name:              name,
		Token:             token,
		Debug:             debug,
		GroupID:           groupID,
		BackBtn:           backBtn,
		locationRequested: make(map[int64]time.Time),
		ReplyCounts:       make(map[int64]int),
		Replies:           make(map[int64][]string),
		Repo:              repo,
	}
}
func (b *Bot) Start() {
	log.Printf("Starting bot %s...", b.Name)
	// Initialize the bot
	botAPI, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Printf("Failed to initialize bot %s: %v", b.Name, err)
		return
	}
	log.Printf("Bot %s initialized successfully", b.Name)
	// Set debug mode
	botAPI.Debug = b.Debug
	//botGroupID := b.GroupID

	log.Printf("Authorized bot %s (%s)", b.Name, botAPI.Self.UserName)

	// Start listening for updates
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)
	log.Printf("Bot %s is now listening for updates...", b.Name)
	// Handle incoming updates
	for update := range updates {
		// Handle regular messages (only process /start and specific commands)
		if update.Message != nil {
			b.handleMessage(botAPI, update.Message, update)

		}
		if update.EditedMessage != nil && update.EditedMessage.Location != nil {
			b.handleLiveLocationUpdate(botAPI, update.EditedMessage)
		}

	}
}

// handleMessage handles regular messages
func (b *Bot) handleMessage(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message, update tgbotapi.Update) {

	if message.IsCommand() && message.Command() == "start" {
		b.initKeyBoard(botAPI, message)
		return
	}
	// Check if the message is a reply to another message
	if message.ReplyToMessage != nil {
		b.mu.Lock()
		b.ReplyCounts[int64(message.ReplyToMessage.MessageID)]++
		b.Replies[int64(message.ReplyToMessage.MessageID)] = append(b.Replies[int64(message.ReplyToMessage.MessageID)], message.Text)
		b.mu.Unlock()
	}
	b.handleMessageTextDynamic(botAPI, b.GroupID, message)
	//b.HandleUpdate(botAPI, update)

}

func (b *Bot) initKeyBoard_(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var Username = message.From.FirstName + " " + message.From.LastName
	redisServerKey := fmt.Sprintf("request:bot:%s", "setting")
	SettingInfo, err := redis.Get(redisServerKey)
	if err != nil {
		utils.ErrorLog(err.Error(), " Field Get Data from Redis videoSetting:server:%d")
	}

	var botSetting []RequestSettingDTO
	if err := json.Unmarshal([]byte(SettingInfo), &botSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
	}
	var keyboardButtons []tgbotapi.KeyboardButton
	var addButtonLine2 []tgbotapi.KeyboardButton
	for _, keyboard := range botSetting {
		if keyboard.BotName == strconv.FormatInt(b.GroupID, 10) {
			if keyboard.ButtonType == "track" {
				keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
			} else {
				addButtonLine2 = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
			}
		}
	}
	// Add additional buttons to the keyboardButtons slice
	additionalButtons := []tgbotapi.KeyboardButton{}

	// Add additional buttons
	keyboardButtons = append(keyboardButtons, additionalButtons...)
	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
		tgbotapi.NewKeyboardButtonRow(addButtonLine2...),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("⬅️ Back"),
		),
	)
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("📍 CheckIn "),
		),
	)
	locationKeyboard.ResizeKeyboard = true
	locationKeyboard.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
	msg.ReplyMarkup = locationKeyboard
	msg.ReplyMarkup = replyKeyboard
	botAPI.Send(msg)
}
func (b *Bot) initKeyBoard(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	var Username = message.From.FirstName + " " + message.From.LastName
	redisServerKey := fmt.Sprintf("request:bot:%s", "setting")
	SettingInfo, err := redis.Get(redisServerKey)
	if err != nil {
		utils.ErrorLog(err.Error(), "Field Get Data from Redis videoSetting:server:%d")
		return false
	}
	var botSetting []RequestSettingDTO
	if err := json.Unmarshal([]byte(SettingInfo), &botSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
		return false // Early return on error
	}
	utils.InfoLog(botSetting, "botSetting")
	if len(botSetting) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		botAPI.Send(msg)
		return false
	}

	var keyboardButtons []tgbotapi.KeyboardButton
	//var addButtonLine2 []tgbotapi.KeyboardButton

	for _, keyboard := range botSetting {
		utils.InfoLog(keyboard.BotName, strconv.FormatInt(message.Chat.ID, 10))
		if keyboard.BotName == strconv.FormatInt(message.Chat.ID, 10) {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
		}
	}
	if len(keyboardButtons) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		botAPI.Send(msg)
		return false
	}
	var replyKeyboard tgbotapi.ReplyKeyboardMarkup
	if b.BackBtn {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Back"),
			),
		)
	} else {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
		)
	}
	// Set properties after initialization
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	// Choose which keyboard to attach to the message
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
	msg.ReplyMarkup = replyKeyboard // Set the replyKeyboard here

	// If you want to send a message with both keyboards, you cannot set both at the same time
	// You need to decide on one or merge their functionalities based on your requirements
	if _, err := botAPI.Send(msg); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	return true
}

func (b *Bot) handleMessageTextDynamic(botAPI *tgbotapi.BotAPI, botGroupID int64, message *tgbotapi.Message) {
	var Username = message.From.FirstName + " " + message.From.LastName
	var chatId = message.From.ID
	var btnText = "🕗 ClockTime "
	currentTime := time.Now()
	chatIdStr := strconv.FormatInt(chatId, 10)
	groupIdStr := strconv.FormatInt(botGroupID, 10)
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	// Init Case
	redisServerKey := fmt.Sprintf("request:bot:%s", "setting")
	SettingInfo, err := redis.Get(redisServerKey)
	if err != nil {
		utils.ErrorLog(err.Error(), " Field Get Data from Redis videoSetting:server:%d")
	}
	var botSetting []RequestSettingDTO
	if err := json.Unmarshal([]byte(SettingInfo), &botSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
	}
	var keyboardButtons []tgbotapi.KeyboardButton
	//var addButtonLine2 []tgbotapi.KeyboardButton
	for _, keyboard := range botSetting {
		if keyboard.BotName == strconv.FormatInt(message.Chat.ID, 10) {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
			//if keyboard.ButtonType == "track" {
			//	keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
			//} else {
			//	addButtonLine2 = append(keyboardButtons, tgbotapi.NewKeyboardButton(keyboard.Name))
			//}
		}
	}
	// Send Reply with keyboard
	additionalButtons := []tgbotapi.KeyboardButton{}
	// Add additional buttons
	keyboardButtons = append(keyboardButtons, additionalButtons...)
	var replyKeyboard tgbotapi.ReplyKeyboardMarkup
	if b.BackBtn {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Back"),
			),
		)
	} else {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
		)
	}
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	// Convert slice to map for faster lookups
	botSettingMap := make(map[string]RequestSettingDTO, len(botSetting))
	for _, keyboard := range botSetting {
		if keyboard.ButtonType == "track" {
			botSettingMap[keyboard.Name] = keyboard
		}
	}
	keyboard, exists := botSettingMap[message.Text]
	if exists {
		// Convert chatIdStr to int64
		chatId64, errchatId := strconv.ParseInt(chatIdStr, 10, 64)
		if errchatId != nil {
			utils.ErrorLog(errchatId.Error(), "Invalid chat ID format")
			return
		}

		// Capture both return values properly
		lastBreak, err := b.Repo.FilterRequestGetBreakOne(chatId64, "pending")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WarnLog(err.Error(), "Records Not Found")
			} else {
				utils.ErrorLog(err.Error(), "Database Query Error")
			}
			return // Stop execution if there's an error
		}
		warningMessage := ""
		if lastBreak != nil && lastBreak.ID > 0 {

			warningMessage = fmt.Sprintf(
				"⚠️! missing back.\n"+
					"🕗 Start Time: %s\n"+
					"🍽️ Break : %s",
				lastBreak.StartTime, lastBreak.RequestType,
			)
			updateSuccess := b.Repo.UpdateRequestStatus(int(lastBreak.ID), "finished", "⚠️ !missing back")
			if !updateSuccess {
				utils.WarnLog("Failed to update request status", "ID: "+strconv.Itoa(int(lastBreak.ID)))
			}
		}

		// Create a new request
		request := model.UserRequest{
			ChatID:      chatIdStr, // FIX: Use int64 instead of string
			RequestType: keyboard.Name,
			StartTime:   time.Now().Format("2006-01-02 15:04:05"),
			Status:      "pending",
			AccountName: Username,
			BotName:     keyboard.BotName,
		}
		var lineBreak = "\n-------------\n"
		// Save the request
		_, errSave := b.Repo.SaveRequest(&request) // FIX: Avoid variable shadowing
		if errSave != nil {
			log.Println("Failed to acknowledge callback query:", errSave)
			return
		}
		msgText := fmt.Sprintf("%s\n%s%s🍽️ New Break: %s\n🕗 Time: %s",
			Username, warningMessage, lineBreak, keyboard.Name, formattedTime)

		// Send a Telegram message
		//msg := tgbotapi.NewMessage(botGroupID, Username+" :  "+keyboard.Name+", 🕗 Time: "+formattedTime)
		msg := tgbotapi.NewMessage(botGroupID, msgText)
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = replyKeyboard
		botAPI.Send(msg)

	} else if message.Text == "⬅️ Back" {
		requestData, err := b.Repo.FilterRequestOne(chatId)
		utils.WarnLog(requestData, "requestData, err := b.Repo.FilterRequestOne(chatId)")
		if err != nil {
			log.Println("Failed to acknowledge callback query:", err)
			return
		}
		if requestData == nil || requestData.ID == 0 {
			msg := tgbotapi.NewMessage(botGroupID, "🚫 No break...")
			msg.ReplyToMessageID = message.MessageID
			msg.ReplyMarkup = replyKeyboard
			botAPI.Send(msg)
			return
		}
		// Set EndTime as the current time (time.Now())
		endTime := time.Now()
		endTimeStr := endTime.Format("2006-01-02 15:04:05")
		fmt.Println("Back Event", requestData)
		startTimeParsed, _ := time.Parse("2006-01-02 15:04:05", requestData.StartTime)
		endTimeParsed, _ := time.Parse("2006-01-02 15:04:05", endTimeStr)
		totalTime := endTimeParsed.Sub(startTimeParsed)
		totalTimeMinutes := float64(totalTime.Minutes()) // Convert to minutes as an integer
		// Filter Allow time
		filter, errf := b.Repo.FilterRequestSetting(requestData.RequestType)
		utils.WarnLog(filter, "FilterRequestSetting")

		if errf != nil {
			log.Println("Failed to acknowledge callback query:", errf)
			return
		}
		if filter == nil || filter.ID == 0 {
			msg := tgbotapi.NewMessage(botGroupID, "🚫 No break...")
			msg.ReplyToMessageID = message.MessageID
			msg.ReplyMarkup = replyKeyboard
			botAPI.Send(msg)
			return
		}
		var Messagetxt = ""
		if totalTimeMinutes > filter.TimeRest {
			excessMinutes := float64(totalTimeMinutes - filter.TimeRest) // Convert excess to integer
			Messagetxt = fmt.Sprintf("⚠️ Warning: %.2f min", excessMinutes)
		} else {
			Messagetxt = " "
		}
		if filter.ButtonType == "normal" || filter.ButtonType == "" {
			Messagetxt = " "
		}
		request := model.UserRequest{
			ID:          requestData.ID,
			ChatID:      requestData.ChatID,
			AccountName: requestData.AccountName,
			StartTime:   requestData.StartTime,
			EndTime:     time.Now().Format("2006-01-02 15:04:05"), // Current time as start time
			Status:      "finished",
			AllowTime:   filter.TimeRest,
			TotalTime:   totalTimeMinutes,
			RequestType: requestData.RequestType, // Default status
			Message:     Messagetxt,
			BotName:     requestData.BotName,
			CreatedAt:   requestData.CreatedAt,
		}
		_, err = b.Repo.SaveRequest(&request)
		if err != nil {
			log.Println("Failed to acknowledge callback query:", err)
		}
		var lineBreak = "\n-------------\n"
		msg := tgbotapi.NewMessage(botGroupID, Username+lineBreak+" : 🔙  time [ "+requestData.RequestType+" ], 🕗 Time: "+formattedTime+", "+Messagetxt)
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = replyKeyboard
		botAPI.Send(msg)
		return
	} else if containsKeyword(message.Text, "check") || containsKeyword(message.Text, "CheckIn") || containsKeyword(message.Text, "Scan") || containsKeyword(message.Text, "Check") || containsKeyword(message.Text, "Time") {
		// Step 1: Send the location request button
		locationKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation(btnText),
			),
		)
		locationKeyboard.ResizeKeyboard = true
		locationKeyboard.OneTimeKeyboard = true
		if message.Chat.Type == "private" {
			//msg := tgbotapi.NewMessage(message.From.ID, "Share Live Location (15 min)")
			//msg.ReplyMarkup = locationKeyboard
			//botAPI.Send(msg)
		} else {

			privateChatLink := fmt.Sprintf("https://t.me/%s", b.Name)
			msg := tgbotapi.NewMessage(b.GroupID, fmt.Sprintf(
				"%s  📣 Please [Click here](%s) to [ Check In/Out ]", Username, privateChatLink,
			))
			msg.ParseMode = "Markdown"
			msg.ReplyMarkup = replyKeyboard
			botAPI.Send(msg)
		}
	} else if containsKeyword(message.Text, "/GetPhone") || containsKeyword(message.Text, "🖨 Phone") || containsKeyword(message.Text, "🖨 Get") || containsKeyword(message.Text, "📠 Phone") || containsKeyword(message.Text, "🖨 Phone +5") {
		var lineBreak = "\n------------------\n"
		phoneLists, err := b.Repo.FilterGet5Numbers(b.Name)
		if err != nil {
			log.Printf("Error retrieving phone numbers: %v", err)
			return
		}
		var phoneListString strings.Builder
		var phoneIDs []uint
		for _, phone := range phoneLists {
			phoneIDs = append(phoneIDs, phone.ID)
			phoneListString.WriteString(phone.Phone + "\n")
		}
		if len(phoneIDs) > 0 {
			bulkRequest := dto.BulkRequestDTO{
				ID:            phoneIDs,
				Requester:     Username,
				RequesterDate: time.Now().Format("2006-01-02 15:04:05"), // Use current timestamp
			}
			b.Repo.BulkUpdatePhone(bulkRequest)
		}

		//msg := tgbotapi.NewMessage(message.Chat.ID, Username+" Work hard and try your best ✅ :"+lineBreak+"\n"+escapeMarkdownV2(phoneListString.String()))
		msg := tgbotapi.NewMessage(botGroupID, escapeMarkdownV2(Username+".Work hard and try your best✅ :"+lineBreak+"```\n"+phoneListString.String()+"```"))
		msg.ReplyToMessageID = int(message.MessageID)
		msg.ParseMode = "MarkdownV2"
		// Add inline keyboard with "View Replies" button
		//inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		//	tgbotapi.NewInlineKeyboardRow(
		//		tgbotapi.NewInlineKeyboardButtonData("View Replies", "view_replies"),
		//	),
		//)
		//msg.ReplyMarkup = inlineKeyboard
		//msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
		//msg.ReplyMarkup = replyKeyboard

		botAPI.Send(msg)

	} else if containsKeyword(message.Text, "/register") {
		var createdAt time.Time
		var lineBreak = "\n-------------\n"
		profilePictureURL, err := getUserProfilePictureURL(botAPI, chatId)
		if err != nil {
			log.Println("Error fetching profile picture:", err)
		}
		register := luckyModel.TelegramUsers{
			ChatID:      chatIdStr,
			Username:    message.From.UserName,
			AccountName: Username,
			Profile:     profilePictureURL,
			BotName:     b.Name,
			GroupID:     groupIdStr,
			CreatedAt:   createdAt,
		}
		utils.InfoLog(register, "register")
		_, err = b.Repo.SaveRegister(&register)
		if err != nil {
			log.Println("Failed to acknowledge callback query:", err)
			utils.ErrorLog(err.Error(), "Failed to reguster account ")
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, Username+lineBreak+"Your has registerd! ")
		msg.ReplyToMessageID = message.MessageID
		// Create inline keyboard button
		inlineBtn := tgbotapi.NewInlineKeyboardButtonURL("🎁 Check Balance", "t.me/XYan_LuckyDraw_bot/Binger")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(inlineBtn),
		)
		msg.ReplyMarkup = inlineKeyboard
		botAPI.Send(msg)
		return
	} else if containsKeyword(message.Text, "Reward") || containsKeyword(message.Text, "🎰 Point") || containsKeyword(message.Text, "🎁 Point") {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Click to open Mini App")
		webAppInfo := tgbotapi.WebAppInfo{
			URL: fmt.Sprintf("https://track.igflexs.com?bot=Tg8899bot&group=%d", message.Chat.ID),
		}

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonWebApp("Open App", webAppInfo),
			),
		)
		msg.ReplyMarkup = inlineKeyboard
		botAPI.Send(msg)

		return
	}

}
func containsKeyword(text, keyword string) bool {
	return strings.Contains(text, keyword)
}
func getUserProfilePictureURL(bot *tgbotapi.BotAPI, userID int64) (string, error) {
	// Fetch user profile photos
	profilePhotos, err := bot.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{
		UserID: userID,
		Limit:  1, // Get the latest photo
	})

	if err != nil {
		return "", fmt.Errorf("error fetching user profile photos: %w", err)
	}

	if profilePhotos.TotalCount > 0 {
		// Get the file ID of the most recent photo
		fileID := profilePhotos.Photos[0][0].FileID

		// Use the file ID to get the file information
		file, err := bot.GetFile(tgbotapi.FileConfig{
			FileID: fileID,
		})

		if err != nil {
			return "", fmt.Errorf("error getting file information: %w", err)
		}

		// Construct the URL for the photo
		photoURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		return photoURL, nil
	}

	return "", nil // No profile picture found
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`*`, `\*`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`, // Escape hyphen (-)
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`,
		`!`, `\!`,
	)
	return replacer.Replace(text)
}
func (b *Bot) RequestLiveLocations(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {

	// Step 1: Send the location request button
	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(" Clock Time"), // Button text
		),
	)
	locationKeyboard.ResizeKeyboard = true
	locationKeyboard.OneTimeKeyboard = true

	if message.Chat.Type == "private" {
		// In private chat, send the live location request directly
		msg := tgbotapi.NewMessage(message.From.ID,
			"Please share your *live location* by clicking the button below and selecting 'Share My Live Location'.",
		)
		msg.ReplyMarkup = locationKeyboard
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
	} else {
		// In group chat, prompt the user to go to private chat
		privateChatLink := fmt.Sprintf("https://t.me/%s", b.Name)
		msg := tgbotapi.NewMessage(b.GroupID, fmt.Sprintf(
			"%s 📣 Please [click here](%s) to share your *live location* for [Check In/Out].",
			message.From.UserName, privateChatLink,
		))
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
	}
}
func (b *Bot) requestLiveLocation(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {
	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation("🕗 Share Live Location (15 min)"),
		),
	)
	locationKeyboard.ResizeKeyboard = true
	locationKeyboard.OneTimeKeyboard = true

	if message.Chat.Type == "private" {
		b.mu.Lock()
		b.locationRequested[message.Chat.ID] = time.Now()
		b.mu.Unlock()

		msg := tgbotapi.NewMessage(message.Chat.ID,
			"Please click the button below and select *'Share My Live Location'* for *15 minutes* to clock your time.",
		)
		msg.ReplyMarkup = locationKeyboard
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
	} else {
		privateChatLink := fmt.Sprintf("https://t.me/%s", b.Name)
		msg := tgbotapi.NewMessage(b.GroupID, fmt.Sprintf(
			"%s 📣 Please [click here](%s) to share your live location for 15 minutes.",
			message.From.UserName, privateChatLink,
		))
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
	}
}

// HandleUpdate processes the initial location, rejecting manual attachments
func (b *Bot) HandleUpdate(botAPI *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Location.LivePeriod > 0 {
		utils.WarnLog(update.Message.Location, "Live location")
	} else {
		if update.Message.Location.LivePeriod == 0 {
			utils.WarnLog(update.Message.Location, "Send location LivePeriod == 0")
		} else {
			utils.WarnLog(update.Message.Location, "Send location LivePeriod == nil")
		}

	}
	if update.Message.Venue != nil {
		utils.WarnLog(update.Message.Venue, "User sent a specific place")
	}

	if update.Message == nil {
		return
	}

	message := update.Message
	if message.Location == nil {
		b.requestLiveLocation(botAPI, message)
		return
	}

	if message.Chat.Type == "private" {
		b.mu.Lock()
		requestTime, requested := b.locationRequested[message.Chat.ID]
		b.mu.Unlock()

		// Reject if not from a recent button request
		if !requested || time.Since(requestTime) > 30*time.Second {
			msg := tgbotapi.NewMessage(message.Chat.ID,
				"⚠️ Manual location attachments are not allowed. Please use the 'Share Live Location (15 min)' button.",
			)
			msg.ParseMode = "Markdown"
			botAPI.Send(msg)
			b.requestLiveLocation(botAPI, message)
			return
		}

		// Acknowledge initial location, wait for live confirmation
		latitude := message.Location.Latitude
		longitude := message.Location.Longitude
		msg := tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("Location received, waiting to confirm it’s live for 15 minutes...\nLatitude: %.6f\nLongitude: %.6f", latitude, longitude),
		)
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
	}
}

func (b *Bot) HandleUpdatesss(botAPI *tgbotapi.BotAPI, update tgbotapi.Update) {
	utils.InfoLog(update, "Live location")
	var Username = update.Message.From.FirstName + " " + update.Message.From.LastName
	var btnText = "🕗 ClockTime "
	currentTime := time.Now()
	chatIdStr := strconv.FormatInt(update.Message.From.ID, 10)
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	if update.Message == nil {
		return // Ignore updates without a new message (e.g., EditedMessage)
	}
	message := update.Message
	if message.Location != nil {
		// Send message private
		locationKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation(btnText),
			),
		)
		locationKeyboard.ResizeKeyboard = true
		locationKeyboard.OneTimeKeyboard = true
		request := model.ClockTime{
			ChatID:      chatIdStr,
			RequestType: "ClockTime",
			StartTime:   time.Now().Format("2006-01-02 15:04:05"),
			AccountName: Username,
			Lat:         fmt.Sprintf("%f", message.Location.Latitude),  // Convert latitude to string
			Long:        fmt.Sprintf("%f", message.Location.Longitude), // Convert longitude to string
			BotName:     b.Name,
		}

		_, err := b.Repo.SaveClockTimeRequest(&request)
		if err != nil {
			log.Println("Failed to save clock time request:", err)
		}

		// Get Location from redis
		redisServerKey := fmt.Sprintf("location:bot:%s", b.Name)
		SettingInfo, err := redis.Get(redisServerKey)
		if err != nil {
			utils.ErrorLog(err.Error(), " Field Get Data from Redis videoSetting:server:%d")
		}

		var locationSetting LocationSettingDTO
		if err := json.Unmarshal([]byte(SettingInfo), &locationSetting); err != nil {
			utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
		}
		//utils.InfoLog(locationSetting, "locationSetting")
		staticLatitudeStr := locationSetting.Lat
		staticLongitudeStr := locationSetting.Long
		// Convert string to float64
		staticLatitude, err := strconv.ParseFloat(staticLatitudeStr, 64)
		if err != nil {
			log.Printf("Error converting static latitude to float64: %v\n", err)
			return
		}

		staticLongitude, err := strconv.ParseFloat(staticLongitudeStr, 64)
		if err != nil {
			log.Printf("Error converting static longitude to float64: %v\n", err)
			return
		}
		latitude := message.Location.Latitude
		longitude := message.Location.Longitude
		// Calculate distance to the static location
		distance := Haversine(staticLatitude, staticLongitude, latitude, longitude)
		// Default value if Allow is null or empty
		allow := 200.0 // Default value

		if locationSetting.Allow != 0 {
			allow = locationSetting.Allow
		}
		//utils.InfoLog(distance, "Clock distance")
		//utils.InfoLog(allow, "Allow Clock")
		//utils.InfoLog(latitude, "User latitude")
		//utils.InfoLog(longitude, "User longitude")
		if distance > allow {
			replyMsg := " ⚠️ Warning !You are outside location"
			msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
			msg.ReplyMarkup = locationKeyboard
			msg.ParseMode = "Markdown"
			botAPI.Send(msg)
			return
		}
		replyMsg := fmt.Sprintf(
			btnText+" !\n🕗 Time: %s\nLatitude: %.6f\nLongitude: %.6f\n",
			formattedTime, latitude, longitude,
		)

		msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
		msg.ReplyMarkup = locationKeyboard
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)

		// Forward location to the group
		groupMsg := tgbotapi.NewMessage(
			b.GroupID,
			fmt.Sprintf("%s : %s, 🕗 Time: %s, Latitude: %.6f, Longitude: %.6f\n",
				Username, btnText, formattedTime, latitude, longitude),
		)
		groupMsg.ParseMode = "Markdown"
		groupMsg.ReplyMarkup = locationKeyboard
		botAPI.Send(groupMsg)
	}
}

// handleLiveLocationUpdate processes the first live update
func (b *Bot) handleLiveLocationUpdate(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {
	b.mu.Lock()
	requestTime, requested := b.locationRequested[message.Chat.ID]
	if requested {
		delete(b.locationRequested, message.Chat.ID) // Clear state
	}
	b.mu.Unlock()

	if !requested || time.Since(requestTime) > 30*time.Second {
		return
	}

	// Process live location
	username := message.From.FirstName + " " + message.From.LastName
	btnText := "🕗 ClockTime "
	currentTime := time.Now()
	chatIdStr := strconv.FormatInt(message.From.ID, 10)
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	latitude := message.Location.Latitude
	longitude := message.Location.Longitude

	// Save to repository
	request := model.ClockTime{
		ChatID:      chatIdStr,
		RequestType: "ClockTime",
		StartTime:   formattedTime,
		AccountName: username,
		Lat:         fmt.Sprintf("%f", latitude),
		Long:        fmt.Sprintf("%f", longitude),
		BotName:     b.Name,
	}
	_, err := b.Repo.SaveClockTimeRequest(&request)
	if err != nil {
		log.Println("Failed to save clock time request:", err)
	}

	// Get location settings from Redis
	redisServerKey := fmt.Sprintf("location:bot:%s", b.Name)
	SettingInfo, err := redis.Get(redisServerKey) // Assuming redis.Get is defined
	if err != nil {
		utils.ErrorLog(err.Error(), "Failed to get data from Redis location:bot:%s")
	}

	var locationSetting LocationSettingDTO
	if err := json.Unmarshal([]byte(SettingInfo), &locationSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
	}
	//utils.InfoLog(locationSetting, "locationSetting")

	staticLatitudeStr := locationSetting.Lat
	staticLongitudeStr := locationSetting.Long
	staticLatitude, err := strconv.ParseFloat(staticLatitudeStr, 64)
	if err != nil {
		log.Printf("Error converting static latitude to float64: %v\n", err)
		return
	}
	staticLongitude, err := strconv.ParseFloat(staticLongitudeStr, 64)
	if err != nil {
		log.Printf("Error converting static longitude to float64: %v\n", err)
		return
	}

	// Calculate distance
	distance := Haversine(staticLatitude, staticLongitude, latitude, longitude)
	allow := 200.0 // Default value
	if locationSetting.Allow != 0 {
		allow = locationSetting.Allow
	}
	//utils.InfoLog(distance, "Clock distance")
	//utils.InfoLog(allow, "Allow Clock")
	//utils.InfoLog(latitude, "User latitude")
	//utils.InfoLog(longitude, "User longitude")

	// Prepare keyboard for potential re-request
	locationKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(btnText),
		),
	)
	locationKeyboard.ResizeKeyboard = true
	locationKeyboard.OneTimeKeyboard = true

	if distance > allow {
		replyMsg := "⚠️ Warning! You are outside the allowed location range."
		msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
		msg.ReplyMarkup = locationKeyboard
		msg.ParseMode = "Markdown"
		botAPI.Send(msg)
		return
	}

	// Send confirmation to user
	replyMsg := fmt.Sprintf(
		"%s!\n🕗 Time: %s\nLatitude: %.6f\nLongitude: %.6f\n*Live location confirmed for 15 minutes*",
		btnText, formattedTime, latitude, longitude,
	)
	msg := tgbotapi.NewMessage(message.Chat.ID, replyMsg)
	msg.ParseMode = "Markdown"
	botAPI.Send(msg)

	// Forward to group
	groupMsg := tgbotapi.NewMessage(
		b.GroupID,
		fmt.Sprintf("%s: %s, 🕗 Time: %s, Latitude: %.6f, Longitude: %.6f\n*Live location (15 min)*",
			username, btnText, formattedTime, latitude, longitude),
	)
	groupMsg.ParseMode = "Markdown"
	botAPI.Send(groupMsg)
}

// Haversine function to calculate the distance between two points
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	// Radius of the Earth in kilometers
	r := 6371.0

	// Convert degrees to radians
	lat1 = degToRad(lat1)
	lon1 = degToRad(lon1)
	lat2 = degToRad(lat2)
	lon2 = degToRad(lon2)

	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Distance in kilometers
	return r * c
}

// degToRad converts degrees to radians
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}
