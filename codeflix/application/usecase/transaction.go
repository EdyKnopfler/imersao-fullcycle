package usecase

import (
	"errors"
	"log"

	"derso.com/imersao-fullcycle/codepix-go/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.TransactionRepositoryInterface
	PixRepository model.PixKeyRepositoryInterface
}

func (t *TransactionUseCase) Register(accountId string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Save(transaction)
	if transaction.ID == "" {
		return nil, errors.New("unable to process this transaction")
	}
	
	return transaction, nil
}

func (t *TransactionUseCase) Confirm(transactionId string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transaction not found", transactionId)
		return nil, err
	}

	err = transaction.Confirm()
	if err != nil {
		log.Println("Unable to confirm transaction")
		return nil, err
	}

	err = t.TransactionRepository.Save(transaction)
	if err != nil {
		log.Println("Unable to save transaction")
		return nil, err
	}

	return transaction, nil
}

// TODO fluxo para as outras regras de negócio (possivelmente refatoradinho)