package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/go-playground/validator/v10"
)

// The function `ContainsFile` checks if a given file name exists in a slice of `fs.DirEntry` objects.
func ContainsFile(files []fs.DirEntry, name string) bool {
	for _, file := range files {
		if file.Name() == name {
			return true
		}
	}
	return false
}

// The function "ExtractInformation" extracts information from a given error message and returns a
// formatted string describing the error.
func ExtractInformation(err error) string {
	errMsg := fmt.Sprintf("%s ", err)
	castringError := regexp.MustCompile(`cannot unmarshal (.*?) into Go struct field (.*?) of type (.*?) `)

	if castringError.MatchString(errMsg) {
		field := strings.Split(castringError.FindStringSubmatch(errMsg)[2], ".")[1]
		givenType := castringError.FindStringSubmatch(errMsg)[1]
		expectedType := castringError.FindStringSubmatch(errMsg)[3]
		return fmt.Sprintf("the '%s' attribute should be of type %s, but a %s value was given", field, expectedType, givenType)
	}
	return ""
}

func ConvertValidationErrors(err error) []types.APIError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]types.APIError, len(ve))
		for i, fe := range ve {
			out[i] = types.APIError{
				Field:   strings.ToLower(fe.Field()),
				Message: msgForTag(fe.Tag()),
			}
			return out
		}
	}
	return nil

}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "opps! this field is required"
	case "email":
		return "opps! email is invalid"
	case "min":
		return "opps! this field is too short"
	case "max":
		return "opps! this field is too long"
	case "gt":
		return "opps! this field is less than required"
	case "lt":
		return "opps! this field is greater than required"
	case "alphanum":
		return "opps! this field should be alphanumeric"
	case "alpha":
		return "opps! this field should be alphabetic"
	case "numeric":
		return "opps! this field should be numeric"
	}
	return ""
}
