package database

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"slices"
	"strings"
)

func AddCoalesceToCompoundFields(str []string, compoundFields []string) []string {
	var newStr = slices.Clone(str)
	for i, s := range newStr {
		for _, field := range compoundFields {
			if strings.Contains(s, field) {
				newStr[i] = strings.Replace(s, field, "COALESCE("+field+",0)", 1)
			}
		}
	}
	return newStr
}

func Paginate(tx *gorm.DB, currentPage *int, pageSize *int) *gorm.DB {
	if currentPage == nil || pageSize == nil {
		return tx
	}
	return tx.Offset((*currentPage - 1) * *pageSize).Limit(*pageSize)
}

func OrderBy(tx *gorm.DB, orderBy []string) *gorm.DB {
	for _, o := range orderBy {
		tx = tx.Order(o)
	}
	return tx
}

func SearchBy(tx *gorm.DB, searchBy []string) *gorm.DB {
	for _, s := range searchBy {
		condition, search := splitSearch(s)
		if reflect.TypeOf(search).String() == "string" && search == "" {
			tx.Where(condition)
		} else {
			tx.Where(condition, search)
		}
	}
	return tx
}

func splitSearch(str string) (string, any) {
	split := strings.SplitN(str, " ", 3)

	property := split[0]
	operator := split[1]
	searchValue := split[2]

	condition := fmt.Sprintf("(%s %s (?))", property, operator)

	if strings.ToLower(operator) == "is" {
		return str, ""
	}
	if strings.ToLower(operator) == "in" {
		var values []string
		for _, val := range strings.Split(searchValue, ",") {
			values = append(values, strings.TrimSpace(val))
		}
		return condition, values
	}

	return condition, searchValue
}
