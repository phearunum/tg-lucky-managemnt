package model

import (
	"time"

	"gorm.io/gorm"
)

type TelegramMessage struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ChatID      string    `gorm:"not null"`
	Text        string    `gorm:"not null"`
	MessageType string    `gorm:"not null"` // e.g., text, photo, etc.
	SentAt      time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type TelegramAccount struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName string     `gorm:"type:varchar(20);not null" json:"account_name"`
	PhoneNumber string     `gorm:"type:varchar(20);not null" json:"phone_number"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type UserRequest struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	ChatID      string  `gorm:"type:varchar(50)" json:"chat_id"`
	AccountName string  `gorm:"type:varchar(100)" json:"account_name"`
	RequestType string  `gorm:"type:varchar(50)" json:"request_type"`
	StartTime   string  `gorm:"type:varchar(50)" json:"start_time"`
	EndTime     string  `gorm:"type:varchar(50)" json:"end_time"`
	Status      string  `gorm:"type:varchar(20)" json:"status"`
	AllowTime   float64 `gorm:"type:varchar(20)" json:"allow_time"`
	TotalTime   float64 `gorm:"type:varchar(20)" json:"total_time"`
	Message     string  `gorm:"type:varchar(250)" json:"message"`
	BotName     string  `gorm:"type:varchar(30)" json:"bot_name"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
type LocationSetting struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(50)" json:"name"`
	Lat       string    `gorm:"type:varchar(50)" json:"lat"`
	Long      string    `gorm:"type:varchar(50)" json:"long"`
	BotToken  string    `gorm:"type:varchar(200)" json:"token_bot"`
	Allow     float64   `gorm:"column:allow" json:"allow"`
	Status    string    `gorm:"type:varchar(20);default:'yes'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
type RequestSetting struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"type:varchar(50)" json:"name"`
	TimeRest   float64   `gorm:"type:varchar(50)" json:"time_rest"`
	ButtonType string    `gorm:"type:varchar(50)" json:"button_type"`
	BotName    string    `gorm:"type:varchar(50)" json:"bot_name"`
	Status     string    `gorm:"type:varchar(20)" json:"status"`
	OrderNo    int       `gorm:"type:varchar(20)" json:"order_int"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
type ClockTime struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ChatID      string    `gorm:"type:varchar(50)" json:"chat_id"`
	AccountName string    `gorm:"type:varchar(50)" json:"account_name"`
	RequestType string    `gorm:"type:varchar(50)" json:"request_type"`
	StartTime   string    `gorm:"type:varchar(50)" json:"start_time"`
	Message     string    `gorm:"type:varchar(250)" json:"message"`
	Lat         string    `gorm:"type:varchar(50)" json:"lat"`
	Long        string    `gorm:"type:varchar(50)" json:"long"`
	BotName     string    `gorm:"type:varchar(200)" json:"bot_name"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type PhoneLists struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Phone       string    `gorm:"type:varchar(50);uniqueIndex" json:"phone"`
	BotName     string    `gorm:"type:varchar(100)" json:"bot_name"`
	Requester   string    `gorm:"type:varchar(100)" json:"requester"`
	RequestDate string    `gorm:"type:varchar(100)" json:"request_date"`
	Status      string    `gorm:"type:varchar(20);default:'yes'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (PhoneLists) TableName() string {
	return "phone_list"
}
func (ClockTime) TableName() string {
	return "clock_time"
}
func (LocationSetting) TableName() string {
	return "LocationSetting"
}
func (UserRequest) TableName() string {
	return "TelegramUserRequest"
}
func (u *UserRequest) BeforeSave(tx *gorm.DB) (err error) {
	// If RoleID is 0 (not set), set it to 0

	return nil
}
func (RequestSetting) TableName() string {
	return "RequestSetting"
}
func (TelegramAccount) TableName() string {
	return "TelegramAccount"
}

// MigrateUsers automates the user table migration
func MigrateTelegramAccount(db *gorm.DB) {
	db.AutoMigrate(&TelegramAccount{})
}

func (TelegramMessage) TableName() string {
	return "TelegramMessage"
}

// MigrateUsers automates the user table migration
func MigrateTelegramMessage(db *gorm.DB) {
	db.AutoMigrate(&TelegramMessage{})
	db.AutoMigrate(&UserRequest{})
	db.AutoMigrate(&RequestSetting{})
	db.AutoMigrate(&LocationSetting{})
	db.AutoMigrate(&ClockTime{})
	db.AutoMigrate(&PhoneLists{})
}
