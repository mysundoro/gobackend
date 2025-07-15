package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gobackend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Ambil dari ENV atau bisa di-hardcode untuk awal
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}

	// Ping untuk tes koneksi
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Minute * 5)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	DB = db
	fmt.Println("✅ Koneksi database berhasil")

	// Auto migrate
	autoMigrate()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.ActiveSession{},
		&models.PasswordReset{},
		&models.AccountUnlockToken{},
		&models.BlockedIP{},
		&models.Setting{},
	)
	if err != nil {
		log.Fatalf("Gagal migrate table: %v", err)
	}
	fmt.Println("✅ Auto migrate selesai")
}
