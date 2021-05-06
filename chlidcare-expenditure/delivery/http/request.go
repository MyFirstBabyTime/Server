package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type expenditureRegistration struct {
	ParentUUID string   `json:"parent_uuid" validate:"required,uuid=parent"`
	BabyUUIDs  []string `json:"baby_uuids" validate:"required"`
	Name       string   `json:"name" validate:"required"`
	Amount     int64    `json:"amount" validate:"required,range=0~1000000000"`
	Rating     int64    `json:"rating" validate:"required,range=0~5"`
	Link       string   `json:"link"`
}

func (r *expenditureRegistration) BindFrom(c *gin.Context) error {
	return errors.Wrap(c.BindJSON(r), "failed to BindJSON")
}
