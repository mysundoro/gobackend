package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RegisterRoutes(router fiber.Router) {
	// Rate limit: max 3 request per 15 menit per IP
	resetLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 15 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Terlalu banyak percobaan. Coba lagi nanti.",
			})
		},
	})

	group := router.Group("/auth")

	group.Post("/login", resetLimiter, Login)

	router.Post("/forgot-password", resetLimiter, ForgotPassword)
	router.Post("/reset-password", ResetPassword)
	group.Get("/unlock-account", UnlockAccount)
}
