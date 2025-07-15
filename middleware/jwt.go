package middleware

import (
	"os"
	"strings"

	"gobackend/database"
	"gobackend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Decode token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid"})
	}

	// Ambil user_id dari JWT
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"]

	// ✅ Validasi token dari DB
	var session models.ActiveSession
	err = database.DB.Where("token = ? AND user_id = ? AND is_active = true", tokenStr, userID).First(&session).Error
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token tidak terdaftar atau sudah logout",
		})
	}

	// Simpan userID ke context
	c.Locals("userID", userID)
	return c.Next()
}
