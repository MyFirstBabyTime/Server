package domain

import (
	"context"
	"github.com/MyFirstBabyTime/Server/tx"
	"mime/multipart"
	"time"
)

// ChildrenUsecase is interface about usecase layer using in delivery layer
type ChildrenUsecase interface {
	CreateNewChildren(ctx context.Context, c *Children, profile *multipart.FileHeader) (uuid string, err error)
}

// ChildrenRepository is repository interface about Children model
type ChildrenRepository interface {
	Store(ctx tx.Context, c *Children) error
}

// Children is model represent parent children using in children domain
type Children struct {
	UUID       *string    `db:"uuid"`
	ParentUUID *string    `db:"parent_uuid"`
	Name       *string    `db:"name"`
	Birth      *time.Time `db:"birth"`
	Sex        *string    `db:"sex"`
	ProfileUri *string    `db:"profile_uri"`
}

// Schema return rdbms schema about Children model
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

// TableName return table name about Children model
func (_ Children) TableName() string {
	return "children"
}
