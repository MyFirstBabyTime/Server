package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
)

// authUsecase is used for usecase layer which implement domain.AuthUsecase interface
type authUsecase struct {

	// parentAuthRepository is repository interface about domain.ParentAuth model
	parentAuthRepository domain.ParentAuthRepository

	// parentPhoneCertifyRepository is repository interface about domain.ParentPhoneCertify model
	parentPhoneCertifyRepository domain.ParentPhoneCertifyRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler TxHandler

	// messageAgency is used as agency about message command
	messageAgency messageAgency
}

// AuthUsecase return implementation of domain.AuthUsecase
func AuthUsecase(
	par domain.ParentAuthRepository,
	ppr domain.ParentPhoneCertifyRepository,
	th TxHandler,
	ma messageAgency,
) domain.AuthUsecase {
	return &authUsecase{
		parentAuthRepository:         par,
		parentPhoneCertifyRepository: ppr,

		txHandler:     th,
		messageAgency: ma,
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

// messageAgency is agency that agent various command about message
type messageAgency interface {
	// SendSMSToOne method send SMS message to one receiver
	SendSMSToOne(receiver, content string) (err error)
}

// SendCertifyCodeToPhone is implement domain.AuthUsecase interface
func (au *authUsecase) SendCertifyCodeToPhone(ctx context.Context, pn string) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	switch err.(type) {
	case nil:
		if ppc.ParentUUID.Valid {
			err = conflictErr{errors.New("this phone number is already in use"), -101}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		ppc.CertifyCode = ppc.GenerateCertifyCode()
		ppc.Certified = sql.NullBool{Bool: false, Valid: true}
		switch err = au.parentPhoneCertifyRepository.Update(_tx, &ppc); err.(type) {
		case nil:
			break
		default:
			err = internalServerErr{errors.Wrap(err, "phone Update return unexpected error")}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	case rowNotExistErr:
		ppc = domain.ParentPhoneCertify{PhoneNumber: pn}
		ppc.CertifyCode = ppc.GenerateCertifyCode()
		switch err = au.parentPhoneCertifyRepository.Store(_tx, &ppc); err.(type) {
		case nil:
			break
		default:
			err = internalServerErr{errors.Wrap(err, "phone Store return unexpected error")}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	default:
		err = internalServerErr{errors.Wrap(err, "GetByPhoneNumber return unexpected error")}
		_ = au.txHandler.Rollback(_tx)
		return
	}
}
