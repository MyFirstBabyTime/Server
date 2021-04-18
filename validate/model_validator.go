package validate

// modelValidator is struct embedding *validator.Validate which is registered custom validation
type modelValidator struct {
	*validator.Validate
}

func (mv *modelValidator) ValidateStruct(s interface{}) error {
	return mv.Struct(s)
}
