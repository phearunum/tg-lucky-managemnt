package dto

type TelegramAccountCreateDTO struct {
	AccountName string `json:"account_name"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type TelegramAccountUpdateDTO struct {
	ID          uint   `json:"id" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	AccountName string `json:"account_name"`
}

type TelegramAccountDeleteDTO struct {
	ID uint `json:"id"`
}

type TelegramMessageDTO struct {
	ChatID  int64  `json:"chat_id"`
	Message string `json:"message"`
}
type BulkRequestDTO struct {
	ID            []uint `json:"id"`
	Requester     string `json:"requester"`
	RequesterDate string `json:"request_date"`
}
