package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// migrator is struct that migrate to mysql repository
type migrator struct{}

// MigrateModel method migrate model to db received from parameter
func (m migrator) MigrateModel(db *sqlx.DB, model interface {
	TableName() string // TableName return table name about model
	Schema() string    // Schema return schema SQL about model
}) (err error) {
	sql, _, _ := squirrel.Select("*").From(model.TableName()).ToSql()
	switch _, err = db.Query(sql); tErr := err.(type) {
	case nil:
		break
	case *mysql.MySQLError:
		switch tErr.Number {
		case mysqlerr.ER_NO_SUCH_TABLE:
			_, err = db.Exec(model.Schema())
			err = errors.Wrapf(err, "failed to exec %s model schema", model.TableName())
		default:
			err = errors.Wrapf(err, "check table query returns unexpected mysql error code")
		}
	default:
		err = errors.Wrapf(err, "check table query returns unexpected error type")
	}

	return
}
