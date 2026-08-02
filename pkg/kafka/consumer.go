package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer handles reading event messages from Apache Kafka topics
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer initializes a Kafka Reader instance for a specific consumer group
func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10,   // 10 Bytes
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})

	return &Consumer{reader: reader}
}

// ReadMessage fetches the next message from the Kafka topic
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("failed to read kafka message: %w", err)
	}
	return msg, nil
}

// Close gracefully shuts down the Kafka Consumer connection
func (c *Consumer) Close() error {
	return c.reader.Close()
}
