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

type ChancePointRepository struct {
	db *gorm.DB
}

// NewChancePointRepository creates a new instance of ChancePointRepository
func NewChancePointRepository(db *gorm.DB) *ChancePointRepository {
	return &ChancePointRepository{db: db}
}

// GetChancePointList retrieves a paginated list of ChancePoints
func (r *ChancePointRepository) GetChancePointList(page, limit int, query string, tg_group string) ([]*models.ChancePoint, int, error) {
	var points []*models.ChancePoint
	var total int64
	offset := (page - 1) * limit

	//utils.InfoLog(tg, "TG Group Id")
	var groupIDs []int64
	for _, idStr := range strings.Split(tg_group, ",") {
		idStr = strings.TrimSpace(idStr)
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			groupIDs = append(groupIDs, id)
		}
	}
	baseQuery := r.db.Model(&models.ChancePoint{}).
		Joins("JOIN telegram_users ON telegram_users.chat_id = telegram_chance.chat_id").
		Preload("Member").
		Where("telegram_users.group_id IN (?)", groupIDs).
		Order("telegram_chance.created_at DESC")

	if query != "" {
		baseQuery = baseQuery.Where(
			r.db.Where("telegram_chance.account_name LIKE ?", "%"+query+"%").
				Or("telegram_chance.chat_id LIKE ?", "%"+query+"%"),
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

// CreateChancePoint adds a new ChancePoint
func (r *ChancePointRepository) CreateChancePoint(createDTO dto.CreateChancePointDTO) (*models.ChancePoint, error) {
	point := mapper.ToChancePointModel(createDTO)

	if err := r.db.Save(&point).Error; err != nil {
		utils.LoggerRepository(err, "CreateChancePoint: Insert failed")
		return nil, err
	}
	return &point, nil
}

// UpdateChancePoint modifies an existing ChancePoint
func (r *ChancePointRepository) UpdateChancePoint(updateDTO dto.UpdateChancePointDTO) (*models.ChancePoint, error) {
	existingPoint, err := r.GetChancePointByID(updateDTO.ID)
	if err != nil {
		utils.LoggerRepository(err, "UpdateChancePoint: Fetch failed")
		return nil, err
	}
	mapper.ApplyUpdateToChancePoint(existingPoint, updateDTO)

	if err := r.db.Save(existingPoint).Error; err != nil {
		utils.LoggerRepository(err, "UpdateChancePoint: Update failed")
		return nil, err
	}

	return existingPoint, nil
}

// DeleteChancePoint removes a ChancePoint by ID
func (r *ChancePointRepository) DeleteChancePoint(id uint) (bool, error) {
	if _, err := r.GetChancePointByID(id); err != nil {
		utils.LoggerRepository(err, "DeleteChancePoint: Fetch failed")
		return false, err
	}

	if err := r.db.Where("id = ?", id).Delete(&models.ChancePoint{}).Error; err != nil {
		utils.LoggerRepository(err, "DeleteChancePoint: Delete failed")
		return false, err
	}
	return true, nil
}

// GetChancePointByID retrieves a ChancePoint by ID
func (r *ChancePointRepository) GetChancePointByID(id uint) (*models.ChancePoint, error) {
	var point models.ChancePoint
	err := r.db.Where("id = ?", id).First(&point).Error
	if err != nil {
		utils.LoggerRepository(err, "GetChancePointByID: Fetch failed")
		return nil, err
	}
	return &point, nil
}
func (r *ChancePointRepository) GetChanceCollection(chat_id string) ([]*models.ChancePoint, error) {
	var points []*models.ChancePoint
	err := r.db.Where("chat_id = ? AND status = 'available'", chat_id).Find(&points).Error
	if err != nil {
		utils.LoggerRepository(err, "GetChanceCollection(chat_id string): Fetch failed")
		return nil, err
	}
	return points, nil
}
