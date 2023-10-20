package main

import (
	"fmt"
	"os"

	"derso.com/imersao-fullcycle/codepix-go/application/kafka"
	"derso.com/imersao-fullcycle/codepix-go/infrastructure/db"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	fmt.Println("Produzindo mensagem...")

	producer := kafka.NewKafkaProducer()
	deliveryChannel := make(chan ckafka.Event)
	topic := "teste" // Cria pelo control center ou linha de comando :)

	err := kafka.Publish("Hellô Cáfica", topic, producer, deliveryChannel)
	if err != nil {
		panic(err)
	}

	go kafka.DeliveryReport(deliveryChannel)

	database := db.ConnectDB(os.Getenv("env"))
	processor := kafka.NewKafkaProcessor(database, producer, deliveryChannel)
	processor.Consume() // loop (segura também o loop do DeliveryReport)
}
