package test

import (
	"os"
	"testing"
	"time"

	"github.com/jr0dbet/weather-api-go.git/parser"
	"github.com/stretchr/testify/assert"
)

func createTempWeatherFile(t *testing.T, mockContent string) string {
	tmpFile, err := os.CreateTemp("", "weather_test_*.dat")
	assert.NoError(t, err)

	_, err = tmpFile.WriteString(mockContent)
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func TestParseWeatherData_ValidData(t *testing.T) {
	mockContent := "2023-01-01 80 22.5\n2023-01-02 65 18.3"
	filepath := createTempWeatherFile(t, mockContent)
	defer os.Remove(filepath)

	data, err := parser.ParseWeatherData(filepath)
	assert.NoError(t, err)
	assert.Len(t, data, 2)

	expectedDate, _ := time.Parse("2006-01-02", "2023-01-01")
	assert.Equal(t, expectedDate, data[0].Date)
	assert.Equal(t, 80, data[0].Humidity)
	assert.Equal(t, 22.5, data[0].Temperature)
}

func TestParseWeatherData_InvalidLineSkipped(t *testing.T) {
	mockContent := "2023-01-01 80 22.5\nInvalid line\n2023-01-02 65 18.3"
	filepath := createTempWeatherFile(t, mockContent)
	defer os.Remove(filepath)

	data, err := parser.ParseWeatherData(filepath)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestParseWeatherData_BadDateFormat(t *testing.T) {
	mockContent := "01-01-2023 80 22.5\n2023-01-02 65 18.3"
	filepath := createTempWeatherFile(t, mockContent)
	defer os.Remove(filepath)

	data, err := parser.ParseWeatherData(filepath)
	assert.NoError(t, err)
	assert.Len(t, data, 1)

	expectedDate, _ := time.Parse("2006-01-02", "2023-01-02")
	assert.Equal(t, expectedDate, data[0].Date)
}
