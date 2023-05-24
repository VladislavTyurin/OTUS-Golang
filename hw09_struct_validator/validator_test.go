package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	errs "github.com/VladislavTyurin/OTUS-Golang/hw09_struct_validator/errors"
	"github.com/stretchr/testify/require"
)

type UserRole string

type errSlice []error

func (e errSlice) Error() string {
	var builder strings.Builder
	for _, err := range e {
		builder.WriteString(err.Error() + "\n")
	}
	return builder.String()
}

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func errorIs(errs ValidationErrors, target error) bool {
	for _, e := range errs {
		if errors.Is(e.Err, target) {
			return true
		}
	}
	return false
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{},
			expectedErr: errSlice{
				errs.ErrValueLessThanMin,
				errs.ErrRegexpField,
				errs.ErrValueNotFoundInSet,
			},
		},
		{
			in: User{
				ID:     strings.Repeat("i", 37),
				Age:    18,
				Email:  "user@mail.com",
				Role:   "stuff",
				Phones: []string{""},
			},
			expectedErr: errSlice{errs.ErrStringTooLong},
		},
		{
			in: User{
				ID:     "id",
				Age:    17,
				Email:  "user@mail.com",
				Role:   "stuff",
				Phones: []string{""},
			},
			expectedErr: errSlice{errs.ErrValueLessThanMin},
		},
		{
			in: User{
				ID:     "id",
				Age:    51,
				Email:  "user@mail.com",
				Role:   "stuff",
				Phones: []string{""},
			},
			expectedErr: errSlice{errs.ErrValueGreaterThanMax},
		},
		{
			in: User{
				ID:     "id",
				Age:    20,
				Email:  "usermail.com",
				Role:   "stuff",
				Phones: []string{""},
			},
			expectedErr: errSlice{errs.ErrRegexpField},
		},
		{
			in: User{
				ID:     "id",
				Age:    20,
				Email:  "user@mail.com",
				Role:   "notstuff",
				Phones: []string{""},
			},
			expectedErr: errSlice{errs.ErrValueNotFoundInSet},
		},
		{
			in: User{
				ID:     "id",
				Age:    20,
				Email:  "user@mail.com",
				Role:   "stuff",
				Phones: []string{strings.Repeat("1", 12)},
			},
			expectedErr: errSlice{errs.ErrStringTooLong},
		},
		{
			in: App{
				Version: "00000000",
			},
			expectedErr: errSlice{errs.ErrStringTooLong},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1",
				Age:    20,
				Email:  "user@mail.com",
				Role:   "admin",
				Phones: []string{"00000000000"},
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			t.Log(err, tt.expectedErr)
			if tt.expectedErr != nil {
				require.Error(t, err)
				require.ErrorAs(t, err, &ValidationErrors{})
				require.Equal(t, len(tt.expectedErr.(errSlice)), len(err.(ValidationErrors))) //nolint:errorlint
				for _, e := range tt.expectedErr.(errSlice) {                                 //nolint:errorlint
					require.True(t, errorIs(err.(ValidationErrors), e)) //nolint:errorlint
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIsStruct(t *testing.T) {
	v := "not a struct"
	v1 := User{}
	require.True(t, isStruct(v1))
	require.False(t, isStruct(v))
}
