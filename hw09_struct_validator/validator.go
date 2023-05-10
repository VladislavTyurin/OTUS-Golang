package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/common"
)

var ErrNotStruct = errors.New("data type is not a struct")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func isStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func Validate(v interface{}) error {
	if !isStruct(v) {
		return fmt.Errorf("%w: %v", ErrNotStruct, v)
	}

	errs := make(ValidationErrors, 0)
	for i := 0; i < reflect.ValueOf(v).NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)
		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		v := common.GetValidator(field)
		if v != nil {
			err := v.Validate(tag)
			if err != nil {
				if v.ValidationError(err) {
					errs = append(errs, ValidationError{
						Field: field.Name,
						Err:   err,
					})
				} else {
					return err
				}
			}
		}
	}
	return errs
}
