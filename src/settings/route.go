package settings

import (
	"gobackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router) {
	r := router.Group("/settings")

	r.Get("/", GetAll)
	r.Get("/:key", GetByKey)
	r.Post("/", middleware.RequireAuth, Create)
	r.Put("/:key", middleware.RequireAuth, Update)
}
