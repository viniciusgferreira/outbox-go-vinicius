package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const consumerGroup = "outbox-consumer-group-1"

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          consumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	err = c.Subscribe("outbox.outbox-db.outbox_messages", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	defer c.Close()

	fmt.Println("Started consumer...")

	var msgCount int

	start := time.Now()
	for {
		msg, err := c.ReadMessage(-1) // wait indefinitely
		if err != nil {
			log.Printf("Consumer error: %v", err)
			continue
		}

		msgCount++

		fmt.Printf("Message #%d | Offset: %v | Total Time: %v\n",
			msgCount, msg.TopicPartition.Offset, time.Since(start))
	}
}
