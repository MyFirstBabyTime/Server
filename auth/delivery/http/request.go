package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
