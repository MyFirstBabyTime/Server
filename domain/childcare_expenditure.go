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
