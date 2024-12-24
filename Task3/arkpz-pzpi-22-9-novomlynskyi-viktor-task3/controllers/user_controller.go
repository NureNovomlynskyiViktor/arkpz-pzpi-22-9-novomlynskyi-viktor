package controllers

import (
	"log"
	"my_project/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// RegisterUser - реєстрація нового користувача
func RegisterUser(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database connection error"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request data"})
	}

	// Перевіряємо, чи існує вже користувач з таким email
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		log.Println("User with this email already exists")
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User with this email already exists"})
	}

	// Хешуємо пароль перед збереженням
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error hashing password"})
	}

	user.PasswordHash = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error saving user to database"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully", "user": user})
}

// LoginUser - логін користувача
func LoginUser(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database connection error"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request data"})
	}

	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		log.Println("User not found:", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	// Перевірка пароля
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		log.Println("Incorrect password:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Incorrect password"})
	}

	userResponse := struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{
		ID:        existingUser.ID,
		Name:      existingUser.Name,
		Email:     existingUser.Email,
		Role:      existingUser.Role,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful", "user": userResponse})
}


// UpdateUser - оновлення даних користувача
func UpdateUser(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database connection error"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request data"})
	}

	userID := c.Params("id")
	if err := db.Model(&models.User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
		log.Println("Database update error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

// DeleteUser - видалення користувача
func DeleteUser(c *fiber.Ctx) error {
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Database connection error"})
	}

	userID := c.Params("id")
	if err := db.Delete(&models.User{}, userID).Error; err != nil {
		log.Println("Database delete error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}


