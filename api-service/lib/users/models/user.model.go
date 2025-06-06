// user models
package users

import (
	roles "api-service/lib/roles/models"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// User represents the model for a user
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	FirstName string         `json:"fname" gorm:"column:first_name"`                                               // Use snake_case for GORM column name
	LastName  string         `json:"lname" gorm:"column:last_name"`                                                // Use snake_case for GORM column name
	Username  string         `json:"username" gorm:"column:username"`                                              // Use snake_case for GORM column name
	Password  string         `json:"password" gorm:"column:password"`                                              // Use snake_case for GORM column name
	RoleID    uint           `json:"role_id" gorm:"column:role_id"`                                                // Use snake_case for GORM column name
	Role      roles.Role     `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL,OnUpdate:SET NULL"` // Ensure correct Role import
	CompanyID int            `json:"company_id" gorm:"column:company_id"`                                          // Use snake_case for GORM column name
	Sex       string         `json:"sex" gorm:"column:sex"`                                                        // Use snake_case for GORM column name
	Phone     string         `json:"phone" gorm:"column:phone"`                                                    // Use snake_case for GORM column name
	Status    string         `json:"status" gorm:"column:status"`                                                  // Consider changing to int for enum status
	TGgroup   string         `json:"tg_group" gorm:"column:tg_group"`
	Token     string         `json:"token" gorm:"column:token"`                         // Use snake_case for GORM column name
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // Change pointer to time.Time
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // Change pointer to time.Time                          // Automatically set update time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "system_user"
}

// BeforeSave hook to handle setting RoleID
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// If RoleID is 0 (not set), set it to 0
	if u.RoleID == 0 {
		u.RoleID = 0
	}
	return nil
}

// MigrateUsers automates the user table migration
func MigrateUsers(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

// Handle if User has no Role return Role {}
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User

	// Check if Role is empty
	emptyRole := roles.Role{}
	isEmptyRole := u.Role == emptyRole

	return json.Marshal(&struct {
		*Alias
		Role interface{} `json:"Role,omitempty"`
	}{
		Alias: (*Alias)(u),
		Role: func() interface{} {
			if isEmptyRole {
				return struct{}{}
			}
			return u.Role
		}(),
	})
}
