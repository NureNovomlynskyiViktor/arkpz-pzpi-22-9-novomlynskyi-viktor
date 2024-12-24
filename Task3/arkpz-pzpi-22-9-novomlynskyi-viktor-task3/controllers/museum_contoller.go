package controllers

import (
	"log"
	"my_project/models" 
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateMuseum - функція для створення музею
func CreateMuseum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	var museum models.Museum
	if err := c.BodyParser(&museum); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Create(&museum).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving museum to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Museum created successfully",
		"museum":  museum,
	})
}
// GetAllMuseums - отримує всі музеї з бази даних
func GetAllMuseums(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Список музеїв
	var museums []models.Museum

	// Отримуємо всі музеї з бази
	if err := db.Find(&museums).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching museums",
		})
	}

	// Повертаємо список музеїв
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"museums": museums,
	})
}
// GetMuseumByID - отримує музей за ID
func GetMuseumByID(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID музею з параметрів URL
	museumID := c.Params("id")

	// Створюємо структуру для музею
	var museum models.Museum

	// Знаходимо музей по ID
	if err := db.First(&museum, museumID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Museum not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching museum",
		})
	}

	// Повертаємо знайдений музей
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"museum": museum,
	})
}
// GetMuseumStats - бізнес-логіка для отримання статистики музеїв
func GetMuseumStats(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection error")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Структура для отримання статистики
	var museumStats []struct {
		MuseumName   string `json:"museum_name"`
		ExhibitCount int64  `json:"exhibit_count"`
	}

	// SQL запит для отримання статистики
	err := db.Table("exhibits").
		Select("museums.name as museum_name, count(exhibits.id) as exhibit_count").
		Joins("left join museums on exhibits.museum_id = museums.id").
		Group("museums.name").
		Scan(&museumStats).Error
	if err != nil {
		log.Println("Error fetching museum stats:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching museum statistics",
		})
	}

	// Відповідь зі статистикою
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Museum statistics retrieved successfully",
		"data":    museumStats,
	})
}
// GetMuseum - функція для отримання музею за ID
func GetMuseum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	museumID := c.Params("id")
	var museum models.Museum
	if err := db.First(&museum, museumID).Error; err != nil {
		log.Println("Museum not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Museum not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"museum": museum,
	})
}

// UpdateMuseum - функція для оновлення музею
func UpdateMuseum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	museumID := c.Params("id")
	var museum models.Museum
	if err := db.First(&museum, museumID).Error; err != nil {
		log.Println("Museum not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Museum not found",
		})
	}

	var updatedMuseum models.Museum
	if err := c.BodyParser(&updatedMuseum); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Save(&updatedMuseum).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating museum",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Museum updated successfully",
		"museum":  updatedMuseum,
	})
}

// DeleteMuseum - функція для видалення музею
func DeleteMuseum(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	museumID := c.Params("id")
	if err := db.Delete(&models.Museum{}, museumID).Error; err != nil {
		log.Println("Error deleting museum:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting museum",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Museum deleted successfully",
	})
}


