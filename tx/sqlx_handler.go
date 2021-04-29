package tx

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// sqlxHandler is struct that handle transaction with sqlx package
type sqlxHandler struct{ db *sqlx.DB }

func NewSqlxHandler(db *sqlx.DB) *sqlxHandler { return &sqlxHandler{db} }

// sqlxTxKey is used for key for transaction value in tx context
type sqlxTxKey struct{}

// BeginTx method start transaction (get option from ctx)
func (sh *sqlxHandler) BeginTx(ctx context.Context, opts interface{}) (txCtx Context, err error) {
	var tx *sqlx.Tx
	switch opts.(type) {
	case *sql.TxOptions:
		tx, err = sh.db.BeginTxx(ctx, opts.(*sql.TxOptions))
	default:
		tx, err = sh.db.BeginTxx(ctx, nil)
	}

	if err != nil {
		err = errors.Wrap(err, "failed to begin sqlx transaction")
		return
	}

	txCtx = &txContext{
		Context: context.Background(),
		txKey:   sqlxTxKey{},
	}
	txCtx.SetTx(tx)
	return
}

// Commit method commit transaction
func (sh *sqlxHandler) Commit(ctx Context) (err error) {
	return ctx.Tx().(*sqlx.Tx).Commit()
}

// Rollback method rollback transaction
func (sh *sqlxHandler) Rollback(ctx Context) (err error) {
	return ctx.Tx().(*sqlx.Tx).Rollback()
}
