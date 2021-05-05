package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

//expenditureHandler represent the http handler for article
type expenditureHandler struct {
	eUsecase  domain.ExpenditureUsecase
	validator validator
	jwtHandler jwtHandler
}

// validator if interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// jwtHandler is interface of jwt handler
type jwtHandler interface {
	// ParseUUIDFromToken parse token & return token payload and type
	ParseUUIDFromToken (c *gin.Context)
}

// NewExpenditureHandler vil initialize the expenditure endpoint
func NewExpenditureHandler(r *gin.Engine, eu domain.ExpenditureUsecase, v validator, jh jwtHandler) {
	h := &expenditureHandler{
		eUsecase:  eu,
		validator: v,
		jwtHandler: jh,
	}

	r.POST("expenditure/registration", h.jwtHandler.ParseUUIDFromToken, h.ExpenditureRegistration)
}

}