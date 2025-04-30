package models

import (
	"time"
)

type WeatherData struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"uniqueIndex" json:"Date"`
	Humidity    int       `json:"humidity"`
	Temperature float64   `json:"temperature"`
}
