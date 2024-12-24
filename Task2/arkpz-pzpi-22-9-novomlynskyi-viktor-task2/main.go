package main

import (
	"log"
	"my_project/config"
	"my_project/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Ініціалізація підключення до бази даних
	if err := config.InitDB(); err != nil {
		log.Fatalf("Помилка ініціалізації бази даних: %v", err)
	}

	// Створення нового Fiber-додатка
	app := fiber.New()

	// Middleware для додавання бази даних у контекст
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("db", config.DB)
		return ctx.Next()
	})

	// Налаштування маршрутів додатка
	routes.SetupRoutes(app)

	// Запуск HTTP-сервера
	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Помилка запуску сервера: %v", err)
	}
}









