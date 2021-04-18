package validate

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
)

// modelValidatorInstance is global variable returned in ModelValidator function
var modelValidatorInstance *modelValidator

// initialize modelValidatorInstance with registering custom validation
func init() {
	v := validator.New()

	_ = v.RegisterValidation("uuid", isValidateUUID)
	_ = v.RegisterValidation("range", isWithinRange)

	v.RegisterCustomTypeFunc(sqlNullStringType, sql.NullString{})

	modelValidatorInstance = &modelValidator{v}
}

// ModelValidator function return modelValidatorInstance global variable
func ModelValidator() *modelValidator {
	return modelValidatorInstance
}

// modelValidator is struct embedding *validator.Validate which is registered custom validation
type modelValidator struct {
	*validator.Validate
}

func (mv *modelValidator) ValidateStruct(s interface{}) error {
	return mv.Struct(s)
}
