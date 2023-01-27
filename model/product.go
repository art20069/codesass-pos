package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	SKU        string   `gorm:"uniqueIndex;type:varchar(100);not null"`
	Name       string   `gorm:"type:varchar(100);not null"`
	Desc       string   `gorm:"type:varchar(255)"`
	Price      float64  `gorm:"not null"`
	Status     uint     `gorm:"not null"`
	Image      string   `gorm:"type:varchar(255);not null"`
	CategoryID uint     `gorm:"not null"`
	Category   Category //ผูกความสัมพันธ์ของตาราง
}

// SKU จะซ้ำไม่ได้บันทึกเข้ามาได้ค่าเดียว
//ผูกความสัมพันธ์ของตาราง
// Category <--> Product (tulip)
// tulip.Category => Category(Name:Flower)
