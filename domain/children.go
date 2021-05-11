package domain

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MyFirstBabyTime/Server/tx"
)

// ChildrenUsecase is interface about usecase layer using in delivery layer
type ChildrenUsecase interface {
	CreateNewChildren(ctx context.Context, c *Children, profile []byte) (uuid string, err error)
}

// ChildrenRepository is repository interface about Children model
type ChildrenRepository interface {
	GetByUUID(ctx tx.Context, uuid string) (children Children, err error)
	GetAvailableUUID(ctx tx.Context) (*string, error)
	Store(ctx tx.Context, c *Children) error
}

// Children is model represent parent children using in children domain
type Children struct {
	UUID       *string    `db:"uuid" validate:"required,uuid=children"`
	ParentUUID *string    `db:"parent_uuid" validate:"required,uuid=parent"`
	Name       *string    `db:"name" validate:"required,min=1,max=10"`
	Birth      *time.Time `db:"birth" validate:"required"`
	Sex        *string    `db:"sex" validate:"required,oneof=male female"`
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
		FOREIGN KEY (parent_uuid)
			REFERENCES parent_auth (uuid)
			ON DELETE CASCADE
	)
`
}

// TableName return table name about Children model
func (_ Children) TableName() string {
	return "children"
}

// GenerateRandomUUID generate & return random uuid value
func (c Children) GenerateRandomUUID() string {
	rand.Seed(time.Now().UnixNano())
	is := []rune("0123456789")
	random := make([]rune, 10)
	for i := range random {
		random[i] = is[rand.Intn(len(is))]
	}
	return fmt.Sprintf("c%s", string(random))
}

// GenerateProfileUri method return ProfileUri value with field value
func (c Children) GenerateProfileUri() string {
	return fmt.Sprintf("/profiles/children/uuid/%s", StringValue(c.UUID))
}
