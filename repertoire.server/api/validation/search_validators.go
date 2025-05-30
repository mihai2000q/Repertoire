package validation

import "strings"

// order

func validateSearchOrderElem(order string) bool {
	split := strings.Split(order, ":")
	return len(split) == 2 &&
		len(strings.Split(order, " ")) == 1 &&
		validateOrderType(strings.ToLower(split[1]))
}

// filter

func validateSearchFilterElem(filter string) bool {
	filter = strings.TrimSpace(filter)

	// Handle empty string
	if filter == "" {
		return false
	}

	// Remove optional outer parentheses
	if strings.HasPrefix(filter, "(") && strings.HasSuffix(filter, ")") {
		filter = strings.TrimSpace(filter[1 : len(filter)-1])
	}

	// Check for OR conditions
	if strings.Contains(filter, " OR ") {
		conditions := strings.Split(filter, " OR ")
		for _, cond := range conditions {
			if !validateAtomicCondition(strings.TrimSpace(cond)) {
				return false
			}
		}
		return true
	}

	// Single condition case
	return validateAtomicCondition(filter)
}

var searchOperators = []string{"=", "!=", ">", "<", ">=", "<=", "IS EMPTY", "EXISTS", "IS NULL", "IS NOT NULL", "IN"}

func validateAtomicCondition(condition string) bool {
	condition = strings.TrimSpace(condition)

	// Handle NOT prefix
	if strings.HasPrefix(condition, "NOT ") {
		condition = strings.TrimSpace(condition[4:])
	}

	// Handle quoted attributes
	var attr string
	// Unquoted attribute
	spaceIdx := strings.Index(condition, " ")
	if spaceIdx == -1 {
		return false
	}
	attr = condition[:spaceIdx]
	condition = strings.TrimSpace(condition[spaceIdx:])

	// Validate attribute
	if strings.ContainsAny(strings.Trim(attr, `"`), " \t") {
		return false
	}

	for _, op := range searchOperators {
		if strings.HasPrefix(condition, op) {
			val := strings.TrimSpace(condition[len(op):])

			switch op {
			case "IS NULL", "IS NOT NULL", "IS EMPTY", "EXISTS":
				return val == ""
			case "IN":
				return validateInValue(val)
			default:
				return validateSimpleValue(val)
			}
		}
	}
	return false
}

func validateInValue(val string) bool {
	return strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]")
}

func validateSimpleValue(val string) bool {
	isQuoted := (strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) ||
		(strings.HasPrefix(val, `'`) && strings.HasSuffix(val, `'`))
	if isQuoted {
		return true // Quoted strings may contain spaces
	}
	return !strings.ContainsAny(val, " \t") // Unquoted string cannot contain spaces
}
