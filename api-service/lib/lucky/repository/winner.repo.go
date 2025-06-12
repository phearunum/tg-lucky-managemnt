package repository

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/mapper"
	"api-service/lib/lucky/models"
	"api-service/utils"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WinnerRepository struct {
	db *gorm.DB
}

// NewWinnerRepository creates a new instance of WinnerRepository
func NewWinnerRepository(db *gorm.DB) *WinnerRepository {
	return &WinnerRepository{db: db}
}

// GetWinnerList retrieves a paginated list of Winners
func (r *WinnerRepository) GetWinnerList(page, limit int, query string, status string, tg_group string) ([]*models.LuckWinner, int, error) {
	var points []*models.LuckWinner
	var total int64
	offset := (page - 1) * limit
	var groupIDs []int64
	for _, idStr := range strings.Split(tg_group, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			groupIDs = append(groupIDs, id)
		}
	}
	//baseQuery := r.db.Model(&models.LuckWinner{}).Preload("Member").Order("created_at DESC")
	baseQuery := r.db.Model(&models.LuckWinner{}).
		Joins("JOIN telegram_users ON telegram_users.chat_id = telegram_winner.chat_id").
		Preload("Member").
		Where("telegram_users.group_id IN (?) AND telegram_winner.status = ?", groupIDs, status).
		Order("telegram_winner.created_at DESC")

	if query != "" {
		baseQuery = baseQuery.Where(
			r.db.Where("telegram_winner.account_name LIKE ?", "%"+query+"%").
				Or("telegram_winner.chat_id LIKE ?", "%"+query+"%"),
		)
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := baseQuery.Offset(offset).Limit(limit).Find(&points).Error; err != nil {
		return nil, 0, err
	}

	return points, int(total), nil
}

// CreateWinner adds a new Winner
func (r *WinnerRepository) CreateWinner(createDTO dto.CreateWinnerDTO) (*models.LuckWinner, error) {
	point := mapper.ToWinnerModel(createDTO)

	if err := r.db.Save(&point).Error; err != nil {
		utils.LoggerRepository(err, "CreateWinner: Insert failed")
		return nil, err
	}
	return &point, nil
}

// UpdateWinner modifies an existing Winner
func (r *WinnerRepository) UpdateWinner(updateDTO dto.UpdateWinnerDTO) (*models.LuckWinner, error) {
	existingPoint, err := r.GetWinnerByID(updateDTO.ID)
	if err != nil {
		utils.LoggerRepository(err, "UpdateWinner: Fetch failed")
		return nil, err
	}
	mapper.ApplyUpdateToWinner(existingPoint, updateDTO)

	if err := r.db.Save(existingPoint).Error; err != nil {
		utils.LoggerRepository(err, "UpdateWinner: Update failed")
		return nil, err
	}

	return existingPoint, nil
}

// DeleteWinner removes a Winner by ID
func (r *WinnerRepository) DeleteWinner(id uint) (bool, error) {
	if _, err := r.GetWinnerByID(id); err != nil {
		utils.LoggerRepository(err, "DeleteWinner: Fetch failed")
		return false, err
	}

	if err := r.db.Where("id = ?", id).Delete(&models.LuckWinner{}).Error; err != nil {
		utils.LoggerRepository(err, "DeleteWinner: Delete failed")
		return false, err
	}
	return true, nil
}

// GetWinnerByID retrieves a Winner by ID
func (r *WinnerRepository) GetWinnerByID(id uint) (*models.LuckWinner, error) {
	var point models.LuckWinner
	err := r.db.Where("id = ?", id).First(&point).Error
	if err != nil {
		utils.LoggerRepository(err, "GetWinnerByID: Fetch failed")
		return nil, err
	}
	return &point, nil
}
func (r *WinnerRepository) GetChanceCollection(chat_id string) ([]*models.LuckWinner, error) {
	var points []*models.LuckWinner
	err := r.db.Where("chat_id = ? AND status = 'available'", chat_id).Find(&points).Error
	if err != nil {
		utils.LoggerRepository(err, "GetChanceCollection(chat_id string): Fetch failed")
		return nil, err
	}
	return points, nil
}

func (r *WinnerRepository) TaskExcuteWinner() bool {
	var setting models.TelegramSettingWinner // Use a value, not a pointer, unless required
	err := r.db.Order("created_at DESC").First(&setting).Error
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no settings exist (e.g., use defaults)
			setting = models.TelegramSettingWinner{TotalWin: 5} // Default value
			utils.InfoLog("No settings found, using default TotalWin=5", "TaskExcuteWinner")
		} else {
			// Log other database errors and fail
			utils.LoggerRepository(err, "TaskExcuteWinner: Failed to fetch Setting Lucky Winner")
			return false
		}
	}
	var members []*models.TelegramUsers
	utils.InfoLog(setting, "setting")
	if setting.TotalWin < 0 {
		setting.TotalWin = 10 // Set a default if negative
	}
	errr := r.db.Order("total_point DESC").Limit(int(setting.TotalWin)).Find(&members).Error
	if err != nil {
		utils.LoggerRepository(errr, "TaskExcuteWinner: Failed to fetch TelegramUsers")
		return false
	}

	if len(members) == 0 {
		utils.LoggerRepository(nil, "TaskExcuteWinner: No TelegramUsers found")
		return false
	}

	var winners []models.LuckWinner
	today := time.Now()
	rankDate := today.Format("20060102")
	for i, m := range members {
		winner := models.LuckWinner{
			ChatID:      m.ChatID,
			TotalPoints: m.TotalPoints,
			Exprired:    "no",
			Status:      "active",
			RankID:      fmt.Sprintf("%s%d", rankDate, i),
		}
		if m.TotalPoints > 0 {
			winners = append(winners, winner)
		}

	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "rank_id"}},
		UpdateAll: true,
	}).Create(&winners).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return false
	}

	return true
}
func (r *WinnerRepository) TaskExcuteWinnerTask() bool {
	var settings []models.TelegramSettingWinner
	err := r.db.
		Where("status = ?", "collected").
		Order("created_at DESC").
		Find(&settings).Error

	if err != nil {
		utils.LoggerRepository(err, "TaskExcuteWinner: Failed to fetch Setting Lucky Winner")
		return false
	}

	if len(settings) == 0 {
		utils.InfoLog("No settings found, using default TotalWin=5", "TaskExcuteWinner")
		return false
	}

	utils.InfoLog(settings, "setting")
	for _, s := range settings {
		utils.InfoLog(s, fmt.Sprintf("Processing Setting for Group: %s (TotalWin: %s)", s.GroupID, s.TotalWin))

		var members []*models.TelegramUsers
		//utils.InfoLog(s, "Setting Group")
		errr := r.db.
			Where("group_id = ?", s.GroupID).
			Order("total_point DESC").
			Limit(int(s.TotalWin)).
			Find(&members).Error

		if errr != nil {
			utils.LoggerRepository(errr, "TaskExcuteWinner: Failed to fetch TelegramUsers")
			continue
		}

		if len(members) == 0 {
			utils.LoggerRepository(nil, "TaskExcuteWinner: No TelegramUsers found")
			continue
		}

		var winners []models.LuckWinner
		today := time.Now()
		rankDate := today.Format("20060102")

		for i, m := range members {
			if m.TotalPoints > 0 {
				winner := models.LuckWinner{
					ChatID:      m.ChatID,
					TotalPoints: m.TotalPoints,
					Exprired:    "no",
					Status:      "active",
					RankID:      fmt.Sprintf("%s_%s_%d", rankDate, s.GroupID, i),
				}
				winners = append(winners, winner)
			}
		}

		if err := r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "rank_id"}},
			UpdateAll: true,
		}).Create(&winners).Error; err != nil {
			utils.LoggerRepository(err, "Execute")
			continue
		}
	}

	return true
}

func (r *WinnerRepository) TaskExcutePointReset() bool {
	// Reset chance to exprired
	err := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Model(&models.ChancePoint{}).
		Updates(map[string]interface{}{
			//"status":  "expired",
			"expired": "yes",
		}).Error
	if err != nil {
		utils.LoggerRepository(err, "TaskExcuteChancePointReset: Failed to reset points")
		return false
	}
	// Execute Reset Point on all Member accounts
	exerr := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}). // Allow global update
									Model(&models.TelegramUsers{}).
									Updates(map[string]interface{}{
			"total_chance": 0,
			"total_point":  0,
		}).Error

	if exerr != nil {
		utils.LoggerRepository(exerr, "TaskExcuteReset Point on Account: Failed to reset points")
		return false
	}

	// Reset Winner
	winerr := r.db.Session(&gorm.Session{AllowGlobalUpdate: true}). // Allow global update
									Model(&models.LuckWinner{}).
									Updates(map[string]interface{}{
			"status":  "false",
			"expired": "yes",
		}).Error

	if winerr != nil {
		utils.LoggerRepository(winerr, "TaskExcuteReset Winner Account: Failed to reset points")
		return false
	}

	return true
}

func (r *WinnerRepository) GetWinnerToday(ids []int, limit int) ([]*models.LuckWinner, error) {
	var winners []*models.LuckWinner

	// Ensure ids is not empty
	if len(ids) == 0 {
		return winners, nil
	}
	if limit == 0 {
		limit = 1
	}
	err := r.db.
		Model(&models.LuckWinner{}).Preload("Member").Order("total_point DESC").
		Where("id IN (?) and expired='no'", ids).
		Find(&winners).Limit(limit).Error
	if err != nil {
		return nil, err
	}
	return winners, nil
}
func (r *WinnerRepository) GetLuckyWinnerSetting(id string) (*models.TelegramSettingWinner, error) {
	var setting models.TelegramSettingWinner
	err := r.db.Where("group_id = ?", id).First(&setting).Error
	if err != nil {
		utils.LoggerRepository(err, "GetWinnerByID: Fetch failed")
		return nil, err
	}
	return &setting, nil
}
