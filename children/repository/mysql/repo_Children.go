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

// Store is implement Store method of domain.ChildrenRepository interface
func (cr *childrenRepository) Store(ctx tx.Context, c *domain.Children) (err error) {
	if domain.StringValue(c.UUID) == "" {
		if c.UUID, err = cr.GetAvailableUUID(ctx); err != nil {
			return errors.Wrap(err, "failed to GetAvailableUUID")
		}
	}

	if err = cr.validator.ValidateStruct(c); err != nil {
		return domain.ErrInvalidModel{RepoErr: errors.Wrap(err, "failed to validate domain.Children")}
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Insert("children").Columns("uuid", "parent_uuid", "name", "birth", "sex", "profile_uri").
		Values(c.UUID, c.ParentUUID, c.Name, c.Birth, c.Sex, c.ProfileUri).ToSql()

	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_REFERENCED_ROW_2:
			err = errors.Wrap(err, "failed to insert children")
			fk := cr.sqlMsgParser.NoReferencedRow(tErr.Message)
			err = domain.ErrNoReferencedRow{RepoErr: err, ForeignKey: fk}
		default:
			err = errors.Wrap(err, "insert children return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "insert children return unexpected error type")
	}
	return
}

// GetByUUID is implement GetByUUID method of domain.ChildrenRepository interface
func (cr *childrenRepository) GetByUUID(ctx tx.Context, uuid string) (children domain.Children, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("*").From("children").Where("uuid = ?", uuid).ToSql()

	switch err = _tx.Get(&children, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = domain.ErrRowNotExist{RepoErr: errors.Wrap(err, "failed to select children")}
	default:
		err = errors.Wrap(err, "select children return unexpected error")
	}
	return
}

// GetAvailableUUID method return available uuid of children table
func (cr *childrenRepository) GetAvailableUUID(ctx tx.Context) (*string, error) {
	pa := new(domain.Children)

	for {
		uuid := pa.GenerateRandomUUID()
		_, err := cr.GetByUUID(ctx, uuid)

		if err == nil {
			continue
		} else if _, ok := err.(domain.ErrRowNotExist); ok {
			return &uuid, nil
		} else {
			return nil, errors.Wrap(err, "failed to GetByUUID")
		}
	}
}
