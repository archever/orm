package orm

import (
	"database/sql"
)

type Session struct {
	db *sql.DB
}

type TxSession struct {
	db *sql.DB
	tx *sql.Tx
}

func Open(driverName, dataSourceName string) (*Session, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Session{db: db}, nil
}

func (s *Session) Begin() (*TxSession, error) {
	var err error
	ret := new(TxSession)
	ret.db = s.db
	ret.tx, err = s.db.Begin()
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (s *Session) MustBegin() *TxSession {
	var err error
	ret := new(TxSession)
	ret.tx, err = s.db.Begin()
	if err != nil {
		panic(err)
	}
	return ret
}

func (s *Session) Exec(sql string, arg ...interface{}) *stmt {
	ret := new(stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.db
	ret.tx = nil
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *Session) Table(t string) *action {
	ret := new(action)
	ret.db = s.db
	ret.tx = nil
	ret.table = t
	return ret
}

func (s *TxSession) Exec(sql string, arg ...interface{}) *stmt {
	ret := new(stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.db
	ret.tx = s.tx
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *TxSession) Table(t string) *action {
	ret := new(action)
	ret.db = s.db
	ret.tx = s.tx
	ret.table = t
	return ret
}

func (s *TxSession) Commit() error {
	return s.tx.Commit()
}

func (s *TxSession) RollBack() error {
	return s.tx.Rollback()
}
