package mysql

import "github.com/jmoiron/sqlx"

// childrenRepository is implementation of domain.ChildrenRepository using mysql
type childrenRepository struct {
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