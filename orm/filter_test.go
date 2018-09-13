// test where filter

package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquel(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Equal("id", 1)).Sql()
	tsql := "select * from test where id=?"
	targs := []interface{}{1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestNotEquel(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(NotEqual("id", 1)).Sql()
	tsql := "select * from test where id!=?"
	targs := []interface{}{1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestLike(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Like("id", "1%")).Sql()
	tsql := "select * from test where id like ?"
	targs := []interface{}{"1%"}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestNotLike(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(NotLike("id", "1%")).Sql()
	tsql := "select * from test where id not like ?"
	targs := []interface{}{"1%"}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}
func TestGte(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Gte("id", 1)).Sql()
	tsql := "select * from test where id>=?"
	targs := []interface{}{1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestGt(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Gt("id", 1)).Sql()
	tsql := "select * from test where id>?"
	targs := []interface{}{1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestAnd(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Gt("id", 1), Gte("id", 1)).Sql()
	tsql := "select * from test where id>? and id>=?"
	targs := []interface{}{1, 1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)

	sql, args, err = o.Table("test").Select().Where(And(Gt("id", 1), Gte("id", 1))).Sql()
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}
func TestOr(t *testing.T) {
	sql, args, err := o.Table("test").Select().Where(Or(Gt("id", 1), Gte("id", 1))).Sql()
	tsql := "select * from test where id>? or id>=?"
	targs := []interface{}{1, 1}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}
func TestLimit(t *testing.T) {
	sql, args, err := o.Table("test").Select().Limit(10).Sql()
	tsql := "select * from test limit ?"
	targs := []interface{}{int64(10)}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestOffset(t *testing.T) {
	sql, args, err := o.Table("test").Select().Offset(1).Sql()
	tsql := "select * from test offset ?"
	targs := []interface{}{int64(1)}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestOrderBy(t *testing.T) {
	sql, args, err := o.Table("test").Select().OrderBy("id", "desc").OrderBy("name").Sql()
	tsql := "select * from test order by id desc, name"
	var targs []interface{}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestGroupBy(t *testing.T) {
	sql, args, err := o.Table("test").Select().GroupBy("id").OrderBy("name", "desc").Sql()
	tsql := "select * from test group by id order by name desc"
	var targs []interface{}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestPage(t *testing.T) {
	sql, args, err := o.Table("test").Select().Page(1, 10).Sql()
	tsql := "select * from test limit ?"
	targs := []interface{}{int64(10)}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}

func TestPageOffset(t *testing.T) {
	sql, args, err := o.Table("test").Select().Page(2, 10).Sql()
	tsql := "select * from test limit ? offset ?"
	targs := []interface{}{int64(10), int64(10)}
	assert.NoError(t, err)
	assert.Equal(t, tsql, sql)
	assert.Equal(t, targs, args)
}
