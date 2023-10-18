package repository

import (
	"fmt"

	"derso.com/imersao-fullcycle/codepix-go/domain/model"
	"gorm.io/gorm"
)

type PixKeyRepositoryDb struct {
	DB *gorm.DB
}

func (r PixKeyRepositoryDb) AddBank(bank *model.Bank) error {
	return r.DB.Create(bank).Error
}

func (r PixKeyRepositoryDb) AddAccount(account *model.Account) error {
	return r.DB.Create(account).Error
}

func (r PixKeyRepositoryDb) RegisterKey(pixKey *model.PixKey) (*model.PixKey, error) {
	err := r.DB.Create(pixKey).Error

	if err != nil {
		return nil, err
	}

	return pixKey, nil
}

func (r PixKeyRepositoryDb) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey
	r.DB.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)

	if pixKey.ID == "" {
		return nil, fmt.Errorf("no key found")
	}

	return &pixKey, nil
}

func (r PixKeyRepositoryDb) FindAccount(id string) (*model.Account, error) {
	var account model.Account
	r.DB.Preload("Bank").First(&account, "id = ?", id)

	if account.ID == "" {
		return nil, fmt.Errorf("no account found")
	}

	return &account, nil
}

func (r PixKeyRepositoryDb) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank
	r.DB.First(&bank, "id = ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("no bank found")
	}

	return &bank, nil
}