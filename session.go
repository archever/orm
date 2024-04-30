package orm

import (
	"context"
	"database/sql"
	"fmt"
)

type Session struct {
	*sql.DB
	Driver string
}

type TxSession struct {
	s  *Session
	Tx *sql.Tx
}

func Open(driverName, dataSourceName string) (*Session, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Session{DB: db, Driver: driverName}, nil
}

func (s *Session) Begin() (*TxSession, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &TxSession{
		s:  s,
		Tx: tx,
	}, nil
}

func (s *Session) MustBegin() *TxSession {
	ret, err := s.Begin()
	if err != nil {
		panic(err)
	}
	return ret
}

func (s *Session) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TxSession, error) {
	tx, err := s.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TxSession{
		s:  s,
		Tx: tx,
	}, nil
}

func (s *Session) MustBeginTx(ctx context.Context, opts *sql.TxOptions) *TxSession {
	ret, err := s.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return ret
}

func (s *Session) Exec(sql string, arg ...interface{}) *Stmt {
	ret := new(Stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.DB
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *Session) Table(schema Schema) *Action {
	ret := new(Action)
	ret.db = s.DB
	// ret.table = fmt.Sprintf(t, arg...)
	return ret
}

func (s *TxSession) Exec(sql string, arg ...interface{}) *Stmt {
	ret := new(Stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.s.DB
	ret.tx = s.Tx
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *TxSession) Table(t string, arg ...interface{}) *Action {
	ret := new(Action)
	ret.db = s.s.DB
	ret.tx = s.Tx
	ret.table = fmt.Sprintf(t, arg...)
	return ret
}

func (s *TxSession) Commit() error {
	return s.Tx.Commit()
}

func (s *TxSession) RollBack() error {
	return s.Tx.Rollback()
}
