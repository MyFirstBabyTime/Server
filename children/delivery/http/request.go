package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
)

// createNewChildrenRequest is request for childrenHandler.CreateNewChildren
type createNewChildrenRequest struct {
	ParentUUID    string                `uri:"parent_uuid" validate:"required"`
	Name          string                `form:"name" json:"name" validate:"required,max=20"`
	Birth         string                `form:"birth" json:"birth" validate:"required,max=20"`
	Sex           string                `form:"sex" json:"sex" validate:"required,max=20,oneof=male female"`
	Profile       *multipart.FileHeader `form:"profile"`
	ProfileBase64 string                `json:"profile_base64"`
}

func (r *createNewChildrenRequest) BindFrom(c *gin.Context) error {
	if err := c.BindUri(r); err != nil {
		return errors.Wrap(err, "failed to BindUri")
	}

	switch c.ContentType() {
	case "application/json":
		return errors.Wrap(c.BindJSON(r), "failed to BindJSON")
	default:
		return errors.Wrap(c.Bind(r), "failed to Bind")
	}
}
