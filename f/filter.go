package f

import (
	"fmt"
	"reflect"
	"strings"
)

type M map[string]interface{}

type FilterItem struct {
	Where string
	Args  []interface{}
}

func FieldWapper(field string) string {
	return fmt.Sprintf("%s%s%s", "`", field, "`")
}

func filter(o, field string, arg interface{}) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s%s?", FieldWapper(field), o),
		Args:  []interface{}{arg},
	}
}

func S(cond string, arg ...interface{}) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s", cond),
		Args:  arg,
	}
}

func IsNull(field string) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s is null", FieldWapper(field)),
		Args:  []interface{}{},
	}
}

func IsNotNull(field string) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s is not null", FieldWapper(field)),
		Args:  []interface{}{},
	}
}

func Equal(field string, arg interface{}) *FilterItem {
	return filter("=", field, arg)
}

func NotEqual(field string, arg interface{}) *FilterItem {
	return filter("!=", field, arg)
}

func Lte(field string, arg interface{}) *FilterItem {
	return filter("<=", field, arg)
}

func Lt(field string, arg interface{}) *FilterItem {
	return filter("<", field, arg)
}

func Gte(field string, arg interface{}) *FilterItem {
	return filter(">=", field, arg)
}

func Gt(field string, arg interface{}) *FilterItem {
	return filter(">", field, arg)
}

func Between(field string, left, right interface{}) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s between ? and ?", FieldWapper(field)),
		Args:  []interface{}{left, right},
	}
}

func nin(o, field string, arg interface{}) *FilterItem {
	argS := []string{}
	argV := []interface{}{}
	v := reflect.ValueOf(arg)
	for i := 0; i < v.Len(); i++ {
		argS = append(argS, "?")
		argV = append(argV, v.Index(i).Interface())
	}
	return &FilterItem{
		Where: fmt.Sprintf("%s %s (%s)", FieldWapper(field), o, strings.Join(argS, ",")),
		Args:  argV,
	}
}

func In(field string, arg interface{}) *FilterItem {
	return nin("in", field, arg)
}

func NotIn(field string, arg interface{}) *FilterItem {
	return nin("not in", field, arg)
}

func Like(field string, arg interface{}) *FilterItem {
	return filter(" like ", field, arg)
}

func NotLike(field string, arg interface{}) *FilterItem {
	return filter(" not like ", field, arg)
}

func And(f ...*FilterItem) *FilterItem {
	if len(f) == 1 {
		return f[0]
	}
	whereS := []string{}
	args := []interface{}{}
	for _, i := range f {
		whereS = append(whereS, i.Where)
		args = append(args, i.Args...)
	}
	return &FilterItem{
		Where: fmt.Sprintf("(%s)", strings.Join(whereS, " and ")),
		Args:  args,
	}
}

func Or(f ...*FilterItem) *FilterItem {
	if len(f) == 1 {
		return f[0]
	}
	whereS := []string{}
	args := []interface{}{}
	for _, i := range f {
		whereS = append(whereS, i.Where)
		args = append(args, i.Args...)
	}
	return &FilterItem{
		Where: fmt.Sprintf("(%s)", strings.Join(whereS, " or ")),
		Args:  args,
	}
}
