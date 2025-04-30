package repositories

import (
	"time"

	"github.com/jr0dbet/weather-api-go.git/config"
	"github.com/jr0dbet/weather-api-go.git/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetAllWeather() ([]models.WeatherData, error) {
	var data []models.WeatherData
	err := config.DB.Order("date").Find(&data).Error
	return data, err
}

func InsertWeather(data models.WeatherData) error {
	return config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "date"}},
		DoNothing: true,
	}).Create(&data).Error
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

func DeleteAllWeather() error {
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("1 = 1").Delete(&models.WeatherData{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("TRUNCATE TABLE weather_data RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
