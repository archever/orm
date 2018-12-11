package orm

import (
	"database/sql"
)

type Action struct {
	db    *sql.DB
	tx    *sql.Tx
	table string
}

func (o *Action) errStmt(err error) *Stmt {
	return &Stmt{
		err: err,
	}
}

func (o *Action) passStmt(sql string, args ...interface{}) *Stmt {
	return &Stmt{
		db:    o.db,
		tx:    o.tx,
		sql:   sql,
		args:  args,
		table: o.table,
	}
}

func (o *Action) Select(field ...string) *Stmt {
	sql := sqlSelect(o.table, true, field...)
	return o.passStmt(sql)
}

func (o *Action) SelectS(field ...string) *Stmt {
	sql := sqlSelect(o.table, false, field...)
	return o.passStmt(sql)
}

func (o *Action) Update(data M) *Stmt {
	sql, args, err := sqlUpdate(o.table, data)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *Action) Insert(row interface{}) *Stmt {
	sql, args, err := sqlInsert(o.table, row)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *Action) InsertMany(rows interface{}) *Stmt {
	sql, args, err := sqlInsertMany(o.table, rows)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *Action) Replace(row interface{}) *Stmt {
	sql, args, err := sqlReplace(o.table, row)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *Action) ReplaceMany(rows interface{}) *Stmt {
	sql, args, err := sqlReplaceMany(o.table, rows)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql, args...)
}

func (o *Action) Delete() *Stmt {
	sql, err := sqlDelete(o.table)
	if err != nil {
		return o.errStmt(err)
	}
	return o.passStmt(sql)
}
