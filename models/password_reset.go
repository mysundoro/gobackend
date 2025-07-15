package models

import "gorm.io/gorm"

type PasswordReset struct {
	gorm.Model
	Email     string `gorm:"size:100;not null"`
	Token     string `gorm:"size:255;not null;unique"`
	ExpiredAt int64  `gorm:"not null"`
}
