package orm

import (
	"context"
	"database/sql"
)

type Session struct {
	db ExecutorIfc
}

// func (s *Session) Exec(sql string, arg ...interface{}) *Stmt {
// 	ret := new(Stmt)
// 	sql, args := sqlExec(sql, arg...)
// 	ret.db = s.DB
// 	ret.sql = sql
// 	ret.args = args
// 	return ret
// }

func (s *Session) Table(schema Schema) *Action {
	return &Action{
		session: s,
		schema:  schema,
	}
}

func (s *Session) queryPayload(ctx context.Context, stmt *Stmt, payload PayloadIfc) error {
	payload.Bind()
	fields := payload.Fields()
	stmt.selectExpr.fields = fields
	sqlRaw, argsRaw, err := stmt.SQL()
	if err != nil {
		return err
	}
	rows, err := s.db.QueryContext(ctx, sqlRaw, argsRaw...)
	if err != nil {
		return err
	}
	return ScanQueryFields(fields, rows)
}

func (s *Session) exec(ctx context.Context, stmt *Stmt) (sql.Result, error) {
	sqlRaw, argsRaw, err := stmt.SQL()
	if err != nil {
		return nil, err
	}
	return s.db.ExecContext(ctx, sqlRaw, argsRaw...)
}
