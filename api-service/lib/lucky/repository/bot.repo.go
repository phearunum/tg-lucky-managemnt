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

type TelegramSettingWinnerRepository struct {
	db *gorm.DB
}

// NewTelegramSettingWinnerRepository creates a new instance of TelegramSettingWinnerRepository
func NewTelegramSettingWinnerRepository(db *gorm.DB) *TelegramSettingWinnerRepository {
	return &TelegramSettingWinnerRepository{db: db}
}

// GetTelegramSettingWinnerList retrieves a paginated list of Telegram setting winners
func (r *TelegramSettingWinnerRepository) GetTelegramSettingWinnerList(page, limit int, query string, tg_group string) ([]*models.TelegramSettingWinner, int, error) {
	var settings []*models.TelegramSettingWinner
	var total int64
	offset := (page - 1) * limit
	var groupIDs []int64
	for _, idStr := range strings.Split(tg_group, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			groupIDs = append(groupIDs, id)
		}
	}
	baseQuery := r.db.Model(&models.TelegramSettingWinner{}).Where("group_id IN (?)", groupIDs).Order("created_at DESC")
	/*baseQuery := r.db.Model(&models.LuckWinner{}).
	Joins("JOIN telegram_users ON telegram_users.chat_id = telegram_win_setting.chat_id").
	Preload("Member").
	Where("telegram_users.group_id IN (?) ", groupIDs).
	Order("telegram_win_setting.created_at DESC")*/
	if query != "" {
		//baseQuery = baseQuery.Where("message LIKE ? OR group_id LIKE ?", "%"+query+"%", "%"+query+"%")
		baseQuery = baseQuery.Where(
			r.db.Where("telegram_win_setting.name LIKE ?", "%"+query+"%").
				Or("telegram_win_setting.group_id LIKE ?", "%"+query+"%"),
		)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&settings).Error; err != nil {
		return nil, 0, err
	}

	return settings, int(total), nil
}

// CreateTelegramSettingWinner adds a new Telegram setting winner
func (r *TelegramSettingWinnerRepository) CreateTelegramSettingWinner(createDTO dto.CreateTelegramSettingWinnerDTO) (*models.TelegramSettingWinner, error) {
	setting := mapper.ToTelegramSettingWinnerModel(createDTO)

	if err := r.db.Create(&setting).Error; err != nil {
		utils.LoggerRepository(err, "CreateTelegramSettingWinner: Insert failed")
		return nil, err
	}
	return &setting, nil
}

// UpdateTelegramSettingWinner modifies an existing Telegram setting winner
func (r *TelegramSettingWinnerRepository) UpdateTelegramSettingWinner(updateDTO dto.UpdateTelegramSettingWinnerDTO) (*models.TelegramSettingWinner, error) {
	existingSetting, err := r.GetTelegramSettingWinnerByID(updateDTO.ID)
	if err != nil {
		utils.LoggerRepository(err, "UpdateTelegramSettingWinner: Fetch failed")
		return nil, err
	}
	mapper.ApplyUpdateToTelegramSettingWinner(existingSetting, updateDTO)

	if err := r.db.Save(existingSetting).Error; err != nil {
		utils.LoggerRepository(err, "UpdateTelegramSettingWinner: Update failed")
		return nil, err
	}

	return existingSetting, nil
}

// DeleteTelegramSettingWinner removes a setting by ID
func (r *TelegramSettingWinnerRepository) DeleteTelegramSettingWinner(id uint) (bool, error) {
	if _, err := r.GetTelegramSettingWinnerByID(id); err != nil {
		utils.LoggerRepository(err, "DeleteTelegramSettingWinner: Fetch failed")
		return false, err
	}

	if err := r.db.Where("id = ?", id).Delete(&models.TelegramSettingWinner{}).Error; err != nil {
		utils.LoggerRepository(err, "DeleteTelegramSettingWinner: Delete failed")
		return false, err
	}
	return true, nil
}

// GetTelegramSettingWinnerByID retrieves a setting by ID
func (r *TelegramSettingWinnerRepository) GetTelegramSettingWinnerByID(id uint) (*models.TelegramSettingWinner, error) {
	var setting models.TelegramSettingWinner
	err := r.db.Where("id = ?", id).First(&setting).Error
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramSettingWinnerByID: Fetch failed")
		return nil, err
	}
	return &setting, nil
}

// GetTelegramSettingWinnerByGroupID retrieves a setting by group ID
func (r *TelegramSettingWinnerRepository) GetTelegramSettingWinnerByGroupID(groupID string) (*models.TelegramSettingWinner, error) {
	var setting models.TelegramSettingWinner
	err := r.db.Where("group_id = ?", groupID).First(&setting).Error
	if err != nil {
		utils.LoggerRepository(err, "GetTelegramSettingWinnerByGroupID: Fetch failed")
		return nil, err
	}
	return &setting, nil
}
