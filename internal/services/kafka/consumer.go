package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	saving_book "SavingBooks/internal/saving-book"
	"SavingBooks/internal/services/kafka/event"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	broker       string
	groupId      string
	readers      map[string]*kafka.Reader
	SavingBookUC saving_book.UseCase
}

func NewKafkaConsumer(broker string, groupId string, savingBookUC saving_book.UseCase) *KafkaConsumer {
	return &KafkaConsumer{broker: broker, groupId: groupId, readers: make(map[string]*kafka.Reader),
		SavingBookUC: savingBookUC}
}

func (kc *KafkaConsumer) FetchTopics() ([]string, error) {
	conn, err := kafka.Dial("tcp", kc.broker)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}
	topicSet := make(map[string]struct{})
	for _, partition := range partitions {
		topicSet[partition.Topic] = struct{}{}
	}
	topics := make([]string, 0, len(topicSet))
	for topic := range topicSet {
		topics = append(topics, topic)
	}
	return topics, nil
}
func (kc *KafkaConsumer) StartListening() error {
	topics, err := kc.FetchTopics()
	if err != nil {
		return err
	}
	fmt.Println("Listening to topics:")
	for _, topic := range topics {
		if topic != "__consumer_offsets" {
			fmt.Println(" -", topic)

			kc.readers[topic] = kc.createReader(topic)
			go kc.listenToTopic(topic)
		}
	}
	return nil
}
func (kc *KafkaConsumer) createReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kc.broker},
		Topic:   topic,
		GroupID: kc.groupId,
	})
}
func (kc *KafkaConsumer) listenToTopic(topic string) {
	r := kc.readers[topic]
	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message from topic %s: %v\n", topic, err)
			continue
		}
			switch topic {
			case CaptureOrderTopic:
				var withDrawEvent event.WithDrawEvent
				err = json.Unmarshal(msg.Value, &withDrawEvent)
				if err != nil {
					log.Printf("Error unmarshaling message from topic %s: %v\n", topic, err)
					continue
				}
				err = kc.SavingBookUC.HandleWithdraw(context.Background(), &withDrawEvent)
				if err != nil {
					log.Printf("Error handling withdraw: %v\n", err)
				}

			}

		log.Printf("Message from topic %s: %s = %s\n", msg.Topic, string(msg.Key), string(msg.Value))
	}
}
