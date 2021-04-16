package domain

import (
	"context"
	"github.com/MyFirstBabyTime/Server/tx"
)

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	// SendPhoneNumberCertifyCode method send phone number certify code
	SendPhoneNumberCertifyCode(ctx context.Context)

	// CertifyPhoneNumber method certify phone number with certify code
	CertifyPhoneNumber(ctx context.Context)

	// SignUpParent method create new parent auth with parent phone number
	SignUpParent(ctx context.Context)
}

// AuthRepository is abstract interface about repository layer using in usecase layer
//type AuthRepositoryTXHandler interface {
	//parentAuthRepository

	//BeginTx(context.Context, *sql.TxOptions) (*sqlx.Tx, error) // BeginTx method start transaction
	//Commit(tx *sqlx.Tx) (err error)                            // Commit method commit transaction
	//Rollback(tx *sqlx.Tx) (err error)                          // Rollback method rollback transaction
//}

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
	Store(ctx tx.Context, pa *ParentAuth) error
}

// ParentPhoneNumberRepository is repository interface about ParentPhoneCertify model
type ParentPhoneNumberRepository interface {
	GetByPhoneNumber(ctx tx.Context, pn string) (ParentPhoneCertify, error)
	Store(ctx tx.Context, ppc *ParentPhoneCertify) error
	Update(ctx tx.Context, ppc *ParentPhoneCertify)
}
