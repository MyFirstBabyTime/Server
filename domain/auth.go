package domain

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	SignUpParent(ctx gin.Context)
}

// AuthRepository is abstract interface about repository layer using in usecase layer
type AuthRepository interface {
	parentAuthRepository

	BeginTx(context.Context, *sql.TxOptions) (*sqlx.Tx, error) // BeginTx method start transaction
	Commit(tx *sqlx.Tx) (err error)                            // Commit method commit transaction
	Rollback(tx *sqlx.Tx) (err error)                          // Rollback method rollback transaction
}

type parentAuthRepository interface {
	CreateParentAuth(tx *sqlx.Tx, auth *ParentAuth) error
}
