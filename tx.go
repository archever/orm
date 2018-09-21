package orm

import (
	"database/sql"
	"errors"
)

type OrmTx struct {
	db    *sql.DB
	tx    *sql.Tx
	table string
}

var _ TxI = &OrmTx{}
var _ TabledTxI = &OrmTx{}

func NewTx(db *sql.DB) (TxI, error) {
	return &OrmTx{
		db: db,
	}, nil
}

func toTabledTxI(i interface{}) TabledTxI {
	return i.(TabledTxI)
}

func (o *OrmTx) reset() {
	o.table = ""
}

func (o *OrmTx) errAction(err error) ActionTxI {
	defer o.reset()
	return &ActionTx{
		err: err,
	}
}

func (o *OrmTx) passAction(sql string, args ...interface{}) ActionTxI {
	defer o.reset()
	return &ActionTx{
		db:   o.db,
		tx:   o.tx,
		sql:  sql,
		args: args,
	}
}

func (o *OrmTx) Table(t string) TabledTxI {
	o.table = t
	return toTabledTxI(o)
}

func (o *OrmTx) Exec(sql string, arg ...interface{}) ActionTxI {
	return o.passAction(sql, arg...)
}

func (o *OrmTx) Begin() error {
	var err error
	if o.tx != nil {
		o.tx.Rollback()
	}
	o.tx, err = o.db.Begin()
	return err
}

func (o *OrmTx) Commit() error {
	if o.tx == nil {
		return errors.New("not begin tx")
	}
	return o.tx.Commit()
}

func (o *OrmTx) RollBack() error {
	if o.tx == nil {
		return errors.New("not begin tx")
	}
	return o.tx.Rollback()
}

func (o *OrmTx) Select(field ...string) ActionTxI {
	sql := _select(o.table, field...)
	return o.passAction(sql)
}

func (o *OrmTx) Update(data map[string]interface{}) ActionTxI {
	sql, args, err := update(o.table, data)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *OrmTx) Insert(data ...interface{}) ActionTxI {
	sql, args, err := insert(o.table, data...)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *OrmTx) Replace(data ...interface{}) ActionTxI {
	sql, args, err := replace(o.table, data...)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *OrmTx) Delete() ActionTxI {
	sql, err := delete(o.table)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql)
}
