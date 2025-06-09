package utils

import (
	"aws-s3-bucket/config/interfaces/mocks"
	"aws-s3-bucket/models/dto"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
)

func Test_ExtractBase64(t *testing.T) {

	type expected struct {
		raw         string
		contentType string
		formatType  string
		err         error
	}
	tests := []struct {
		name     string
		input    string
		expected expected
	}{
		{
			name:  "Valid Base64 String",
			input: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUA",
			expected: expected{
				err:         nil,
				formatType:  "png",
				contentType: "image/png",
				raw:         "iVBORw0KGgoAAAANSUhEUgAAAAUA",
			},
		},
		{
			name:  "Invalid Base64 String - get content type",
			input: "dataimage/pngbase64,iVBORw0KGgoAAAANSUhEUgAAAAUA",
			expected: expected{
				err: fmt.Errorf("invalid base64 string format"),
			},
		},
		{
			name:  "Invalid Base64 String - get format type",
			input: "data:imagepngbase64,iVBORw0KGgoAAAANSUhEUgAAAAUA",
			expected: expected{
				err: fmt.Errorf("invalid base64 string format"),
			},
		},
		{
			name:  "Empty String",
			input: "",
			expected: expected{
				err: fmt.Errorf("invalid base64 string format"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, contentType, formatType, decodedBytes, err := ExtractBase64(tt.input)
			var decoded []byte
			if raw != "" {
				decoded, _ = base64.StdEncoding.DecodeString(raw)
			}

			require.Equal(t, tt.expected.err, err)
			require.Equal(t, decodedBytes, decoded)
			require.Equal(t, tt.expected.formatType, formatType)
			require.Equal(t, tt.expected.contentType, contentType)
			require.Equal(t, tt.expected.raw, raw)

		})
	}

}

func TestFormatMessageValidator(t *testing.T) {
	tests := []struct {
		name     string
		mockErr  validator.FieldError
		expected string
	}{
		{
			name: "required",
			mockErr: mocks.MockFieldError{
				Fields: "Email",
				Tags:   "required",
			},
			expected: "Parameter email Required field",
		},
		{
			name: "numeric",
			mockErr: mocks.MockFieldError{
				Fields: "Age",
				Tags:   "numeric",
			},
			expected: "Parameter age accepted:format=number field",
		},
		{
			name: "gt",
			mockErr: mocks.MockFieldError{
				Fields: "Score",
				Tags:   "gt",
				Params: "10",
			},
			expected: "Parameter score accepted:gt=10 field",
		},
		{
			name: "email",
			mockErr: mocks.MockFieldError{
				Fields: "EmailAddress",
				Tags:   "email",
			},
			expected: "Parameter email_address accepted:format=email field",
		},
		{
			name: "dateformat",
			mockErr: mocks.MockFieldError{
				Fields: "BirthDate",
				Tags:   "dateformat",
			},
			expected: "Parameter birth_date accepted:format=YYYY-MM-DD field",
		},
		{
			name: "oneof",
			mockErr: mocks.MockFieldError{
				Fields: "Status",
				Tags:   "oneof",
				Params: "active inactive",
			},
			expected: "Parameter status accepted:value=active inactive field",
		},
		{
			name: "gte",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "gte",
				Params: "10",
			},
			expected: "Parameter score accepted:gte=10 field",
		},
		{
			name: "lt",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "lt",
				Params: "10",
			},
			expected: "Parameter score accepted:lt=10 field",
		},
		{
			name: "lte",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "lte",
				Params: "10",
			},
			expected: "Parameter score accepted:lte=10 field",
		},
		{
			name: "min",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "min",
				Params: "10",
			},
			expected: "Parameter score accepted:min=10 field",
		},
		{
			name: "max",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "max",
				Params: "10",
			},
			expected: "Parameter score accepted:max=10 field",
		},
		{
			name: "len",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "len",
				Params: "10",
			},
			expected: "Parameter score accepted:len=10 field",
		},
		{
			name: "eq",
			mockErr: mocks.MockFieldError{
				Fields: "score",
				Tags:   "eq",
				Params: "10",
			},
			expected: "Parameter score accepted:eq=10 field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMessageValidator(tt.mockErr)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUnwrapValidation(t *testing.T) {
	mockErr := validator.ValidationErrors{
		mocks.MockFieldError{Fields: "Email", Tags: "required", Params: ""},
		mocks.MockFieldError{Fields: "Age", Tags: "gte", Params: "18"},
	}

	err := error(mockErr)

	expected := []dto.ErrorValidation{
		{
			Message:   "Parameter email Required field",
			Parameter: "email",
		},
		{
			Message:   "Parameter age accepted:gte=18 field",
			Parameter: "age",
		},
	}

	result := UnwrapValidation(err)

	require.Equal(t, expected, result)
}
