package main

import (
	"fmt"

	"derso.com/imersao-fullcycle/codepix-go/application/kafka"
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

	/*
		Para a mensagem ser entregue este loop precisa ficar rodando e reportar.
		Se chamamos como goroutine, a publicação de uma única mensagem não acontece pois ela
		morre junto com a main :()
	*/
	kafka.DeliveryReport(deliveryChannel)
}
