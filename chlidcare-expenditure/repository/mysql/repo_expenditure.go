package mysql

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

type expenditureRepository struct {
	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
	validator    validator
}

// sqlMsgParser is interface used for parse sql result message
type sqlMsgParser interface {
	EntryDuplicate(msg string) (entry, key string)
	NoReferencedRow(msg string) (fk string)
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// ExpenditureRepository return implementation of domain.ExpenditureRepository using mysql
func ExpenditureRepository(
	db *sqlx.DB,
	sp sqlMsgParser,
	v validator,
) domain.ExpenditureRepository {
	repo := &expenditureRepository{
		db:           db,
		sqlMsgParser: sp,
		validator:    v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.Expenditure{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.ExpenditureBabyTag{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}
