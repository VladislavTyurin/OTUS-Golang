package validator

import (
	"reflect"
)

var (
	intValidator = integerValidator{}
	strValidator = stringValidator{}
	slValidator  = sliceValidator{}
)

type Validator interface {
	Validate(tag string) error
	ValidationError(err error) bool
}

func GetValidator(fieldValue reflect.Value, fieldType reflect.Type) Validator {
	switch fieldType.Kind() { //nolint:exhaustive
	case reflect.Int:
		intValidator.fieldValue = fieldValue
		return &intValidator
	case reflect.String:
		strValidator.fieldValue = fieldValue
		return &strValidator
	case reflect.Slice:
		slValidator.fieldValue = fieldValue
		return &slValidator
	default:
		return nil
	}
}
