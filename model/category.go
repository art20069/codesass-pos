package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;type:varchar(100);not null"`
}
