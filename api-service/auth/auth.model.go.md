package auth

import (
"time"

    "gorm.io/gorm"

)

type Role struct {
ID uint `gorm:"primaryKey" json:"id"`
RoleName string `gorm:"type:varchar(100)" json:"roleName"`
RoleStatus int `json:"roleStatus"`
RoleKey string `gorm:"type:varchar(100)" json:"roleKey"`
CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
return "Roles"
}

type User struct {
ID uint `json:"id" gorm:"primaryKey"`
FirstName string `json:"fname"`
LastName string `json:"lname"`
Username string `json:"username"`
Password string `json:"password"`
RoleID uint `json:"roleId"`
Role Role `json:"Role" gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL,OnUpdate:SET NULL"`

    CompanyID int            `json:"companyId"`
    Sex       string         `json:"sex"`
    Phone     string         `json:"phone"`
    Status    string         `json:"status"`
    Token     string         `json:"token"`
    CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
    UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

}

func (User) TableName() string {
return "Users"
}
