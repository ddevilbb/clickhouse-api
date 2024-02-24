package helper

import (
	"fmt"
	"strings"
)

func FieldIsRequired(field string) string {
	return fmt.Sprintf("Field '%s' is required.", field)
}

func AnyOfFieldsIsRequired(fields []string) string {
	return fmt.Sprintf("Any of the fields:`%s`; is required", strings.Join(fields, "`, `"))
}

func InvalidFieldValue(field string) string {
	return fmt.Sprintf("Field '%s' value invalid", field)
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}
