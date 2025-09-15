package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type TemperatureReading struct {
	SensorID    string    `json:"sensor_id"`
	Temperature float64   `json:"temperature"`
	Unit        string    `json:"unit"`
	Location    string    `json:"location"`
	Timestamp   time.Time `json:"timestamp"`
}

func main() {
	credsPath := filepath.Join("..", "Credentials", "NGS-Premium-CLI.creds")
	natsURL := "tls://connect.ngs.global"
	streamName := "Temperatures"
	subject := "telemetry.sensors.temperature"

	if _, err := os.Stat(credsPath); os.IsNotExist(err) {
		log.Fatalf("Credentials file not found: %s", credsPath)
	}

	opts := []nats.Option{
		nats.UserCredentials(credsPath),
		nats.Name("Temperature Sensor JetStream Publisher"),
	}

	nc, err := nats.Connect(natsURL, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Printf("Successfully connected to NATS at %s", natsURL)

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Failed to create JetStream context: %v", err)
	}

	ctx := context.Background()

	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		log.Printf("Stream '%s' not found, creating it...", streamName)

		streamConfig := jetstream.StreamConfig{
			Name:      streamName,
			Subjects:  []string{"telemetry.sensors.>"},
			Retention: jetstream.LimitsPolicy,
			MaxAge:    time.Hour * 24 * 7,
			Storage:   jetstream.FileStorage,
			Replicas:  1,
			MaxBytes:  10 * 1024 * 1024, // 10 MB
		}

		stream, err = js.CreateStream(ctx, streamConfig)
		if err != nil {
			log.Fatalf("Failed to create stream '%s': %v", streamName, err)
		}
		log.Printf("Successfully created stream '%s'", streamName)
	} else {
		log.Printf("Successfully connected to existing stream '%s'", streamName)
	}
	log.Printf("Publishing temperature data to subject: %s", subject)
	log.Println("Publishing temperature readings every 5 seconds...")

	locations := []string{"Server Room", "Office", "Warehouse", "Loading Dock", "Rooftop"}
	sensorIDs := []string{"TEMP-001", "TEMP-002", "TEMP-003", "TEMP-004", "TEMP-005"}

	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	publishReading := func() {
		idx := rand.Intn(len(locations))

		baseTemp := 22.0
		if locations[idx] == "Server Room" {
			baseTemp = 18.0
		} else if locations[idx] == "Rooftop" {
			baseTemp = 15.0
		}

		variation := (rand.Float64() - 0.5) * 8
		temperature := baseTemp + variation

		reading := TemperatureReading{
			SensorID:    sensorIDs[idx],
			Temperature: temperature,
			Unit:        "Celsius",
			Location:    locations[idx],
			Timestamp:   time.Now(),
		}

		data, err := json.Marshal(reading)
		if err != nil {
			log.Printf("Error marshalling temperature reading: %v", err)
			return
		}

		ack, err := js.Publish(ctx, subject, data)
		if err != nil {
			log.Printf("Error publishing to JetStream: %v", err)
			return
		}

		log.Printf("Published: Sensor %s at %s - %.1f°C (Stream: %s, Seq: %d)",
			reading.SensorID, reading.Location, reading.Temperature,
			ack.Stream, ack.Sequence)
	}

	info, err := stream.Info(ctx)
	if err == nil {
		log.Printf("Stream info - Messages: %d, Bytes: %d",
			info.State.Msgs, info.State.Bytes)
	}

	publishReading()

	for range ticker.C {
		publishReading()
	}
}
