package validation

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
	
	// Custom validation for username
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		
		// Check length
		if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 20 {
			return false
		}
		
		// Check allowed characters (letters, numbers, underscore, hyphen)
		matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
		return matched
	})
	
	// Custom validation for strong password
	validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		
		if len(password) < 8 {
			return false
		}
		
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[ \]{};':"\\|,.<>\/?]`).MatchString(password)
		
		return hasUpper && hasLower && hasNumber && hasSpecial
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
	case "email":
		return "Invalid email format"
	case "username":
		return "Username must be 3-20 characters and contain only letters, numbers, underscore, and hyphen"
	case "strong_password":
		return "Password must be at least 8 characters with uppercase, lowercase, number, and special character"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}