package mysql

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
)

// parentPhoneCertifyRepository is implementation of domain.AuthRepository using mysql
type parentPhoneCertifyRepository struct {
	domain.ParentPhoneCertifyRepository

	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
	validator    validator
}

// ParentPhoneCertifyRepository return implementation of domain.ParentPhoneCertifyRepository using mysql
func ParentPhoneCertifyRepository(db *sqlx.DB, sp sqlMsgParser, v validator) domain.ParentPhoneCertifyRepository {
	repo := &parentPhoneCertifyRepository{
		db:           db,
		sqlMsgParser: sp,
		validator:    v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.ParentPhoneCertify{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent phone certify").Error())
	}
	return repo
}

func (pp *parentPhoneCertifyRepository) GetByPhoneNumber(ctx tx.Context, pn string) (ppc domain.ParentPhoneCertify, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").From("parent_phone_certify").Where("phone_number = ?", pn).ToSql()

	switch err = _tx.Get(&ppc, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = rowNotExistErr{errors.Wrap(err, "failed to select parent phone certify")}
	default:
		err = errors.Wrap(err, "select parent phone certify return unexpected error")
	}
	return
}

func (pp *parentPhoneCertifyRepository) Store(ctx tx.Context, ppc *domain.ParentPhoneCertify) (err error) {
	if ppc.CertifyCode == 0 {
		ppc.CertifyCode = ppc.GenerateCertifyCode()
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Insert("parent_phone_certify").
		Columns("parent_uuid", "phone_number", "certify_code").
		Values(ppc.ParentUUID, ppc.PhoneNumber, ppc.CertifyCode).ToSql()

	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_DUP_ENTRY:
			err = errors.Wrap(err, "failed to insert parent phone certify")
			_, key := pp.sqlMsgParser.EntryDuplicate(tErr.Message)
			err = entryDuplicateErr{err, key}
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to insert parent phone certify")
			fk := pp.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = noReferencedRowErr{err, fk}
		default:
			err = errors.Wrap(err, "insert parent auth return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "insert parent auth return unexpected error type")
	}
	return
}
