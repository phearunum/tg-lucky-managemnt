package videosnap

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Videosnap struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProcessId uint           `gorm:"column:process_id" json:"process_id"`
	Code      string         `gorm:"column:tran_id" json:"tran_id"`
	GameNo    string         `gorm:"column:game_no" json:"period"`
	Streamkey string         `gorm:"column:stream_key" json:"rtmpurl"`
	Rtmp      string         `gorm:"column:rtmp" json:"rtmp"`
	ImageURL  string         `gorm:"column:image_url" json:"image_url"`
	VideoURL  string         `gorm:"column:video_url" json:"video_url"`
	StorePath string         `gorm:"column:store_path" json:"store_path"`
	Status    bool           `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Videosnap) TableName() string {
	return "video_records"
}
func (v *Videosnap) BeforeSave(tx *gorm.DB) (err error) {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now().In(loc) // Set CreatedAt to current local time if not set
	}
	v.UpdatedAt = time.Now().In(loc) // Always set UpdatedAt to current local time

	return nil
}
func (v *Videosnap) BeforeCreate(tx *gorm.DB) (err error) {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	v.CreatedAt = time.Now().In(loc) // Set CreatedAt to current local time
	v.UpdatedAt = time.Now().In(loc) // Set UpdatedAt to current local time

	return nil
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&Videosnap{})
}
