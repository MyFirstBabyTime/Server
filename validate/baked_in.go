package validate

import (
	"database/sql/driver"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
	"strings"
)

// isValidateUUID function return if uuid format is validate
func isValidateUUID(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	switch fl.Param() {
	case "parent":
		return parentUUIDRegex.MatchString(fl.Field().String())
	case "item":
		return itemUUIDRegex.MatchString(fl.Field().String())
	case "children":
		return childrenRegex.MatchString(fl.Field().String())
	}
	return false
}

// isWithinRange function return if field value is within range
func isWithinRange(fl validator.FieldLevel) bool {
	_range := strings.Split(fl.Param(), "~")
	if len(_range) != 2 {
		return false
	}

	start, err := strconv.Atoi(_range[0])
	if err != nil {
		return false
	}
	end, err := strconv.Atoi(_range[1])
	if err != nil {
		return false
	}

	field := int(fl.Field().Int())
	return field >= start && field <= end
}

func isNotEmptyValue(fl validator.FieldLevel) bool {
	if fl.Field().Interface() == nil {
		return false
	}

	switch fl.Field().Kind() {
	case reflect.String:
		if fl.Field().String() == "" {
			return false
		}
	case reflect.Int64:
		if fl.Field().Int() == 0 {
			return false
		}
	default:
		return false
	}
	return true
}

// sqlNullStringTypeConverter function assert driver.Valuer & return Value()
func sqlNullStringTypeConverter(field reflect.Value) (v interface{}) {
	v = ""
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		if value, err := valuer.Value(); err == nil && value != nil {
			v = value
		}
	}
	return
}
