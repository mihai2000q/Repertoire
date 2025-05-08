package validation

import (
	"regexp"
	"repertoire/server/internal/enums"
	"slices"
	"strings"
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

func DifficultyEnum(fl validator.FieldLevel) bool {
	difficulties := []enums.Difficulty{enums.Easy, enums.Medium, enums.Hard, enums.Impossible}

	difficulty, ok := fl.Field().Interface().(enums.Difficulty)
	if !ok {
		return false
	}
	return slices.Contains(difficulties, difficulty)
}

func SearchTypeEnum(fl validator.FieldLevel) bool {
	searchTypes := []enums.SearchType{enums.Artist, enums.Album, enums.Song, enums.Playlist}

	searchType, ok := fl.Field().Interface().(enums.SearchType)
	if !ok {
		return false
	}
	return slices.Contains(searchTypes, searchType)
}

func YoutubeLink(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^(https?://)?(www\.)?(youtube\.com|youtu\.be)/(watch\?v=|embed/|v/|.+\?v=)?([^&=%\?]{11})`)
	return regex.MatchString(fl.Field().String())
}

func Color(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
	return regex.MatchString(fl.Field().String())
}

func OrderBy(fl validator.FieldLevel) bool {
	orderBy, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, o := range orderBy {
		if !validateOrderByElem(o) {
			return false
		}
	}
	return true
}

func SearchOrder(fl validator.FieldLevel) bool {
	order, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, o := range order {
		if !validateSearchOrderElem(o) {
			return false
		}
	}
	return true
}

func SearchBy(fl validator.FieldLevel) bool {
	searchBy, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, s := range searchBy {
		if !validateSearchByElem(s) {
			return false
		}
	}
	return true
}

// private functions

func validateOrderByElem(orderBy string) bool {
	split := strings.Split(orderBy, " ")
	if len(split) == 1 {
		return true
	}
	if len(split) == 2 {
		return validateOrderType(strings.ToLower(split[1]))
	}
	if len(split) == 3 {
		return validateOrderNullability(strings.ToLower(split[1]), strings.ToLower(split[2]))
	}
	if len(split) == 4 {
		return validateOrderType(strings.ToLower(split[1])) &&
			validateOrderNullability(strings.ToLower(split[2]), strings.ToLower(split[3]))
	}
	return false
}

func validateOrderType(str string) bool {
	return str == "asc" || str == "desc"
}

func validateOrderNullability(str1 string, str2 string) bool {
	return str1 == "nulls" && (str2 == "last" || str2 == "first")
}

func validateSearchOrderElem(order string) bool {
	split := strings.Split(order, ":")
	return len(split) == 2 &&
		len(strings.Split(order, " ")) == 1 &&
		validateOrderType(strings.ToLower(split[1]))
}

var filterOperators = []string{"=", "!=", "<>", "<", ">", "<=", ">=", "is", "in"}

func validateSearchByElem(searchBy string) bool {
	split := strings.SplitN(searchBy, " ", 3)
	if len(split) != 3 {
		return false
	}

	operator := strings.ToLower(split[1])
	searchValue := split[2]

	if operator == "is" {
		searchValue = strings.ToLower(searchValue)
		return searchValue == "null" || searchValue == "not null"
	}
	return slices.Contains(filterOperators, operator)
}
