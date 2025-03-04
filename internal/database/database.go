package database

import (
	"log"
	"rateLimiter/internal/config"
	"rateLimiter/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	// auto migrate
	err = DB.AutoMigrate(&models.RateLimit{}, &models.TokenBucket{})
	if err != nil {
		log.Fatalf("Failed to migrate database : %v", err)
	}

	log.Println("Database Initialized Succesfully")
}

func GetDB() *gorm.DB {
	return DB
}
