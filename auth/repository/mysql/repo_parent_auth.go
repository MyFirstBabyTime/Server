package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
)

// mysqlAuthRepository is implementation of domain.AuthRepository using mysql
type mysqlAuthRepository struct {
	db *sqlx.DB
}

// NewMysqlAuthRepository return implementation of domain.AuthRepository using mysql
func NewMysqlAuthRepository(db *sqlx.DB) domain.ParentAuthRepository {
	repo := &mysqlAuthRepository{db}

	if err := repo.migrateModel(domain.ParentAuth{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}

	if err := repo.migrateModel(domain.ParentPhoneCertify{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent phone number certify").Error())
	}
	return repo
}

// Store is implement domain.AuthRepository interface
func (ar *mysqlAuthRepository) Store(ctx tx.Context, pa *domain.ParentAuth) (err error) {
	return
}

// migrateModel method migrate model received from parameter to this repository
func (ar *mysqlAuthRepository) migrateModel(model interface {
	TableName() string // TableName return table name about model
	Schema() string    // Schema return schema SQL about model
}) (err error) {
	_sql, _, _ := squirrel.Select("*").From(model.TableName()).ToSql()
	switch _, err = ar.db.Query(_sql); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_SUCH_TABLE:
			_, err = ar.db.Exec(model.Schema())
			err = errors.Wrapf(err, "failed to exec %s model schema", model.TableName())
		default:
			err = errors.Wrapf(err, "check table query returns unexpected mysql error code")
		}
	default:
		err = errors.Wrapf(err, "check table query returns unexpected error type")
	}

	return
}
