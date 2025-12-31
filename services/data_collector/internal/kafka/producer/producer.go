package producer

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer() *KafkaProducer {
	broker := os.Getenv("KAFKA_ADDR")
	topic := os.Getenv("KAFKA_TOPIC")

	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(broker),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (p *KafkaProducer) Close() {
	p.writer.Close()
}

func (p *KafkaProducer) Send(ctx context.Context, key string, msg any) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: b,
		Time:  time.Now(),
	})
}
