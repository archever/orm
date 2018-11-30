package orm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/archever/orm/f"
	"github.com/archever/orm/logger"
)

type stmt struct {
	db       *sql.DB
	tx       *sql.Tx
	sql      string
	table    string
	args     []interface{}
	limit    int
	orderby  []string
	groupby  []string
	offset   int
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

func (a *stmt) finish() (string, []interface{}, error) {
	if a.err != nil {
		return "", nil, a.err
	}
	sqls := a.sql
	rawArgs := a.args[:]
	args := []interface{}{}
	if len(a.filters) != 0 {
		filter := f.And(a.filters...)
		sqls += " where " + filter.Where
		rawArgs = append(rawArgs, filter.Args...)
	}
	if len(a.groupby) > 0 {
		sqls += fmt.Sprintf(" group by %s", strings.Join(a.groupby, ", "))
	}
	if len(a.orderby) > 0 {
		sqls += fmt.Sprintf(" order by %s", strings.Join(a.orderby, ", "))
	}
	if a.limit > 0 {
		sqls += " limit ?"
		rawArgs = append(rawArgs, a.limit)
	}
	if a.offset > 0 {
		sqls += " offset ?"
		rawArgs = append(rawArgs, a.offset)
	}
	for _, i := range rawArgs {
		m := ITOMarshaler(i)
		if m != nil {
			data, err := m.MarshalSQL()
			if err == ErrNull {
				args = append(args, nil)
			} else if err != nil {
				return "", nil, err
			} else {
				args = append(args, data)
			}
		} else {
			args = append(args, i)
		}
	}
	if Echo {
		logger.Info.Printf("sql: %s, %v", sqls, args)
	}
	return sqls, args, nil
}

// SQL return the sql info for testing
func (a *stmt) SQL() (string, []interface{}, error) {
	sqls, args, err := a.finish()
	if err != nil {
		return "", nil, err
	}
	return sqls, args, nil
}

func (a *stmt) MustDo() (rowID, rowCount int64) {
	rowID, rowCount, err := a.Do()
	if err != nil {
		panic(err)
	}
	return rowID, rowCount
}

func (a *stmt) MustGet(dest interface{}) {
	err := a.Get(dest)
	if err != nil && err != ErrNotFund {
		panic(err)
	}
}

func (a *stmt) MustOne(dest interface{}) {
	err := a.One(dest)
	if err != nil && err != ErrNotFund {
		panic(err)
	}
}

// Do executing sql
func (a *stmt) Do() (rowID, rowCount int64, err error) {
	sqls, args, err := a.finish()
	if err != nil {
		return 0, 0, err
	}
	var res sql.Result
	if a.isTx() {
		res, err = a.tx.Exec(sqls, args...)
	} else {
		res, err = a.db.Exec(sqls, args...)
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

func (a *stmt) Count() (int64, error) {
	tmp := a.sql
	defer func() {
		a.sql = tmp
	}()
	a.sql = "select count(*) as cnt from " + a.table
	sqls, args, err := a.finish()
	if err != nil {
		return 0, err
	}
	dest := f.M{}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.Query(sqls, args...)
	} else {
		rows, err = a.db.Query(sqls, args...)
	}
	if err != nil {
		return 0, err
	}
	err = ScanQueryOne(&dest, rows)
	if err == ErrNotFund {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return dest["cnt"].(int64), nil
}

// Get executing sql and fetch the data and restore to dest
func (a *stmt) Get(dest interface{}) error {
	sqls, args, err := a.finish()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.Query(sqls, args...)
	} else {
		rows, err = a.db.Query(sqls, args...)
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
	sqls, args, err := a.finish()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.Query(sqls, args...)
	} else {
		rows, err = a.db.Query(sqls, args...)
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
	a.groupby = append(a.groupby, strings.Join(o, ", "))
	return a
}

// Limit set sql limit
func (a *stmt) Limit(l int) *stmt {
	a.limit = l
	return a
}

// Offset set sql offset
func (a *stmt) Offset(o int) *stmt {
	a.offset = o
	return a
}

// Page set sql limit and offset
func (a *stmt) Page(page, psize int) *stmt {
	if page < 1 {
		page = 1
	}
	a.limit = psize
	a.offset = (page - 1) * psize
	return a
}
