package domain

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID       string         `db:"uuid" validate:"required,uuid=parent"`
	ID         string         `db:"id" validate:"required,min=4,max=20"`
	PW         string         `db:"pw" validate:"required"`
	Name       string         `db:"name" validate:"required,max=20"`
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

// ParentPhoneCertify is model represent parent phone number using in auth domain
type ParentPhoneCertify struct {
	ParentUUID  sql.NullString `db:"parent_uuid" validate:"uuid=parent"`
	PhoneNumber string         `db:"phone_number" validate:"required,len=11"`
	CertifyCode int            `db:"certify_code" validate:"required,range=100000~999999"`
	Certified   bool           `db:"certified"`
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
func (pn *ParentPhoneCertify) GenerateCertifyCode() int {
	rand.Seed(time.Now().UnixNano())
	is := []rune("0123456789")
	random := make([]rune, 6)
	for i := range random {
		random[i] = is[rand.Intn(len(is))]
	}
	v, _ := strconv.Atoi(string(random))
	return v
}
