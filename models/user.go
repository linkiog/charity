package models

import "gorm.io/gorm"

const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"-"`
	Email    string `gorm:"unique;not null" json:"-"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null;default:'user'" json:"-"`

	Mosques []Mosque `gorm:"foreignKey:AdminID" json:"-"`
}
