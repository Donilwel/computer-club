package generate

import (
	"computer-club/internal/repository/models"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

var adjectives = []string{
	"Fast", "Crazy", "Cool", "Brave", "Smart", "Lucky", "Wild", "Slowed", "Bad", "Good",
	"Sick", "Punished", "Elite", "Sweet", "Fierce", "Mighty", "Stealthy", "Cunning", "Fearless",
	"Epic", "Bold", "Stormy", "Loyal", "Daring", "Powerful", "Savage",
}

var nouns = []string{
	"Tiger", "Eagle", "Wolf", "Shark", "Panther", "Hawk", "Dragon", "Chicken", "Pow", "Dog",
	"Cat", "Pig", "Lion", "Cobra", "Phoenix", "Grizzly", "Rhino", "Jaguar", "Falcon",
	"Raven", "Bison", "Scorpion", "Viper", "Bear", "Fox", "Whale",
}

func GenerateUsername(db *gorm.DB) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		adj := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		number := rand.Intn(9000) + 1000
		username := fmt.Sprintf("%s%s%d", adj, noun, number)

		var count int64
		db.Model(&models.User{}).Where("name = ?", username).Count(&count)
		if count == 0 {
			return username
		}
	}
}
