package database

import (
	"gorm.io/gorm"
	"strings"
)

func AddCoalesceToCompoundFields(str *[]string, compoundFields []string) {
	for i, s := range *str {
		for _, field := range compoundFields {
			if strings.Contains(s, field) {
				(*str)[i] = strings.Replace(s, field, "COALESCE("+field+", 0)", 1)
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
		tx.Order(o)
	}
	return tx
}

func SearchBy(tx *gorm.DB, searchBy []string) *gorm.DB {
	for _, s := range searchBy {
		tx.Where(s)
		// TODO: SECURITY ZERO SO CHANGE
		/*property, search := splitSearch(s)
		if search == "" {
			tx = tx.Where(property)
		} else {
			tx = tx.Where(property, search)
		}*/
	}
	return tx
}

func splitSearch(str string) (string, string) {
	spaces := 0
	firstIndex := 0
	lastIndex := 0

	for i, r := range str {
		if r == ' ' {
			spaces = spaces + 1
		}
		if spaces == 1 && firstIndex == 0 {
			firstIndex = i
		}
		if spaces == 2 {
			lastIndex = i
			break
		}
	}

	// What if the condition is:
	// Property IS NULL
	// Property IS NOT NULL
	// Property = IS NULL
	symbol := str[firstIndex+1 : lastIndex-1]
	if (symbol != "=" && symbol != "<>" && symbol != "<" && symbol != ">" && symbol != "=<" && symbol != ">=") &&
		(strings.HasSuffix(str, "IS NULL") || strings.HasSuffix(str, "IS NOT NULL")) {
		return str, ""
	}

	return str[:lastIndex] + " ?", str[lastIndex+1:]
}
