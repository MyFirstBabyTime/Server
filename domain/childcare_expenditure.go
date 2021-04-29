package domain

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type ExpenditureUsecase interface {
}

type ExpenditureRepository interface {
}

// Expenditure is model represent expenditure using in child_expenditure domain
type Expenditure struct {
	UUID       string         `db:"uuid" validate:"uuid=item"`
	ParentUUID string         `db:"parent_uuid" validate:"required uuid=parent"`
	BabyUUID   string         `db:"baby_uuid" validate:"required"`
	Name       string         `db:"name" validate:"required"`
	Amount     int64          `db:"amount" validate:"required"`
	Rating     int64          `db:"rating" validate:"range=0~5"`
	Link       sql.NullString `db:"link"`
}

// TableName return table name about Expenditure model
func (e Expenditure) TableName() string {
	return "expenditure"
}

func (e Expenditure) Schema() string {
	return `CREATE TABLE expenditure (
		uuid 		CHAR(11) 	NOT NULL,
		parent_uuid CHAR(11)	NOT NULL,
		baby_uuid	CHAR(11)    NOT NULL,
		name 		VARCHAR(20) NOT NULL,
		amount 		INT(15) 	NOT NULL,
		rating 		INT(1) 		NOT NULL,
		link 		VARCHAR(100),
		PRIMARY KEY (uuid)
		FOREIGN KEY (parent_uuid)
			REFERENCES parent_auth(uuid)
			ON DELETE CASCADE
	);`
}
