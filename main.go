package main

import (
	"context"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/startup"
	cfg "github.com/ZMS-DevOps/search-service/startup/config"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"log"
	"time"
)

func initJaegerTracer(jaegerHost string) (*sdktrace.TracerProvider, error) {
	log.Printf("Initializing tracing to jaeger at %s\n", jaegerHost)
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerHost)))
	if err != nil {
		return nil, err
	}
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(domain.ServiceName),
		)),
	), nil
}

func initPromtailClient(lokiHost string) (promtail.Client, error) {
	labels := "{source=\"" + domain.ServiceName + "\",service_name=\"" + "\"}"
	conf := promtail.ClientConfig{
		PushURL:            lokiHost,
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.ERROR,
	}

	return promtail.NewClientJson(conf)
}

func main() {
	config := cfg.NewConfig()

	var err error
	tp, err := initJaegerTracer(config.JaegerHost)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	loki, err := initPromtailClient(config.LokiHost)

	if err != nil {
		log.Printf("promtail.NewClient: %s\n", err)
	}

	server := startup.NewServer(config, tp, loki)
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
	loki.Shutdown()
}
