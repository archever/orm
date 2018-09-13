package orm

import "fmt"

func Sum(field string) string {
	return fmt.Sprintf("sum(%s)", field)
}

func Distinct(field string) string {
	return fmt.Sprintf("distinct(%s)", field)
}

func Count(field ...string) string {
	fieldV := "*"
	if len(field) != 0 {
		fieldV = field[0]
	}
	return fmt.Sprintf("count(%s)", fieldV)

}

type M map[string]interface{}
