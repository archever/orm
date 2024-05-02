package orm

import (
	"context"
	"database/sql"
)

type Session struct {
	db ExecutorIfc
}

func (s *Session) Table(schema Schema) *Action {
	return &Action{
		session: s,
		schema:  schema,
	}
}

func (s *Session) queryPayload(ctx context.Context, stmt *Stmt, payload PayloadIfc) error {
	payload.Bind()
	fields := payload.Fields()
	stmt.selectField = fields
	expr, err := stmt.completeSelect()
	if err != nil {
		return err
	}
	sqlRaw, argsRaw := expr.Expr()
	rows, err := s.db.QueryContext(ctx, sqlRaw, argsRaw...)
	if err != nil {
		return err
	}
	return ScanQueryFields(fields, rows)
}

func (s *Session) exec(ctx context.Context, stmt *Stmt) (sql.Result, error) {
	expr, err := stmt.complete()
	if err != nil {
		return nil, err
	}
	sqlRaw, argsRaw := expr.Expr()
	return s.db.ExecContext(ctx, sqlRaw, argsRaw...)
}
