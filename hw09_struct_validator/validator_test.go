package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "User case [positive]",
			in: User{
				ID:     "12345_12345_12345_12345_12345_12345_",
				Name:   "Vasya",
				Age:    42,
				Email:  "vasya@pochta.com",
				Role:   "admin",
				Phones: []string{"89165556677", "78074445522"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			name: "App case [positive]",
			in: App{
				Version: "1.0.3",
			},
			expectedErr: nil,
		},
		{
			name: "Token case [positive]",
			in: Token{
				Header:    []byte("un"),
				Payload:   []byte("validated"),
				Signature: []byte("struct"),
			},
			expectedErr: nil,
		},
		{
			name: "Response case [positive]",
			in: Response{
				Code: 200,
				Body: "some body",
			},
			expectedErr: nil,
		},
		{
			name:        "Arbitrary case [negative]",
			in:          42,
			expectedErr: ErrType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("unexpected error: want %s, but got %s", tt.expectedErr, err)
			}
		})
	}
}

func TestValidateWithValidationErrors(t *testing.T) {
	const (
		userErr = `ID: [validation error] invalid "len"
Age: [validation error] invalid "min"
Email: [validation error] invalid "regexp"
Role: [validation error] invalid "in"
Phones[0]: [validation error] invalid "len"
Phones[1]: [validation error] invalid "len"`

		appErr      = "Version: [validation error] invalid \"len\""
		responseErr = "Code: [validation error] invalid \"in\""
	)

	tests := []struct {
		name        string
		in          interface{}
		expectedErr string
	}{
		{
			name: "User case [negative]",
			in: User{
				ID:     "12345_12345",
				Name:   "smallVasya",
				Age:    12,
				Email:  "vasya&pochta.com",
				Role:   "user",
				Phones: []string{"8916", "2"},
				meta:   nil,
			},
			expectedErr: userErr,
		},
		{
			name: "App case [negative]",
			in: App{
				Version: "0.1",
			},
			expectedErr: appErr,
		},
		{
			name: "Response case [negative]",
			in: Response{
				Code: 301,
				Body: "another some body",
			},
			expectedErr: responseErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err == nil {
				t.Fatal("expected error, but got nil")
			}
			if err.Error() != tt.expectedErr {
				t.Fatalf("unexpected error: want %s, but got %s", tt.expectedErr, err)
			}
		})
	}
}
