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

// parentAuthRepository is implementation of domain.ParentAuthRepository using mysql
type parentAuthRepository struct {
	myCfg parentAuthRepositoryConfig

	db           *sqlx.DB
	migrator     migrator
	sqlMsgParser sqlMsgParser
	validator    validator
}

// parentAuthRepositoryConfig is interface get config value for parent auth repository
type parentAuthRepositoryConfig interface{}

// sqlMsgParser is interface used for parse sql result message
type sqlMsgParser interface {
	EntryDuplicate(msg string) (entry, key string)
	NoReferencedRow(msg string) (fk string)
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// ParentAuthRepository return implementation of domain.ParentAuthRepository using mysql
func ParentAuthRepository(
	cfg parentAuthRepositoryConfig,
	db *sqlx.DB,
	sp sqlMsgParser,
	v validator,
) domain.ParentAuthRepository {
	repo := &parentAuthRepository{
		myCfg:        cfg,
		db:           db,
		sqlMsgParser: sp,
		validator:    v,
	}

	if err := repo.migrator.MigrateModel(repo.db, domain.ParentAuth{}); err != nil {
		log.Fatal(errors.Wrap(err, "failed to migrate parent auth model").Error())
	}
	return repo
}

// GetByUUID is implement domain.ParentAuthRepository interface
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
		err = domain.ErrRowNotExist{RepoErr: errors.Wrap(err, "failed to select parent auth")}
	default:
		err = errors.Wrap(err, "select parent auth return unexpected error")
	}
	return
}

// GetByID is implement domain.ParentAuthRepository interface
func (ar *parentAuthRepository) GetByID(ctx tx.Context, id string) (auth struct {
	domain.ParentAuth
	domain.ParentPhoneCertify
}, err error) {
	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Select("parent_auth.*, IF(phone_number IS NULL, '', phone_number) AS phone_number").
		From("parent_auth").
		LeftJoin("parent_phone_certify ON parent_auth.uuid = parent_phone_certify.parent_uuid").
		Where("parent_auth.id = ?", id).ToSql()

	switch err = _tx.Get(&auth, _sql, args...); err {
	case nil:
		break
	case sql.ErrNoRows:
		err = domain.ErrRowNotExist{RepoErr: errors.Wrap(err, "failed to select parent auth")}
	default:
		err = errors.Wrap(err, "select parent auth return unexpected error var")
	}
	return
}

// Store is implement domain.ParentAuthRepository interface
func (ar *parentAuthRepository) Store(ctx tx.Context, pa *domain.ParentAuth) (err error) {
	if domain.StringValue(pa.UUID) == "" {
		if uuid, err := ar.GetAvailableUUID(ctx); err != nil {
			err = errors.Wrap(err, "failed to GetAvailableUUID")
			return err
		} else {
			pa.UUID = domain.String(uuid)
		}
	}

	if err = ar.validator.ValidateStruct(pa); err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.Wrap(err, "failed to validate domain.ParentAuth")}
		return
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, _ := squirrel.Insert("parent_auth").
		Columns("uuid", "id", "pw", "name", "profile_uri").
		Values(pa.UUID, pa.ID, pa.PW, pa.Name, pa.ProfileUri).ToSql()

	switch _, err = _tx.Exec(_sql, args...); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_DUP_ENTRY:
			err = errors.Wrap(err, "failed to insert parent auth")
			_, key := ar.sqlMsgParser.EntryDuplicate(tErr.Message)
			err = domain.ErrEntryDuplicate{RepoErr: err, DuplicateKey: key}
		default:
			err = errors.Wrap(err, "insert parent auth return unexpected code return")
		}
	default:
		err = errors.Wrap(err, "insert parent auth return unexpected error type")
	}
	return
}

// GetAvailableUUID method return available uuid of parent auth table
func (ar *parentAuthRepository) GetAvailableUUID(ctx tx.Context) (string, error) {
	pa := new(domain.ParentAuth)

	for {
		uuid := pa.GenerateRandomUUID()
		_, err := ar.GetByUUID(ctx, uuid)

		if err == nil {
			continue
		} else if _, ok := err.(domain.ErrRowNotExist); ok {
			return uuid, nil
		} else {
			return "", errors.Wrap(err, "failed to GetByUUID")
		}
	}
}

// Update method update tuple of domain.ParentAuth model by UUID field value
func (ar *parentAuthRepository) Update(ctx tx.Context, pa *domain.ParentAuth) (err error) {
	if domain.StringValue(pa.UUID) == "" {
		err = errors.New("UUID(PK) value in model must be set")
		return
	}

	if err = ar.validator.ValidateStruct(pa.GenerateValidModel()); err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.Wrap(err, "failed to validate domain.ParentAuth")}
		return
	}

	b := squirrel.Update("parent_auth").Where("uuid = ?", pa.UUID)
	if pa.ID != nil {
		if *pa.ID == "" {
			pa.ID = nil
		}
		b = b.Set("id", pa.ID)
	}
	if pa.PW != nil {
		if *pa.PW == "" {
			pa.PW = nil
		}
		b = b.Set("pw", pa.PW)
	}
	if pa.Name != nil {
		if *pa.Name == "" {
			pa.Name = nil
		}
		b = b.Set("name", pa.Name)
	}
	if pa.ProfileUri != nil {
		if *pa.ProfileUri == "" {
			pa.ProfileUri = nil
		}
		b = b.Set("profile_uri", pa.ProfileUri)
	}

	_tx, _ := ctx.Tx().(*sqlx.Tx)
	_sql, args, err := b.ToSql()
	if err != nil {
		err = domain.ErrInvalidModel{RepoErr: errors.New("update statements must have at least one")}
		return
	}

	if _, err = _tx.Exec(_sql, args...); err != nil {
		err = errors.Wrap(err, "failed to update parent auth")
	}
	return
}
