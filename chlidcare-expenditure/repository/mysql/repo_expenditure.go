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

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

func ExpenditureRepository(
	db *sqlx.DB,
	v validator,
) domain.ExpenditureRepository {
	repo := &expenditureRepository{
		db:        db,
		validator: v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.Expenditure{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}
