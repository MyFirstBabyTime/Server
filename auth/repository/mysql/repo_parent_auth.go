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
func (ar *parentAuthRepository) GetByUUID(ctx tx.Context, uuid string) (auth struct {
	domain.ParentAuth
	domain.ParentPhoneCertify
}, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("parent_auth.*, IF(phone_number IS NULL, '', phone_number) AS phone_number").
		From("parent_auth").
		LeftJoin("parent_phone_certify ON parent_auth.uuid = parent_phone_certify.parent_uuid").
		Where("parent_auth.uuid = ?", uuid).ToSql()

	switch err = _tx.Get(&auth, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = rowNotExistErr{errors.Wrap(err, "failed to select parent auth")}
	default:
		err = errors.Wrap(err, "select parent auth return unexpected error var")
	}
	return
}

// GetByID is implement domain.AuthRepository interface
func (ar *parentAuthRepository) GetByID(ctx tx.Context, id string) (auth domain.ParentAuth, err error) {
	return
}

// Store is implement domain.AuthRepository interface
func (ar *parentAuthRepository) Store(ctx tx.Context, pa *domain.ParentAuth) (err error) {
	if pa.UUID == "" {
		pa.UUID = ar.getAvailableUUID(ctx)
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Insert(pa.TableName()).Columns("uuid", "id", "pw", "name", "profile_uri").
		Values(pa.UUID, pa.ID, pa.PW, pa.Name, pa.ProfileUri).ToSql()
	
	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_DUP_ENTRY:
			err = errors.Wrap(err, "failed to insert parent auth")
			err = entryDuplicateErr{err, ar.sqlMsgParser.EntryDuplicate(tErr.Message)}
		default:
			err = errors.Wrap(err, "insert parent auth return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "insert parent auth return unexpected error type")
	}
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
