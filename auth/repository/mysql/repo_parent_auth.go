package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/MyFirstBabyTime/Server/tx"
)

// mysqlAuthRepository is implementation of domain.AuthRepository using mysql
type parentAuthRepository struct {
	migrator migrator
	db *sqlx.DB
}

// ParentAuthRepository return implementation of domain.ParentAuthRepository using mysql
func ParentAuthRepository(db *sqlx.DB) domain.ParentAuthRepository {
	repo := &parentAuthRepository{db: db}
	if err := repo.migrator.MigrateModel(repo.db, domain.ParentAuth{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}

// Store is implement domain.AuthRepository interface
func (ar *parentAuthRepository) Store(ctx tx.Context, pa *domain.ParentAuth) (err error) {
	return
}

// getAvailableUUID method return available uuid of parent auth table
func (ar *parentAuthRepository) getAvailableUUID(ctx tx.Context) (uuid string) {
	auth := new(domain.ParentAuth)

	for {
		auth.SetRandomUUID()
		_, err := ar.GetByUUID(ctx, auth.UUID)
		if isRowNotExist(err) {
			return auth.UUID
		}
	}
}
