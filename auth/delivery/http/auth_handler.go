package http

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

}