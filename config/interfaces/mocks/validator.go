package mocks

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/stretchr/testify/mock"
)

// MockValidator mocks the validator
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(req interface{}) error {
	args := m.Called(req)
	return args.Error(0)
}

// MockFieldError implements validator.FieldError
type MockFieldError struct {
	Fields string
	Tags   string
	Params string
}

func (m MockFieldError) Field() string      { return m.Fields }
func (m MockFieldError) Tag() string        { return m.Tags }
func (m MockFieldError) Param() string      { return m.Params }
func (m MockFieldError) ActualTag() string  { return "" }
func (m MockFieldError) Kind() reflect.Kind { return reflect.String }
func (m MockFieldError) Type() reflect.Type { return nil }
func (m MockFieldError) Value() interface{} { return nil }
func (m MockFieldError) StructField() string {
	return ""
}
func (m MockFieldError) Translate(_ ut.Translator) string {
	return "translated"
}

func (m MockFieldError) StructNamespace() string { return "" }
func (m MockFieldError) Namespace() string       { return "" }
