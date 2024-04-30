package orm

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func _sqlInsert(table string, action string, row interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var args []interface{}
	var keyS []string
	var argS []string

	if table == "" {
		return "", nil, errors.New("table not set")
	}

	value, err := IToMap(reflect.ValueOf(row))
	if err != nil {
		return sql, args, err
	}
	for k := range value {
		keyS = append(keyS, k)
	}
	sort.Strings(keyS)
	for i, k := range keyS {
		keyS[i] = FieldWrapper(k)
		args = append(args, value[k])
		argS = append(argS, "?")
	}
	sql = fmt.Sprintf("%s into %s(%s) values (%s)", action, FieldWrapper(table),
		strings.Join(keyS, ", "),
		strings.Join(argS, ", "),
	)
	return sql, args, err
}

func _sqlInsertMany(table string, action string, rows interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var args []interface{}
	if table == "" {
		err = fmt.Errorf("table not set")
		return sql, args, err
	}
	init := false
	keys := []string{}
	checkkeys := map[string]bool{}
	wappedKeys := []string{}
	insertData := []string{}

	rv := reflect.ValueOf(rows)
	switch rv.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if rv.Len() == 0 {
			err = fmt.Errorf("save empty data")
			break
		}
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
					checkkeys[k] = true
				}
				sort.Strings(keys)
				for _, k := range keys {
					wappedKeys = append(wappedKeys, FieldWrapper(k))
				}
				init = true
			} else {
				if len(value) != len(keys) {
					return "", nil, errors.New("can not save many with different data field")
				}
				for k := range value {
					if _, ok := checkkeys[k]; !ok {
						return "", nil, errors.New("can not save many with different data field")
					}
				}
			}
			argS := []string{}
			for _, k := range keys {
				argS = append(argS, "?")
				args = append(args, value[k])
			}
			insertData = append(insertData, fmt.Sprintf("(%s)", strings.Join(argS, ", ")))
		}
	default:
		err = fmt.Errorf("invalid data type: %T", rows)
	}
	if err == nil {
		sql = fmt.Sprintf("%s into %s(%s) values %s", action, FieldWrapper(table),
			strings.Join(wappedKeys, ", "),
			strings.Join(insertData, ", "),
		)
	}
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
		err = fmt.Errorf("table not set")
		return sql, err
	}
	sql += " from " + FieldWrapper(table)
	return sql, err
}

func sqlExec(sql string, args ...interface{}) (string, []interface{}) {
	return sql, args
}
