package common

import (
	"errors"
	"reflect"
)

var intValidator = integerValidator{}
var strValidator = stringValidator{}

var ErrTagInvalid = errors.New("invalid tag")
var ErrValueNotFoundInSet = errors.New("field value not found in set")

type Validator interface {
	Validate(tag string) error
	ValidationError(err error) bool
}

func GetValidator(field reflect.StructField) Validator {
	switch field.Type.Kind() {
	case reflect.Int:
		intValidator.field = field
		return &intValidator
	case reflect.String:
		strValidator.field = field
		return &strValidator
	default:
		return nil
	}
}
