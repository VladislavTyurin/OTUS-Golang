package integervalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var tagIntMinMaxPattern = regexp.MustCompile(`(min|max):(\d+)`)
var tagIntInPattern = regexp.MustCompile(`(in):(\d+(,\d+)*)`)

var ErrTagInvalid = errors.New("invalid tag for type 'int'")
var ErrValueLessThanMin = errors.New("field value is less than min")
var ErrValueGreaterThanMax = errors.New("field value is greater than max")
var ErrValueNotFoundInSet = errors.New("field value not found in set")

func Validate(tag string, field reflect.StructField) error {
	tags := strings.Split(tag, "|")
	for _, t := range tags {
		if tagIntMinMaxPattern.MatchString(t) {
			groups := tagIntMinMaxPattern.FindAllStringSubmatch(t, -1)
			for _, g := range groups {
				if g[1] == "min" {
					if err := checkMin(g[2], field); err != nil {
						return err
					}
				} else {
					if err := checkMax(g[2], field); err != nil {
						return err
					}
				}
			}
		} else if tagIntInPattern.MatchString(t) {
			groups := tagIntMinMaxPattern.FindAllStringSubmatch(t, -1)
			for _, g := range groups {
				if err := checkIn(g[2], field); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("%w: %s", ErrTagInvalid, t)
		}
	}
	return nil
}

func checkMin(tagValue string, field reflect.StructField) error {
	minValue, err := getValueFromTag(tagValue, field)
	if err != nil {
		return err
	}
	fieldValue := reflect.ValueOf(field).Int()
	if fieldValue < int64(minValue) {
		return fmt.Errorf("%w: %d < %d", ErrValueLessThanMin, fieldValue, minValue)
	}
	return nil
}

func checkMax(tagValue string, field reflect.StructField) error {
	maxValue, err := getValueFromTag(tagValue, field)
	if err != nil {
		return err
	}
	fieldValue := reflect.ValueOf(field).Int()
	if fieldValue > int64(maxValue) {
		return fmt.Errorf("%w: %d > %d", ErrValueGreaterThanMax, fieldValue, maxValue)
	}
	return nil
}

func checkIn(tagValue string, field reflect.StructField) error {
	values := strings.Split(tagValue, ",")
	fieldValue := reflect.ValueOf(field).Int()
	for _, v := range values {
		if vInt, err := getValueFromTag(v, field); err != nil {
			return err
		} else {
			if fieldValue == int64(vInt) {
				return nil
			}
		}
	}
	return fmt.Errorf("%w: %d not in %v", ErrValueNotFoundInSet, fieldValue, values)
}

func getValueFromTag(tagValue string, field reflect.StructField) (int, error) {
	value, err := strconv.Atoi(tagValue)
	if err != nil {
		return 0, err
	}
	return value, nil
}
