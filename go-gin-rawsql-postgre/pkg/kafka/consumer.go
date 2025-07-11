package kafka

import (
	"context"
	"errors"
	db "film-rental/pkg/db/gorm"
	"film-rental/pkg/monitoring/model"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries       = 2
	retryDelayMillis = 2000
)

func StartFilmConsumer(consumerName string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicFilmEvents,
		GroupID:  "film-consumer-group", // same group for all
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	log.Printf("[%s] Starting consumer...", consumerName)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[%s] Error reading message: %v", consumerName, err)
			continue
		}
		log.Printf("[%s] Received message from partition %d: %s", consumerName, msg.Partition, string(msg.Value))

		success := false
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err := processMessage(msg)
			if err != nil {
				log.Printf("[%s] Retry %d/%d failed: %v", consumerName, attempt, maxRetries, err)
				time.Sleep(retryDelayMillis * time.Millisecond)
			} else {
				success = true
				break
			}
		}
		if !success {
			logEvent(consumerName, "kafka_messages_failed_total", string(msg.Value))
		} else {
			logEvent(consumerName, "kafka_messages_processed_total", string(msg.Value))
		}
	}
}

func logEvent(service, message, context string) {
	eventLog := model.EventLog{
		Service: service,
		Message: message,
		Context: context,
	}
	if err := db.DB.Create(&eventLog).Error; err != nil {
		log.Printf("[%s] Failed to insert event log: %v", service, err)
	}
}

// Simulated business logic
func processMessage(msg kafka.Message) error {
	log.Printf("Processing message: %s", string(msg.Value))

	// Simulate random failure for testing
	if rand.Intn(2) == 0 {
		return errors.New("simulated failure")
	}
	return nil
}

func StartMetricsServer() {
	server := &http.Server{
		Addr: ":9090",
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		errorCounts := CountErrorLogsByEvent()
		for _, e := range errorCounts {
			fmt.Fprintf(w, "%s %d\n", e.Message, e.Count)
		}
	})

	log.Println("Metrics available at http://localhost:9090/metrics")

	// Run in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Metrics server error: %v", err)
		}
	}()
}

type ErrorCountByEvent struct {
	Message string
	Count   int64
}

func CountErrorLogsByEvent() []ErrorCountByEvent {
	var results []ErrorCountByEvent
	if err := db.DB.
		Model(&model.EventLog{}).
		Select("message, COUNT(*) as count").
		Group("message").
		Scan(&results).Error; err != nil {
		log.Printf("Error counting error logs by event: %v", err)
	}
	return results
}
