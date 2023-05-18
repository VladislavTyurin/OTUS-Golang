package validator

import (
	"reflect"
)

type sliceValidator struct {
	fieldValue   reflect.Value
	subValidator Validator
}

func (sv *sliceValidator) Validate(tag string) error {
	slice := sv.fieldValue
	elemType := sv.fieldValue.Type().Elem().Kind()
	switch elemType {
	case reflect.Int:
		sv.subValidator = &intValidator
		for i := 0; i < slice.Len(); i++ {
			intValidator.fieldValue = sv.fieldValue.Index(i)
			return intValidator.Validate(tag)
		}
	case reflect.String:
		sv.subValidator = &strValidator
		for i := 0; i < slice.Len(); i++ {
			strValidator.fieldValue = sv.fieldValue.Index(i)
			return strValidator.Validate(tag)
		}
	}
	return nil
}

func (sv *sliceValidator) ValidationError(err error) bool {
	return sv.subValidator.ValidationError(err)
}
