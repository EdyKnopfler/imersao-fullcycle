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
	AccountFromID string `gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount float64 `gorm:"type:decimal(10,2)" json:"amount" valid:"notnull"`
	PixKeyTo *PixKey `valid:"-"`
	PixKeyIdTo string `gorm:"column:pix_key_id_to;type:uuid;" valid:"notnull"`
	Status string `gorm:"type:varchar(20)" json:"status" valid:"notnull"`
	Description string `gorm:"type:varchar(255)" json:"description" valid:"notnull"`
	CancelDescription string `gorm:"type:varchar(255)" json:"cancelDescription" valid:"-"`
}

func NewTransaction(
	accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	
	transaction := Transaction{
		AccountFrom: accountFrom,
		AccountFromID: accountFrom.ID,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		PixKeyIdTo: pixKeyTo.ID,
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

func (t *Transaction) Complete() error {
	return t.atualizarStatusEValidar(TransactionCompleted)
}

func (t *Transaction) Cancel(cancelDescription string) error {
	t.CancelDescription = cancelDescription
	return t.atualizarStatusEValidar(TransactionError)
}

func (t *Transaction) Confirm() error {
	return t.atualizarStatusEValidar(TransactionConfirmed)
}

func (t *Transaction) atualizarStatusEValidar(novoStatus string) error {
	t.Status = novoStatus
	t.UpdatedAt = time.Now()
	err := t.isValid()
	return err
}
