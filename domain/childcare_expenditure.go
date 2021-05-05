package domain

import (
	"context"
	"fmt"
	"github.com/MyFirstBabyTime/Server/tx"
	"math/rand"
	"time"
)

type ExpenditureUsecase interface {
	ExpenditureRegistration(ctx context.Context, req *Expenditure, babyUUIDs *[]string) (err error)}

type ExpenditureRepository interface {
	Store(ctx tx.Context, e *Expenditure, babyUUIDs *[]string) (err error)
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

// GenerateRandomUUID method return random UUID value
func (e Expenditure) GenerateRandomUUID() string {
	rand.Seed(time.Now().UnixNano())
	is := []rune("0123456789")
	random := make([]rune, 10)
	for i := range random {
		random[i] = is[rand.Intn(len(is))]
	}
	return fmt.Sprintf("e%s", string(random))
}
