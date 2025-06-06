package models

import (
	"time"

	"gorm.io/gorm"
)

type TelegramUsers struct {
	ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string  `gorm:"type:varchar(50);column:username" json:"username"`
	AccountName string  `gorm:"type:varchar(50);column:account_name" json:"account_name"`
	EmployeeID  string  `gorm:"type:varchar(50);column:employee_id" json:"employee_id"`
	Profile     string  `gorm:"type:varchar(250);column:profile" json:"profile"`
	ChatID      string  `gorm:"type:varchar(30);uniqueIndex;column:chat_id" json:"chat_id"`
	TotalChance float64 `gorm:"column:total_chance" json:"total_chance"`
	TotalPoints float64 `gorm:"column:total_point" json:"total_point"`
	//ChancePoints []ChancePoint `gorm:"foreignKey:ChatID;references:ChatID" json:"chance_points"`
	BotName   string    `gorm:"type:varchar(30)" json:"bot_name"`
	GroupID   string    `gorm:"type:varchar(30)" json:"gtoup_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
type ChancePoint struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName string        `gorm:"type:varchar(50);column:account_name" json:"account_name"`
	EmployeeID  string        `gorm:"type:varchar(50);column:employee_id" json:"employee_id"`
	Profile     string        `gorm:"type:varchar(250);column:profile" json:"profile"`
	ChatID      string        `gorm:"type:varchar(30);column:chat_id" json:"chat_id"`
	Member      TelegramUsers `gorm:"foreignKey:ChatID;references:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"member"`
	TotalChance float64       `gorm:"column:total_chance" json:"total_chance"`
	Status      string        `gorm:"type:varchar(20);column:status" json:"status"`
	Exprired    string        `gorm:"type:varchar(20);column:expired;default:'no'" json:"expired"`

	Note      string     `gorm:"type:varchar(250);column:note" json:"note"`
	CreateBy  string     `gorm:"type:varchar(20);column:create_by" json:"create_by"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type LuckWinner struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	RankID      string        `gorm:"type:varchar(30);uniqueIndex;column:rank_id" json:"rank_id"`
	ChatID      string        `gorm:"type:varchar(30);column:chat_id" json:"chat_id"`
	Member      TelegramUsers `gorm:"foreignKey:ChatID;references:ChatID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"member"`
	TotalPoints float64       `gorm:"column:total_point" json:"total_point"`
	Status      string        `gorm:"type:varchar(20);column:status" json:"status"`
	Exprired    string        `gorm:"type:varchar(20);column:expired;default:'no'" json:"expired"`
	Note        string        `gorm:"type:varchar(250);column:note" json:"note"`
	CreateBy    string        `gorm:"type:varchar(20);column:create_by" json:"create_by"`
	CreatedAt   *time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   *time.Time    `gorm:"column:updated_at" json:"updatedAt"`
}

type TelegramSettingWinner struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	TotalWin  float64    `gorm:"column:total_win" json:"total_win"`
	MaxPoint  float64    `gorm:"column:max_point" json:"max_point"`
	MinPoint  float64    `gorm:"column:min_point" json:"min_point"`
	GroupID   string     `gorm:"type:varchar(30);column:group_id" json:"group_id"`
	Token     string     `gorm:"type:varchar(100);column:token" json:"token"`
	Message   string     `gorm:"type:text;column:message" json:"message"`
	Image     string     `gorm:"type:varchar(250);column:image" json:"image"`
	Type      string     `gorm:"type:varchar(30);column:bot_type" json:"bot_type"`
	Status    string     `gorm:"type:varchar(30);column:status" json:"status"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (u *ChancePoint) AfterSave(tx *gorm.DB) (err error) {
	// Step 1: Sum all ChancePoints for this ChatID
	var totalPoints float64
	err = tx.Model(&ChancePoint{}).
		Where("chat_id = ? AND status = ? and expired =?", u.ChatID, "collected", "no").
		Select("COALESCE(SUM(total_chance), 0)").
		Scan(&totalPoints).Error
	if err != nil {
		return err
	}
	// Step 2: Update the TelegramUser's TotalPoints field
	err = tx.Model(&TelegramUsers{}).
		Where("chat_id = ?", u.ChatID).
		Update("total_point", totalPoints).Error
	if err != nil {
		return err
	}

	return nil
}

func (TelegramUsers) TableName() string {
	return "telegram_users"
}
func (ChancePoint) TableName() string {
	return "telegram_chance"
}
func (TelegramSettingWinner) TableName() string {
	return "telegram_win_setting"
}
func (LuckWinner) TableName() string {
	return "telegram_winner"
}

// MigrateUsers automates the user table migration
func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&TelegramUsers{})
	db.AutoMigrate(&ChancePoint{})
	db.AutoMigrate(&TelegramSettingWinner{})
	db.AutoMigrate(&LuckWinner{})
}
