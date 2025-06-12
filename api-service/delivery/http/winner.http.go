package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	dto "api-service/lib/lucky/dto"
	service "api-service/lib/lucky/service"
	"api-service/utils"

	"github.com/gorilla/mux"
)

type WinnerController struct {
	WinnerService service.WinnerServiceInterface
}

// NewWinnerController creates a new instance of WinnerController
func NewWinnerController(service service.WinnerServiceInterface) *WinnerController {
	return &WinnerController{WinnerService: service}
}

// GetWinnerListHandler retrieves a paginated list of Winners
func (cc *WinnerController) GetWinnerListHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")
	status := r.URL.Query().Get("status")
	tg_group := r.Header.Get("tg_group")

	// Parse page number
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Prepare request DTO
	requestDto := utils.PaginationRequestFilterDTO{
		Page:    page,
		Limit:   limit,
		Query:   query,
		Status:  status,
		TGgroup: tg_group,
	}

	// Call the service to get the list of Winners
	WinnerListResponse := cc.WinnerService.WinnerServiceGetList(requestDto)

	// Send the response back to the client
	utils.NewHttpPaginationResponse(w, WinnerListResponse.Data, WinnerListResponse.Meta.Total, WinnerListResponse.Meta.Page, WinnerListResponse.Meta.LastPage, int(WinnerListResponse.Status), WinnerListResponse.Message)
}

// WinnerCreateHandler creates a new Winner
func (cc *WinnerController) WinnerCreateHandler(w http.ResponseWriter, r *http.Request) {
	var createWinnerDTO dto.CreateWinnerDTO
	err := json.NewDecoder(r.Body).Decode(&createWinnerDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	response, err := cc.WinnerService.WinnerServiceCreate(createWinnerDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// WinnerUpdateHandler updates an existing Winner
func (cc *WinnerController) WinnerUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var updateWinnerDTO dto.UpdateWinnerDTO
	err := json.NewDecoder(r.Body).Decode(&updateWinnerDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}
	response, err := cc.WinnerService.WinnerServiceUpdate(updateWinnerDTO)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// WinnerGetByIDHandler retrieves a Winner by ID
func (cc *WinnerController) WinnerGetByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid ID")
		return
	}
	requestDto := dto.WinnerFilterDTO{
		ID: uint(id),
	}
	response, err := cc.WinnerService.WinnerServiceGetById(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// WinnerGetByIDHandler retrieves a Winner by ID
func (cc *WinnerController) WinnerGetByChatIDHandler(w http.ResponseWriter, r *http.Request) {
	chatIdStr := r.URL.Query().Get("chat")

	if chatIdStr == "" {
		utils.HttpSuccessResponse(w, nil, http.StatusBadRequest, "Invalid chat ID")
		return
	}

	requestDto := dto.WinnerFilterDTO{
		ChatID: chatIdStr,
	}
	response, err := cc.WinnerService.WinnerServiceGetByChatID(requestDto)
	if err != nil {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, err.Error())
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// Task Excute GET Winner
func (cc *WinnerController) WinnerTaskExcuteHandler(w http.ResponseWriter, r *http.Request) {
	response := cc.WinnerService.WinnerServiceTaskExcuteGetWiner()
	if !response {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, "False WinnerTaskExcuteHandler ")
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

// Task Excute Reset point
func (cc *WinnerController) TaskExcuteResetPointHandler(w http.ResponseWriter, r *http.Request) {
	response := cc.WinnerService.ServiceTaskExcuteResetPoint()
	if !response {
		utils.HttpSuccessResponse(w, nil, http.StatusNotFound, "False WinnerTaskExcuteHandler ")
		return
	}

	utils.HttpSuccessResponse(w, response, http.StatusOK, string(utils.SuccessMessage))
}

func (cc *WinnerController) TaskSendNotificationWinnerHandler_(w http.ResponseWriter, r *http.Request) {
	token := "5692442281:AAFKWDPe1O4t1nTvLIxXUeREt3foywstcIY" // Replace with the correct token for @mytgbooter_bot
	chatID := int64(-1002617162989)
	message := "🎉 <b>Congratulations to new lucky winner!</b>\n"

	var winnerDto dto.WinnerFilterIDsDTO
	err := json.NewDecoder(r.Body).Decode(&winnerDto)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}

	winners, err := cc.WinnerService.ServiceGetWinner(winnerDto)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	for _, winner := range winners {
		points := strconv.FormatFloat(winner.TotalPoints, 'f', 2, 64)
		id, err := strconv.ParseInt(winner.ChatID, 10, 64)
		if err != nil {
			fmt.Printf("Invalid ChatID for user %s: %v\n", winner.Member.Username, err)
			continue
		}
		message += fmt.Sprintf(" - <b>%s</b> <a href=\"tg://user?id=%d\">@%s</a> [🎁 %s]\n",
			winner.Member.AccountName,
			id,
			winner.Member.Username,
			points,
		)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	telegramMessage := map[string]interface{}{
		"chat_id":    chatID,
		"text":       message,
		"parse_mode": "HTML",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": "🌐 Visit Website", // Temporary workaround
						"url":  "https://track.igflexs.com/",
					},
				},
			},
		},
	}

	messageData, err := json.Marshal(telegramMessage)
	if err != nil {
		fmt.Printf("JSON Marshal error: %v\n", err)
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(messageData))
	if err != nil {
		fmt.Printf("HTTP Post error: %v\n", err)
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			utils.HttpSuccessResponse(w, "Failed to read Telegram API response", http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
			return
		}
		fmt.Printf("Telegram API error: %s\n", string(body))
		utils.HttpSuccessResponse(w, string(body), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	utils.HttpSuccessResponse(w, nil, http.StatusOK, string(utils.SuccessMessage))
}
func (cc *WinnerController) TaskSendNotificationWinnerHandler(w http.ResponseWriter, r *http.Request) {
	//token := "5692442281:AAFKWDPe1O4t1nTvLIxXUeREt3foywstcIY"
	//chatID := int64(-1002617162989)
	message := "🎉 <b>Congratulations to new lucky winner!</b>\n"

	var winnerDto dto.WinnerFilterIDsDTO
	err := json.NewDecoder(r.Body).Decode(&winnerDto)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}
	utils.WarnLog(winnerDto, "dto.WinnerFilterIDsDTO")
	settingDto := dto.TelegramSettingWinnerFilterDTO{
		GroupID: &winnerDto.ChatID,
	}
	luckysetting, errsw := cc.WinnerService.LuckySettingService(settingDto)
	if errsw != nil {
		utils.ErrorLog(errsw, "error  cc.WinnerService.LuckySettingService")
		utils.HttpSuccessResponse(w, errsw.Error(), http.StatusBadRequest, string(utils.ErrorMessage))
		return
	}
	winnerDto.Limit = int(luckysetting.TotalWin)
	winners, err := cc.WinnerService.ServiceGetWinner(winnerDto)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	for _, winner := range winners {
		points := strconv.FormatFloat(winner.TotalPoints, 'f', 2, 64)
		id, err := strconv.ParseInt(winner.ChatID, 10, 64)
		if err != nil {
			fmt.Printf("Invalid ChatID for user %s: %v\n", winner.Member.Username, err)
			continue
		}
		message += fmt.Sprintf(" - <b>%s</b> <a href=\"tg://user?id=%d\">@%s</a> 🎁 %s\n",
			winner.Member.AccountName,
			id,
			winner.Member.Username,
			points,
		)
	}
	message += "\n " + luckysetting.Message

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendAnimation", luckysetting.Token)
	telegramPayload := map[string]interface{}{
		"chat_id":    luckysetting.GroupID,
		"animation":  luckysetting.Image,
		"caption":    message,
		"parse_mode": "HTML",
		"reply_markup": map[string]interface{}{
			"inline_keyboard": [][]map[string]interface{}{
				{
					{
						"text": "🎁 Check Balance",
						"url":  "t.me/XYan_LuckyDraw_bot/Binger",
					},
				},
			},
		},
	}

	messageData, err := json.Marshal(telegramPayload)
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(messageData))
	if err != nil {
		utils.HttpSuccessResponse(w, err.Error(), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		utils.HttpSuccessResponse(w, string(body), http.StatusInternalServerError, "False TaskSendNotificationWinnerHandler")
		return
	}

	utils.HttpSuccessResponse(w, nil, http.StatusOK, string(utils.SuccessMessage))
}
