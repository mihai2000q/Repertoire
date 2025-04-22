package database

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func AddCoalesceToCompoundFields(str *[]string, compoundFields []string) {
	for i, s := range *str {
		for _, field := range compoundFields {
			if strings.Contains(s, field) {
				(*str)[i] = strings.Replace(s, field, "COALESCE("+field+",0)", 1)
			}
		}
	}
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
	startIndexOfOperator := 0
	startIndexOfSearchValue := 0

	spaces := 0
	for i, s := range str {
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

	property := str[:startIndexOfOperator]
	operator := str[startIndexOfOperator+1 : startIndexOfSearchValue]
	searchValue := str[startIndexOfSearchValue+1:]

	condition := fmt.Sprintf("(%s %s (?))", property, operator)

	if operator == "IS" {
		return str, ""
	}
	if operator == "IN" {
		var values []string
		for _, val := range strings.Split(searchValue, ",") {
			values = append(values, strings.TrimSpace(val))
		}
		return condition, values
	}

	return condition, searchValue
}
