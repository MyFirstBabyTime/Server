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
	myCfg parentPhoneCertifyRepositoryConfig

	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
	validator    validator
}

// ParentPhoneCertifyRepository return implementation of domain.ParentPhoneCertifyRepository using mysql
func ParentPhoneCertifyRepository(
	cfg parentPhoneCertifyRepositoryConfig,
	db *sqlx.DB,
	sp sqlMsgParser,
	v validator,
) domain.ParentPhoneCertifyRepository {
	repo := &parentPhoneCertifyRepository{
		myCfg:        cfg,
		db:           db,
		sqlMsgParser: sp,
		validator:    v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.ParentPhoneCertify{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent phone certify").Error())
	}
	return repo
}

// parentPhoneCertifyRepositoryConfig is interface get config value for parent phone certify repository
type parentPhoneCertifyRepositoryConfig interface{}

// GetByPhoneNumber is implement domain.ParentPhoneCertifyRepository interface
func (pp *parentPhoneCertifyRepository) GetByPhoneNumber(ctx tx.Context, pn string) (ppc domain.ParentPhoneCertify, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").From("parent_phone_certify").Where("phone_number = ?", pn).ToSql()

	switch err = _tx.Get(&ppc, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = domain.ErrRowNotExist{RepoErr: errors.Wrap(err, "failed to select parent phone certify")}
	default:
		err = errors.Wrap(err, "select parent phone certify return unexpected error")
	}
	return
}

// Store is implement domain.ParentPhoneCertifyRepository interface
func (pp *parentPhoneCertifyRepository) Store(ctx tx.Context, ppc *domain.ParentPhoneCertify) (err error) {
	if domain.Int64Value(ppc.CertifyCode) == 0 {
		ppc.CertifyCode = domain.Int64(ppc.GenerateCertifyCode())
	}

	if err = pp.validator.ValidateStruct(ppc); err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.Wrap(err, "failed to validate domain.ParentPhoneCertify")}
		return
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
			err = domain.ErrEntryDuplicate{RepoErr: err, DuplicateKey: key}
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to insert parent phone certify")
			fk := pp.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = domain.ErrNoReferencedRow{RepoErr: err, ForeignKey: fk}
		default:
			err = errors.Wrap(err, "insert parent auth return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "insert parent auth return unexpected error type")
	}
	return
}

// Update is implement domain.ParentPhoneCertifyRepository interface
// where -> PK, set -> field with value set (so, cannot set to NULL in this method)
func (pp *parentPhoneCertifyRepository) Update(ctx tx.Context, ppc *domain.ParentPhoneCertify) (err error) {
	if domain.StringValue(ppc.PhoneNumber) == "" {
		err = errors.New("PhoneNumber(PK) value in model must be set")
		return
	}

	if err = pp.validator.ValidateStruct(ppc.GenerateValidModel()); err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.Wrap(err, "failed to validate domain.ParentPhoneCertify")}
		return
	}

	b := squirrel.Update("parent_phone_certify").Where("phone_number = ?", ppc.PhoneNumber)
	if ppc.ParentUUID != nil {
		b = b.Set("parent_uuid", ppc.ParentUUID)
	}
	if ppc.CertifyCode != nil {
		b = b.Set("certify_code", ppc.CertifyCode)
	}
	if ppc.Certified != nil {
		b = b.Set("certified", ppc.Certified)
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, err := b.ToSql()
	if err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.New("update statements must have at least one")}
		return
	}

	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to update parent phone certify")
			fk := pp.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = domain.ErrNoReferencedRow{RepoErr: err, ForeignKey: fk}
		default:
			err = errors.Wrap(err, "update parent auth return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "update parent auth return unexpected error type")
	}
	return
}
