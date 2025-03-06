package repository

import (
	"computer-club/internal/models"
	"fmt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Session{}, &models.Computer{})

	// Проверяем, есть ли компьютеры в базе
	var count int64
	db.Model(&models.Computer{}).Count(&count)
	if count == 0 {
		fmt.Println("Создаем 7 компьютеров...")
		for i := 1; i <= 7; i++ {
			db.Create(&models.Computer{
				PCNumber: i,
				Status:   models.Free,
			})
		}
	}
}
