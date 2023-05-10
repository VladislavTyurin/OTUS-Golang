package common

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tagStringLenPattern = regexp.MustCompile(`len:(\d+)`)
var tagStringInPattern = regexp.MustCompile(`in:(\w+(,\w+)*)`)
var tagRegexpPattern = regexp.MustCompile(`regexp:\\d\+`)

var ErrStringTooLong = errors.New("field value too long")
var ErrRegexpField = errors.New("incorrect field")

type stringValidator struct {
	field reflect.StructField
}

func (sv *stringValidator) Validate(tag string) error {
	tags := strings.Split(tag, "|")
	for _, t := range tags {
		if tagStringLenPattern.MatchString(t) {
			if err := sv.checkLen(t); err != nil {
				return err
			}
		} else if tagStringInPattern.MatchString(t) {
			if err := sv.checkIn(t); err != nil {
				return err
			}
		} else if tagRegexpPattern.MatchString(t) {
			if err := sv.checkRegexp(t); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w for type 'int': %s", ErrTagInvalid, t)
		}
	}
	return nil
}

func (sv *stringValidator) ValidationError(err error) bool {
	return errors.Is(err, ErrTagInvalid) || errors.Is(err, ErrStringTooLong) ||
		errors.Is(err, ErrRegexpField) || errors.Is(err, ErrValueNotFoundInSet)
}

func (sv *stringValidator) checkLen(tag string) error {
	groups := tagStringLenPattern.FindStringSubmatch(tag)
	maxLen, err := strconv.Atoi(groups[1])
	if err != nil {
		return err
	}
	fieldValue := reflect.ValueOf(sv.field).String()
	if len(fieldValue) > maxLen {
		return fmt.Errorf("%w: %s with len %d, but max len is %d",
			ErrStringTooLong, fieldValue, len(fieldValue), maxLen)
	}
	return nil
}

func (sv *stringValidator) checkIn(tag string) error {
	groups := tagStringInPattern.FindStringSubmatch(tag)
	values := strings.Split(groups[1], ",")
	fieldValue := reflect.ValueOf(sv.field).String()
	for _, v := range values {
		if fieldValue == v {
			return nil
		}
	}
	return fmt.Errorf("%w: %s not in %v", ErrValueNotFoundInSet, fieldValue, values)
}

func (sv *stringValidator) checkRegexp(tag string) error {
	r, err := regexp.Compile(tag)
	if err != nil {
		return err
	}
	fieldValue := reflect.ValueOf(sv.field).String()

	if r.MatchString(fieldValue) {
		return fmt.Errorf("%w: %s", ErrRegexpField, fieldValue)
	}
	return nil
}
