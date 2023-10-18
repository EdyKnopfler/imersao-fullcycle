package model

import (
	"time"
	uuid "github.com/satori/go.uuid"
	govalidator "github.com/asaskevich/govalidator"
)

type Account struct {
	Base `valid:"required"`
	OwnerName string `json:"ownerName" valid:"notnull"`
	Bank *Bank `valid:"-"`
	Number string `json:"number" valid:"-"`
	PixKeys []*PixKey `valid:"-"`
}

func NewAccount(bank *Bank, number string, ownerName string) (*Account, error) {
	account := Account{
		OwnerName: ownerName,
		Bank: bank,
		Number: number,
	}

	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	err := account.isValid()
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (account *Account) isValid() error {
	_, err := govalidator.ValidateStruct(account)
	return err
}