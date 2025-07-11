package mqtt

import (
	"log"
	"time"

	"film-rental/pkg/kafka"

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
				log.Printf("Failed to publish to Kafka: %v", err)
			}
		}).SetCleanSession(false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	// Subscribe to your topic
	topic := "film/mqtt"
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic %s: %v", topic, token.Error())
	}

	log.Printf("MQTT subscriber is listening on topic %s", topic)
}

// docker run --rm eclipse-mosquitto mosquitto_pub   -h host.docker.internal   -t film/mqtt   -m "test from CLI" // send message to mqtt broker
