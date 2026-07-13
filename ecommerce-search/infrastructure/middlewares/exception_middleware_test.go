package middlewares

import (
	"errors"
	"testing"

	"go-service/infrastructure/exceptions"
)

func TestParseBindingErrors_InvalidJSON(t *testing.T) {
	result := parseBindingErrors(errors.New("unexpected EOF"))
	assertValidationError(t, result, "body", "The request body is invalid JSON.")
}

func assertValidationError(t *testing.T, errors []exceptions.ValidationErrorItem, key string, value string) {
	t.Helper()

	for _, errItem := range errors {
		if errItem.Key == key && errItem.Value == value {
			return
		}
	}

	t.Fatalf("expected validation error %s=%q, got %#v", key, value, errors)
}
