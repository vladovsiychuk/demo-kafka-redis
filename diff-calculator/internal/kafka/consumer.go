package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Next(ctx context.Context) (*Event, bool)
	Close() error
}

type kafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, group string) Consumer {
	return &kafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{broker},
			GroupID:        group,
			Topic:          topic,
			MinBytes:       1,    // Smallest allowed batch
			MaxBytes:       10e6, // 10MB max batch size
			StartOffset:    kafka.FirstOffset,
			CommitInterval: time.Second,
		}),
	}
}

func (c *kafkaConsumer) Next(ctx context.Context) (*Event, bool) {
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, false
	}
	var ev Event
	if err := json.Unmarshal(m.Value, &ev); err != nil {
		fmt.Printf("kafka consumer: failed to unmarshal event: %v\n", err)
		return nil, true
	}
	return &ev, true
}

func (c *kafkaConsumer) Close() error {
	return c.reader.Close()
}
