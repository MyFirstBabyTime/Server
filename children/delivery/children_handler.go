package delivery

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/MyFirstBabyTime/Server/domain"
)

// childrenHandler represent the http handler for children
type childrenHandler struct {
	cUsecase   domain.ChildrenUsecase
	validator  validator
	jwtHandler jwtHandler
}

// jwtHandler is interface of jwt handler
type jwtHandler interface {
	// ParseUUIDFromToken parse token & return token payload and type
	ParseUUIDFromToken(c *gin.Context)
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// NewChildrenHandler will initialize the children resources endpoint
func NewChildrenHandler(r *gin.Engine, cu domain.ChildrenUsecase, v validator, jh jwtHandler) {
	h := &childrenHandler{
		cUsecase:   cu,
		validator:  v,
		jwtHandler: jh,
	}

	r.POST("parents/uuid/:parent_uuid/children", h.jwtHandler.ParseUUIDFromToken, h.CreateNewChildren)
}
