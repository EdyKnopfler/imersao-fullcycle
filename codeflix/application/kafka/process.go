package kafka

import (
	"fmt"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/gorm"
)

type KafkaProcessor struct {
	Database        *gorm.DB
	Producer        *ckafka.Producer
	DeliveryChannel chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		Database:        database,
		Producer:        producer,
		DeliveryChannel: deliveryChan,
	}
}

func (k *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092", // TODO externalizar
		"group.id":          "consumergroup",
		"auto.offset.reset": "earliest",
	}
	consumer, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{"teste"}
	consumer.SubscribeTopics(topics, nil)
	fmt.Println("Kafka consumer has been started")

	for {
		timeout := time.Duration(-1)
		msg, err := consumer.ReadMessage(timeout)
		if err == nil {
			fmt.Println("Message received:", string(msg.Value))
		}
	}
}
