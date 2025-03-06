package repository

import (
	"computer-club/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewPostgresDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}
