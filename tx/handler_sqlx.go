package tx

import (
	"database/sql"
)

// sqlxHandler is struct that handle transaction with sqlx package
type sqlxHandler struct { db *sqlx.DB }
func NewSqlxHandler(db *sqlx.DB) *sqlxHandler { return &sqlxHandler{db} }
