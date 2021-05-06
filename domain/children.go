package domain

import "time"

// Children is model represent parent children using in children domain
type Children struct {
	UUID       *string    `db:"uuid"`
	ParentUUID *string    `db:"parent_uuid"`
	Name       *string    `db:"name"`
	Birth      *time.Time `db:"birth"`
	Sex        *string    `db:"sex"`
	ProfileUri *string    `db:"profile_uri"`
}

func (_ Children) Schema() string {
	return `CREATE TABLE children (
		uuid        CHAR(11) NOT NULL,
		parent_uuid CHAR(11) NOT NULL,
		name        VARCHAR(10) NOT NULL,
		birth       DATETIME NOT NULL,
		sex         VARCHAR(10) NOT NULL,
		profile_uri VARCHAR(100),
		PRIMARY KEY (uuid),
		FOREIGN KEY (parent_uuid) REFERENCE parent_auth (uuid)
	)
`
}

func (_ Children) TableName() string {
	return "children"
}
