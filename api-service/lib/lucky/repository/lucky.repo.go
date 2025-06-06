package repository

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/mapper"
	"api-service/lib/lucky/models"
	"api-service/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type TelegramUserRepository struct {
	db *gorm.DB
}

// NewTelegramUserRepository creates a new instance of TelegramUserRepository
func NewTelegramUserRepository(db *gorm.DB) *TelegramUserRepository {
	return &TelegramUserRepository{db: db}
}

// GetTelegramUserList retrieves a paginated list of Telegram users
func (r *TelegramUserRepository) GetTelegramUserList(page, limit int, query string, tg_group string) ([]*models.TelegramUsers, int, error) {
	var users []*models.TelegramUsers
	var total int64
	offset := (page - 1) * limit

	var groupIDs []int64
	for _, idStr := range strings.Split(tg_group, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			groupIDs = append(groupIDs, id)
		}
	}
	// Build base query
	baseQuery := r.db.Model(&models.TelegramUsers{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("account_name LIKE ? OR username LIKE ? AND group_id in (?)", "%"+query+"%", "%"+query+"%", tg_group)
	}
	baseQuery = baseQuery.Where("group_id in (?)", groupIDs)

	// Use transactions for better performance when counting and fetching
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// CreateTelegramUser adds a new Telegram user
func (r *TelegramUserRepository) CreateTelegramUser(createDTO dto.CreateTelegramUserDTO) (*models.TelegramUsers, error) {
	user := mapper.ToTelegramUserModel(createDTO)

	// Insert record (pass pointer)
	if err := r.db.Create(&user).Error; err != nil {
		utils.LoggerRepository(err, "CreateTelegramUser: Insert failed")
		return nil, err
	}
	return &user, nil
}

// UpdateTelegramUser modifies an existing Telegram user
func (r *TelegramUserRepository) UpdateTelegramUser(updateDTO dto.UpdateTelegramUserDTO) (*models.TelegramUsers, error) {
	existingUser, err := r.GetTelegramUserByID(updateDTO.ID)
	if err != nil {
		utils.LoggerRepository(err, "UpdateTelegramUser: Fetch failed")
		return nil, err
	}
	mapper.ApplyUpdateToTelegramUser(existingUser, updateDTO)

	if err := r.db.Save(existingUser).Error; err != nil {
		utils.LoggerRepository(err, "UpdateTelegramUser: Update failed")
		return nil, err
	}

	return existingUser, nil
}

// DeleteTelegramUser removes a Telegram user by ID
func (r *TelegramUserRepository) DeleteTelegramUser(id uint) (bool, error) {
	// Check if the user exists
	if _, err := r.GetTelegramUserByID(id); err != nil {
		utils.LoggerRepository(err, "DeleteTelegramUser: Fetch failed")
		return false, err
	}
	// Delete the user from the database
	if err := r.db.Where("id = ?", id).Delete(&models.TelegramUsers{}).Error; err != nil {
		utils.LoggerRepository(err, "DeleteTelegramUser: Delete failed")
		return false, err
	}
	return true, nil
}

// GetTelegramUserByID retrieves a Telegram user by ID
func (r *TelegramUserRepository) GetTelegramUserByID(id uint) (*models.TelegramUsers, error) {
	var user models.TelegramUsers
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramUserByID: Fetch failed")
		return nil, err
	}
	return &user, nil
}
func (r *TelegramUserRepository) GetTelegramUserByChatID(chatId string) (*models.TelegramUsers, error) {
	var user models.TelegramUsers
	err := r.db.Where("chat_id = ?", chatId).First(&user).Error
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramUserByChatID: Fetch failed")
		return nil, err
	}
	return &user, nil
}

// Get Lucky Bot List
func (r *TelegramUserRepository) GetLuckyBotList(page, limit int, query string) ([]*models.TelegramSettingWinner, int, error) {
	var users []*models.TelegramSettingWinner
	var total int64
	offset := (page - 1) * limit
	baseQuery := r.db.Model(&models.TelegramSettingWinner{}).Order("created_at DESC")
	if query != "" {
		baseQuery = baseQuery.Where("group_id LIKE ? ", "%"+query+"%")
	}
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}
