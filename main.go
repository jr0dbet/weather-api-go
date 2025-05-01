package main

import (
	"log"

	"github.com/jr0dbet/weather-api-go.git/cmd/server"
)

func main() {
	log.Println("🌦️  Starting Weather API Server...")

	if err := server.Run(); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}
