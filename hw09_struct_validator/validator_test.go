package hw09_struct_validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

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

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Igor",
				Age:    20,
				Email:  "test@mail.ru",
				Role:   "admin",
				Phones: []string{"7", "7", "7", "7", "7", "7", "7", "7", "7", "7", "7"},
			},
			nil,
		},
		{
			App{
				Version: "12345",
			},
			nil,
		},
		{
			Token{
				Header:    []byte("qwerty"),
				Payload:   []byte("qwerty"),
				Signature: []byte("qwerty"),
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.NoError(t, err, tt.expectedErr)
		})
	}
}

func TestNegativeValidation(t *testing.T) {
	user := User{
		ID:     "12",
		Name:   "Igor",
		Age:    15,
		Email:  "test@mail.",
		Role:   "adm",
		Phones: []string{"7", "7", "7", "7", "7", "7", "7", "7", "7", "7"},
	}

	err := Validate(user)
	var actualErrors ValidationErrors
	errors.As(err, &actualErrors)
	require.Error(t, actualErrors)
	errorFields := []string{"ID", "Name", "Age", "Email", "Role", "Phones"}
	for i := range len(actualErrors) {
		assert.Contains(t, errorFields, actualErrors[i].Field)
	}
}
