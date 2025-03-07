package repository

import (
	"computer-club/internal/models"
	"fmt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Session{}, &models.Computer{}, &models.Tariff{}, &models.Wallet{}, &models.Transaction{})

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

	var countTariffs int64
	db.Model(&models.Tariff{}).Count(&countTariffs)
	if countTariffs == 0 {
		tariffs := []models.Tariff{
			{ID: 1, Name: "1 час", Price: 100, Duration: 60},
			{ID: 2, Name: "3 часа", Price: 250, Duration: 180},
			{ID: 3, Name: "5 часов", Price: 400, Duration: 300},
			{ID: 4, Name: "8 часов", Price: 600, Duration: 480},
			{ID: 5, Name: "Всю ночь", Price: 500, Duration: 480},
			{ID: 6, Name: "Ночной (4ч)", Price: 350, Duration: 240},
			{ID: 7, Name: "Ночной (6ч)", Price: 450, Duration: 360},
		}

		for _, tariff := range tariffs {
			db.Create(&tariff)
		}
	}
}
