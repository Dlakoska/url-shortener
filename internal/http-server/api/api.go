package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"url-shortener/internal/http-server/middleware"
	"url-shortener/internal/service"
)

type Routers struct {
	Service service.Service
}

// NewRouters - конструктор для настройки API
func NewRouters(r *Routers, token string) *fiber.App {
	app := fiber.New()

	// Настройка CORS (разрешенные методы, заголовки, авторизация)
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET, POST, PUT, DELETE",
		AllowHeaders:  "Accept, Authorization, Content-Type, X-CSRF-Token, X-REQUEST-ID",
		ExposeHeaders: "Link",
		MaxAge:        300,
	}))

	// Группа маршрутов с авторизацией
	apiGroup := app.Group("/v1", middleware.Authorization(token))

	// Роут для создания задачи
	apiGroup.Post("/new-alias", r.Service.CreateAlias)
	apiGroup.Get("/:alias", r.Service.GetURL)
	apiGroup.Delete("/:alias", r.Service.DeleteUrl)

	return app
}
