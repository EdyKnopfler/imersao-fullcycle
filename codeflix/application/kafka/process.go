package kafka

import (
	"fmt"
	"os"
	"time"

	"derso.com/imersao-fullcycle/codepix-go/application/dto"
	"derso.com/imersao-fullcycle/codepix-go/application/factory"
	"derso.com/imersao-fullcycle/codepix-go/application/usecase"
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
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	consumer, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{
		os.Getenv("kafkaTransactionTopic"),
		os.Getenv("kafkaTransactionConfirmationTopic"),
	}
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
	transactionsTopic := os.Getenv("kafkaTransactionTopic")
	transactionConfirmationTopic := os.Getenv("kafkaTransactionConfirmationTopic")

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic:", topic, ", msg:", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transactionDto, err := dto.NewTransactionDTO(msg.Value)
	if err != nil {
		fmt.Println("error parsing transaction JSON", err)
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
		// TODO o banco de origem tem que ser notificado em caso de erro, né não? rs
		// Isso exigiria alguns refactorings para tornar isto um pouco mais flexível.
		// Por enquanto somente pegando a ideia geral está bom.
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

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transactionDto, err := dto.NewTransactionDTO(msg.Value)
	if err != nil {
		fmt.Println("error parsing transaction JSON", err)
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.Database)

	if transactionDto.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(transactionDto, transactionUseCase)
		if err != nil {
			fmt.Println("error when confirming transaction ", err)
			return err
		}
	} else if transactionDto.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transactionDto.ID)
		if err != nil {
			fmt.Println("error when completing transaction ", err)
			return err
		}
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(
	transaction *dto.TransactionDTO,
	transactionUseCase usecase.TransactionUseCase,
) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	err = Publish(string(transactionJson), topic, k.Producer, k.DeliveryChannel)
	if err != nil {
		return err
	}

	return nil
}
