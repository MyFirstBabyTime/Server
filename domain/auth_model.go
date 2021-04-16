package domain

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID       string         `db:"uuid" validate:"required"`
	ID         string         `db:"id" validate:"required"`
	PW         string         `db:"pw" validate:"required"`
	Name       string         `db:"name" validate:"required"`
	ProfileUri sql.NullString `db:"profile_uri"`
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

// SetRandomUUID method set UUID field to random value
func (pa *ParentAuth) SetRandomUUID() {
	rand.Seed(time.Now().UnixNano())
	intLetters := []rune("0123456789")
	randomRuneArr := make([]rune, 10)
	for i := range randomRuneArr {
		randomRuneArr[i] = intLetters[rand.Intn(len(intLetters))]
	}
	pa.UUID = fmt.Sprintf("p%s", string(randomRuneArr))
}

// ParentPhoneCertify is model represent parent phone number using in auth domain
type ParentPhoneCertify struct {
	ParentUUID  sql.NullString `db:"parent_uuid"`
	PhoneNumber string         `db:"phone_number" validate:"required"`
	CertifyCode int            `db:"certify_code" validate:"required"`
	Certified   bool           `db:"certified"`
}

// TableName return table name about ParentPhoneNumber model
func (pn ParentPhoneCertify) TableName() string {
	return "parent_phone_certify"
}

// Schema return schema SQL about ParentPhoneNumber model
func (pn ParentPhoneCertify) Schema() string {
	return `CREATE TABLE parent_phone_certify (
		parent_uuid  CHAR(11),
		phone_number CHAR(11) NOT NULL,
		certify_code INT(11)  NOT NULL,
		certified    TINYINT  NOT NULL DEFAULT 0,
		PRIMARY KEY (phone_number),
		FOREIGN KEY (parent_uuid)
        	REFERENCES parent_auth(uuid)
        	ON DELETE CASCADE
	);`
}
