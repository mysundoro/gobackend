package users

import (
	"gobackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {
	r := router.Group("/users", middleware.RequireAuth)

	r.Get("/", GetAll)
	r.Post("/", Create)
	r.Get("/:id", GetByID)
	r.Put("/:id", Update)
	r.Delete("/:id", Delete)
}
