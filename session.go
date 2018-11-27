package orm

import (
	"database/sql"
)

// 是否打印sql日志 默认是
var Echo = true

func SetEcho(echo bool) {
	Echo = echo
}

type Session struct {
	*sql.DB
}

type TxSession struct {
	DB *sql.DB
	TX *sql.Tx
}

func Open(driverName, dataSourceName string) (*Session, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Session{DB: db}, nil
}

func (s *Session) Begin() (*TxSession, error) {
	var err error
	ret := new(TxSession)
	ret.DB = s.DB
	ret.TX, err = s.DB.Begin()
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (s *Session) MustBegin() *TxSession {
	var err error
	ret := new(TxSession)
	ret.TX, err = s.DB.Begin()
	if err != nil {
		panic(err)
	}
	return ret
}

func (s *Session) Exec(sql string, arg ...interface{}) *stmt {
	ret := new(stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.DB
	ret.tx = nil
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *Session) Table(t string) *action {
	ret := new(action)
	ret.db = s.DB
	ret.tx = nil
	ret.table = t
	return ret
}

func (s *TxSession) Exec(sql string, arg ...interface{}) *stmt {
	ret := new(stmt)
	sql, args := sqlExec(sql, arg...)
	ret.db = s.DB
	ret.tx = s.TX
	ret.sql = sql
	ret.args = args
	return ret
}

func (s *TxSession) Table(t string) *action {
	ret := new(action)
	ret.db = s.DB
	ret.tx = s.TX
	ret.table = t
	return ret
}

func (s *TxSession) Commit() error {
	return s.TX.Commit()
}

func (s *TxSession) RollBack() error {
	return s.TX.Rollback()
}
