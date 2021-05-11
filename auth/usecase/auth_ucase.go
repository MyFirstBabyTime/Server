package usecase

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

	// s3Agency is used as agency about aws s3 API
	s3Agency s3Agency
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
	sa s3Agency,
) domain.AuthUsecase {
	return &authUsecase{
		myCfg: cfg,

		parentAuthRepository:         par,
		parentPhoneCertifyRepository: ppr,

		txHandler:     th,
		messageAgency: ma,
		hashHandler:   hh,
		jwtHandler:    jh,
		s3Agency:      sa,
	}
}

// authUsecaseConfig is interface get config value for auth usecase
type authUsecaseConfig interface {
	// AccessTokenDuration return access token valid duration
	AccessTokenDuration() time.Duration

	// ParentProfileS3Bucket return aws s3 bucket name for parent profile
	ParentProfileS3Bucket() string
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

// s3Agency is agency that agent various API about aws s3
type s3Agency interface {
	// PutObject method put(insert or update) object to s3
	PutObject(input *s3.PutObjectInput) (output *s3.PutObjectOutput, err error)
}

// SendCertifyCodeToPhone implement SendCertifyCodeToPhone method of domain.AuthUsecase interface
func (au *authUsecase) SendCertifyCodeToPhone(ctx context.Context, pn string) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	switch err.(type) {
	case nil:
		if domain.StringValue(ppc.ParentUUID) != "" {
			err = errors.New("this phone number is already in use")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyInUse}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		ppc.CertifyCode = domain.Int64(ppc.GenerateCertifyCode())
		ppc.Certified = domain.Bool(false)
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
		ppc = domain.ParentPhoneCertify{
			PhoneNumber: domain.String(pn),
			CertifyCode: domain.Int64(ppc.GenerateCertifyCode()),
		}
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

	content := fmt.Sprintf("[육아는 처음이지 인증 번호]\n회원가입 인증 번호: %d", domain.Int64Value(ppc.CertifyCode))
	if err = au.messageAgency.SendSMSToOne(domain.StringValue(ppc.PhoneNumber), content); err != nil {
		err = errors.Wrap(err, "SendSMSToOne return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	_ = au.txHandler.Commit(_tx)
	return nil
}

// CertifyPhoneWithCode implement CertifyPhoneWithCode method of domain.AuthUsecase interface
func (au *authUsecase) CertifyPhoneWithCode(ctx context.Context, pn string, code int64) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, pn)
	switch err.(type) {
	case nil:
		if domain.BoolValue(ppc.Certified) == true {
			err = errors.New("this phone number is already certified")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyCertified}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		if code != domain.Int64Value(ppc.CertifyCode) {
			err = errors.New("incorrect certify code to that phone number")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.IncorrectCertifyCode}
			_ = au.txHandler.Rollback(_tx)
			return
		}
		ppc.Certified = domain.Bool(true)
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

// SignUpParent implement SignUpParent method of domain.AuthUsecase interface
func (au *authUsecase) SignUpParent(ctx context.Context, pi struct {
	*domain.ParentAuth
	*domain.ParentPhoneCertify
}, profile []byte) (uuid string, err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	ppc, err := au.parentPhoneCertifyRepository.GetByPhoneNumber(_tx, domain.StringValue(pi.PhoneNumber))
	if err == nil && domain.BoolValue(ppc.Certified) == true {
		if domain.StringValue(ppc.ParentUUID) != "" {
			err = errors.New("this phone number is already in use")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusConflict, Code: domain.PhoneAlreadyInUse}
			_ = au.txHandler.Rollback(_tx)
			return
		}

		if hash, err := au.hashHandler.GenerateHashWithMinSalt(domain.StringValue(pi.PW)); err != nil {
			err = errors.Wrap(err, "failed to GenerateHashWithMinSalt")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return "", err
		} else {
			pi.PW = domain.String(hash)
		}

		if uuid, err = au.parentAuthRepository.GetAvailableUUID(_tx); err != nil {
			pi.UUID = domain.String(pi.GenerateRandomUUID())
		} else {
			pi.UUID = domain.String(uuid)
		}
		if profile != nil && string(profile) != "" {
			pi.ProfileUri = domain.String(pi.ParentAuth.GenerateProfileUri())
		}

		switch err = au.parentAuthRepository.Store(_tx, pi.ParentAuth); tErr := err.(type) {
		case nil:
			break
		case domain.ErrInvalidModel:
			err = errors.Wrap(err, "parent auth Store return invalid model")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		case domain.ErrEntryDuplicate:
			switch tErr.DuplicateKey {
			case "id", "parent_auth.id":
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

	ppc.ParentUUID = domain.String(domain.StringValue(pi.UUID))
	if err = au.parentPhoneCertifyRepository.Update(_tx, &ppc); err != nil {
		err = errors.Wrap(err, "phone Update return unexpected error")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	if profile != nil && string(profile) != "" {
		if _, err = au.s3Agency.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(au.myCfg.ParentProfileS3Bucket()),
			Key:    aws.String(pi.ParentAuth.GenerateProfileUri()),
			Body:   bytes.NewReader(profile),
			ACL:    aws.String("public-read"),
		}); err != nil {
			err = errors.Wrap(err, "s3 PutObject return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	}

	uuid = domain.StringValue(pi.UUID)
	err = nil
	_ = au.txHandler.Commit(_tx)
	return
}

// LoginParentAuth implement LoginParentAuth method of domain.AuthUsecase interface
func (au *authUsecase) LoginParentAuth(ctx context.Context, id, pw string) (uuid, token string, err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	pa, err := au.parentAuthRepository.GetByID(_tx, id)
	switch err.(type) {
	case nil:
		switch err = au.hashHandler.CompareHashAndPW(domain.StringValue(pa.PW), pw); err.(type) {
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

	uuid = domain.StringValue(pa.UUID)
	token, err = au.jwtHandler.GenerateUUIDJWT(uuid, "access_token", au.myCfg.AccessTokenDuration())
	err = nil

	_ = au.txHandler.Commit(_tx)
	return
}

// GetParentInformByID implement GetParentInformByID method of domain.AuthUsecase interface
func (au *authUsecase) GetParentInformByID(ctx context.Context, id string) (pi struct {
	domain.ParentAuth
	domain.ParentPhoneCertify
}, err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}

	pi, err = au.parentAuthRepository.GetByID(_tx, id)
	switch err.(type) {
	case nil:
	case domain.ErrRowNotExist:
		err = domain.UsecaseError{UsecaseErr: errors.New("not exist parent auth with that ID"), Status: http.StatusNotFound}
		_ = au.txHandler.Rollback(_tx)
	}

	_ = au.txHandler.Commit(_tx)
	return pi, err
}

// UpdateParentInform implement UpdateParentInform method of domain.AuthUsecase interface
func (au *authUsecase) UpdateParentInform(ctx context.Context, uuid string, pa *domain.ParentAuth, profile []byte) (err error) {
	_tx, err := au.txHandler.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "failed to begin transaction")
		return
	}
	pa.UUID = domain.String(uuid)

	if profile != nil && len(profile) != 0 {
		pa.ProfileUri = domain.String(pa.GenerateProfileUri())
	}

	if err = au.parentAuthRepository.Update(_tx, pa); err != nil {
		err = errors.Wrap(err, "failed to Update")
		err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
		_ = au.txHandler.Rollback(_tx)
		return
	}

	if profile != nil && len(profile) != 0 {
		if _, err = au.s3Agency.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(au.myCfg.ParentProfileS3Bucket()),
			Key:    aws.String(domain.StringValue(pa.ProfileUri)),
			Body:   bytes.NewReader(profile),
			ACL:    aws.String("public-read"),
		}); err != nil {
			err = errors.Wrap(err, "s3 PutObject return unexpected error")
			err = domain.UsecaseError{UsecaseErr: err, Status: http.StatusInternalServerError}
			_ = au.txHandler.Rollback(_tx)
			return
		}
	}

	_ = au.txHandler.Commit(_tx)
	return nil
}
