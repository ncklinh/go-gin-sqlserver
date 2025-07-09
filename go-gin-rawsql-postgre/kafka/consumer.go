package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	maxRetries       = 3
	retryDelayMillis = 2000
)

var (
	processedCount int64
	failedCount    int64
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
			log.Printf("[%s] Failed to process message after %d retries: %s", consumerName, maxRetries, string(msg.Value))
			// Optional: push to dead-letter topic
		}
		if success {
			atomic.AddInt64(&processedCount, 1)
		} else {
			atomic.AddInt64(&failedCount, 1)
		}
	}
}

// Simulated business logic
func processMessage(msg kafka.Message) error {
	log.Printf("Processing message: %s", string(msg.Value))

	// Simulate random failure for testing
	if rand.Intn(5) == 0 {
		return errors.New("simulated failure")
	}
	return nil
}

func StartMetricsServer() {
	server := &http.Server{
		Addr: ":9090",
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "kafka_messages_processed_total %d\n", atomic.LoadInt64(&processedCount))
		fmt.Fprintf(w, "kafka_messages_failed_total %d\n", atomic.LoadInt64(&failedCount))
	})

	log.Println("Metrics available at http://localhost:9090/metrics")

	// Run in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Metrics server error: %v", err)
		}
	}()
}
