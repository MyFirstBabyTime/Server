package domain

import "database/sql"

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID       string         `db:"uuid" validate:"required"`
	ID         string         `db:"id" validate:"required"`
	PW         string         `db:"pw" validate:"required"`
	Name       string         `db:"name" validate:"required"`
	ProfileUri sql.NullString `db:"profile_uri"`
}

// TableName return table name about model
func (pa ParentAuth) TableName() string {
	return "parent_auth"
}

// Schema return schema SQL about model
func (pa ParentAuth) Schema() string {
	return `CREATE TABLE parent_auth (
		uuid        CHAR(11)  NOT NULL,
		id          VARCHAR(20)  NOT NULL,
		pw          VARCHAR(100) NOT NULL,
		name        VARCHAR(10)  NOT NULL,
		profile_uri VARCHAR(100),
		PRIMARY KEY (uuid)
	);`
}

// ParentPhoneNumber is model represent parent phone number using in auth domain
type ParentPhoneNumber struct {
	ParentUUID  string `db:"parent_uuid" validate:"required"`
	PhoneNumber string `db:"phone_number" validate:"required"`
	CertifyCode int    `db:"certify_code" validate:"required"`
	Certified   bool   `db:"certified"`
}
