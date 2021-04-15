package domain

// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID       string         `db:"uuid" validate:"required"`
	ID         string         `db:"id" validate:"required"`
	PW         string         `db:"pw" validate:"required"`
	Name       string         `db:"name" validate:"required"`
	ProfileUri sql.NullString `db:"profile_uri"`
}

// TableName return table name about model
func (pa ParentAuth) TableName() string {
	return "parent_auth"
}

// Schema return schema SQL about model
func (pa ParentAuth) Schema() string {
	return `CREATE TABLE parent_auth (
		uuid VARCHAR(11) NOT NULL,
		id VARCHAR(20) NOT NULL,
		pw VARCHAR(100) NOT NULL
	);`
}
