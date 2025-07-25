package mqtt

import (
	"fmt"
	"log"
	"time"

	"film-rental/pkg/kafka"
	"film-rental/pkg/monitoring"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartMQTTSubscriber() {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetClientID("film_mqtt_subscriber").
		SetKeepAlive(60 * time.Second).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			log.Printf("Received MQTT message: topic=%s payload=%s", msg.Topic(), msg.Payload())

			err := kafka.PublishFilmEvent(string(msg.Payload()))
			if err != nil {
				monitoring.SendEmailAlert("Failed to publish to Kafka", fmt.Sprintf("Failed to publish to Kafka: %v", err))
				log.Printf("Failed to publish to Kafka: %v", err)
			}
		}).SetCleanSession(false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		monitoring.SendEmailAlert("Failed to connect to MQTT broker", fmt.Sprintf("Failed to connect to MQTT broker: %v", token.Error()))
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	// Subscribe to your topic
	topic := "film/mqtt"
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("MQTT subscriber is listening on topic %s", topic)
}    
