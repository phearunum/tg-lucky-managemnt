package mapper

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/models"
)

// Mapper functions
func ToTelegramUserModel(dto dto.CreateTelegramUserDTO) models.TelegramUsers {
	return models.TelegramUsers{
		Username:    dto.Username,
		AccountName: dto.AccountName,
		EmployeeID:  dto.EmployeeID,
		Profile:     dto.Profile,
		ChatID:      dto.ChatID,
		TotalChance: dto.TotalChance,
		TotalPoints: dto.TotalPoints,
		BotName:     dto.BotName,
		GroupID:     dto.GroupID,
	}
}

func ToResponseTelegramUserDTO(user models.TelegramUsers) dto.ResponseTelegramUserDTO {
	return dto.ResponseTelegramUserDTO{
		ID:          user.ID,
		Username:    user.Username,
		AccountName: user.AccountName,
		EmployeeID:  user.EmployeeID,
		Profile:     user.Profile,
		ChatID:      user.ChatID,
		TotalChance: user.TotalChance,
		TotalPoints: user.TotalPoints,
		BotName:     user.BotName,
		GroupID:     user.GroupID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func ApplyUpdateToTelegramUser(user *models.TelegramUsers, dto dto.UpdateTelegramUserDTO) {
	if dto.AccountName != nil {
		user.AccountName = *dto.AccountName
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.EmployeeID != nil {
		user.EmployeeID = *dto.EmployeeID
	}
	if dto.Profile != nil {
		user.Profile = *dto.Profile
	}
	if dto.ChatID != nil {
		user.ChatID = *dto.ChatID
	}
	if dto.TotalChance != nil {
		user.TotalChance = *dto.TotalChance
	}
	if dto.TotalPoints != nil {
		user.TotalPoints = *dto.TotalPoints
	}
}
