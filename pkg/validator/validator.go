package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator instance (bisa digunakan ulang)
var validate = validator.New()

// ValidateStruct mengecek validitas struct berdasarkan tag
func ValidateStruct(data interface{}) map[string]string {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	}
	return errors
}
