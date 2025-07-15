package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name             string `gorm:"size:100;not null"`
	Email            string `gorm:"size:100;unique;not null"`
	Password         string `gorm:"size:255;not null"`
	Role             string `gorm:"size:50;default:user"`
	FailedLoginCount int    `gorm:"default:0"`
	LockedUntil      *time.Time
}
