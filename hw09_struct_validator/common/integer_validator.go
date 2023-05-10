package common

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tagIntMinMaxPattern = regexp.MustCompile(`(min|max):(\d+)`)
var tagIntInPattern = regexp.MustCompile(`in:(\d+(,\d+)*)`)

var ErrValueLessThanMin = errors.New("field value is less than min")
var ErrValueGreaterThanMax = errors.New("field value is greater than max")

type integerValidator struct {
	field reflect.StructField
}

func (iv *integerValidator) Validate(tag string) error {
	tags := strings.Split(tag, "|")
	for _, t := range tags {
		if tagIntMinMaxPattern.MatchString(t) {
			if err := iv.checkMinMax(t); err != nil {
				return err
			}
		} else if tagIntInPattern.MatchString(t) {
			if err := iv.checkIn(t); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w for type 'int': %s", ErrTagInvalid, t)
		}
	}
	return nil
}

func (iv *integerValidator) ValidationError(err error) bool {
	return errors.Is(err, ErrTagInvalid) ||
		errors.Is(err, ErrValueGreaterThanMax) ||
		errors.Is(err, ErrValueLessThanMin) ||
		errors.Is(err, ErrValueNotFoundInSet)
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
	fieldValue := reflect.ValueOf(iv.field).Int()
	if fieldValue < int64(minValue) {
		return fmt.Errorf("%w: %d < %d", ErrValueLessThanMin, fieldValue, minValue)
	}
	return nil
}

func (iv *integerValidator) checkMax(tag string) error {
	maxValue, err := iv.getValueFromTag(tag)
	if err != nil {
		return err
	}
	fieldValue := reflect.ValueOf(iv.field).Int()
	if fieldValue > int64(maxValue) {
		return fmt.Errorf("%w: %d > %d", ErrValueGreaterThanMax, fieldValue, maxValue)
	}
	return nil
}

func (iv *integerValidator) checkIn(tag string) error {
	groups := tagIntInPattern.FindStringSubmatch(tag)
	values := strings.Split(groups[1], ",")
	fieldValue := reflect.ValueOf(iv.field).Int()
	for _, v := range values {
		if vInt, err := iv.getValueFromTag(v); err != nil {
			return err
		} else {
			if fieldValue == int64(vInt) {
				return nil
			}
		}
	}
	return fmt.Errorf("%w: %d not in %v", ErrValueNotFoundInSet, fieldValue, values)
}

func (iv *integerValidator) getValueFromTag(tag string) (int, error) {
	value, err := strconv.Atoi(tag)
	if err != nil {
		return 0, err
	}
	return value, nil
}
