package mysql

import (
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
