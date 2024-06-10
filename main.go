package main

import (
	"github.com/ZMS-DevOps/search-service/startup"
	cfg "github.com/ZMS-DevOps/search-service/startup/config"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"security.protocol": "sasl_plaintext",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "user1",
		"sasl.password":     config.KafkaAuthPassword,
		"group.id":          "hotel-service",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	consumer.SubscribeTopics([]string{"accommodation-rating.changed"}, nil)
	topicHandlers := map[string]func(*kafka.Message){
		"accommodation-rating.changed": server.AccommodationHandler.OnRatingChanged,
	}

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			if msg == nil {
				log.Printf("Received nil message")
				continue
			}

			handlerFunc, ok := topicHandlers[*msg.TopicPartition.Topic]
			if !ok {
				log.Printf("No handler for topic: %s\n", *msg.TopicPartition.Topic)
				continue
			}
			if handlerFunc == nil {
				log.Printf("Handler function for topic %s is nil", *msg.TopicPartition.Topic)
				continue
			}

			handlerFunc(msg)
		}
	}()

	server.Start()
}
