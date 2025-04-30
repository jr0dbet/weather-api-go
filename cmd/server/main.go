package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/jr0dbet/weather-api-go.git/config"
	"github.com/jr0dbet/weather-api-go.git/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("⚠️ Cannot load .env file")
	}

	config.InitDB()

	r := router.InitRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(": "+port, r))
}
