package models

import "gorm.io/gorm"

type AccountUnlockToken struct {
	gorm.Model
	UserID    uint
	Token     string `gorm:"size:255;uniqueIndex"`
	ExpiredAt int64
}
