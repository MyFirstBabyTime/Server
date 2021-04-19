package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

// authHandler represent the http handler for article
type authHandler struct {
	aUsecase  domain.AuthUsecase
	validator validator
}

// validator is interface used for validating struct value
type validator interface {
	ValidateStruct(s interface{}) (err error)
}

// NewAuthHandler will initialize the auth/ resources endpoint
func NewAuthHandler(e *gin.Engine, au domain.AuthUsecase, v validator) {
	h := &authHandler{
		aUsecase:  au,
		validator: v,
	}

	e.POST("phones/phone-number/:phone_number/certify-code", h.SendCertifyCodeToPhone)
	//e.GET("/articles/:id", handler.GetByID)
}

// SendCertifyCodeToPhone is implement domain.AuthUsecase interface
func (ah *authHandler) SendCertifyCodeToPhone(c *gin.Context) {
	req := new(sendCertifyCodeToPhoneRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	switch err := ah.aUsecase.SendCertifyCodeToPhone(c.Request.Context(), req.PhoneNumber); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to send certify code to phone")
		c.JSON(http.StatusOK, resp)
	case usecaseErr:
		resp := defaultResp(tErr.Status(), tErr.Code(), tErr.Error())
		c.JSON(tErr.Status(), resp)
	default:
		msg := errors.Wrap(err, "SendCertifyCodeToPhone return unexpected error").Error()
		resp := defaultResp(http.StatusInternalServerError, 0, msg)
		c.JSON(http.StatusInternalServerError, resp)
	}
	return
}

// CertifyPhoneWithCode is implement domain.AuthUsecase interface
func (ah *authHandler) CertifyPhoneWithCode(c *gin.Context) {
	req := new(certifyPhoneWithCodeRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	switch err := ah.aUsecase.CertifyPhoneWithCode(c.Request.Context(), req.PhoneNumber, req.CertifyCode); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to send certify code to phone")
		c.JSON(http.StatusOK, resp)
	case usecaseErr:
		c.JSON(tErr.Status(), defaultResp(tErr.Status(), tErr.Code(), tErr.Error()))
	default:
		msg := errors.Wrap(err, "CertifyPhoneWithCode return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// bindRequest method bind *gin.Context to request having BindFrom method
func (ah *authHandler) bindRequest(req interface{ BindFrom(ctx *gin.Context) error }, c *gin.Context) error {
	if err := req.BindFrom(c); err != nil {
		return errors.Wrap(err, "failed to bind req")
	}
	if err := ah.validator.ValidateStruct(req); err != nil {
		return errors.Wrap(err, "invalid request")
	}
	return nil
}

// defaultResp return response have status, code, message inform
func defaultResp(status, code int, msg string) (resp struct{
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}) {
	resp.Status, resp.Code, resp.Message = status, code, msg
	return
}
