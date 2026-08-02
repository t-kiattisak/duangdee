package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer handles publishing events to Apache Kafka topics
type Producer struct {
	writer *kafka.Writer
}

// EventMessage defines the standardized JSON envelope for all Kafka events
type EventMessage struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// NewProducer initializes a new Kafka Writer instance
func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return &Producer{writer: writer}
}

// Publish sends a structured event payload to a specified Kafka topic
func (p *Producer) Publish(ctx context.Context, topic string, key string, event EventMessage) error {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal kafka event payload: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write kafka message to topic %s: %w", topic, err)
	}

	return nil
}

// Close gracefully stops the Kafka Producer connection
func (p *Producer) Close() error {
	return p.writer.Close()
}
