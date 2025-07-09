package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const TopicFilmEvents = "film-events"

var filmWriter *kafka.Writer

func InitKafkaProducer() {
	filmWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    TopicFilmEvents,
		Balancer: &kafka.LeastBytes{}, // or RoundRobin
	})
}

func PublishFilmEvent(msg string) {
	err := filmWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
			Value: []byte(msg),
		},
	)
	if err != nil {
		log.Printf("Failed to publish film event: %v", err)
	}
}
