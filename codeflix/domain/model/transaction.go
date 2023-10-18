package model

import (
	"time"
	"errors"
	uuid "github.com/satori/go.uuid"
	govalidator "github.com/asaskevich/govalidator"
)

const (
	TransactionPending string = "pending"
	TransactionCompleted string = "completed"
	TransactionError string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base `valid:"required"`
	AccountFrom *Account `valid:"-"`
	Amount float64 `json:"amount" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	Status string `json:"status" valid:"notnull"`
	Description string `json:"description" valid:"notnull"`
	CancelDescription string `json:"cancelDescription" valid:"-"`
}

func NewTransaction(
	accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		Status: TransactionPending,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(t)

	if err != nil {
		return err
	}

	if t.Amount <= 0 {
		return errors.New("amount must be greather than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError  && t.Status != TransactionConfirmed {
		return errors.New("invalid status")
	}

	if t.PixKeyTo.AccountID == t.AccountFrom.ID {
		return errors.New("same account")
	}

	return err
}

// TODO Daria para refatorar em um método "privado atualizarStatus"
func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Cancel(cancelDescription string) error {
	t.Status = TransactionError
	t.CancelDescription = cancelDescription
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}