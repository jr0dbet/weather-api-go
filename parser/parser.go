package parser

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jr0dbet/weather-api-go.git/models"
)

func ParseWeatherData(filepath string) ([]models.WeatherData, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []models.WeatherData
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		if len(parts) != 3 {
			log.Println("Ignored line, incorrect number of fields:", line)
			continue
		}

		date, err := time.Parse("2006-01-02", parts[0])
		if err != nil {
			log.Println("Error parsing date:", parts[0], err)
			continue
		}

		humidityFloat, _ := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Println("Error parsing humidity:", parts[1], err)
			continue
		}
		humidity := int64(math.Round(humidityFloat))

		temperature, _ := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Println("Error parsing temperature:", parts[2], err)
			continue
		}

		results = append(results, models.WeatherData{
			Date:        date,
			Humidity:    humidity,
			Temperature: temperature,
		})

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
