package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jr0dbet/weather-api-go.git/httphelpers"
	"github.com/jr0dbet/weather-api-go.git/models"
	"github.com/jr0dbet/weather-api-go.git/repositories"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var wsClients = make(map[*websocket.Conn]bool)

func HandleAllWeatherRecords(w http.ResponseWriter, r *http.Request) {
	data, err := repositories.GetAllWeather()
	if err != nil {
		httphelpers.RespondError(w, http.StatusInternalServerError, "Error retrieving weather data: "+err.Error())
		return
	}
	httphelpers.RespondJSON(w, http.StatusOK, data)
}

func HandleWeatherIngest(w http.ResponseWriter, r *http.Request) {
	var entry models.WeatherData

	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		httphelpers.RespondError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := repositories.InsertWeather(entry); err != nil {
		httphelpers.RespondError(w, http.StatusInternalServerError, "Error while saving data: "+err.Error())
	}

	broadcastWeather(entry)
	httphelpers.RespondJSON(w, http.StatusCreated, map[string]string{"status": "Weather data save sucessfully"})
}

func HandleWeatherByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := httphelpers.ParseDateParam(dateStr)
	if err != nil {
		httphelpers.RespondError(w, http.StatusBadRequest, "Invalid date format: "+err.Error())
		return
	}

	data, err := repositories.GetWeatherByDate(date)
	if err != nil {
		if err.Error() == "record not found" {
			httphelpers.RespondError(w, http.StatusNotFound, "No data found for the given date")
		} else {
			httphelpers.RespondError(w, http.StatusInternalServerError, "Error fetching data: "+err.Error())
		}
		return
	}
	httphelpers.RespondJSON(w, http.StatusOK, data)
}

func HandleWeatherByRange(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err := httphelpers.ParseDateParam(fromStr)
	if err != nil {
		httphelpers.RespondError(w, http.StatusBadRequest, "Invalid date format 'from': "+err.Error())
		return
	}

	to, err := httphelpers.ParseDateParam(toStr)
	if err != nil {
		httphelpers.RespondError(w, http.StatusBadRequest, "Invalid date format 'to': "+err.Error())
		return
	}

	data, err := repositories.GetWeatherByRange(from, to)
	if err != nil {
		httphelpers.RespondError(w, http.StatusInternalServerError, "Error fetching data: "+err.Error())
		return
	}

	httphelpers.RespondJSON(w, http.StatusOK, data)
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		httphelpers.RespondError(w, http.StatusInternalServerError, "Error establishing WebSocket"+err.Error())
		return
	}
	defer conn.Close()

	wsClients[conn] = true

	for {
		if _, _, err := conn.NextReader(); err != nil {
			delete(wsClients, conn)
			break
		}
	}
}

func broadcastWeather(data models.WeatherData) {
	for client := range wsClients {
		if err := client.WriteJSON(data); err != nil {
			client.Close()
			delete(wsClients, client)
		}
	}
}
