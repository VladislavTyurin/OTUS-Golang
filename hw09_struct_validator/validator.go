package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	errs "github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/errors"
	"github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/validator"
)

var m = sync.Mutex{}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	if len(v) != 0 {
		builder.WriteString("Validation Error:\n")
	}
	for _, err := range v {
		builder.WriteString("Field: " + err.Field + ": " + err.Err.Error() + "\n")
	}
	return builder.String()
}

func isStruct(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func Validate(v interface{}) error {
	m.Lock()
	defer m.Unlock()
	if !isStruct(v) {
		return fmt.Errorf("%w: %v", errs.ErrNotStruct, v)
	}

	errs := make(ValidationErrors, 0)
	for i := 0; i < reflect.ValueOf(v).NumField(); i++ {
		fieldValue := reflect.ValueOf(v).Field(i)
		fieldType := reflect.TypeOf(v).Field(i)
		tag, ok := fieldType.Tag.Lookup("validate")
		if !ok {
			continue
		}
		v := validator.GetValidator(fieldValue, fieldType.Type)
		if v != nil {
			err := v.Validate(tag)
			if err != nil {
				if v.ValidationError(err) {
					errs = append(errs, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				} else {
					return err
				}
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
