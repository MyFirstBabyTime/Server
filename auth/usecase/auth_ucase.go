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
func AuthUsecase() domain.AuthUsecase {
	return &authUsecase{}
}
