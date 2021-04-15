package mysql

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
)

// mysqlAuthRepository is implementation of domain.AuthRepository using mysql
type mysqlAuthRepository struct {
	db *sqlx.DB
}

// NewMysqlAuthRepository return implementation of domain.AuthRepository using mysql
func NewMysqlAuthRepository(db *sqlx.DB) domain.AuthRepository {
	repo := &mysqlAuthRepository{db}
	if err := repo.migrateModel(domain.ParentAuth{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}

// BeginTx is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) { return ar.db.BeginTxx(ctx, opts) }

// Commit is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) Commit(tx *sqlx.Tx) error { return tx.Commit() }

// Rollback is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) Rollback(tx *sqlx.Tx) error { return tx.Rollback() }

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
