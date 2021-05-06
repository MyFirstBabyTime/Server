package mysql

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
)

// childrenRepository is implementation of domain.ChildrenRepository using mysql
type childrenRepository struct {
	domain.ChildrenRepository
	myCfg childrenRepositoryConfig

	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
	validator    validator
}

// childrenRepositoryConfig is interface get config value for children repository
type childrenRepositoryConfig interface{}

// sqlMsgParser is interface used for parse sql result message
type sqlMsgParser interface {
	EntryDuplicate(msg string) (entry, key string)
	NoReferencedRow(msg string) (fk string)
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// ChildrenRepository return implementation of domain.ChildrenRepository using mysql
func ChildrenRepository(
	cfg childrenRepositoryConfig,
	db *sqlx.DB,
	sp sqlMsgParser,
	v validator,
) domain.ChildrenRepository {
	repo := &childrenRepository{
		myCfg:        cfg,
		db:           db,
		sqlMsgParser: sp,
		validator:    v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.Children{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent children model").Error())
	}
	return repo
}
