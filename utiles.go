package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type M map[string]interface{}

func sqlSelect(table string, field ...string) string {
	var sql string
	if len(field) == 0 {
		sql = "select *"
	} else {
		sql = "select " + strings.Join(field, ", ")
	}
	if table != "" {
		sql += " from " + table
	}
	return sql
}

func sqlUpdate(table string, data map[string]interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var set []string
	var args []interface{}
	if len(data) == 0 {
		err = errors.New("no data set")
		return sql, args, err
	}
	if table == "" {
		err = errors.New("No table set")
		return sql, args, err
	}
	for k, v := range data {
		set = append(set, k+"=?")
		args = append(args, v)
	}
	sql = "update " + table + " set " + strings.Join(set, ", ")
	return sql, args, err
}

func _sqlInsert(table string, action string, row interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var args []interface{}
	var keys []string
	var argS []string

	value, err := IToMap(reflect.ValueOf(row))
	if err != nil {
		return sql, args, err
	}
	for k := range value {
		keys = append(keys, k)
		args = append(args, value[k])
		argS = append(argS, "?")
	}
	sql = fmt.Sprintf("%s into %s(%s) values (%s)", action, table,
		strings.Join(keys, ", "),
		strings.Join(argS, ", "),
	)
	return sql, args, err
}

func _sqlInsertMany(table string, action string, rows interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var args []interface{}
	if table == "" {
		err = errors.New("No table set")
		return sql, args, err
	}
	init := false
	keys := []string{}
	insertData := []string{}

	rv := reflect.ValueOf(rows)
	switch rv.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			ele := rv.Index(i)
			value, err := IToMap(ele)
			if err != nil {
				return sql, args, err
			}
			// init keys
			if !init {
				for k := range value {
					keys = append(keys, k)
				}
				init = true
			}
			argS := []string{}
			for _, k := range keys {
				v, _ := value[k]
				argS = append(argS, "?")
				args = append(args, v)
			}
			insertData = append(insertData, fmt.Sprintf("(%s)", strings.Join(argS, ", ")))
		}
	default:
		err = errors.New("insertmany rows not iterable")
	}

	sql = fmt.Sprintf("%s into %s(%s) values %s", action, table,
		strings.Join(keys, ", "),
		strings.Join(insertData, ", "),
	)
	return sql, args, err
}

func sqlInsert(table string, row interface{}) (string, []interface{}, error) {
	return _sqlInsert(table, "insert", row)
}

func sqlInsertMany(table string, rows interface{}) (string, []interface{}, error) {
	return _sqlInsertMany(table, "insert", rows)
}

func sqlReplace(table string, row interface{}) (string, []interface{}, error) {
	return _sqlInsert(table, "replace", row)
}

func sqlReplaceMany(table string, rows interface{}) (string, []interface{}, error) {
	return _sqlInsertMany(table, "replace", rows)
}

func sqlDelete(table string) (string, error) {
	var err error
	sql := "delete"
	if table == "" {
		err = errors.New("No table set")
		return sql, err
	}
	sql += " form " + table
	return sql, err
}

func sqlExec(sql string, args ...interface{}) (string, []interface{}) {
	return sql, args
}
