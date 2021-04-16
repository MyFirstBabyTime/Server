package domain

import (
	"context"
	"github.com/gin-gonic/gin"
)

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	// SendPhoneNumberCertifyCode method send phone number certify code
	SendPhoneNumberCertifyCode(ctx gin.Context)

	// CertifyPhoneNumber method certify phone number with certify code
	CertifyPhoneNumber(ctx gin.Context)

	// SignUpParent method create new parent auth with parent phone number
	SignUpParent(ctx gin.Context)
}

// AuthRepository is abstract interface about repository layer using in usecase layer
//type AuthRepositoryTXHandler interface {
	//parentAuthRepository

	//BeginTx(context.Context, *sql.TxOptions) (*sqlx.Tx, error) // BeginTx method start transaction
	//Commit(tx *sqlx.Tx) (err error)                            // Commit method commit transaction
	//Rollback(tx *sqlx.Tx) (err error)                          // Rollback method rollback transaction
//}

// ParentAuthRepository is interface only about ParentAuth model
type ParentAuthRepository interface {
	Store(ctx context.Context, pa *ParentAuth) error
}

type ParentPhoneNumberRepository interface {
	GetByPhoneNumber(ctx context.Context, pn string) (ParentPhoneCertify, error)
	Store(ctx context.Context, ppc *ParentPhoneCertify) error
	Update(ctx context.Context, ppc *ParentPhoneCertify)
}
