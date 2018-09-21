// orm test init

package orm

import (
	"database/sql"
)

// Orm entrance
type Orm struct {
	db    *sql.DB
	table string
}

var _ OrmI = &Orm{}
var _ TabledI = &Orm{}

func New(db *sql.DB) OrmI {
	return &Orm{
		db: db,
	}
}

func toTabledI(i interface{}) TabledI {
	return i.(TabledI)
}

// Reset the orm instance
func (o *Orm) Reset() {
	o.table = ""
}

func (o *Orm) errAction(err error) ActionI {
	defer o.Reset()
	return &Action{
		err: err,
	}
}
func (o *Orm) passAction(sql string, args ...interface{}) ActionI {
	defer o.Reset()
	return &Action{
		db:   o.db,
		sql:  sql,
		args: args,
	}
}

func (o *Orm) Table(t string) TabledI {
	o.table = t
	return toTabledI(o)
}

func (o *Orm) Exec(sql string, arg ...interface{}) ActionI {
	return o.passAction(sql, arg...)
}

func (o *Orm) Select(field ...string) ActionI {
	sql := _select(o.table, field...)
	return o.passAction(sql)
}

func (o *Orm) Update(data map[string]interface{}) ActionI {
	sql, args, err := update(o.table, data)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *Orm) Insert(data ...interface{}) ActionI {
	sql, args, err := insert(o.table, data...)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *Orm) Replace(data ...interface{}) ActionI {
	sql, args, err := replace(o.table, data...)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql, args...)
}

func (o *Orm) Delete() ActionI {
	sql, err := delete(o.table)
	if err != nil {
		return o.errAction(err)
	}
	return o.passAction(sql)
}
