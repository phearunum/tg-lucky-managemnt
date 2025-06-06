package dto

import (
	"time"
)

// CreateWinnerDTO represents the DTO for creating a Winner entry
type CreateWinnerDTO struct {
	ChatID      string  `json:"chat_id"`
	Status      string  `json:"status"`
	TotalPoints float64 `json:"total_point"`
	Note        string  `json:"note"`
	CreateBy    string  `json:"create_by"`
	TotalChance float64 `json:"total_chance"`
}

// UpdateWinnerDTO represents the DTO for updating a Winner entry
type UpdateWinnerDTO struct {
	ID uint `json:"id"`

	ChatID      *string  `json:"chat_id,omitempty"`
	Status      *string  `json:"status,omitempty"`
	TotalPoints float64  `json:"total_point"`
	Note        *string  `json:"note,omitempty"`
	CreateBy    *string  `json:"create_by,omitempty"`
	TotalChance *float64 `json:"total_chance,omitempty"`
}

// ResponseWinnerDTO represents the DTO for responding with a Winner entry
type ResponseWinnerDTO struct {
	ID          uint                    `json:"id"`
	Member      ResponseTelegramUserDTO `json:"member"`
	ChatID      string                  `json:"chat_id"`
	TotalPoints float64                 `json:"total_point"`
	Status      string                  `json:"status"`
	Note        string                  `json:"note"`
	Exprired    string                  `json:"expired"`
	CreateBy    string                  `json:"create_by"`
	CreatedAt   *time.Time              `json:"created_at"`
	UpdatedAt   *time.Time              `json:"updated_at"`
}

// WinnerFilterDTO represents the DTO for filtering Winner entries
type WinnerFilterDTO struct {
	ID uint `json:"id"`

	ChatID      string  `json:"chat_id"`
	TotalPoints float64 `json:"total_point"`
	Status      *string `json:"status,omitempty"`
	CreateBy    *string `json:"create_by,omitempty"`
	Exprired    string  `json:"expired"`
}
type WinnerFilterIDsDTO struct {
	ID          []int   `json:"id"`
	ChatID      string  `json:"chat_id"`
	TotalPoints float64 `json:"total_point"`
	Status      *string `json:"status,omitempty"`
	CreateBy    *string `json:"create_by,omitempty"`
	Exprired    string  `json:"expired"`
	Img         string  `json:"img"`
	Token       string  `json:"token"`
	Limit       int     `json:"limit"`
}
