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

type DefaultExecutor struct {
	*sql.DB
}

func NewDefaultExecutor(db *sql.DB) *DefaultExecutor {
	return &DefaultExecutor{
		DB: db,
	}
}

func (e *DefaultExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	Log.Printf("exec: %s, args: %v", query, args)
	return e.DB.ExecContext(ctx, query, args...)
}

func (e *DefaultExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	Log.Printf("query: %s, args: %v", query, args)
	return e.DB.QueryContext(ctx, query, args...)
}
