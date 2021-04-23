package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
)

// authUsecase is used for usecase layer which implement domain.AuthUsecase interface
type authUsecase struct {
	// myCfg is used for get config value for auth usecase
	myCfg authUsecaseConfig

	// parentAuthRepository is repository interface about domain.ParentAuth model
	parentAuthRepository domain.ParentAuthRepository

	// parentPhoneCertifyRepository is repository interface about domain.ParentPhoneCertify model
	parentPhoneCertifyRepository domain.ParentPhoneCertifyRepository

	// txHandler is used for handling transaction to begin & commit or rollback
	txHandler txHandler

	// messageAgency is used as agency about message API
	messageAgency messageAgency

	// messageAgency is used as handler about hashing
	hashHandler hashHandler

	// jwtHandler is used as handler about jwt
	jwtHandler jwtHandler
}

// AuthUsecase return implementation of domain.AuthUsecase
func AuthUsecase(
	cfg authUsecaseConfig,
	par domain.ParentAuthRepository,
	ppr domain.ParentPhoneCertifyRepository,
	th txHandler,
	ma messageAgency,
	hh hashHandler,
	jh jwtHandler,
) domain.AuthUsecase {
	return &authUsecase{
		myCfg: cfg,

		parentAuthRepository:         par,
		parentPhoneCertifyRepository: ppr,

		txHandler:     th,
		messageAgency: ma,
		hashHandler:   hh,
		jwtHandler:    jh,
	}
}

// authUsecaseConfig is interface get config value for auth usecase
type authUsecaseConfig interface {
	// AccessTokenDuration return access token valid duration
	AccessTokenDuration() time.Duration
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

// messageAgency is agency that agent various API about message
type messageAgency interface {
	// SendSMSToOne method send SMS message to one receiver
	SendSMSToOne(receiver, content string) (err error)
}

// hashHandler is interface about hash handler
type hashHandler interface {
	// GenerateHashWithMinSalt generate & return hashed value from password with minimum salt
	GenerateHashWithMinSalt(pw string) (hash string, err error)

	// CompareHashAndPW compare hashed value and password & return error
	CompareHashAndPW(hash, pw string) (err error)
}

// jwtHandler is interface about JWT handler
type jwtHandler interface {
	// GenerateUUIDJWT generate & return JWT UUID token with type & time
	GenerateUUIDJWT(uuid, _type string, t time.Duration) (token string, err error)
}

// SendCertifyCodeToPhone is implement domain.AuthUsecase interface
func (au *authUsecase) SendCertifyCodeToPhone(ctx context.Context, pn string) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	switch err.(type) {
	case nil:
		if ppc.ParentUUID.Valid {
			err = errors.New("this phone number is already in use")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyInUse}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		ppc.CertifyCode = ppc.GenerateCertifyCode()
		ppc.Certified = sql.NullBool{Bool: false, Valid: true}
		switch err = au.parentPhoneCertifyRepository.Update(_tx, &ppc); err.(type) {
		case nil:
			break
		default:
			err = errors.Wrap(err, "phone Update return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	case domain.ErrRowNotExist:
		ppc = domain.ParentPhoneCertify{PhoneNumber: pn}
		ppc.CertifyCode = ppc.GenerateCertifyCode()
		switch err = au.parentPhoneCertifyRepository.Store(_tx, &ppc); err.(type) {
		case nil:
			break
		default:
			err = errors.Wrap(err, "phone Store return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	default:
		err = errors.Wrap(err, "GetByPhoneNumber return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	content := fmt.Sprintf("[육아는 처음이지 인증 번호]\n회원가입 인증 번호: %d", ppc.CertifyCode)
	if err = au.messageAgency.SendSMSToOne(ppc.PhoneNumber, content); err != nil {
		err = errors.Wrap(err, "SendSMSToOne return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	_ = au.txHandler.Commit(_tx)
	return nil
}

// CertifyPhoneWithCode is implement domain.AuthUsecase interface
func (au *authUsecase) CertifyPhoneWithCode(ctx context.Context, pn string, code int64) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	switch err.(type) {
	case nil:
		if ppc.Certified.Valid && ppc.Certified.Bool {
			err = errors.New("this phone number is already certified")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyCertified}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		if code != ppc.CertifyCode {
			err = errors.New("incorrect certify code to that phone number")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.IncorrectCertifyCode}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		ppc.Certified = sql.NullBool{Bool: true, Valid: true}
		switch err = au.parentPhoneCertifyRepository.Update(_tx, &ppc); err.(type) {
		case nil:
			break
		default:
			err = errors.Wrap(err, "phone Update return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	case domain.ErrRowNotExist:
		err = errors.New("not exist phone number")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusNotFound}
		_ = au.txHandler.Rollback(_tx)
		return
	default:
		err = errors.Wrap(err, "GetByPhoneNumber return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	_ = au.txHandler.Commit(_tx)
	return nil
}

// SignUpParent is implement domain.AuthUsecase interface
func (au *authUsecase) SignUpParent(ctx context.Context, pa *domain.ParentAuth, pn string) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	if err == nil && ppc.Certified.Valid && ppc.Certified.Bool {
		if ppc.ParentUUID.Valid {
			err = errors.New("this phone number is already in use")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyInUse}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		if pa.PW, err = au.hashHandler.GenerateHashWithMinSalt(pa.PW); err != nil {
			err = errors.Wrap(err, "failed to GenerateHashWithMinSalt")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		if pa.UUID, err = au.parentAuthRepository.GetAvailableUUID(_tx); err != nil {
			pa.UUID = pa.GenerateRandomUUID()
		}
		switch err = au.parentAuthRepository.Store(_tx, pa); tErr := err.(type) {
		case nil:
			break
		case domain.ErrInvalidModel:
			err = errors.Wrap(err, "parent auth Store return invalid model")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		case domain.ErrEntryDuplicate:
			switch tErr.DuplicateKey {
			case "id":
				err = errors.New("this parent ID is already in use")
				err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.ParentIDAlreadyInUse}
				_ = au.txHandler.Rollback(_tx)
				return
			default:
				err = errors.Wrap(err, "parent auth Store return unexpected duplicate error")
				err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
				_ = au.txHandler.Rollback(_tx)
				return
			}
		default:
			err = errors.Wrap(err, "parent auth Store return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	} else {
		if _, ok := err.(domain.ErrRowNotExist); err == nil || ok {
			err = errors.New("this phone number is not certified")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.UncertifiedPhone}
			_ = au.txHandler.Rollback(_tx)
			return
		} else {
			err = errors.Wrap(err, "GetByPhoneNumber return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	}

	ppc.ParentUUID = sql.NullString{String: pa.UUID, Valid: true}
	if err = au.parentPhoneCertifyRepository.Update(_tx, &ppc); err != nil {
		err = errors.Wrap(err, "phone Update return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	_ = au.txHandler.Commit(_tx)
	return nil
}

// LoginParentAuth is implement domain.AuthUsecase interface
func (au *authUsecase) LoginParentAuth(ctx context.Context, id, pw string) (uuid, token string, err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	pa, err := au.parentAuthRepository.GetByID(_tx, id)
	switch err.(type) {
	case nil:
		switch err = au.hashHandler.CompareHashAndPW(pa.PW, pw); err.(type) {
		case nil:
			break
		case interface{ Mismatch() }:
			err = errors.New("incorrect password")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.IncorrectParentPW}
			_ = au.txHandler.Rollback(_tx)
			return
		default:
			err = errors.Wrap(err, "CompareHashAndPW return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	case domain.ErrRowNotExist:
		err = errors.New("not exist parent ID")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.NotExistParentID}
		_ = au.txHandler.Rollback(_tx)
		return
	default:
		err = errors.Wrap(err, "GetByID return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	uuid = pa.UUID
	token, err = au.jwtHandler.GenerateUUIDJWT(pa.UUID, "access_token", au.myCfg.AccessTokenDuration())
	err = nil

	_ = au.txHandler.Commit(_tx)
	return
}
