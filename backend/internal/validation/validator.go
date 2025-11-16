package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()

	// Настраиваем валидацию под URL
	validate.RegisterValidation("url", func(fl validator.FieldLevel) bool {
		url := fl.Field().String()
		matched, _ := regexp.MatchString(`^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(/\S*)?$`, url)
		return matched
	})
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	if validate == nil {
		Init()
	}

	err := validate.Struct(s)
	if err != nil {
		var errors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getErrorMessage(err),
			})
		}
		return errors
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "url":
		return "Invalid URL format"
	default:
		return "Invalid value"
	}
}
