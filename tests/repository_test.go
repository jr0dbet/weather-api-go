package test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jr0dbet/weather-api-go.git/api"
	"github.com/jr0dbet/weather-api-go.git/config"
	"github.com/jr0dbet/weather-api-go.git/models"
	"github.com/jr0dbet/weather-api-go.git/repositories"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error while opening a mock database connection: %v", err)
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the mock database: %v", err)
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func wrapAsGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func TestGetAllWeather(t *testing.T) {
	type testCase struct {
		name         string
		mockSetup    func(sqlmock.Sqlmock)
		expectedCode int
		expectedBody string
	}

	mockWeatherData := []models.WeatherData{
		{ID: 1, Date: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC), Humidity: 80, Temperature: 20.5},
		{ID: 2, Date: time.Date(2023, time.January, 2, 0, 0, 0, 0, time.UTC), Humidity: 75, Temperature: 21.8},
	}

	tests := []testCase{
		{
			name: "Sucessful get",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mockRows := sqlmock.NewRows([]string{"id", "date", "humidity", "temperature"}).
					AddRow(mockWeatherData[0].ID, mockWeatherData[0].Date, mockWeatherData[0].Humidity, mockWeatherData[0].Temperature).
					AddRow(mockWeatherData[1].ID, mockWeatherData[1].Date, mockWeatherData[1].Humidity, mockWeatherData[1].Temperature)
				mock.ExpectQuery(`SELECT \* FROM "weather_data" ORDER BY date`).WillReturnRows(mockRows)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":1,"date":"2023-01-01T00:00:00Z","humidity":80,"temperature":20.5},{"id":2,"date":"2023-01-02T00:00:00Z","humidity":75,"temperature":21.8}]`,
		},
		{
			name: "db error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "weather_data" ORDER BY date`).WillReturnError(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error": "Error retrieving weather data: db error"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := createMockDB()
			require.NoError(t, err)

			tc.mockSetup(mock)
			config.DB = db

			router := gin.Default()
			router.GET("/weather", wrapAsGinHandler(api.HandleAllWeatherRecords))

			req := httptest.NewRequest(http.MethodGet, "/weather", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)
			require.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}

}

func TestInsertWeather(t *testing.T) {
	type testCase struct {
		name        string
		prepareMock func(mock sqlmock.Sqlmock, data models.WeatherData)
		expectedErr bool
	}

	mockDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	mockWeatherData := models.WeatherData{
		Date:        mockDate,
		Humidity:    80,
		Temperature: 22.5,
	}

	tests := []testCase{
		{
			name: "Sucessful insert",
			prepareMock: func(mock sqlmock.Sqlmock, data models.WeatherData) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "weather_data".*RETURNING "id"`).
					WithArgs(data.Date, data.Humidity, data.Temperature).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedErr: false,
		},
		{
			name: "Conflict - no insert",
			prepareMock: func(mock sqlmock.Sqlmock, data models.WeatherData) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "weather_data".*RETURNING "id"`).
					WithArgs(data.Date, data.Humidity, data.Temperature).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
				mock.ExpectCommit()
			},
			expectedErr: false,
		},
		{
			name: "db error",
			prepareMock: func(mock sqlmock.Sqlmock, data models.WeatherData) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "weather_data".*RETURNING "id"`).
					WithArgs(data.Date, data.Humidity, data.Temperature).
					WillReturnError(fmt.Errorf("db insert failed"))
				mock.ExpectRollback()
			},
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gormDB, mock, err := createMockDB()
			require.NoError(t, err)

			tc.prepareMock(mock, mockWeatherData)

			config.DB = gormDB
			err = repositories.InsertWeather(mockWeatherData)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}

}
