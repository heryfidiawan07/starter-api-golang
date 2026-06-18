package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(s interface{}) []ValidationError {
	var errs []ValidationError
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: formatMessage(e),
			})
		}
	}
	return errs
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", strings.ToLower(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", strings.ToLower(e.Field()), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", strings.ToLower(e.Field()), e.Param())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", strings.ToLower(e.Field()))
	default:
		return fmt.Sprintf("%s is invalid", strings.ToLower(e.Field()))
	}
}
