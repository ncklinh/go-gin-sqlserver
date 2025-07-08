package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var TopicRentalEvents = "rental-events"

func PublishRentalEvent(message string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicRentalEvents,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("film"),
			Value: []byte(message),
		},
	)
	if err != nil {
		log.Printf("failed to write message: %v", err)
	}
}
