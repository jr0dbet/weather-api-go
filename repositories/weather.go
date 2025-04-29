package repositories

import (
	"time"

	"github.com/jr0dbet/weather-api-go.git/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&models.WeatherData{})
}

func InsertWeather(data models.WeatherData) error {
	return db.Create(&data).Error
}

func GetWeatherByDate(date time.Time) (models.WeatherData, error) {
	var data models.WeatherData
	err := db.Where("date = ?", date).First(&data).Error
	return data, err
}

func GetWeatherByRange(from time.Time, to time.Time) ([]models.WeatherData, error) {
	var data []models.WeatherData
	err := db.Where("date BETWEEN ? AND ?", from, to).Find(&data).Error
	return data, err
}
