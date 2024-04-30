package orm

import "fmt"

func FieldWrapper(field string) string {
	return fmt.Sprintf("%s%s%s", "`", field, "`")
}
