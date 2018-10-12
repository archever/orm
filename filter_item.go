// filter generate the where condition sql

package orm

import (
	"fmt"
	"reflect"
	"strings"
)

// FilterItem build where condition
type FilterItem struct {
	Where string
	Args  []interface{}
}

func filter(o, field string, arg interface{}) *FilterItem {
	return &FilterItem{
		Where: fmt.Sprintf("%s%s?", field, o),
		Args:  []interface{}{arg},
	}
}

func S(cond string, arg ...interface{}) *FilterItem {
	return &FilterItem{
		Where: cond,
		Args:  arg,
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

func nin(o, field string, arg interface{}) *FilterItem {
	argS := []string{}
	argV := []interface{}{}
	v := reflect.ValueOf(arg)
	for i := 0; i < v.Len(); i++ {
		argS = append(argS, "?")
		argV = append(argV, v.Index(i).Interface())
	}
	return &FilterItem{
		Where: fmt.Sprintf("%s %s (%s)", field, o, strings.Join(argS, ",")),
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
	whereS := []string{}
	args := []interface{}{}
	for _, i := range f {
		whereS = append(whereS, i.Where)
		args = append(args, i.Args...)
	}
	return &FilterItem{
		Where: strings.Join(whereS, " and "),
		Args:  args,
	}
}

func Or(left, right *FilterItem) *FilterItem {
	whereS := []string{}
	args := []interface{}{}
	for _, i := range [...]*FilterItem{
		left, right,
	} {
		whereS = append(whereS, i.Where)
		args = append(args, i.Args...)
	}
	return &FilterItem{
		Where: strings.Join(whereS, " or "),
		Args:  args,
	}
}
