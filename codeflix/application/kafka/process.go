package kafka

import (
	"fmt"
	"time"

	"derso.com/imersao-fullcycle/codepix-go/application/dto"
	"derso.com/imersao-fullcycle/codepix-go/application/factory"
	"derso.com/imersao-fullcycle/codepix-go/domain/model"
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
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions" // TODO em variável de ambiente
	transactionConfirmationTopic := "transaction-confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
	default:
		fmt.Println("not a valid topic", topic, "msg:", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transactionDto, err := dto.NewTransactionDTO(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)
	createdTransaction, err := transactionUseCase.Register(
		transactionDto.AccountID,
		transactionDto.Amount,
		transactionDto.PixKeyTo,
		transactionDto.PixKeyKindTo,
		transactionDto.Description,
	)

	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	destinationBankTopic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transactionDto.ID = createdTransaction.ID
	transactionDto.Status = model.TransactionPending
	transactionJson, err := transactionDto.ToJson()

	if err != nil {
		fmt.Println("error generating JSON from transaction", err)
		return err
	}

	err = Publish(string(transactionJson), destinationBankTopic, k.Producer, k.DeliveryChannel)
	if err != nil {
		fmt.Println("error sending message to destination bank", err)
		return err
	}

	return nil
}
