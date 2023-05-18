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

var (
	tagIntMinMaxPattern = regexp.MustCompile(`(min|max):(\d+)`)
	tagIntInPattern     = regexp.MustCompile(`in:(\d+(,\d+)*)`)
)

type integerValidator struct {
	fieldValue reflect.Value
}

func (iv *integerValidator) Validate(tag string) error {
	tags := strings.Split(tag, "|")
	for _, t := range tags {
		switch {
		case tagIntMinMaxPattern.MatchString(t):
			if err := iv.checkMinMax(t); err != nil {
				return err
			}
		case tagIntInPattern.MatchString(t):
			if err := iv.checkIn(t); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w for type 'int': %s", errs.ErrTagInvalid, t)
		}
	}
	return nil
}

func (iv *integerValidator) ValidationError(err error) bool {
	return errors.Is(err, errs.ErrTagInvalid) ||
		errors.Is(err, errs.ErrValueGreaterThanMax) ||
		errors.Is(err, errs.ErrValueLessThanMin) ||
		errors.Is(err, errs.ErrValueNotFoundInSet)
}

func (iv *integerValidator) checkMinMax(tag string) error {
	groups := tagIntMinMaxPattern.FindStringSubmatch(tag)

	if groups[1] == "min" {
		if err := iv.checkMin(groups[2]); err != nil {
			return err
		}
	} else {
		if err := iv.checkMax(groups[2]); err != nil {
			return err
		}
	}

	return nil
}

func (iv *integerValidator) checkMin(tag string) error {
	minValue, err := iv.getValueFromTag(tag)
	if err != nil {
		return err
	}
	fieldValue := iv.fieldValue.Int()
	if fieldValue < int64(minValue) {
		return fmt.Errorf("%w: %d < %d", errs.ErrValueLessThanMin, fieldValue, minValue)
	}
	return nil
}

func (iv *integerValidator) checkMax(tag string) error {
	maxValue, err := iv.getValueFromTag(tag)
	if err != nil {
		return err
	}

	if iv.fieldValue.Int() > int64(maxValue) {
		return fmt.Errorf("%w: %d > %d", errs.ErrValueGreaterThanMax, iv.fieldValue.Int(), maxValue)
	}
	return nil
}

func (iv *integerValidator) checkIn(tag string) error {
	groups := tagIntInPattern.FindStringSubmatch(tag)
	values := strings.Split(groups[1], ",")
	fieldValue := iv.fieldValue.Int()
	for _, v := range values {
		if vInt, err := iv.getValueFromTag(v); err != nil {
			return err
		} else if fieldValue == int64(vInt) {
			return nil
		}
	}
	return fmt.Errorf("%w: %d not in %v", errs.ErrValueNotFoundInSet, fieldValue, values)
}

func (iv *integerValidator) getValueFromTag(tag string) (int, error) {
	value, err := strconv.Atoi(tag)
	if err != nil {
		return 0, err
	}
	return value, nil
}
