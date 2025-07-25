package main

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := "tcp://localhost:1883"
	topic := "film/mqtt"
	clientID := "test_publisher"

	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetCleanSession(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	var wg sync.WaitGroup
	concurrency := 1000 // number of concurrent publishers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			msg := fmt.Sprintf("test message %d at %s", id, time.Now().Format(time.RFC3339))
			token := client.Publish(topic, 1, false, msg)
			token.Wait()
			fmt.Printf("Published: %s\n", msg)
		}(i)
	}

	wg.Wait()
	fmt.Println("All messages sent")
}
