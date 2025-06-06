package dto

import (
	"time"
)

// CreateTelegramUserDTO represents the DTO for creating a Telegram user
type CreateTelegramUserDTO struct {
	Username    string  `json:"username"`
	AccountName string  `json:"account_name"`
	EmployeeID  string  `json:"employee_id"`
	Profile     string  `json:"profile"`
	ChatID      string  `json:"chat_id"`
	TotalChance float64 `json:"total_chance"`
	TotalPoints float64 `json:"total_point"`
	BotName     string  `json:"bot_name"`
	GroupID     string  `json:"group_id"`
}

// UpdateTelegramUserDTO represents the DTO for updating a Telegram user
type UpdateTelegramUserDTO struct {
	ID          uint     `json:"id"`
	Username    *string  `json:"username"`
	AccountName *string  `json:"account_name,omitempty"`
	EmployeeID  *string  `json:"employee_id,omitempty"`
	Profile     *string  `json:"profile,omitempty"`
	ChatID      *string  `json:"chat_id,omitempty"`
	TotalChance *float64 `json:"total_chance,omitempty"`
	TotalPoints *float64 `json:"total_point,omitempty"`
	BotName     string   `json:"bot_name"`
	GroupID     string   `json:"group_id"`
}

// ResponseTelegramUserDTO represents the DTO for responding with a Telegram user
type ResponseTelegramUserDTO struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	AccountName string    `json:"account_name"`
	EmployeeID  string    `json:"employee_id"`
	Profile     string    `json:"profile"`
	ChatID      string    `json:"chat_id"`
	TotalChance float64   `json:"total_chance"`
	TotalPoints float64   `json:"total_point"`
	BotName     string    `json:"bot_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     string    `json:"group_id"`
}

type TelegramUserFilterDTO struct {
	ID          uint     `json:"id"`
	AccountName *string  `json:"account_name,omitempty"`
	EmployeeID  *string  `json:"employee_id,omitempty"`
	Profile     *string  `json:"profile,omitempty"`
	BotName     string   `json:"bot_name"`
	ChatID      *string  `json:"chat_id,omitempty"`
	GroupID     string   `json:"group_id"`
	TotalChance *float64 `json:"total_chance,omitempty"`
	TotalPoints *float64 `json:"total_point,omitempty"`
}
