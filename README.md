# Weather API

A simple and efficient RESTful API built with Go to manage historical weather data.

## Features

- Parse weather data from a `.dat` file
- Store, retrieve, and delete weather records via RESTful endpoints
- Simulate real-time data ingestion with a one-second interval
- Modular architecture with clear separation of concerns

## Project Structure

```
weather-api-go/
├── cmd/
│   └── importer/         # CLI script for importing .dat data via HTTP POST
├── config/               # Database configuration
├── handlers/             # HTTP request handlers
├── httphelpers/          # Utility functions for HTTP responses
├── models/               # Data models
├── parser/               # Parses .dat files into model structs
├── repositories/         # Database operations
├── routes/               # Route registration
├── server/               # Server startup logic
├── weather.dat           # Input data file
└── main.go               # Entry point for the API server
```

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/jr0dbet/weather-api-go.git
   cd weather-api-go
   ```

2. Start the API server:
   ```bash
   go run main.go
   ```

3. Import data (requires the server to be running):
   ```bash
   go run cmd/importer/import_data.go
   ```

## API Endpoints

- `GET /weather`: Retrieve all weather records
- `POST /weather`: Insert a new record (expects JSON)
- `GET /weather/range?from=YYYY-MM-DD&to=YYYY-MM-DD`: Filter by date range
- `DELETE /weather`: Delete all records and reset ID sequence

## Notes

- Ensure PostgreSQL is running and properly configured in `config/database.go`
- Timestamps must be in `RFC3339` format (`2006-01-02T15:04:05Z07:00`) when posting JSON data

## License

MIT