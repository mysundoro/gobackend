package settings

import (
	"gobackend/database"
	"gobackend/models"

	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	var settings []models.Setting
	search := c.Query("search")

	query := database.DB.Model(&models.Setting{})
	if search != "" {
		query = query.Where("key LIKE ? OR value LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&settings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}

	return c.JSON(settings)
}

func GetByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	var setting models.Setting

	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Setting tidak ditemukan"})
	}

	return c.JSON(setting)
}

func Create(c *fiber.Ctx) error {
	var input models.Setting

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	// Cek duplikasi
	var exists models.Setting
	if err := database.DB.Where("key = ?", input.Key).First(&exists).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Key sudah ada"})
	}

	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan setting"})
	}

	return c.JSON(input)
}

func Update(c *fiber.Ctx) error {
	key := c.Params("key")
	var setting models.Setting

	if err := database.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Setting tidak ditemukan"})
	}

	var update models.Setting
	if err := c.BodyParser(&update); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse input"})
	}

	setting.Value = update.Value
	setting.Group = update.Group

	if err := database.DB.Save(&setting).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update setting"})
	}

	return c.JSON(setting)
}
