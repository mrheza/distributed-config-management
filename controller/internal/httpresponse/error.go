package httpresponse

import (
	"controller/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ValidationFieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Fields  []ValidationFieldError `json:"fields,omitempty"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ValidationErrorResponse struct {
	Error ValidationErrorDetail `json:"error"`
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func InternalServerError(c *gin.Context, err error) {
	log.Printf("internal error: %v", err)
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func Unauthorized(c *gin.Context) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "resource not found"
	}
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

func FromError(c *gin.Context, err error) {
	if errors.Is(err, repository.ErrConfigNotFound) {
		NotFound(c, "config not found")
		return
	}
	InternalServerError(c, err)
}

func ValidationError(c *gin.Context, err error, requestStruct interface{}) {
	fields := extractValidationFields(err, requestStruct)
	c.JSON(http.StatusBadRequest, ValidationErrorResponse{
		Error: ValidationErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
			Fields:  fields,
		},
	})
}

func extractValidationFields(err error, requestStruct interface{}) []ValidationFieldError {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		fields := make([]ValidationFieldError, 0, len(ve))
		for _, fe := range ve {
			jsonName := getJSONFieldName(requestStruct, fe.Field())
			fields = append(fields, ValidationFieldError{
				Field:   jsonName,
				Message: validationMessage(fe),
			})
		}
		return fields
	}

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		field := ute.Field
		if field == "" {
			field = "body"
		} else {
			field = getJSONFieldName(requestStruct, field)
		}
		return []ValidationFieldError{
			{
				Field:   field,
				Message: fmt.Sprintf("must be %s", ute.Type.String()),
			},
		}
	}

	var se *json.SyntaxError
	if errors.As(err, &se) {
		return []ValidationFieldError{
			{
				Field:   "body",
				Message: "invalid JSON format",
			},
		}
	}

	return nil
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "url":
		return "must be a valid URL"
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	default:
		return "is invalid"
	}
}

func getJSONFieldName(requestStruct interface{}, structFieldName string) string {
	t := reflect.TypeOf(requestStruct)

	// handle pointer
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, ok := t.FieldByName(structFieldName)

	if !ok {
		return toSnakeCase(structFieldName)
	}

	jsonTag := field.Tag.Get("json")

	if jsonTag == "" {
		return toSnakeCase(structFieldName)
	}

	name := strings.Split(jsonTag, ",")[0]

	if name == "" {
		return toSnakeCase(structFieldName)
	}

	return name
}

func toSnakeCase(s string) string {
	var result []rune

	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}

		result = append(result, r)
	}

	return strings.ToLower(string(result))
}
