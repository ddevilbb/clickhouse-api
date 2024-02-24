package helper

import (
	"fmt"
	"strings"
)

func PrepareSqlWhere(params map[string][]string) (string, error) {
	where := ""
	operation := "="
	if len(params) > 0 {
		for key, values := range params {
			value := fmt.Sprintf("'%s'", strings.Join(values, "', '"))
			if len(value) == 0 {
				continue
			}
			if where != "" {
				where += " AND "
			} else {
				where += "WHERE "
			}
			operation = "="
			if len(values) > 1 {
				operation = "IN"
				value = fmt.Sprintf("(%s)", value)
			}
			where += fmt.Sprintf("%s %s %s", key, operation, value)
		}
	}

	return where, nil
}

func PrepareSqlOrder(params map[string]string) (string, error) {
	orderBy := ""
	if len(params) > 0 {
		var orderValue string
		for key, value := range params {
			orderValue = fmt.Sprintf("`%s` %s", key, value)
			if orderBy != "" {
				orderBy += ", "
			} else {
				orderBy += "ORDER BY "
			}
			orderBy += orderValue
		}
	} else {
		orderBy = "ORDER BY created_at DESC"
	}

	return orderBy, nil
}
