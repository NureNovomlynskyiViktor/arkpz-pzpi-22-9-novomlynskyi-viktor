package routes

import (
	"my_project/controllers"
	"github.com/gofiber/fiber/v2"
)



func SetupRoutes(app *fiber.App) {

	
	app.Post("/register", controllers.RegisterUser)

	
	app.Get("/login", controllers.LoginUser)

	
	app.Put("/users/:id", controllers.UpdateUser)

	
	app.Delete("/admin/users/:id", controllers.DeleteUser)

	// Маршрути для музеїв (museums)
	app.Post("/admin/museums", controllers.CreateMuseum) 
	app.Get("/admin/museums", controllers.GetAllMuseums) 
	app.Get("/admin/museums/:id", controllers.GetMuseumByID) 
	app.Put("/admin/museums/:id", controllers.UpdateMuseum) 
	app.Delete("/admin/museums/:id", controllers.DeleteMuseum) 

	// Маршрути для експонатів (exhibits)
	app.Post("/admin/exhibits", controllers.CreateExhibit) 
	app.Get("/admin/exhibits", controllers.GetAllExhibits) 
	app.Get("/admin/exhibits/:id", controllers.GetExhibitByID) 
	app.Put("/admin/exhibits/:id", controllers.UpdateExhibit) 
	app.Delete("/admin/exhibits/:id", controllers.DeleteExhibit) 

	// Маршрути для сповіщень (alerts)
	app.Post("/admin/alerts", controllers.CreateAlert) 
	app.Get("/admin/alerts", controllers.GetAllAlerts) 
	app.Get("/admin/alerts/:id", controllers.GetAlertByID) 
	app.Put("/admin/alerts/:id", controllers.UpdateAlert) 
	app.Delete("/admin/alerts/:id", controllers.DeleteAlert) 
	
	// Маршрут для отримання експонатів за ID музею
	app.Get("/admin/museums/:id/exhibits", controllers.GetExhibitsByMuseumID)

	// Бізнес-логіка: Статистика музеїв
	/*app.Get("/admin/museums/stats", controllers.GetMuseumStats) // Статистика по музеях

	// Бізнес-логіка: Статистика експонатів
	app.Get("/admin/exhibits/stats", controllers.GetExhibitStats) // Статистика по експонатах

	// Бізнес-логіка: Пошук експонатів за параметрами
	app.Get("/admin/exhibits/search", controllers.SearchExhibits) // Пошук за матеріалом, ціною, датою створення

	// Бізнес-логіка: Пошук сповіщень за статусом
	app.Get("/admin/alerts/status/:status", controllers.GetAlertsByStatus) // Пошук сповіщень за статусом

	// Бізнес-логіка: Пошук сповіщень за типом
	app.Get("/admin/alerts/type/:type", controllers.GetAlertsByType) // Пошук сповіщень за типом

	// Бізнес-логіка: Історія сповіщень
	app.Get("/admin/alerts/history", controllers.GetAlertHistory) // Отримати історію сповіщень

	// Бізнес-логіка: Експонати під ремонтом
	app.Get("/admin/exhibits/under-maintenance", controllers.GetExhibitsUnderMaintenance) // Експонати на ремонті

	// Бізнес-логіка: Пошук експонатів по вартості
	app.Get("/admin/exhibits/price-range", controllers.GetExhibitsByPriceRange) // Пошук за діапазоном вартості

	// Бізнес-логіка: Пошук експонатів за матеріалом
	app.Get("/admin/exhibits/material/:material", controllers.GetExhibitsByMaterial) // Пошук за матеріалом

	// Бізнес-логіка: Фільтрація сповіщень по даті
	app.Get("/admin/alerts/date-range", controllers.GetAlertsByDateRange) // Пошук сповіщень за діапазоном дат

	// Бізнес-логіка: Термінові сповіщення
	app.Get("/admin/alerts/urgent", controllers.GetUrgentAlerts) // Термінові сповіщення*/

}

