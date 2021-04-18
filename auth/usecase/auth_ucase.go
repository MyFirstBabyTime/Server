package usecase

import "github.com/MyFirstBabyTime/Server/domain"

// authUsecase is used for usecase layer which implement domain.AuthUsecase interface
type authUsecase struct {

}

// AuthUsecase return implementation of domain.AuthUsecase
func AuthUsecase() domain.AuthUsecase {
	return &authUsecase{}
}
