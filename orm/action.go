// action is a sql generater after orm, usally to execut sql

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// Action sql executor implement
type Action struct {
	db       *sql.DB
	sql      string
	args     []interface{}
	limit    int64
	orderby  []string
	groupby  []string
	offset   int64
	isOffset bool
	filter   *Filter
	err      error
}

var _ ActionI = &Action{}

// ActionI sql executor interface
type ActionI interface {
	Get(dest interface{}) error
	One(dest interface{}) error
	Do() (int64, int64, error)
	Where(f ...*Filter) ActionI
	OrderBy(o ...string) ActionI
	GroupBy(o ...string) ActionI
	Limit(l int64) ActionI
	Offset(o int64) ActionI
	Page(page, psize int64) ActionI
	Sql() (string, []interface{}, error)
}

func (a *Action) finish() error {
	if a.err != nil {
		return a.err
	}
	if a.filter != nil {
		a.sql += " where " + a.filter.Where
		a.args = append(a.args, a.filter.Args...)
	}
	if len(a.groupby) > 0 {
		a.sql += fmt.Sprintf(" group by %s", strings.Join(a.groupby, ", "))
	}
	if len(a.orderby) > 0 {
		a.sql += fmt.Sprintf(" order by %s", strings.Join(a.orderby, ", "))
	}
	if a.limit > 0 {
		a.sql += " limit ?"
		a.args = append(a.args, a.limit)
	}
	if a.offset > 0 {
		a.sql += " offset ?"
		a.args = append(a.args, a.offset)
	}
	log.Printf("sql: %s, %v", a.sql, a.args)
	return nil
}

// Sql return the sql info for testing
func (a *Action) Sql() (string, []interface{}, error) {
	err := a.finish()
	if err != nil {
		return "", nil, err
	}
	return a.sql, a.args, nil
}

// Do executing sql
func (a *Action) Do() (int64, int64, error) {
	err := a.finish()
	if err != nil {
		return 0, 0, err
	}
	res, err := a.db.Exec(a.sql, a.args...)
	if err != nil {
		return 0, 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	return id, count, nil
}

// Get executing sql and fetch the data and restore to dest
func (a *Action) Get(dest interface{}) error {
	err := a.finish()
	if err != nil {
		return err
	}
	rows, err := a.db.Query(a.sql, a.args...)
	if err != nil {
		return err
	}
	err = scanQueryRows(dest, rows)
	return err
}

// One executing sql fetch one data and restore to dest
func (a *Action) One(dest interface{}) error {
	err := a.finish()
	if err != nil {
		return err
	}
	a.limit = 1
	rows, err := a.db.Query(a.sql, a.args...)
	if err != nil {
		return err
	}
	err = scanQueryOne(dest, rows)
	return err
}

// Where generate where condition
func (a *Action) Where(f ...*Filter) ActionI {
	if len(f) == 0 {
		a.err = errors.New("where can not be empty")
		return a
	}
	filter := And(f...)
	a.filter = filter
	return a
}

// OrderBy set sql order by
func (a *Action) OrderBy(o ...string) ActionI {
	if len(o) == 0 {
		a.err = errors.New("order by empty")
		return a
	}
	a.orderby = append(a.orderby, strings.Join(o, " "))
	return a
}

// GroupBy set sql group by
func (a *Action) GroupBy(o ...string) ActionI {
	if len(o) == 0 {
		a.err = errors.New("group by empty")
		return a
	}
	a.groupby = append(a.groupby, strings.Join(o, " "))
	return a
}

// Limit set sql limit
func (a *Action) Limit(l int64) ActionI {
	a.limit = l
	return a
}

// Offset set sql offset
func (a *Action) Offset(o int64) ActionI {
	a.offset = o
	return a
}

// Page set sql limit and offset
func (a *Action) Page(page, psize int64) ActionI {
	if page < 1 {
		page = 1
	}
	a.limit = psize
	a.offset = (page - 1) * psize
	return a
}
