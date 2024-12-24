package controllers

import (
	"log"
	"my_project/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateAlert - функція для створення сповіщення
func CreateAlert(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	var alert models.Alert
	if err := c.BodyParser(&alert); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Create(&alert).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving alert to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alert created successfully",
		"alert":   alert,
	})
}
// GetAllAlerts - отримує всі попередження з бази даних
func GetAllAlerts(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Список попереджень
	var alerts []models.Alert

	// Отримуємо всі попередження з бази
	if err := db.Find(&alerts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching alerts",
		})
	}

	// Повертаємо список попереджень
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alerts": alerts,
	})
}
// GetAlertByID - отримує попередження за ID
func GetAlertByID(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID попередження з параметрів URL
	alertID := c.Params("id")

	// Створюємо структуру для попередження
	var alert models.Alert

	// Знаходимо попередження по ID
	if err := db.First(&alert, alertID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Alert not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching alert",
		})
	}

	// Повертаємо знайдене попередження
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alert": alert,
	})
}
// GetAlertsByStatus - отримує попередження за статусом
func GetAlertsByStatus(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо параметр статусу з query string
	status := c.Query("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status query parameter is required",
		})
	}

	// Створюємо структуру для результатів
	var alerts []models.Alert

	// Виконуємо запит для отримання попереджень за статусом
	if err := db.Where("status = ?", status).Find(&alerts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching alerts",
		})
	}

	// Повертаємо знайдені попередження
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alerts": alerts,
	})
}
// GetAlertsByType - отримує попередження за типом
func GetAlertsByType(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо параметр типу з query string
	alertType := c.Query("type")
	if alertType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Type query parameter is required",
		})
	}

	// Створюємо структуру для результатів
	var alerts []models.Alert

	// Виконуємо запит для отримання попереджень за типом
	if err := db.Where("type = ?", alertType).Find(&alerts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching alerts",
		})
	}

	// Повертаємо знайдені попередження
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alerts": alerts,
	})
}
// GetAlertHistory - отримує історію попереджень для експонату
func GetAlertHistory(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо параметр ID експонату або музею
	exhibitID := c.Query("exhibit_id")
	museumID := c.Query("museum_id")

	// Створюємо структуру для результатів
	var alerts []models.Alert

	// Виконуємо запит для отримання історії попереджень
	if exhibitID != "" {
		if err := db.Where("exhibit_id = ?", exhibitID).Find(&alerts).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error fetching alert history",
			})
		}
	} else if museumID != "" {
		if err := db.Where("museum_id = ?", museumID).Find(&alerts).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error fetching alert history",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Either exhibit_id or museum_id query parameter is required",
		})
	}

	// Повертаємо історію попереджень
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alert_history": alerts,
	})
}

// GetAlert - функція для отримання сповіщення за ID
func GetAlert(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	alertID := c.Params("id")
	var alert models.Alert
	if err := db.First(&alert, alertID).Error; err != nil {
		log.Println("Alert not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Alert not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alert": alert,
	})
}

// UpdateAlert - функція для оновлення сповіщення
func UpdateAlert(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	alertID := c.Params("id")
	var alert models.Alert
	if err := db.First(&alert, alertID).Error; err != nil {
		log.Println("Alert not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Alert not found",
		})
	}

	var updatedAlert models.Alert
	if err := c.BodyParser(&updatedAlert); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Save(&updatedAlert).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating alert",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alert updated successfully",
		"alert":   updatedAlert,
	})
}

// DeleteAlert - функція для видалення сповіщення
func DeleteAlert(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	alertID := c.Params("id")
	if err := db.Delete(&models.Alert{}, alertID).Error; err != nil {
		log.Println("Error deleting alert:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting alert",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Alert deleted successfully",
	})
}


