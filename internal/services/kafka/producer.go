package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Broker string
}

func NewKafkaProducer(broker string) *KafkaProducer {
	return &KafkaProducer{
		Broker: broker,
	}
}

func (p *KafkaProducer) SendMessage(topic string, value []byte) error {

	const maxRetries = 3
	const retryDelay = 1 * time.Second

	writer := &kafka.Writer{
		Addr:         kafka.TCP(p.Broker),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 5 * time.Second,
		BatchTimeout: 50 * time.Millisecond,
		BatchSize:    50,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := writer.WriteMessages(ctx,
			kafka.Message{
				Value: value,
			},
		)

		if err == nil {
			return nil
		}

		log.Printf("Failed to write message to topic %s: %v", topic, err)

		if attempt == maxRetries-1 {
			return err
		}

		time.Sleep(retryDelay)
	}

	return nil
}
