package utils

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
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
		return fmt.Sprintf("the \\'%s\\' attribute should be a %s, but a %s value was given", field, expectedType, givenType)
	}
	return ""
}
