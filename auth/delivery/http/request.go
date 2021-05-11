package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
)

// sendCertifyCodeToPhoneRequest is request for authHandler.SendCertifyCodeToPhone
type sendCertifyCodeToPhoneRequest struct {
	PhoneNumber string `uri:"phone_number" validate:"required,len=11"`
}

func (r *sendCertifyCodeToPhoneRequest) BindFrom(c *gin.Context) error {
	return errors.Wrap(c.BindUri(r), "failed to BindUri")
}

// certifyPhoneWithCodeRequest is request for authHandler.CertifyPhoneWithCode
type certifyPhoneWithCodeRequest struct {
	PhoneNumber string `uri:"phone_number" validate:"required"`
	CertifyCode int64  `json:"certify_code" validate:"required"`
}

func (r *certifyPhoneWithCodeRequest) BindFrom(c *gin.Context) error {
	if err := c.BindUri(r); err != nil {
		return errors.Wrap(err, "failed to BindUri")
	}

	if err := c.BindJSON(r); err != nil {
		return errors.Wrap(err, "failed to BindJSON")
	}
	return nil
}

// signUpParentRequest is request for authHandler.SignUpParent
type signUpParentRequest struct {
	ParentID      string                `form:"id" json:"id" validate:"required,min=4,max=20"`
	ParentPW      string                `form:"pw" json:"pw" validate:"required,min=6,max=20"`
	Name          string                `form:"name" json:"name" validate:"required,max=20"`
	PhoneNumber   string                `form:"phone_number" json:"phone_number" validate:"required,len=11"`
	Profile       *multipart.FileHeader `form:"profile"`
	ProfileBase64 string                `json:"profile_base64"`
}

func (r *signUpParentRequest) BindFrom(c *gin.Context) error {
	switch c.ContentType() {
	case "application/json":
		return errors.Wrap(c.BindJSON(r), "failed to BindJSON")
	default:
		return errors.Wrap(c.Bind(r), "failed to Bind")
	}
}

type loginParentAuthRequest struct {
	ID string `json:"id" validate:"required"`
	PW string `json:"pw" validate:"required"`
}

func (r *loginParentAuthRequest) BindFrom(c *gin.Context) error {
	return errors.Wrap(c.BindJSON(r), "failed to BindJSON")
}

type getParentInformByIDRequest struct {
	ParentID string `uri:"parent_id" validate:"required"`
}

func (r *getParentInformByIDRequest) BindFrom(c *gin.Context) error {
	return errors.Wrap(c.BindUri(r), "failed to BindUri")
}

type updateParentInformRequest struct {
	ParentUUID    string                `uri:"parent_uuid" validate:"required"`
	Name          *string               `form:"name" json:"name" validate:"max=20"`
	Profile       *multipart.FileHeader `form:"profile"`
	ProfileBase64 string                `json:"profile_base64"`
}

func (r *updateParentInformRequest) BindFrom(c *gin.Context) (err error) {
	if err = c.BindUri(r); err != nil {
		return errors.Wrap(err, "failed to BindUri")
	}

	switch c.ContentType() {
	case "application/json":
		err = errors.Wrap(c.BindJSON(r), "failed to BindJSON")
	default:
		err = errors.Wrap(c.Bind(r), "failed to Bind")
	}
	if err != nil {
		return
	}

	if r.Name == nil && r.Profile == nil {
		return errors.New("all field blank is not allowed")
	}
	return
}
