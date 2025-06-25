package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	MosqueID uint `gorm:"not null"`

	Name        string `gorm:"not null"`
	Description string
	Price       float64 `gorm:"not null"`

	Need      int `gorm:"not null"`
	Purchased int `gorm:"not null; default:0"`

	Mosque *Mosque `gorm:"foreignKey:MosqueID" json:"mosque,omitempty"`

	Donations []Donation `gorm:"foreignKey:ProductID"`
}
