package model

import (
	"time"
	"errors"
	uuid "github.com/satori/go.uuid"
	govalidator "github.com/asaskevich/govalidator"
)

type PixKey struct {
	Base `valid:"required"`
	Kind string `json:"kind" valid:"notnull"`
	Key string `json:"key" valid:"notnull"`
	AccountID string `json:"accountId" valid:"notnull"`
	Account *Account `valid:"-"`
	Status string `json:"status" valid:"notnull"`
}

type PixKeyRepositoryInterface interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{
		Kind: kind,
		Account: account,
		AccountId: account.ID,
		Key: key,
		Status: "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.isValid()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}

func (pixKey *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if err != nil {
		return err
	}

	// Uma validação rudimentar
	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid key type")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}

	return err
}