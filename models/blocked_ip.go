package models

import "gorm.io/gorm"

type BlockedIP struct {
	gorm.Model
	IP        string `gorm:"unique"`
	Attempts  int
	BlockedAt int64
	ExpiresAt int64
}
