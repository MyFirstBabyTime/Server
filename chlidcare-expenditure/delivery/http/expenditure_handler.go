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

func (eh *expenditureHandler) ExpenditureRegistration(c *gin.Context) {
	req := new(expenditureRegistration)
	if err := eh.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	err := eh.eUsecase.ExpenditureRegistration(c,
		&domain.Expenditure{
			ParentUUID: &req.ParentUUID,
			Name:      	&req.Name,
			Amount:     &req.Amount,
			Rating:     &req.Rating,
			Link:       &req.Link,
		},
		&req.BabyUUIDs,
	)

	switch tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to registration expenditure")
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "ExpenditureRegistration return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}
