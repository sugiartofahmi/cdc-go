package exceptions

import "net/http"

type Exception struct {
	StatusCode       int
	ErrorMessage     string
	ValidationErrors []ValidationErrorItem
}

type ValidationErrorItem struct {
	Key   string
	Value string
}

func BadRequestException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusBadRequest,
	}
}

func ValidationException(errors map[string]string) *Exception {
	validationErrors := make([]ValidationErrorItem, 0, len(errors))
	for key, value := range errors {
		validationErrors = append(validationErrors, ValidationErrorItem{
			Key:   key,
			Value: value,
		})
	}

	return &Exception{
		ErrorMessage:     "Validation failed",
		StatusCode:       http.StatusBadRequest,
		ValidationErrors: validationErrors,
	}
}

func UnauthenticatedException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusUnauthorized,
	}
}

func ForbiddenException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusForbidden,
	}
}

func NotFoundException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusNotFound,
	}
}

func ConflictException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusConflict,
	}
}

func UnprocessableEntityException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusUnprocessableEntity,
	}
}

func ServerErrorException(err error) *Exception {
	return &Exception{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
}

func ServiceUnavailableException(message string) *Exception {
	return &Exception{
		ErrorMessage: message,
		StatusCode:   http.StatusServiceUnavailable,
	}
}
