package orm

import (
	"context"
	"database/sql"
)

var _ ExecutorIfc = (*sql.DB)(nil)
var _ ExecutorIfc = (*sql.Tx)(nil)
var _ TransactionIfc = (*sql.Tx)(nil)

type ExecutorIfc interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type TransactionIfc interface {
	Rollback() error
	Commit() error
}

// func (a *Stmt) MustDo() (rowID, rowCount int64) {
// 	rowID, rowCount, err := a.Do()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return rowID, rowCount
// }

// // Do executing sql
// func (a *Stmt) Do() (rowID, rowCount int64, err error) {
// 	sqls, args, err := a.complete()
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	var res sql.Result
// 	if a.isTx() {
// 		res, err = a.tx.ExecContext(a.ctx, sqls, args...)
// 	} else {
// 		res, err = a.db.ExecContext(a.ctx, sqls, args...)
// 	}
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	count, err := res.RowsAffected()
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	return id, count, nil
// }

// func (a *Stmt) Count() (int64, error) {
// 	tmp := a.sql
// 	defer func() {
// 		a.sql = tmp
// 	}()
// 	a.sql = "select count(*) as cnt from " + a.table.TableName()
// 	sqls, args, err := a.complete()
// 	if err != nil {
// 		return 0, err
// 	}
// 	dest := M{}
// 	var rows *sql.Rows
// 	if a.isTx() {
// 		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
// 	} else {
// 		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
// 	}
// 	if err != nil {
// 		return 0, err
// 	}
// 	err = ScanQueryOne(&dest, rows)
// 	if err == ErrNotFund {
// 		return 0, nil
// 	} else if err != nil {
// 		return 0, err
// 	}
// 	return dest["cnt"].(int64), nil
// }

// // Get executing sql and fetch the data and restore to dest
// func (a *Stmt) Get(dest PayloadIfc) error {
// 	dest.Bind()
// 	queryFields := dest.Fields()
// 	a.selectFields = queryFields
// 	a.limit = 1
// 	sqls, args, err := a.complete()
// 	if err != nil {
// 		return err
// 	}
// 	var rows *sql.Rows
// 	if a.isTx() {
// 		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
// 	} else {
// 		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	// err = ScanQueryRows(dest, rows)
// 	err = ScanQueryFields(queryFields, rows)
// 	return err
// }

// // One executing sql fetch one data and restore to dest
// func (a *Stmt) One(dest interface{}) error {
// 	a.limit = 1
// 	sqls, args, err := a.complete()
// 	if err != nil {
// 		return err
// 	}
// 	var rows *sql.Rows
// 	if a.isTx() {
// 		rows, err = a.tx.QueryContext(a.ctx, sqls, args...)
// 	} else {
// 		rows, err = a.db.QueryContext(a.ctx, sqls, args...)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	err = ScanQueryOne(dest, rows)
// 	return err
// }
