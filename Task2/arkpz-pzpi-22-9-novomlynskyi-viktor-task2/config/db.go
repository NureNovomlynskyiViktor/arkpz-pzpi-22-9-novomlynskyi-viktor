package config

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

// InitDB ініціалізує з'єднання з базою даних.
func InitDB() *gorm.DB {
	dbUser := "postgres"    
	dbPass := "16072005!"         
	dbName := "ARTGUARD" 
	dbHost := "localhost"     
	dbPort := "5432"          

	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", dbHost, dbUser, dbName, dbPass, dbPort)

	// Підключаємося до бази даних.
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Перевірка з'єднання
	sqlDB, err := DB.DB() // Отримуємо доступ до звичайного SQL драйвера для виконання запиту
	if err != nil {
		log.Fatal("Failed to get DB connection:", err)
	}

	// Перевірка з'єднання з базою даних
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	} else {
		log.Println("Successfully connected to the database.")
	}

	

	DB.AutoMigrate()

	
	return DB
}





