package delivery

import "mime/multipart"

// createNewChildrenRequest is request for childrenHandler.CreateNewChildren
type createNewChildrenRequest struct {
	ParentUUID string                `uri:"parent_uuid" validate:"required"`
	Name       string                `form:"name" validate:"required,max=20"`
	Birth      string                `form:"birth" validate:"required,max=20"`
	Sex        string                `form:"sex" validate:"required,max=20"`
	Profile    *multipart.FileHeader `form:"profile"`
}
