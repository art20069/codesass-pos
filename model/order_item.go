package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	SKU      string  `gorm:"not null"`
	Name     string  `gorm:"not null"`
	Image    string  `gorm:"not null"`
	Price    float64 `gorm:"not null"`
	Quantity uint    `gorm:"not null"`
	OrderID  uint
}
