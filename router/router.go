package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jr0dbet/weather-api-go.git/api"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	r.HandleFunc("/weather", api.HandleWeatherIngest).Methods("POST")
	r.HandleFunc("/weather/all", api.HandleAllWeatherRecords).Methods("GET")
	r.HandleFunc("/weather", api.HandleWeatherByDate).Methods("GET").Queries("date", "{date}")
	r.HandleFunc("/weather", api.HandleWeatherByRange).Methods("GET").Queries("from", "{from}", "to", "{to}")
	r.HandleFunc("/ws", api.HandleWebSocket)

	return r
}
