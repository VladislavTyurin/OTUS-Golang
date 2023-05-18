package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	errs "github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/errors"
)

var tagStringLenPattern = regexp.MustCompile(`len:(\d+)`)
var tagStringInPattern = regexp.MustCompile(`in:(\w+(,\w+)*)`)
var tagRegexpPattern = regexp.MustCompile(`regexp:`)

type stringValidator struct {
	fieldValue reflect.Value
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
			return fmt.Errorf("%w for type 'string': %s", errs.ErrTagInvalid, t)
		}
	}
	return nil
}

func (sv *stringValidator) ValidationError(err error) bool {
	return errors.Is(err, errs.ErrTagInvalid) || errors.Is(err, errs.ErrStringTooLong) ||
		errors.Is(err, errs.ErrRegexpField) || errors.Is(err, errs.ErrValueNotFoundInSet)
}

func (sv *stringValidator) checkLen(tag string) error {
	groups := tagStringLenPattern.FindStringSubmatch(tag)
	maxLen, err := strconv.Atoi(groups[1])
	if err != nil {
		return err
	}

	if len(sv.fieldValue.String()) > maxLen {
		return fmt.Errorf("%w: %s with len %d, but max len is %d",
			errs.ErrStringTooLong, sv.fieldValue.String(), len(sv.fieldValue.String()), maxLen)
	}
	return nil
}

func (sv *stringValidator) checkIn(tag string) error {
	groups := tagStringInPattern.FindStringSubmatch(tag)
	values := strings.Split(groups[1], ",")
	fieldValue := sv.fieldValue.String()
	for _, v := range values {
		if fieldValue == v {
			return nil
		}
	}
	return fmt.Errorf("%w: %s not in %v", errs.ErrValueNotFoundInSet, fieldValue, values)
}

func (sv *stringValidator) checkRegexp(tag string) error {
	tag = strings.TrimLeft(tag, "regexp:")
	r, err := regexp.Compile(tag)
	if err != nil {
		return err
	}
	fieldValue := sv.fieldValue.String()

	if !r.MatchString(fieldValue) {
		return fmt.Errorf("%w: %s", errs.ErrRegexpField, fieldValue)
	}
	return nil
}
