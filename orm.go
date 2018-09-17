// orm test init

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Orm entrance
type Orm struct {
	db    *sql.DB
	table string
}

var _ OrmI = &Orm{}

type Tabler interface {
	Table() string
}

// OrmI interface to build sql
type OrmI interface {
	Table(t string) OrmI
	Select(field ...string) ActionI
	Update(data map[string]interface{}) ActionI
	Insert(data ...interface{}) ActionI
	Replace(data ...interface{}) ActionI
	Delete() ActionI
	Exec(sql string, arg ...interface{}) ActionI
}

// New init a orm
func New(db *sql.DB) OrmI {
	return &Orm{
		db: db,
	}
}

// Reset the orm instance
func (o *Orm) Reset() {
	o.table = ""
}

func (o *Orm) errAction(err error) ActionI {
	defer o.Reset()
	return &Action{
		err: err,
	}
}

func (o *Orm) passAction(sql string, args ...interface{}) ActionI {
	defer o.Reset()
	return &Action{
		db:   o.db,
		sql:  sql,
		args: args,
	}
}

// Table set the table to work on
func (o *Orm) Table(t string) OrmI {
	o.table = t
	return o
}

// Select execute query fields
func (o *Orm) Select(field ...string) ActionI {
	var verb string
	if len(field) == 0 {
		verb = "select *"
	} else {
		verb = "select " + strings.Join(field, ", ")
	}
	if o.table != "" {
		verb += " from " + o.table
	}
	return o.passAction(verb)
}

// Update execute update sql
func (o *Orm) Update(data map[string]interface{}) ActionI {
	var err error
	if len(data) == 0 {
		err = errors.New("No data set")
		return o.errAction(err)
	}
	if o.table == "" {
		err = errors.New("No table set")
		return o.errAction(err)
	}
	set := []string{}
	args := []interface{}{}
	for k, v := range data {
		set = append(set, k+"=?")
		args = append(args, v)
	}
	verb := "update " + o.table + " set " + strings.Join(set, ", ")
	return o.passAction(verb, args...)
}

// Insert TODO: suport struct and multi
func (o *Orm) _insert(action string, data ...interface{}) ActionI {
	var err error
	if len(data) == 0 {
		err = errors.New("No data set")
		return o.errAction(err)
	}
	if o.table == "" {
		err = errors.New("No table set")
		return o.errAction(err)
	}
	init := false
	keys := []string{}
	insertData := []string{}
	args := []interface{}{}
	for index := range data {
		item := data[index]
		value, err := IToMap(item)
		if err != nil {
			return o.errAction(err)
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
				return o.errAction(errors.New("insert sequence must be same type"))
			}
			argS = append(argS, "?")
			args = append(args, v)
		}
		insertData = append(insertData, fmt.Sprintf("(%s)", strings.Join(argS, ", ")))
	}

	verb := fmt.Sprintf("%s into %s(%s) values %s", action, o.table,
		strings.Join(keys, ", "),
		strings.Join(insertData, ", "),
	)
	return o.passAction(verb, args...)
}

// Insert insert data
func (o *Orm) Insert(data ...interface{}) ActionI {
	return o._insert("insert", data...)
}

// Replace replace data
func (o *Orm) Replace(data ...interface{}) ActionI {
	return o._insert("replace", data...)
}

// Delete execute delete sql
func (o *Orm) Delete() ActionI {
	verb := "delete"
	if o.table == "" {
		err := errors.New("No table set")
		return o.errAction(err)
	}
	verb += " form " + o.table
	return o.passAction(verb)
}

// Exec do exec sql
func (o *Orm) Exec(sql string, args ...interface{}) ActionI {
	return o.passAction(sql, args...)
}
