package model

import (
	"time"
	uuid "github.com/satori/go.uuid"
	govalidator "github.com/asaskevich/govalidator"
)

type Bank struct {
	Base `valid:"required"`
	Code string `json:"code" valid:"notnull"`
	Name string `json:"name" valid:"notnull"`
	Accounts []*Account `valid:"-"`
}

func NewBank(code string, name string) (*Bank, error) {
	bank := Bank{
		Code: code,
		Name: name,
	}

	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()

	err := bank.isValid()
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (bank *Bank) isValid() error {
	_, err := govalidator.ValidateStruct(bank)
	return err
}