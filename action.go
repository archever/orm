package orm

import (
	"database/sql"
)

type action struct {
	db    *sql.DB
	tx    *sql.Tx
	table string
}

func (o *action) errStmt(err error) *stmt {
	return &stmt{
		err: err,
	}
}
func (o *action) passStmt(sql string, args ...interface{}) *stmt {
	return &stmt{
		db:    o.db,
		tx:    o.tx,
		sql:   sql,
		args:  args,
		table: o.table,
	}
}

func (o *action) Select(field ...string) *stmt {
	sql := sqlSelect(o.table, true, field...)
	return o.passStmt(sql)
}

func (o *action) SelectS(field ...string) *stmt {
	sql := sqlSelect(o.table, false, field...)
	return o.passStmt(sql)
}

func (o *action) Update(data map[string]interface{}) *stmt {
	sql, args, err := sqlUpdate(o.table, data)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *action) Insert(row interface{}) *stmt {
	sql, args, err := sqlInsert(o.table, row)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *action) InsertMany(rows interface{}) *stmt {
	sql, args, err := sqlInsertMany(o.table, rows)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *action) Replace(row interface{}) *stmt {
	sql, args, err := sqlReplace(o.table, row)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *action) ReplaceMany(rows interface{}) *stmt {
	sql, args, err := sqlReplaceMany(o.table, rows)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *action) Delete() *stmt {
	sql, err := sqlDelete(o.table)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql)
}
