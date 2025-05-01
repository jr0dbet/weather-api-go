package models

import (
	"time"
)

type WeatherData struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Date        time.Time `gorm:"uniqueIndex" json:"date"`
	Humidity    int64     `json:"humidity"`
	Temperature float64   `json:"temperature"`
}
