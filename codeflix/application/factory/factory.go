package factory

import (
	"derso.com/imersao-fullcycle/codepix-go/application/usecase"
	"derso.com/imersao-fullcycle/codepix-go/infrastructure/repository"
	"gorm.io/gorm"
)

func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	transactionRepository := repository.TransactionRepositoryDb{DB: database}
	pixRepository := repository.PixKeyRepositoryDb{DB: database}

	return usecase.TransactionUseCase{
		TransactionRepository: transactionRepository,
		PixRepository:         pixRepository,
	}
}
