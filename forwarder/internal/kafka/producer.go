package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github/vladovsiychuk/demo-kafka-redis-forwarder/internal/reporter"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Send(ctx context.Context, partner string, data *reporter.Data) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(bootstrapServer, topic string) Producer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(bootstrapServer),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (k *kafkaProducer) Send(ctx context.Context, partner string, data *reporter.Data) error {
	event := Event{
		Partner: partner,
		Data:    data,
		SentAt:  time.Now().UTC(),
	}
	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(partner),
		Value: bytes,
	}
	return k.writer.WriteMessages(ctx, msg)
}

func (k *kafkaProducer) Close() error {
	return k.writer.Close()
}

func CreateTopicIfNotExists(broker, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller: %w", err)
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("dial controller: %w", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
	if err != nil && err.Error() != "topic already exists" {
		return fmt.Errorf("create topic: %w", err)
	}
	time.Sleep(time.Second) // let Kafka propagate topic metadata
	return nil
}
