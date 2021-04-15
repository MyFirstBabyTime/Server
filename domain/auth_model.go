package domain

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID string `db:"uuid" validate:"required"`
	ID   string `db:"id" validate:"required"`
	PW   string `db:"pw" validate:"required"`
}
