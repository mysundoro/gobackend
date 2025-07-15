package users

import (
	"fmt"

	"gobackend/database"
	"gobackend/models"

	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	var users []models.User

	// 1. Ambil query parameter
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)
	search := c.Query("search")
	sort := c.Query("sort", "id")
	order := c.Query("order", "desc")

	offset := (page - 1) * limit

	// 2. Query dasar
	query := database.DB.Model(&models.User{})

	// 3. Filter search (nama atau email mengandung keyword)
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ? OR role LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 4. Hitung total data
	var total int64
	query.Count(&total)

	// 5. Ambil data dengan limit, offset, dan sorting
	if err := query.Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data users"})
	}

	// 6. Return dengan pagination meta
	return c.JSON(fiber.Map{
		"data":  users,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func Create(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse request"})
	}

	// Validasi simple
	if input.Email == "" || input.Password == "" || input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Semua field wajib diisi"})
	}

	// Cek duplikasi email
	var exists models.User
	if err := database.DB.Where("email = ?", input.Email).First(&exists).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "Email sudah terdaftar"})
	}

	// Simpan ke DB
	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan user"})
	}

	return c.JSON(input)
}

func GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	return c.JSON(user)
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse data"})
	}

	// Update field yang boleh
	user.Name = updateData.Name
	user.Email = updateData.Email
	user.Role = updateData.Role

	if updateData.Password != "" {
		user.Password = updateData.Password // hash dulu jika perlu
	}

	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update user"})
	}

	return c.JSON(user)
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User tidak ditemukan"})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus user"})
	}

	return c.JSON(fiber.Map{"message": "User berhasil dihapus"})
}
