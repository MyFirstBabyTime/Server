package http

// sendCertifyCodeToPhoneRequest is request for authHandler.SendCertifyCodeToPhone
type sendCertifyCodeToPhoneRequest struct {
	PhoneNumber string `uri:"phone_number" validate:"required,len=11"`
}
