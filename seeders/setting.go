package seeders

import (
	"gobackend/database"
	"gobackend/models"
	"log"
)

func SeedSettings() {
	settings := []models.Setting{
		{Key: "site_name", Value: "MyApp", Group: 1},
		{Key: "logo", Value: "/uploads/logo.png", Group: 1},
		{Key: "favicon", Value: "/uploads/favicon.ico", Group: 1},
		{Key: "frontend_url", Value: "https://yourfrontend.com", Group: 1},

		// Konfigurasi keamanan login
		{Key: "max_failed_login_user", Value: "3", Group: 2},
		{Key: "lock_duration_minutes", Value: "30", Group: 2},
		{Key: "max_failed_login_ip", Value: "5", Group: 2},
		{Key: "ip_block_duration_minutes", Value: "60", Group: 2},

		// Konfigurasi email
		{Key: "mail_from", Value: "admin@yourdomain.com", Group: 3},
		{Key: "mail_host", Value: "smtp.gmail.com", Group: 3},
		{Key: "mail_port", Value: "587", Group: 3},
		{Key: "mail_user", Value: "admin@yourdomain.com", Group: 3},
		{Key: "mail_pass", Value: "aplikasi-password-atau-token", Group: 3},
	}

	for _, setting := range settings {
		var existing models.Setting
		err := database.DB.Where("key = ?", setting.Key).First(&existing).Error
		if err != nil {
			if err := database.DB.Create(&setting).Error; err != nil {
				log.Printf("❌ Gagal menyimpan setting '%s': %v", setting.Key, err)
			} else {
				log.Printf("✅ Setting '%s' berhasil ditambahkan", setting.Key)
			}
		}
	}
}
