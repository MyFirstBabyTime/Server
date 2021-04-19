package domain

import (
	"context"
	"github.com/MyFirstBabyTime/Server/tx"
)

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	// SendCertifyCodeToPhone method send certify code to phone with pn(phone number)
	SendCertifyCodeToPhone(ctx context.Context, pn string) error

	// CertifyPhoneWithCode method certify phone with certify code
	CertifyPhoneWithCode(ctx context.Context, pn string, code int64) error

	// SignUpParent method create new parent auth with ParentAuth & parent phone number
	SignUpParent(ctx context.Context, pa *ParentAuth, pn string) error

	// LoginParentAuth method login parent auth & return logged ParentAuth model, token
	LoginParentAuth(ctx context.Context, id, pw string) (uuid, token string, err error)
}

// ParentAuthRepository is repository interface about ParentAuth model
type ParentAuthRepository interface {
	GetByUUID(ctx tx.Context, uuid string) (struct{
		ParentAuth
		ParentPhoneCertify
	}, error)
	GetByID(ctx tx.Context, id string) (struct{
		ParentAuth
		ParentPhoneCertify
	}, error)
	GetAvailableUUID(ctx tx.Context) (uuid string, err error)
	Store(ctx tx.Context, pa *ParentAuth) error
}

// ParentPhoneCertifyRepository is repository interface about ParentPhoneCertify model
type ParentPhoneCertifyRepository interface {
	GetByPhoneNumber(ctx tx.Context, pn string) (ParentPhoneCertify, error)
	Store(ctx tx.Context, ppc *ParentPhoneCertify) error
	Update(ctx tx.Context, ppc *ParentPhoneCertify) error
}
