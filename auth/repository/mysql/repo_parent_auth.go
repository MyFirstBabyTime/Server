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
	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
}

// sqlMsgParser is interface used for parse sql result message
type sqlMsgParser interface{
	EntryDuplicate(msg string) (entry string)
}

// ParentAuthRepository return implementation of domain.ParentAuthRepository using mysql
func ParentAuthRepository(db *sqlx.DB) domain.ParentAuthRepository {
	repo := &parentAuthRepository{db: db}
	if err := repo.migrator.MigrateModel(repo.db, domain.ParentAuth{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}

// GetByUUID is implement domain.AuthRepository interface
func (ar *parentAuthRepository) GetByUUID(ctx tx.Context, uuid string) (auth domain.ParentAuth, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").
		From(domain.ParentAuth{}.TableName()).
		Where("uuid = ?", uuid).ToSql()

	switch err = _tx.Get(&auth, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = rowNotExistErr{errors.Wrap(err, "failed to select parent auth")}
	default:
		err = errors.Wrap(err, "unexpected failed to select parent auth")
	}
	return
}
	return
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
