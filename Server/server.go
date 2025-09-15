package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	switch r.Method {
	case http.MethodGet:
		locations := []string{"Server Room", "Office", "Warehouse", "Loading Dock", "Rooftop"}
		sensorIDs := []string{"TEMP-001", "TEMP-002", "TEMP-003", "TEMP-004", "TEMP-005"}
		
		idx := rand.Intn(len(locations))
		
		baseTemp := 22.0
		if locations[idx] == "Server Room" {
			baseTemp = 18.0
		} else if locations[idx] == "Rooftop" {
			baseTemp = 15.0
		}
		
		variation := (rand.Float64() - 0.5) * 8
		temperature := baseTemp + variation
		
		tempData := TemperatureData{
			SensorID:    sensorIDs[idx],
			Temperature: temperature,
			Unit:        "Celsius",
			Location:    locations[idx],
			Timestamp:   time.Now(),
		}
		
		log.Printf("Generated temperature data: SensorID=%s, Temp=%.2f%s, Location=%s",
			tempData.SensorID, tempData.Temperature, tempData.Unit, tempData.Location)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tempData)
		
	case http.MethodPost:
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
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/temperature", temperatureHandler)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /             - Home endpoint")
	fmt.Println("  GET  /temperature  - Get random temperature reading")
	fmt.Println("  POST /temperature  - Submit temperature data")

	log.Fatal(http.ListenAndServe(port, mux))
}