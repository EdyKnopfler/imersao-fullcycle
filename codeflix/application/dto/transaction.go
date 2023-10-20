package dto

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
)

type TransactionDTO struct {
	ID           string  `json:"id" valid:"notnull,uuid"`
	AccountID    string  `json:"accountId" valid:"notnull,uuid"`
	Amount       float64 `json:"amount" valid:"notnull,numeric"`
	PixKeyTo     string  `json:"pixKeyTo" valid:"notnull"`
	PixKeyKindTo string  `json:"pixKeyKindTo" valid:"notnull"`
	Description  string  `json:"description" valid:"notnull"`
	Status       string  `json:"status" valid:"notnull"`
	Error        string  `json:"error"`
}

func NewTransactionDTO(data []byte) (*TransactionDTO, error) {
	transactionDto := &TransactionDTO{}

	// TODO a lógica de marshal/unmarshal + validação poderia ser refatorada
	err := json.Unmarshal(data, transactionDto)
	if err != nil {
		return nil, err
	}

	err = transactionDto.isValid()
	if err != nil {
		return nil, err
	}

	return transactionDto, nil
}

func (t *TransactionDTO) isValid() error {
	_, err := govalidator.ValidateStruct(t)
	return err
}

func (t *TransactionDTO) ToJson() ([]byte, error) {
	err := t.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return result, nil
}
