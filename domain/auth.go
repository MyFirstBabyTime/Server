package domain


// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID string `json:"uuid" validate:"required"`
	ID   string `json:"id" validate:"required"`
	PW   string `json:"pw" validate:"required"`
}

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	SignUpParent(ctx gin.Context)
}

// AuthRepository is abstract interface about repository layer using in usecase layer
type AuthRepository interface {
	parentAuthRepository

	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) // BeginTx method start transaction
	Commit(tx *sql.Tx) (err error)                                     // Commit method commit transaction
	Rollback(tx *sql.Tx) (err error)                                   // Rollback method rollback transaction
}

type parentAuthRepository interface {
	CreateParentAuth(tx *sql.Tx, auth *ParentAuth) error
}
