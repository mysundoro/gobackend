package models

import "gorm.io/gorm"

type Setting struct {
	gorm.Model
	Key   string `gorm:"size:100;not null;uniqueIndex"`
	Value string `gorm:"type:text"`
	Group int    `gorm:"default:0"` // Untuk mengelompokkan setting (misal: 1 = sistem, 2 = email, dst)
}
