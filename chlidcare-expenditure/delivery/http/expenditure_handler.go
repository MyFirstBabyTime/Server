package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

//expenditureHandler represent the http handler for article
type expenditureHandler struct {
	eUsecase   domain.ExpenditureUsecase
	validator  validator
	jwtHandler jwtHandler
}

// validator if interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// jwtHandler is interface of jwt handler
type jwtHandler interface {
	// ParseUUIDFromToken parse token & return token payload and type
	ParseUUIDFromToken(c *gin.Context)
}

// NewExpenditureHandler vil initialize the expenditure endpoint
func NewExpenditureHandler(r *gin.Engine, eu domain.ExpenditureUsecase, v validator, jh jwtHandler) {
	h := &expenditureHandler{
		eUsecase:   eu,
		validator:  v,
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

	err := eh.eUsecase.ExpenditureRegistration(c, &domain.Expenditure{
		ParentUUID: domain.String(req.ParentUUID),
		Name:       domain.String(req.Name),
		Amount:     domain.Int64(req.Amount),
		Rating:     domain.Int64(req.Rating),
		Link:       domain.String(req.Link),
	}, req.BabyUUIDs)

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

// bindRequest method bind *gin.Context to request having BindFrom method
func (eh *expenditureHandler) bindRequest(req interface {
	BindFrom(ctx *gin.Context) error
}, c *gin.Context) error {
	if err := req.BindFrom(c); err != nil {
		return errors.Wrap(err, "failed to bind req")
	}
	if err := eh.validator.ValidateStruct(req); err != nil {
		return errors.Wrap(err, "invalid request")
	}
	return nil
}

// defaultResp return response have status, code, message inform
func defaultResp(status, code int, msg string) (resp gin.H) {
	resp = gin.H{}
	resp["status"] = status
	resp["code"] = code
	resp["message"] = msg
	return
}
