package mapper

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/models"
	"time"
)

// ToTelegramSettingWinnerModel maps CreateTelegramSettingWinnerDTO to TelegramSettingWinner model
func ToTelegramSettingWinnerModel(dto dto.CreateTelegramSettingWinnerDTO) models.TelegramSettingWinner {
	return models.TelegramSettingWinner{
		Name:     dto.Name,
		TotalWin: dto.TotalWin,
		MaxPoint: dto.MaxPoint,
		MinPoint: dto.MinPoint,
		GroupID:  dto.GroupID,
		Token:    dto.Token,
		Message:  dto.Message,
		Image:    dto.Image,
		Type:     dto.Type,
		Status:   dto.Status,
	}
}

// ToResponseTelegramSettingWinnerDTO maps TelegramSettingWinner model to ResponseTelegramSettingWinnerDTO
func ToResponseTelegramSettingWinnerDTO(setting models.TelegramSettingWinner) dto.ResponseTelegramSettingWinnerDTO {
	var createdAt time.Time
	var updatedAt time.Time

	if setting.CreatedAt != nil {
		createdAt = *setting.CreatedAt
	}

	if setting.UpdatedAt != nil {
		updatedAt = *setting.UpdatedAt
	}

	return dto.ResponseTelegramSettingWinnerDTO{
		ID:        setting.ID,
		Name:      setting.Name,
		TotalWin:  setting.TotalWin,
		MaxPoint:  setting.MaxPoint,
		MinPoint:  setting.MinPoint,
		GroupID:   setting.GroupID,
		Token:     setting.Token,
		Message:   setting.Message,
		Image:     setting.Image,
		Type:      setting.Type,
		Status:    setting.Status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// ApplyUpdateToTelegramSettingWinner updates an existing TelegramSettingWinner model with values from UpdateTelegramSettingWinnerDTO
func ApplyUpdateToTelegramSettingWinner(setting *models.TelegramSettingWinner, dto dto.UpdateTelegramSettingWinnerDTO) {
	if dto.TotalWin != nil {
		setting.TotalWin = *dto.TotalWin

	}
	if dto.Name != "" {
		setting.Name = dto.Name
	}
	if dto.MaxPoint != nil {
		setting.MaxPoint = *dto.MaxPoint
	}
	if dto.MinPoint != nil {
		setting.MinPoint = *dto.MinPoint
	}
	if dto.GroupID != nil {
		setting.GroupID = *dto.GroupID
	}
	if dto.Token != nil {
		setting.Token = *dto.Token
	}
	if dto.Message != nil {
		setting.Message = *dto.Message
	}
	if dto.Image != nil {
		setting.Image = *dto.Image
	}
	if dto.Type != nil {
		setting.Type = *dto.Type
	}
	if dto.Status != nil {
		setting.Status = *dto.Status
	}
}
