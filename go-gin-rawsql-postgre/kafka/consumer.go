package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartFilmConsumer(consumerName string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicFilmEvents,
		GroupID:  "film-consumer-group", // same group for all
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	log.Printf("[%s] Consumer started", consumerName)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[%s] Error reading message: %v", consumerName, err)
			continue
		}
		log.Printf("[%s] Got from partition %d: %s", consumerName, msg.Partition, string(msg.Value))
	}
}
