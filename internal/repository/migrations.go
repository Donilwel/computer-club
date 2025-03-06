package repository

import (
	"computer-club/internal/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Session{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
