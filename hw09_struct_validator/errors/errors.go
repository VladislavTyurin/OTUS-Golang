package errors

import "errors"

var ErrNotStruct = errors.New("data type is not a struct")
var ErrTagInvalid = errors.New("invalid tag")
var ErrValueNotFoundInSet = errors.New("field value not found in set")
var ErrValueLessThanMin = errors.New("field value is less than min")
var ErrValueGreaterThanMax = errors.New("field value is greater than max")
var ErrStringTooLong = errors.New("field value too long")
var ErrRegexpField = errors.New("incorrect field")
