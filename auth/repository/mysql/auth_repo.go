package mysql

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
)

// mysqlAuthRepository is implementation of domain.AuthRepository using mysql
type mysqlAuthRepository struct {
	conn *sql.DB
}

// NewMysqlAuthRepository return implementation of domain.AuthRepository using mysql
func NewMysqlAuthRepository(conn *sql.DB) domain.AuthRepository {
	repo := &mysqlAuthRepository{conn}
	if err := repo.migrateSchema(domain.AuthMysqlSchema); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate schema").Error())
	}
	return repo
}

// BeginTx is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) { return ar.conn.BeginTx(ctx, opts) }

// Commit is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) Commit(tx *sql.Tx) error { return tx.Commit() }

// Rollback is method of domain.AuthRepository interface
func (ar *mysqlAuthRepository) Rollback(tx *sql.Tx) error { return tx.Rollback() }

// migrate
func (ar *mysqlAuthRepository) migrateSchema(schema string) (err error) {
	_, err = ar.conn.Exec(schema)
	return
}
