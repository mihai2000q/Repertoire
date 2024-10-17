package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func HasUpper(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func HasLower(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func HasDigit(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
