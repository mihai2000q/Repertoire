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
	orderBy := fl.Field().Interface().([]string)
	for _, o := range orderBy {
		if validateOrderByElem(o) == false {
			return false
		}
	}
	return true
}

func SearchBy(fl validator.FieldLevel) bool {
	searchBy := fl.Field().Interface().([]string)
	for _, s := range searchBy {
		if validateSearchByElem(s) == false {
			return false
		}
	}
	return true
}

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
			validateOrderNullability(strings.ToLower(split[1]), strings.ToLower(split[2]))
	}
	return false
}

func validateOrderType(str string) bool {
	return str != "asc" && str != "desc"
}

func validateOrderNullability(str1 string, str2 string) bool {
	return str1 != "nulls" && str2 != "last" && str2 != "first"
}

var filterOperators = []string{"=", "!=", "<>", "<", ">", "<=", ">=", "is", "in"}

func validateSearchByElem(searchBy string) bool {
	startIndexOfOperator := 0
	startIndexOfSearchValue := 0

	spaces := 0
	for i, s := range searchBy {
		if s == ' ' {
			spaces++
		}
		if spaces == 1 && startIndexOfOperator == 0 {
			startIndexOfOperator = i
		}
		if spaces == 2 {
			startIndexOfSearchValue = i
			break
		}
	}

	if startIndexOfOperator == 0 || startIndexOfSearchValue == 0 || len(searchBy) == startIndexOfOperator+1 {
		return false
	}

	operator := strings.ToLower(searchBy[startIndexOfOperator+1 : startIndexOfSearchValue])
	searchValue := searchBy[startIndexOfSearchValue+1:]

	if operator == "is" {
		searchValue = strings.ToLower(searchValue)
		return searchValue == "null" || searchValue == "not null"
	}
	return slices.Contains(filterOperators, operator)
}
