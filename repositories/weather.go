package repositories

import (
	"time"

	"github.com/jr0dbet/weather-api-go.git/config"
	"github.com/jr0dbet/weather-api-go.git/models"
)

func GetAllWeather() ([]models.WeatherData, error) {
	var data []models.WeatherData
	err := config.DB.Order("date").Find(&data).Error
	return data, err
}

func InsertWeather(data models.WeatherData) error {
	return config.DB.Create(&data).Error
}

func GetWeatherByDate(date time.Time) (models.WeatherData, error) {
	var data models.WeatherData
	err := config.DB.Where("date = ?", date).First(&data).Error
	return data, err
}

func GetWeatherByRange(from time.Time, to time.Time) ([]models.WeatherData, error) {
	var data []models.WeatherData
	err := config.DB.Where("date BETWEEN ? AND ?", from, to).Find(&data).Error
	return data, err
}
