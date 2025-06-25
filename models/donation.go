package models

import "gorm.io/gorm"

type Donation struct {
	gorm.Model

	ProductID uint
	Product   Product `gorm:"foreignKey:ProductID"`

	Qty    int     `gorm:"not null"` // сколько единиц товара хотят оплатить, либо купили
	Amount float64 `gorm:"not null"` // Price × Qty (фиксируем на момент оплаты)

	UserID uint
}
