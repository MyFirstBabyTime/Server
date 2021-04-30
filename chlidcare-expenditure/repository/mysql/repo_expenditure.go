package mysql

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

type expenditureRepository struct {
	db        *sqlx.DB
	migrator  migrator
	validator validator
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
