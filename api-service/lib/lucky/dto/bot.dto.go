package dto

import "time"

// CreateTelegramSettingWinnerDTO represents the DTO for creating a Telegram setting winner
type CreateTelegramSettingWinnerDTO struct {
	Name     string  `json:"name"`
	TotalWin float64 `json:"total_win"`
	MaxPoint float64 `json:"max_point"`
	MinPoint float64 `json:"min_point"`
	GroupID  string  `json:"group_id"`
	Token    string  `json:"token"`
	Message  string  `json:"message"`
	Image    string  `json:"image"`
	Type     string  `json:"bot_type"`
	Status   string  `json:"status"`
}

// UpdateTelegramSettingWinnerDTO represents the DTO for updating a Telegram setting winner
type UpdateTelegramSettingWinnerDTO struct {
	Name     string   `json:"name"`
	ID       uint     `json:"id"`
	TotalWin *float64 `json:"total_win,omitempty"`
	MaxPoint *float64 `json:"max_point,omitempty"`
	MinPoint *float64 `json:"min_point,omitempty"`
	GroupID  *string  `json:"group_id,omitempty"`
	Token    *string  `json:"token,omitempty"`
	Message  *string  `json:"message,omitempty"`
	Image    *string  `json:"image,omitempty"`
	Type     *string  `json:"bot_type,omitempty"`
	Status   *string  `json:"status,omitempty"`
}

// ResponseTelegramSettingWinnerDTO represents the DTO for responding with a Telegram setting winner
type ResponseTelegramSettingWinnerDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	TotalWin  float64   `json:"total_win"`
	MaxPoint  float64   `json:"max_point"`
	MinPoint  float64   `json:"min_point"`
	GroupID   string    `json:"group_id"`
	Token     string    `json:"token"`
	Message   string    `json:"message"`
	Image     string    `json:"image"`
	Type      string    `json:"bot_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TelegramSettingWinnerFilterDTO represents the DTO for filtering Telegram setting winners
type TelegramSettingWinnerFilterDTO struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	GroupID  *string  `json:"group_id,omitempty"`
	Token    *string  `json:"token,omitempty"`
	Type     *string  `json:"bot_type,omitempty"`
	Status   *string  `json:"status,omitempty"`
	MinPoint *float64 `json:"min_point,omitempty"`
	MaxPoint *float64 `json:"max_point,omitempty"`
}
