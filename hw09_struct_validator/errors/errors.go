package errors

import "errors"

var (
	ErrNotStruct           = errors.New("data type is not a struct")
	ErrTagInvalid          = errors.New("invalid tag")
	ErrValueNotFoundInSet  = errors.New("field value not found in set")
	ErrValueLessThanMin    = errors.New("field value is less than min")
	ErrValueGreaterThanMax = errors.New("field value is greater than max")
	ErrStringTooLong       = errors.New("field value too long")
	ErrRegexpField         = errors.New("incorrect field")
)
