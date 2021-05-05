package usecase

import (
	"context"
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/pkg/errors"
	"net/http"
)

type expenditureUsecase struct {
	// expenditureRepository is repository interface about domain.ExpenditureRepository
	expenditureRepository domain.ExpenditureRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler txHandler
}

func ExpenditureUsecase(
	er domain.ExpenditureRepository,
	th txHandler,
) *expenditureUsecase {
	return &expenditureUsecase{
		expenditureRepository: er,

		txHandler:  th,
	}
}
