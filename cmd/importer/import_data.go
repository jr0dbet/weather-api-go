package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jr0dbet/weather-api-go.git/parser"
)

const (
	filepath = "weather.dat"
	apiURL   = "http://localhost:8080/weather"
)

func main() {
	weatherData, err := parser.ParseWeatherData(filepath)
	if err != nil {
		log.Fatalf("❌ Error parsing weather data: %v\n", err)
	}

	for _, record := range weatherData {
		log.Printf("📦 Sending: %+v\n", record)
		jsonData, err := json.Marshal(record)
		if err != nil {
			log.Printf("⚠️ Failed to marshal record: %v", err)
			continue
		}

		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("⚠️ Error sending POST request: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusCreated {
			log.Printf("⚠️ Unexpected response code: %d", resp.StatusCode)
		} else {
			log.Printf("✅ Inserted: %s", record.Date.Format("2006-01-02"))
		}

		time.Sleep(1 * time.Second)

	}

}
