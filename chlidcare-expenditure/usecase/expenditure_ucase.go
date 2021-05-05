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

// txHandler is used for handling transaction to begin & commit or rollback
type txHandler interface {
	// BeginTx method start transaction (get option from ctx)
	BeginTx(ctx context.Context, opts interface{}) (tx tx.Context, err error)

	// Commit method commit transaction
	Commit(tx tx.Context) (err error)

	// Rollback method rollback transaction
	Rollback(tx tx.Context) (err error)
}
