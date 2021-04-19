package validate

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
)

// validatorInstance is global variable returned in customValidator function
var validatorInstance *customValidator

// initialize customValidator with registering custom validation
func init() {
	v := validator.New()

	_ = v.RegisterValidation("uuid", isValidateUUID)
	_ = v.RegisterValidation("range", isWithinRange)

	v.RegisterCustomTypeFunc(sqlNullStringType, sql.NullString{})

	validatorInstance = &customValidator{v}
}

// New function return customValidator global variable
func New() *customValidator {
	return validatorInstance
}

// customValidator is struct embedding *validator.Validate which is registered custom validation
type customValidator struct {
	*validator.Validate
}

func (mv *customValidator) ValidateStruct(s interface{}) error {
	return mv.Struct(s)
}
