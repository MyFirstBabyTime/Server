package mysql

type expenditureRepository struct {
	db        *sqlx.DB
	migrator  migrator
	validator validator
}
