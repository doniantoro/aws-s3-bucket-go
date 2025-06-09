package utils

import (
	"aws-s3-bucket/models/dto"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/go-playground/validator"
)

func ExtractBase64(request string) (raw, contentType, formatType string, decodedBytes []byte, err error) {

	data := strings.Split(request, ",")
	if len(data) != 2 {
		err = fmt.Errorf("invalid base64 string format")
		return
	}

	getContentType := strings.Split(data[0], ";")
	getContentType = strings.Split(getContentType[0], ":")
	if len(getContentType) != 2 {
		err = fmt.Errorf("invalid base64 string format")
		return
	}

	getFormatType := strings.Split(getContentType[1], "/")
	if len(getFormatType) != 2 {
		err = fmt.Errorf("invalid base64 string format")
		return
	}

	decodedBytes, err = base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return
	}

	return data[1], getContentType[1], getFormatType[1], decodedBytes, nil
}

func FormatMessageValidator(err validator.FieldError) string {
	param := err.Param()
	message := err.Tag()
	field := err.Field()
	switch err.Tag() {
	case "required":
		message = "Required"
	case "numeric":
		message = "accepted:format=number"
	case "email":
		message = "accepted:format=email"
	case "gt":
		message = fmt.Sprintf("accepted:gt=%s", param)
	case "gte":
		message = fmt.Sprintf("accepted:gte=%s", param)
	case "lt":
		message = fmt.Sprintf("accepted:lt=%s", param)
	case "lte":
		message = fmt.Sprintf("accepted:lte=%s", param)
	case "min":
		message = fmt.Sprintf("accepted:min=%s", param)
	case "max":
		message = fmt.Sprintf("accepted:max=%s", param)
	case "len":
		message = fmt.Sprintf("accepted:len=%s", param)
	case "eq":
		message = fmt.Sprintf("accepted:eq=%s", param)
	case "dateformat":
		message = "accepted:format=YYYY-MM-DD"
	case "oneof":
		message = fmt.Sprintf("accepted:value=%s", param)
	}
	result := fmt.Sprintf("Parameter %s %s field", strcase.ToSnake(field), message)
	return result
}

func UnwrapValidation(err error) []dto.ErrorValidation {
	errors := make([]dto.ErrorValidation, len(err.(validator.ValidationErrors)))
	for i, err := range err.(validator.ValidationErrors) {
		errors[i] = dto.ErrorValidation{
			Message:   FormatMessageValidator(err),
			Parameter: strcase.ToSnake(err.Field()),
		}
	}
	return errors
}
