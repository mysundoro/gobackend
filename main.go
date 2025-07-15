package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"gobackend/database"
	"gobackend/src/users"

	"gobackend/seeders"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file tidak ditemukan, gunakan ENV dari sistem")
	}

	// Koneksi DB
	database.Connect()

	// Tes koneksi database
	var version string
	database.DB.Raw("SELECT VERSION()").Scan(&version)
	fmt.Println("Database version:", version)

	// Seeder awal (jika diperlukan)
	seeders.SeedUsers()
	seeders.SeedSettings()

	// Inisialisasi Fiber
	app := fiber.New()

	// 🌐 Middleware global
	app.Use(recover.New()) // Hindari crash dari panic
	app.Use(helmet.New())  // Tambahkan header keamanan HTTP
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 📦 Routing
	api := app.Group("/api")
	users.RegisterRoutes(api)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ PORT tidak ditemukan di .env")
	}

	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
