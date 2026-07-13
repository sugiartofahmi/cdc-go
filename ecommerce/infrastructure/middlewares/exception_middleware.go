package middlewares

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"

	"go-service/infrastructure/exceptions"
	"go-service/infrastructure/utils"
)

var (
	validationTranslator  ut.Translator
	structPathPartPattern = regexp.MustCompile(`^([A-Za-z0-9_]+)(.*)$`)
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			if name == "" {
				return field.Name
			}
			return name
		})

		english := en.New()
		translator := ut.New(english, english)
		trans, _ := translator.GetTranslator("en")
		_ = enTranslations.RegisterDefaultTranslations(v, trans)
		validationTranslator = trans
	}
}

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handlePanic(c)
		c.Next()
	}
}

func handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		panicException := createPanicException(err)
		errorMessage := panicException.ErrorMessage
		if errorMessage == "" {
			errorMessage = getErrorMessageByStatusCode(panicException.StatusCode)
		}

		var errors *[]utils.ErrorItemResponse
		if panicException.ValidationErrors != nil {
			responseErrors := make([]utils.ErrorItemResponse, len(panicException.ValidationErrors))
			for i, validationError := range panicException.ValidationErrors {
				responseErrors[i] = utils.ErrorItemResponse{
					Key:   validationError.Key,
					Value: validationError.Value,
				}
			}
			errors = &responseErrors
		}

		response := utils.ErrorResponse(panicException.StatusCode, errorMessage, errors)
		c.JSON(panicException.StatusCode, response)
		c.Abort()
	}
}

func createPanicException(err interface{}) exceptions.Exception {
	if ex, ok := err.(*exceptions.Exception); ok {
		return *ex
	}

	if ex, ok := err.(exceptions.Exception); ok {
		return ex
	}

	if bindingErr, ok := err.(gin.Error); ok {
		if bindingErr.Type == gin.ErrorTypeBind {
			validationErrors := parseBindingErrors(bindingErr.Err)
			return exceptions.Exception{
				ErrorMessage:     "Validation failed",
				StatusCode:       http.StatusBadRequest,
				ValidationErrors: validationErrors,
			}
		}
	}

	errMsg := "Internal Server Error"
	if e, ok := err.(error); ok {
		errMsg = e.Error()
	} else if s, ok := err.(string); ok {
		errMsg = s
	}

	return exceptions.Exception{
		ErrorMessage: errMsg,
		StatusCode:   http.StatusInternalServerError,
	}
}

func parseBindingErrors(err error) []exceptions.ValidationErrorItem {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errors := make([]exceptions.ValidationErrorItem, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			key := validationFieldKey(fieldErr)
			errors = append(errors, exceptions.ValidationErrorItem{
				Key:   key,
				Value: validationFieldMessage(key, fieldErr),
			})
		}
		return errors
	}

	return []exceptions.ValidationErrorItem{
		{
			Key:   "body",
			Value: "The request body is invalid JSON.",
		},
	}
}

func validationFieldKey(fieldErr validator.FieldError) string {
	namespace := fieldErr.StructNamespace()
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}

	keys := make([]string, 0, len(parts))
	for _, part := range parts {
		keys = append(keys, structPathPartToJSON(part))
	}
	return strings.Join(keys, ".")
}

func structPathPartToJSON(part string) string {
	matches := structPathPartPattern.FindStringSubmatch(part)
	if len(matches) != 3 {
		return toSnakeCase(part)
	}
	return toSnakeCase(matches[1]) + matches[2]
}

func toSnakeCase(value string) string {
	var builder strings.Builder
	for i, r := range value {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(r))
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

func validationFieldMessage(key string, fieldErr validator.FieldError) string {
	if validationTranslator != nil {
		message := fieldErr.Translate(validationTranslator)
		message = strings.Replace(message, fieldErr.Field(), key, 1)
		return message
	}
	return "The " + key + " field is invalid."
}

func getErrorMessageByStatusCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case http.StatusServiceUnavailable:
		return "Service Unavailable"
	default:
		return "Internal Server Error"
	}
}
