package usecase

import (

	"github.com/MyFirstBabyTime/Server/domain"
)

// authUsecase is used for usecase layer which implement domain.AuthUsecase interface
type authUsecase struct {

	// parentAuthRepository is repository interface about domain.ParentAuth model
	parentAuthRepository domain.ParentAuthRepository

	// parentPhoneCertifyRepository is repository interface about domain.ParentPhoneCertify model
	parentPhoneCertifyRepository domain.ParentPhoneCertifyRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler TxHandler
}

// AuthUsecase return implementation of domain.AuthUsecase
func AuthUsecase(
	par domain.ParentAuthRepository,
	ppr domain.ParentPhoneCertifyRepository,
	th TxHandler,
) domain.AuthUsecase {
	return &authUsecase{
		parentAuthRepository:         par,
		parentPhoneCertifyRepository: ppr,
		txHandler:                    th,
}
}

// TxHandler is used for handling transaction to begin & commit or rollback
type TxHandler interface {
	// BeginTx method start transaction (get option from ctx)
	BeginTx(ctx context.Context, opts interface{}) (tx tx.Context, err error)

	// Commit method commit transaction
	Commit(tx tx.Context) (err error)

	// Rollback method rollback transaction
	Rollback(tx tx.Context) (err error)
}
}
