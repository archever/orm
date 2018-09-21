package orm

import (
	"errors"
	"fmt"
	"strings"
)

func _select(table string, field ...string) string {
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

// Update execute update sql
func update(table string, data map[string]interface{}) (string, []interface{}, error) {
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

// Insert TODO: suport struct and multi
func _insert(table string, action string, data ...interface{}) (string, []interface{}, error) {
	var err error
	var sql string
	var args []interface{}
	if len(data) == 0 {
		err = errors.New("No data set")
		return sql, args, err
	}
	if table == "" {
		err = errors.New("No table set")
		return sql, args, err
	}
	init := false
	keys := []string{}
	insertData := []string{}
	for index := range data {
		item := data[index]
		value, err := IToMap(item)
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
			v, ok := value[k]
			if !ok {
				err = errors.New("insert sequence must be same type")
				return sql, args, err
			}
			argS = append(argS, "?")
			args = append(args, v)
		}
		insertData = append(insertData, fmt.Sprintf("(%s)", strings.Join(argS, ", ")))
	}

	sql = fmt.Sprintf("%s into %s(%s) values %s", action, table,
		strings.Join(keys, ", "),
		strings.Join(insertData, ", "),
	)
	return sql, args, err
}

// Insert insert data
func insert(table string, data ...interface{}) (string, []interface{}, error) {
	return _insert(table, "insert", data...)
}

// Replace replace data
func replace(table string, data ...interface{}) (string, []interface{}, error) {
	return _insert(table, "replace", data...)
}

// Delete execute delete sql
func delete(table string) (string, error) {
	var err error
	sql := "delete"
	if table == "" {
		err = errors.New("No table set")
		return sql, err
	}
	sql += " form " + table
	return sql, err
}

// Exec do exec sql
func exec(sql string, args ...interface{}) (string, []interface{}) {
	return sql, args
}
