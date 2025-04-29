package models

import (
	"time"

	"gorm.io/gorm"
)

type WeatherData struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Date        time.Time      `gorm:"uniqueIndex" json:"Date"`
	Humidity    int            `json:"humidity"`
	Temperature float64        `json:"temperature"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
