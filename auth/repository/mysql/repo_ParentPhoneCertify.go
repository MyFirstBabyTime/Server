package mysql

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
)

// parentPhoneCertifyRepository is implementation of domain.AuthRepository using mysql
type parentPhoneCertifyRepository struct {
	domain.ParentPhoneCertifyRepository

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

func (pp *parentPhoneCertifyRepository) GetByPhoneNumber(ctx tx.Context, pn string) (certify domain.ParentPhoneCertify, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").From("parent_phone_certify").Where("phone_number = ?", pn).ToSql()

	switch err = _tx.Get(&certify, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = rowNotExistErr{errors.Wrap(err, "failed to select parent phone certify")}
	default:
		err = errors.Wrap(err, "select parent phone certify return unexpected error")
	}
	return
}
