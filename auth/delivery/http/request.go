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
