package http

import (
	"github.com/MyFirstBabyTime/Server/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

// authHandler represent the http handler for article
type authHandler struct {
	aUsecase   domain.AuthUsecase
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

// NewAuthHandler will initialize the auth/ resources endpoint
func NewAuthHandler(r *gin.Engine, au domain.AuthUsecase, v validator, jh jwtHandler) {
	h := &authHandler{
		aUsecase:   au,
		validator:  v,
		jwtHandler: jh,
	}

	r.POST("phones/phone-number/:phone_number/certify-code", h.SendCertifyCodeToPhone)
	r.POST("phones/phone-number/:phone_number/certification", h.CertifyPhoneWithCode)
	r.POST("parents", h.SignUpParent)
	r.POST("login/parent", h.LoginParentAuth)
	r.GET("parents/id/:parent_id/existence", h.CheckIfParentIDExist)
	r.PATCH("parents/uuid/:parent_uuid", h.jwtHandler.ParseUUIDFromToken, h.UpdateParentInform)
}

// SendCertifyCodeToPhone deliver data to SendCertifyCodeToPhone of domain.AuthUsecase
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
	case domain.UsecaseError:
		resp := defaultResp(tErr.Status, tErr.Code, tErr.Error())
		c.JSON(tErr.Status, resp)
	default:
		msg := errors.Wrap(err, "SendCertifyCodeToPhone return unexpected error").Error()
		resp := defaultResp(http.StatusInternalServerError, 0, msg)
		c.JSON(http.StatusInternalServerError, resp)
	}
	return
}

// CertifyPhoneWithCode deliver data to CertifyPhoneWithCode of domain.AuthUsecase
func (ah *authHandler) CertifyPhoneWithCode(c *gin.Context) {
	req := new(certifyPhoneWithCodeRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	switch err := ah.aUsecase.CertifyPhoneWithCode(c.Request.Context(), req.PhoneNumber, req.CertifyCode); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to certify phone with certify code")
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "CertifyPhoneWithCode return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// SignUpParent deliver data to SignUpParent of domain.AuthUsecase
func (ah *authHandler) SignUpParent(c *gin.Context) {
	req := new(signUpParentRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	pi := struct {
		*domain.ParentAuth
		*domain.ParentPhoneCertify
	}{
		ParentAuth: &domain.ParentAuth{
			ID:   domain.String(req.ParentID),
			PW:   domain.String(req.ParentPW),
			Name: domain.String(req.Name),
		},
		ParentPhoneCertify: &domain.ParentPhoneCertify{
			PhoneNumber: domain.String(req.PhoneNumber),
		},
	}

	switch uuid, err := ah.aUsecase.SignUpParent(c.Request.Context(), pi, req.Profile); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusCreated, 0, "succeed to sign up new parent auth")
		resp["parent_uuid"] = uuid
		c.JSON(http.StatusCreated, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "SignUpParent return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// LoginParentAuth deliver data to LoginParentAuth of domain.AuthUsecase
func (ah *authHandler) LoginParentAuth(c *gin.Context) {
	req := new(loginParentAuthRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	uuid, token, err := ah.aUsecase.LoginParentAuth(c.Request.Context(), req.ID, req.PW)
	switch tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to login parent auth")
		resp["uuid"], resp["token"] = uuid, token
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "LoginParentAuth return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// CheckIfParentIDExist deliver data to GetParentInformByID of domain.AuthUsecase
func (ah *authHandler) CheckIfParentIDExist(c *gin.Context) {
	req := new(getParentInformByIDRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	switch _, err := ah.aUsecase.GetParentInformByID(c.Request.Context(), req.ParentID); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "parent auth with that ID is exist")
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "GetParentInformByID return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// UpdateParentInform deliver data to UpdateParentInform of domain.AuthUsecase
func (ah *authHandler) UpdateParentInform(c *gin.Context) {
	req := new(updateParentInformRequest)
	if err := ah.bindRequest(req, c); err != nil {
		c.JSON(http.StatusBadRequest, defaultResp(http.StatusBadRequest, 0, err.Error()))
		return
	}

	pa := &domain.ParentAuth{
		Name: domain.String(domain.StringValue(req.Name)),
	}
	switch err := ah.aUsecase.UpdateParentInform(c.Request.Context(), req.ParentUUID, pa, req.Profile); tErr := err.(type) {
	case nil:
		resp := defaultResp(http.StatusOK, 0, "succeed to update parent inform")
		c.JSON(http.StatusOK, resp)
	case domain.UsecaseError:
		c.JSON(tErr.Status, defaultResp(tErr.Status, tErr.Code, tErr.Error()))
	default:
		msg := errors.Wrap(err, "UpdateParentInform return unexpected error").Error()
		c.JSON(http.StatusInternalServerError, defaultResp(http.StatusInternalServerError, 0, msg))
	}
	return
}

// bindRequest method bind *gin.Context to request having BindFrom method
func (ah *authHandler) bindRequest(req interface {
	BindFrom(ctx *gin.Context) error
}, c *gin.Context) error {
	if err := req.BindFrom(c); err != nil {
		return errors.Wrap(err, "failed to bind req")
	}
	if err := ah.validator.ValidateStruct(req); err != nil {
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
