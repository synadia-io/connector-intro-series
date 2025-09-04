package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
}

type TemperatureData struct {
	SensorID    string    `json:"sensor_id"`
	Temperature float64   `json:"temperature"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location,omitempty"`
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message:   "Welcome to the Go HTTP Server!",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tempData TemperatureData
	err := json.NewDecoder(r.Body).Decode(&tempData)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if tempData.Timestamp.IsZero() {
		tempData.Timestamp = time.Now()
	}

	log.Printf("Received temperature data: SensorID=%s, Temp=%.2f%s, Location=%s",
		tempData.SensorID, tempData.Temperature, tempData.Unit, tempData.Location)

	response := map[string]interface{}{
		"status":  "success",
		"message": "Temperature data received",
		"data":    tempData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/temperature", temperatureHandler)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /             - Home endpoint")
	fmt.Println("  POST /temperature  - Submit temperature data")

	log.Fatal(http.ListenAndServe(port, mux))
}