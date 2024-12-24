package controllers

import (
	"log"
	"my_project/models" 
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateSensor - функція для створення сенсора
func CreateSensor(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	var sensor models.Sensor
	if err := c.BodyParser(&sensor); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Create(&sensor).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving sensor to database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Sensor created successfully",
		"sensor":  sensor,
	})
}

// GetSensor - функція для отримання сенсора за ID
func GetSensor(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	sensorID := c.Params("id")
	var sensor models.Sensor
	if err := db.First(&sensor, sensorID).Error; err != nil {
		log.Println("Sensor not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Sensor not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"sensor": sensor,
	})
}

// UpdateSensor - функція для оновлення сенсора
func UpdateSensor(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	sensorID := c.Params("id")
	var sensor models.Sensor
	if err := db.First(&sensor, sensorID).Error; err != nil {
		log.Println("Sensor not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Sensor not found",
		})
	}

	var updatedSensor models.Sensor
	if err := c.BodyParser(&updatedSensor); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	if err := db.Save(&updatedSensor).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating sensor",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sensor updated successfully",
		"sensor":  updatedSensor,
	})
}

// DeleteSensor - функція для видалення сенсора
func DeleteSensor(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	sensorID := c.Params("id")
	if err := db.Delete(&models.Sensor{}, sensorID).Error; err != nil {
		log.Println("Error deleting sensor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting sensor",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sensor deleted successfully",
	})
}

