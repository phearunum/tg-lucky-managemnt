package mapper

import (
	"api-service/lib/lucky/dto"
	"api-service/lib/lucky/models"
)

// ToChancePointModel maps CreateChancePointDTO to ChancePoint model
func ToChancePointModel(dto dto.CreateChancePointDTO) models.ChancePoint {
	return models.ChancePoint{
		AccountName: dto.AccountName,
		EmployeeID:  dto.EmployeeID,
		Profile:     dto.Profile,
		ChatID:      dto.ChatID,
		TotalChance: dto.TotalChance,
		Status:      dto.Status,
		Note:        dto.Note,
		CreateBy:    dto.CreateBy,
	}
}

// ToResponseChancePointDTO maps ChancePoint model to ResponseChancePointDTO
func ToResponseChancePointDTO(chancePoint models.ChancePoint) dto.ResponseChancePointDTO {
	return dto.ResponseChancePointDTO{
		ID:          chancePoint.ID,
		AccountName: chancePoint.Member.AccountName,
		EmployeeID:  chancePoint.Member.EmployeeID,
		Profile:     chancePoint.Member.Profile,
		ChatID:      chancePoint.ChatID,
		TotalChance: chancePoint.TotalChance,
		Status:      chancePoint.Status,
		Note:        chancePoint.Note,
		Exprired:    chancePoint.Exprired,
		CreateBy:    chancePoint.CreateBy,
		CreatedAt:   chancePoint.CreatedAt,
		UpdatedAt:   chancePoint.UpdatedAt,
	}
}

// ApplyUpdateToChancePoint applies UpdateChancePointDTO to an existing ChancePoint model
func ApplyUpdateToChancePoint(chancePoint *models.ChancePoint, dto dto.UpdateChancePointDTO) {
	if dto.Note != nil {
		chancePoint.Note = *dto.Note
	}
	if dto.AccountName != nil {
		chancePoint.AccountName = *dto.AccountName
	}

	if dto.EmployeeID != nil {
		chancePoint.EmployeeID = *dto.EmployeeID
	}
	if dto.TotalChance != nil {
		chancePoint.TotalChance = *dto.TotalChance
	}
	if dto.Profile != nil {
		chancePoint.Profile = *dto.Profile
	}
	if dto.ChatID != nil {
		chancePoint.ChatID = *dto.ChatID
	}
	if dto.Status != nil {
		chancePoint.Status = *dto.Status
	}
	if dto.CreateBy != nil {
		chancePoint.CreateBy = *dto.CreateBy
	}
}
