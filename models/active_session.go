package models

import "gorm.io/gorm"

type ActiveSession struct {
	gorm.Model
	UserID      uint
	Token       string `gorm:"size:500"`
	DeviceID    string `gorm:"size:255"` // wajib dikirim dari frontend
	IP          string `gorm:"size:100"`
	UserAgent   string `gorm:"size:255"` // opsional tapi berguna untuk tracking
	IsActive    bool   `gorm:"default:true"`
	LoggedOutAt *int64 `gorm:""` // waktu logout (jika diperlukan tracking)
}
