package controllers

import (
	"log"
	"my_project/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateExhibit - функція для створення експоната
func CreateExhibit(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	var exhibit models.Exhibit
	if err := c.BodyParser(&exhibit); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Create(&exhibit).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving exhibit to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Exhibit created successfully",
		"exhibit": exhibit,
	})
}
// GetAllExhibits - отримує всі експонати з бази даних
func GetAllExhibits(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Список експонатів
	var exhibits []models.Exhibit

	// Отримуємо всі експонати з бази
	if err := db.Find(&exhibits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits",
		})
	}

	// Повертаємо список експонатів
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibits": exhibits,
	})
}
// GetExhibitByID - отримує експонат за ID
func GetExhibitByID(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID експонату з параметрів URL
	exhibitID := c.Params("id")

	// Створюємо структуру для експонату
	var exhibit models.Exhibit

	// Знаходимо експонат по ID
	if err := db.First(&exhibit, exhibitID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Exhibit not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibit",
		})
	}

	// Повертаємо знайдений експонат
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibit": exhibit,
	})
}
// GetExhibitStats - отримує статистику по експонатах
func GetExhibitStats(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Запит для отримання статистики по експонатах (наприклад, кількість попереджень)
	var stats []struct {
		ExhibitID   uint   `json:"exhibit_id"`
		ExhibitName string `json:"exhibit_name"`
		AlertCount  int    `json:"alert_count"`
	}

	// Виконуємо запит для отримання статистики
	if err := db.Table("exhibits").
		Select("exhibits.id as exhibit_id, exhibits.name as exhibit_name, count(alerts.id) as alert_count").
		Joins("left join alerts on alerts.exhibit_id = exhibits.id").
		Group("exhibits.id").
		Scan(&stats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibit stats",
		})
	}

	// Повертаємо статистику експонатів
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibit_stats": stats,
	})
}
// SearchExhibits - пошук експонатів за критеріями
func SearchExhibits(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо параметри пошуку з query string
	name := c.Query("name")    // Назва експонату
	material := c.Query("material") // Матеріал експонату

	// Створюємо структуру для результатів пошуку
	var exhibits []models.Exhibit

	// Будуємо запит
	query := db.Model(&models.Exhibit{})

	// Додаємо умови пошуку
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if material != "" {
		query = query.Where("material LIKE ?", "%"+material+"%")
	}

	// Виконуємо запит для отримання експонатів
	if err := query.Find(&exhibits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits",
		})
	}

	// Повертаємо знайдені експонати
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibits": exhibits,
	})
}
// GetExhibitsUnderMaintenance - отримує експонати, що знаходяться на технічному обслуговуванні
func GetExhibitsUnderMaintenance(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо структуру для результатів
	var exhibits []models.Exhibit

	// Виконуємо запит для отримання експонатів, що знаходяться на технічному обслуговуванні
	if err := db.Where("maintenance_status = ?", "under_maintenance").Find(&exhibits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits under maintenance",
		})
	}

	// Повертаємо список експонатів на технічному обслуговуванні
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibits_under_maintenance": exhibits,
	})
}
// GetExhibitsByPriceRange - отримує експонати в певному діапазоні цін
func GetExhibitsByPriceRange(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо мінімальну та максимальну ціну з query параметрів
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	if minPrice == "" || maxPrice == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Both min_price and max_price query parameters are required",
		})
	}

	// Створюємо структуру для результатів
	var exhibits []models.Exhibit

	// Виконуємо запит для отримання експонатів в певному діапазоні цін
	if err := db.Where("value BETWEEN ? AND ?", minPrice, maxPrice).Find(&exhibits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits by price range",
		})
	}

	// Повертаємо знайдені експонати
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibits": exhibits,
	})
}
// GetExhibitsByMaterial - отримує експонати за матеріалом
func GetExhibitsByMaterial(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо матеріал з query параметра
	material := c.Query("material")
	if material == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Material query parameter is required",
		})
	}

	// Створюємо структуру для результатів
	var exhibits []models.Exhibit

	// Виконуємо запит для отримання експонатів за матеріалом
	if err := db.Where("material = ?", material).Find(&exhibits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits by material",
		})
	}

	// Повертаємо знайдені експонати
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibits": exhibits,
	})
}
// GetAlertsByDateRange - отримує попередження в певному діапазоні дат
func GetAlertsByDateRange(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо початкову та кінцеву дату з query параметрів
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Both start_date and end_date query parameters are required",
		})
	}

	// Створюємо структуру для результатів
	var alerts []models.Alert

	// Виконуємо запит для отримання попереджень в діапазоні дат
	if err := db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&alerts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching alerts by date range",
		})
	}

	// Повертаємо знайдені попередження
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alerts": alerts,
	})
}
// GetUrgentAlerts - отримує всі термінові попередження
func GetUrgentAlerts(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо структуру для результатів
	var alerts []models.Alert

	// Виконуємо запит для отримання термінових попереджень
	if err := db.Where("status = ?", "urgent").Find(&alerts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching urgent alerts",
		})
	}

	// Повертаємо знайдені термінові попередження
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"urgent_alerts": alerts,
	})
}

// GetExhibit - функція для отримання експоната за ID
func GetExhibit(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	exhibitID := c.Params("id")
	var exhibit models.Exhibit
	if err := db.First(&exhibit, exhibitID).Error; err != nil {
		log.Println("Exhibit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Exhibit not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"exhibit": exhibit,
	})
}

// UpdateExhibit - функція для оновлення експоната
func UpdateExhibit(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	exhibitID := c.Params("id")
	var exhibit models.Exhibit
	if err := db.First(&exhibit, exhibitID).Error; err != nil {
		log.Println("Exhibit not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Exhibit not found",
		})
	}

	var updatedExhibit models.Exhibit
	if err := c.BodyParser(&updatedExhibit); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Save(&updatedExhibit).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating exhibit",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exhibit updated successfully",
		"exhibit": updatedExhibit,
	})
}
// GetExhibitsByMuseumID - бізнес-логіка для отримання всіх експонатів за ID музею
func GetExhibitsByMuseumID(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID музею з параметрів запиту
	museumID := c.Params("id")
	if museumID == "" {
		log.Println("Museum ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Museum ID is required",
		})
	}

	// Структура для збереження результатів запиту
	var exhibits []models.Exhibit

	// Виконуємо SQL-запит, щоб отримати всі експонати для конкретного музею
	err := db.Where("museum_id = ?", museumID).Find(&exhibits).Error
	if err != nil {
		log.Printf("Error fetching exhibits for museum ID %s: %v", museumID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching exhibits",
		})
	}

	// Якщо немає експонатів для цього музею
	if len(exhibits) == 0 {
		log.Printf("No exhibits found for museum ID %s", museumID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No exhibits found for the museum",
		})
	}

	// Повертаємо результат у вигляді JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exhibits retrieved successfully",
		"data":    exhibits,
	})
}

// DeleteExhibit - функція для видалення експоната
func DeleteExhibit(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	exhibitID := c.Params("id")
	if err := db.Delete(&models.Exhibit{}, exhibitID).Error; err != nil {
		log.Println("Error deleting exhibit:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting exhibit",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Exhibit deleted successfully",
	})
}


