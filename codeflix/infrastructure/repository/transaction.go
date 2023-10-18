package repository

import (
	"fmt"

	"derso.com/imersao-fullcycle/codepix-go/domain/model"
	"gorm.io/gorm"
)

type TransactionRepositoryDb struct {
	DB *gorm.DB
}

func (r TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	return r.DB.Create(transaction).Error
}

func (r TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	return r.DB.Save(transaction).Error
}

func (r TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	r.DB.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("no key found")
	}

	return &transaction, nil
}