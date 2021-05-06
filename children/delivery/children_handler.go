package delivery

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
)

// childrenHandler represent the http handler for children
type childrenHandler struct {
	aUsecase   domain.ChildrenUsecase
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
