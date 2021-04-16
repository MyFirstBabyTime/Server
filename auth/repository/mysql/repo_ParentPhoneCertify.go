package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
)

// parentPhoneCertifyRepository is implementation of domain.AuthRepository using mysql
type parentPhoneCertifyRepository struct {
	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
}

// ParentPhoneCertifyRepository return implementation of domain.ParentPhoneCertifyRepository using mysql
func ParentPhoneCertifyRepository(db *sqlx.DB, sp sqlMsgParser) domain.ParentPhoneCertifyRepository {
	repo := &parentPhoneCertifyRepository{
		db:           db,
		sqlMsgParser: sp,
	}
	if err := repo.migrator.MigrateModel(repo.db, domain.ParentPhoneCertify{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent phone certify").Error())
	}
	return repo
}
