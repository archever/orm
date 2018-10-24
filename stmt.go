package orm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/archever/orm/f"
)

type stmt struct {
	db       *sql.DB
	tx       *sql.Tx
	sql      string
	args     []interface{}
	limit    int64
	orderby  []string
	groupby  []string
	offset   int64
	isOffset bool
	filters  []*f.FilterItem
	err      error
}

func (a *stmt) isTx() bool {
	if a.tx == nil {
		return false
	}
	return true
}

func (a *stmt) finish() error {
	if a.err != nil {
		return a.err
	}
	if len(a.filters) != 0 {
		filter := f.And(a.filters...)
		a.sql += " where " + filter.Where
		a.args = append(a.args, filter.Args...)
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

// SQL return the sql info for testing
func (a *stmt) SQL() (string, []interface{}, error) {
	err := a.finish()
	if err != nil {
		return "", nil, err
	}
	return a.sql, a.args, nil
}

// Do executing sql
func (a *stmt) Do() (rowID, rowCount int64, err error) {
	err = a.finish()
	if err != nil {
		return 0, 0, err
	}
	var res sql.Result
	if a.isTx() {
		res, err = a.tx.Exec(a.sql, a.args...)
	} else {
		res, err = a.db.Exec(a.sql, a.args...)
	}
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
func (a *stmt) Get(dest interface{}) error {
	err := a.finish()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.Query(a.sql, a.args...)
	} else {
		rows, err = a.db.Query(a.sql, a.args...)
	}
	if err != nil {
		return err
	}
	err = ScanQueryRows(dest, rows)
	return err
}

// One executing sql fetch one data and restore to dest
func (a *stmt) One(dest interface{}) error {
	a.limit = 1
	err := a.finish()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.Query(a.sql, a.args...)
	} else {
		rows, err = a.db.Query(a.sql, a.args...)
	}
	if err != nil {
		return err
	}
	err = ScanQueryOne(dest, rows)
	return err
}

// Filter generate where condition
func (a *stmt) Filter(filters ...*f.FilterItem) *stmt {
	filter := f.And(filters...)
	a.filters = append(a.filters, filter)
	return a
}

// Where generate where condition
func (a *stmt) Where(cond string, arg ...interface{}) *stmt {
	filter := f.S(cond, arg...)
	a.filters = append(a.filters, filter)
	return a
}

// OrderBy set sql order by
func (a *stmt) OrderBy(field string, reverse bool) *stmt {
	if reverse {
		field += " desc"
	}
	a.orderby = append(a.orderby, field)
	return a
}

// GroupBy set sql group by
func (a *stmt) GroupBy(o ...string) *stmt {
	a.groupby = append(a.groupby, strings.Join(o, " "))
	return a
}

// Limit set sql limit
func (a *stmt) Limit(l int64) *stmt {
	a.limit = l
	return a
}

// Offset set sql offset
func (a *stmt) Offset(o int64) *stmt {
	a.offset = o
	return a
}

// Page set sql limit and offset
func (a *stmt) Page(page, psize int64) *stmt {
	if page < 1 {
		page = 1
	}
	a.limit = psize
	a.offset = (page - 1) * psize
	return a
}
