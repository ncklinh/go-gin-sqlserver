package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartRentalConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   TopicRentalEvents,
		GroupID: "rental-consumer-group",
	})
	log.Println("rental-consumer-group")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}
		log.Printf("Received message: %s = %s", string(msg.Key), string(msg.Value))
	}
}
