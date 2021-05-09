package domain

import (
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/MyFirstBabyTime/Server/tx"
)

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	// SendCertifyCodeToPhone method send certify code to phone with pn(phone number)
	SendCertifyCodeToPhone(ctx context.Context, pn string) error

	// CertifyPhoneWithCode method certify phone with certify code
	CertifyPhoneWithCode(ctx context.Context, pn string, code int64) error

	// SignUpParent method create new parent auth with ParentAuth, ParentPhoneCertify model & profile multipart
	SignUpParent(ctx context.Context, pi struct {
		*ParentAuth
		*ParentPhoneCertify
	}, profile *multipart.FileHeader) (uuid string, err error)

	// LoginParentAuth method login parent auth & return logged ParentAuth model, token
	LoginParentAuth(ctx context.Context, id, pw string) (uuid, token string, err error)

	// GetParentInformByID method get ParentAuth & ParentPhoneCertify model inform by parent ID
	GetParentInformByID(ctx context.Context, id string) (struct {
		ParentAuth
		ParentPhoneCertify
	}, error)

	// UpdateParentInform method update ParentAuth model inform & profile image with parent uuid
	UpdateParentInform(ctx context.Context, uuid string, pa *ParentAuth, profile *multipart.FileHeader) (err error)
}

// ParentAuthRepository is repository interface about ParentAuth model
type ParentAuthRepository interface {
	GetByUUID(ctx tx.Context, uuid string) (struct {
		ParentAuth
		ParentPhoneCertify
	}, error)
	GetByID(ctx tx.Context, id string) (struct {
		ParentAuth
		ParentPhoneCertify
	}, error)
	GetAvailableUUID(ctx tx.Context) (uuid string, err error)
	Store(ctx tx.Context, pa *ParentAuth) error
	Update(ctx tx.Context, pa *ParentAuth) error
}

// ParentPhoneCertifyRepository is repository interface about ParentPhoneCertify model
type ParentPhoneCertifyRepository interface {
	GetByPhoneNumber(ctx tx.Context, pn string) (ParentPhoneCertify, error)
	Store(ctx tx.Context, ppc *ParentPhoneCertify) error
	Update(ctx tx.Context, ppc *ParentPhoneCertify) error
}

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID       *string `db:"uuid" validate:"not_empty,uuid=parent"`
	ID         *string `db:"id" validate:"not_empty,min=4,max=20"`
	PW         *string `db:"pw" validate:"not_empty"`
	Name       *string `db:"name" validate:"not_empty,max=20"`
	ProfileUri *string `db:"profile_uri"`
}

// TableName return table name about ParentAuth model
func (pa ParentAuth) TableName() string {
	return "parent_auth"
}

// Schema return schema SQL about ParentAuth model
func (pa ParentAuth) Schema() string {
	return `CREATE TABLE parent_auth (
		uuid        CHAR(11)  NOT NULL,
		id          VARCHAR(20)  NOT NULL UNIQUE,
		pw          VARCHAR(100) NOT NULL,
		name        VARCHAR(10)  NOT NULL,
		profile_uri VARCHAR(100),
		PRIMARY KEY (uuid)
	);`
}

// GenerateRandomUUID method return random UUID value
func (pa ParentAuth) GenerateRandomUUID() string {
	rand.Seed(time.Now().UnixNano())
	is := []rune("0123456789")
	random := make([]rune, 10)
	for i := range random {
		random[i] = is[rand.Intn(len(is))]
	}
	return fmt.Sprintf("p%s", string(random))
}

// GenerateProfileUri method return ProfileUri value with field value
func (pa ParentAuth) GenerateProfileUri() string {
	return fmt.Sprintf("/profiles/parents/uuid/%s", StringValue(pa.UUID))
}

// GenerateValidModel method return model referenced by value with set valid value
func (pa ParentAuth) GenerateValidModel() ParentAuth {
	var (
		validUUID = String(pa.GenerateRandomUUID())
		validID   = String("validID")
		validPW   = String("validPWHashedString")
		validName = String("validName")
	)

	if pa.UUID == nil {
		pa.UUID = validUUID
	}
	if pa.ID == nil {
		pa.ID = validID
	}
	if pa.PW == nil {
		pa.PW = validPW
	}
	if pa.Name == nil {
		pa.Name = validName
	}

	return pa
}

// ParentPhoneCertify is model represent parent phone number using in auth domain
type ParentPhoneCertify struct {
	ParentUUID  *string `db:"parent_uuid" validate:"uuid=parent"`
	PhoneNumber *string `db:"phone_number" validate:"not_empty,len=11"`
	CertifyCode *int64  `db:"certify_code" validate:"not_empty,range=100000~999999"`
	Certified   *bool   `db:"certified"`
}

// TableName return table name about ParentPhoneNumber model
func (pn ParentPhoneCertify) TableName() string {
	return "parent_phone_certify"
}

// Schema return schema SQL about ParentPhoneNumber model
func (pn ParentPhoneCertify) Schema() string {
	return `CREATE TABLE parent_phone_certify (
		parent_uuid  CHAR(11) UNIQUE,
		phone_number CHAR(11) NOT NULL,
		certify_code INT(11)  NOT NULL,
		certified    TINYINT  NOT NULL DEFAULT 0,
		PRIMARY KEY (phone_number),
		FOREIGN KEY (parent_uuid)
        	REFERENCES parent_auth(uuid)
        	ON DELETE CASCADE
	);`
}

// GenerateCertifyCode method return CertifyCode value
func (pn *ParentPhoneCertify) GenerateCertifyCode() int64 {
	rand.Seed(time.Now().UnixNano())
	is := []rune("0123456789")
	random := make([]rune, 6)
	for i := range random {
		random[i] = is[rand.Intn(len(is))]
	}
	v, _ := strconv.Atoi(string(random))
	return int64(v)
}

// GenerateValidModel method return model referenced by value with set valid value
func (pn ParentPhoneCertify) GenerateValidModel() ParentPhoneCertify {
	var (
		validCertifyCode = Int64(123456)
	)

	if pn.CertifyCode == nil {
		pn.CertifyCode = validCertifyCode
	}

	return pn
}
