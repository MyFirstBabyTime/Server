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

func (er *expenditureRepository) GetByUUID(ctx tx.Context, uuid string) (expenditure domain.Expenditure, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").
		From("expenditure").
		Where("expenditure.uuid = ?", uuid).ToSql()

	switch err = _tx.Get(&expenditure, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = domain.ErrRowNotExist{RepoErr: errors.Wrap(err, "failed to select expenditure uuid")}
	default:
		err = errors.Wrap(err, "select expenditure return unexpected error")
	}
	return
}

func (er *expenditureRepository) Store(ctx tx.Context, e *domain.Expenditure, babyUUIDs *[]string) (err error) {
	if e.UUID == nil {
		if e.UUID, err = er.GetAvailableUUID(ctx); err != nil {
			err = errors.Wrap(err, "failed to getAvailableUUID")
			return
		}
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Insert("expenditure").
		Columns("uuid", "parent_uuid", "name", "amount", "rating", "link").
		Values(e.UUID, e.ParentUUID, e.Name, e.Amount, e.Rating, e.Link).ToSql()

	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to insert expenditure")
			fk := er.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = domain.ErrNoReferencedRow{RepoErr: err, ForeignKey: fk}
			return
		default:
			err = errors.Wrap(err, "insert expenditure unexpected mysql error")
			return
		}
	default:
		err = errors.Wrap(err, "insert expenditure return unexpected error")
		return
	}

	for _, babyUUID := range *babyUUIDs {
		_sql, args, _ = squirrel.Insert("expenditure_baby_tag").
			Columns("expenditure_uuid", "baby_uuid").
			Values(*e.UUID, babyUUID).ToSql()

		_, err = _tx.Exec(_sql, args...)
		if err != nil {
			break
		}
	}

	switch tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to insert expenditure_baby_tag")
			fk := er.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = domain.ErrNoReferencedRow{RepoErr: err, ForeignKey: fk}
			return
		case mysqlerr.ER_DUP_ENTRY:
			err = errors.Wrap(err, "failed to insert expenditure_baby_tag")
			_, key := er.sqlMsgParser.EntryDuplicate(tErr.Message)
			err = domain.ErrEntryDuplicate{RepoErr: err, DuplicateKey: key}
			return
		default:
			err = errors.Wrap(err, "insert expenditure_baby_tag unexpected mysql error")
			return
		}
	default:
		err = errors.Wrap(err, "insert expenditure_baby_tag return unexpected error")
		return
	}

	return
}

func (er *expenditureRepository) GetAvailableUUID(ctx tx.Context) (*string, error) {
	e := new(domain.Expenditure)

	for {
		uuid := e.GenerateRandomUUID()
		_, err := er.GetByUUID(ctx, uuid)

		if err == nil {
			continue
		} else if _, ok := err.(domain.ErrRowNotExist); ok {
			return &uuid, nil
		} else {
			return nil, errors.Wrap(err, "failed to GetByUUID")
		}
	}
}
