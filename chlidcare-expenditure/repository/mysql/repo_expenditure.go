package mysql

type expenditureRepository struct {
	db        *sqlx.DB
	migrator  migrator
	validator validator
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}