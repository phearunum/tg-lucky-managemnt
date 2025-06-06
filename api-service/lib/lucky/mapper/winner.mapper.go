package mapper

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/models"
)

// ToWinnerModel maps CreateWinnerDTO to Winner model
func ToWinnerModel(dto dto.CreateWinnerDTO) models.LuckWinner {
	return models.LuckWinner{

		ChatID:      dto.ChatID,
		TotalPoints: dto.TotalPoints,
		Status:      dto.Status,
		Note:        dto.Note,
		CreateBy:    dto.CreateBy,
	}
}

// ToResponseWinnerDTO maps Winner model to ResponseWinnerDTO
func ToResponseWinnerDTO(Winner models.LuckWinner) dto.ResponseWinnerDTO {
	return dto.ResponseWinnerDTO{
		ID:     Winner.ID,
		ChatID: Winner.ChatID,
		Member: dto.ResponseTelegramUserDTO{ // Using the named struct
			AccountName: Winner.Member.AccountName,
			Username:    Winner.Member.Username,
			EmployeeID:  Winner.Member.EmployeeID,
			Profile:     Winner.Member.Profile,
			TotalPoints: Winner.Member.TotalPoints,
		},
		TotalPoints: Winner.TotalPoints,
		Status:      Winner.Status,
		Note:        Winner.Note,
		Exprired:    Winner.Exprired,
		CreateBy:    Winner.CreateBy,
		CreatedAt:   Winner.CreatedAt,
		UpdatedAt:   Winner.UpdatedAt,
	}
}

// ApplyUpdateToWinner applies UpdateWinnerDTO to an existing Winner model
func ApplyUpdateToWinner(Winner *models.LuckWinner, dto dto.UpdateWinnerDTO) {
	if dto.Note != nil {
		Winner.Note = *dto.Note
	}

	if dto.ChatID != nil {
		Winner.ChatID = *dto.ChatID
	}
	if dto.Status != nil {
		Winner.Status = *dto.Status
	}
	if dto.CreateBy != nil {
		Winner.CreateBy = *dto.CreateBy
	}
}
