package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"

	integervalidator "github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/integer_validator"
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
		switch field.Type.Kind() {
		case reflect.Int:
			err := integervalidator.Validate(tag, field)
			if errors.Is(err, integervalidator.ErrTagInvalid) ||
				errors.Is(err, integervalidator.ErrValueGreaterThanMax) ||
				errors.Is(err, integervalidator.ErrValueLessThanMin) ||
				errors.Is(err, integervalidator.ErrValueNotFoundInSet) {
				errs = append(errs, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		case reflect.String:
			fmt.Println("string", tag)
		case reflect.Slice:
			fmt.Println("slice", tag)
		default:
		}
	}
	return nil
}
