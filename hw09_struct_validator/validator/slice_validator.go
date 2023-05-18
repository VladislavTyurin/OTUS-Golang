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
	switch elemType { //nolint:exhaustive
	case reflect.Int:
		sv.subValidator = &intValidator
		for i := 0; i < slice.Len(); i++ {
			intValidator.fieldValue = sv.fieldValue.Index(i)
			if err := intValidator.Validate(tag); err != nil {
				return err
			}
		}
	case reflect.String:
		sv.subValidator = &strValidator
		for i := 0; i < slice.Len(); i++ {
			strValidator.fieldValue = sv.fieldValue.Index(i)
			if err := strValidator.Validate(tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func (sv *sliceValidator) ValidationError(err error) bool {
	return sv.subValidator.ValidationError(err)
}
