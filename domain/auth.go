package domain


// ParentAuth is model represent parent auth using in auth domain
type ParentAuth struct {
	UUID string `json:"uuid" validate:"required"`
	ID   string `json:"id" validate:"required"`
	PW   string `json:"pw" validate:"required"`
}

// AuthUsecase is abstract interface about usecase layer using in delivery layer
type AuthUsecase interface {
	SignUpParent(ctx gin.Context)
}
