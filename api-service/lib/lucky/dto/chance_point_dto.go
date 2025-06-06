package dto

import (
	"time"
)

// CreateChancePointDTO represents the DTO for creating a ChancePoint entry
type CreateChancePointDTO struct {
	AccountName string  `json:"account_name"`
	EmployeeID  string  `json:"employee_id"`
	Profile     string  `json:"profile"`
	ChatID      string  `json:"chat_id"`
	Status      string  `json:"status"`
	Note        string  `json:"note"`
	CreateBy    string  `json:"create_by"`
	TotalChance float64 `json:"total_chance"`
}

// UpdateChancePointDTO represents the DTO for updating a ChancePoint entry
type UpdateChancePointDTO struct {
	ID          uint     `json:"id"`
	AccountName *string  `json:"account_name,omitempty"`
	EmployeeID  *string  `json:"employee_id,omitempty"`
	Profile     *string  `json:"profile,omitempty"`
	ChatID      *string  `json:"chat_id,omitempty"`
	Status      *string  `json:"status,omitempty"`
	Note        *string  `json:"note,omitempty"`
	CreateBy    *string  `json:"create_by,omitempty"`
	TotalChance *float64 `json:"total_chance,omitempty"`
}

// ResponseChancePointDTO represents the DTO for responding with a ChancePoint entry
type ResponseChancePointDTO struct {
	ID     uint `json:"id"`
	Member struct {
		AccountName string  `json:"account_name"`
		EmployeeID  string  `json:"employee_id"`
		Profile     string  `json:"profile"`
		TotalPoints float64 `json:"total_points"`
	} `json:"member"`
	AccountName string     `json:"account_name"`
	EmployeeID  string     `json:"employee_id"`
	Profile     string     `json:"profile"`
	ChatID      string     `json:"chat_id"`
	TotalChance float64    `json:"total_chance"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
	Exprired    string     `json:"expired"`
	CreateBy    string     `json:"create_by"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// ChancePointFilterDTO represents the DTO for filtering ChancePoint entries
type ChancePointFilterDTO struct {
	ID          uint    `json:"id"`
	AccountName *string `json:"account_name,omitempty"`
	EmployeeID  *string `json:"employee_id,omitempty"`
	Profile     *string `json:"profile,omitempty"`
	ChatID      string  `json:"chat_id"`
	TotalChance float64 `json:"total_chance"`
	Status      *string `json:"status,omitempty"`
	CreateBy    *string `json:"create_by,omitempty"`
	Exprired    string  `json:"expired"`
}
