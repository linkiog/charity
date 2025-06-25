package models

import "gorm.io/gorm"

type Mosque struct {
	gorm.Model `json:"-"`

	Name       string `gorm:"unique;not null"`
	City       string `json:",omitempty"`
	Region     string `json:",omitempty"`
	Requisites string `json:",omitempty"`

	AdminID  uint      `json:"-"`
	Admin    *User     `gorm:"foreignKey:AdminID" json:"-"`
	Products []Product `gorm:"foreignKey:MosqueID" json:"products,omitempty"`
}
