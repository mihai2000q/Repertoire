package validation

import (
	"regexp"
	"repertoire/server/internal/enums"
	"slices"
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

func IsDifficultyEnum(fl validator.FieldLevel) bool {
	difficulties := []enums.Difficulty{enums.Easy, enums.Medium, enums.Hard, enums.Impossible}

	difficulty, ok := fl.Field().Interface().(enums.Difficulty)
	if !ok {
		return false
	}
	return slices.Contains(difficulties, difficulty)
}

func IsYoutubeLink(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/(watch\?v=|embed/|v/|.+\?v=)?([^&=%\?]{11})`)
	return regex.MatchString(fl.Field().String())
}

func IsColor(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
	return regex.MatchString(fl.Field().String())
}
