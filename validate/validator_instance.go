package validate

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// validatorInstance is global variable returned in customValidator function
var validatorInstance *customValidator

// initialize customValidator with registering custom validation
func init() {
	v := validator.New()

	_ = v.RegisterValidation("uuid", isValidateUUID)
	_ = v.RegisterValidation("range", isWithinRange)

	v.RegisterCustomTypeFunc(sqlNullStringTypeConverter, sql.NullString{})

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
	var v reflect.Value
	if reflect.TypeOf(s).Kind() != reflect.Ptr {
		v = reflect.ValueOf(s)
	} else {
		v = reflect.New(reflect.TypeOf(s).Elem()).Elem()
		v.Set(reflect.ValueOf(s).Elem())
	}

	if v.Kind() != reflect.Struct {
		return &validator.InvalidValidationError{Type: v.Type()}
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Type().Kind() != reflect.Ptr {
			continue
		}
		if f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}
	}

	return mv.Struct(v)
}
