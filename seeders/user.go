package seeders

import (
	"golang.org/x/crypto/bcrypt"

	"gobackend/database"
	"gobackend/models"
)

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Gagal hash password: " + err.Error())
	}
	return string(hashed)
}

func SeedUsers() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		database.DB.Create(&[]models.User{
			{
				Name:     "Admin",
				Email:    "admin@gobackend.com",
				Password: hashPassword("[password]"),
				Role:     "admin",
			},
		})
	}
}
