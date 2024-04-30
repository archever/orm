package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type M map[string]interface{}

func FieldWapper(field string) string {
	return fmt.Sprintf("%s%s%s", "`", field, "`")
}

type Stmt struct {
	db  *sql.DB
	tx  *sql.Tx
	ctx context.Context
	sql string

	action       string // select, update, delete, insert, create
	selectFields []FieldIfc
	args         []any
	table        Schema
	conds        []Cond
	orderby      []string
	groupby      []string
	limit        int
	offset       int
	err          error
}

func (a *Stmt) isTx() bool {
	if a.tx == nil {
		return false
	}
	return true
}

func (a *Stmt) complete() (string, []interface{}, error) {
	if a.err != nil {
		return "", nil, a.err
	}
	if a.ctx == nil {
		a.ctx = context.Background()
	}
	var sqls string
	switch a.action {
	case "select":
		cols := []string{}
		for _, field := range a.selectFields {
			cols = append(cols, FieldWapper(field.ColName()))
		}
		sqls = sqlSelect(a.table.TableName(), false, cols...)
	case "update":
		updateSql, args, err := sqlUpdate(a.table.TableName(), a.selectFields)
		if err != nil {
			return "", nil, err
		}
		sqls = updateSql
		a.args = append(a.args, args...)
	}
	rawArgs := a.args[:]
	var args []interface{}
	if len(a.conds) > 0 {
		var cond ExprIfc
		if len(a.conds) == 1 {
			cond = &a.conds[0]
		} else {
			cond = And(a.conds...)
		}
		sqls += " where " + cond.Expr()
		rawArgs = append(rawArgs, cond.Args()...)
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
		m, ok := ITOMarshaler(reflect.ValueOf(&i))
		if ok {
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
	Log.Printf("%s; %v", sqls, args)
	return sqls, args, nil
}

func (a *Stmt) Ctx(ctx context.Context) *Stmt {
	a.ctx = ctx
	return a
}

// SQL return the sql info for testing
func (a *Stmt) SQL() (string, []interface{}, error) {
	sqls, args, err := a.complete()
	if err != nil {
		return "", nil, err
	}
	return sqls, args, nil
}

func (a *Stmt) MustDo() (rowID, rowCount int64) {
	rowID, rowCount, err := a.Do()
	if err != nil {
		panic(err)
	}
	return rowID, rowCount
}

// Do executing sql
func (a *Stmt) Do() (rowID, rowCount int64, err error) {
	sqls, args, err := a.complete()
	if err != nil {
		return 0, 0, err
	}
	var res sql.Result
	if a.isTx() {
		res, err = a.tx.ExecContext(a.ctx, sqls, args...)
	} else {
		res, err = a.db.ExecContext(a.ctx, sqls, args...)
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

func (a *Stmt) Count() (int64, error) {
	tmp := a.sql
	defer func() {
		a.sql = tmp
	}()
	a.sql = "select count(*) as cnt from " + a.table.TableName()
	sqls, args, err := a.complete()
	if err != nil {
		return 0, err
	}
	dest := M{}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
	} else {
		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
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
func (a *Stmt) Get(dest PayloadIfc) error {
	dest.Bind()
	queryFields := dest.Fields()
	a.selectFields = queryFields
	a.limit = 1
	sqls, args, err := a.complete()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
	} else {
		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
	}
	if err != nil {
		return err
	}
	// err = ScanQueryRows(dest, rows)
	err = ScanQueryFields(queryFields, rows)
	return err
}

// One executing sql fetch one data and restore to dest
func (a *Stmt) One(dest interface{}) error {
	a.limit = 1
	sqls, args, err := a.complete()
	if err != nil {
		return err
	}
	var rows *sql.Rows
	if a.isTx() {
		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
	} else {
		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
	}
	if err != nil {
		return err
	}
	err = ScanQueryOne(dest, rows)
	return err
}

// Where generate where condition
func (a *Stmt) Where(cond ...Cond) *Stmt {
	a.conds = append(a.conds, cond...)
	return a
}

// OrderBy set sql order by
func (a *Stmt) OrderBy(field string, reverse bool) *Stmt {
	field = FieldWapper(field)
	if reverse {
		field += " desc"
	}
	a.orderby = append(a.orderby, field)
	return a
}

// GroupBy set sql group by
func (a *Stmt) GroupBy(o ...string) *Stmt {
	for i := range o {
		o[i] = FieldWapper(o[i])
	}
	a.groupby = append(a.groupby, strings.Join(o, ", "))
	return a
}

// Limit set sql limit
func (a *Stmt) Limit(l int) *Stmt {
	a.limit = l
	return a
}

// Offset set sql offset
func (a *Stmt) Offset(o int) *Stmt {
	a.offset = o
	return a
}

// Page set sql limit and offset
func (a *Stmt) Page(page, psize int) *Stmt {
	if page < 1 {
		page = 1
	}
	a.limit = psize
	a.offset = (page - 1) * psize
	return a
}
