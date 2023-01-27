package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Name     string
	Email    string
	Tel      string
	Products []OrderItem `gorm:"foreignKey:OrderID"`
}

// somchaiOrder.OrrderItems =>
// Products.OrrderItems =>
