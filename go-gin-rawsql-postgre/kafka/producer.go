package kafka

import (
	"context"

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

func PublishFilmEvent(message string) error {
	msg := kafka.Message{
		Key:   []byte("film-key"),
		Value: []byte(message),
	}
	return filmWriter.WriteMessages(context.Background(), msg)
}
