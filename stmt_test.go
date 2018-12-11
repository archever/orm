package orm

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.Lshortfile)
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(drop)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(table)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`insert into test(name, data, birth) values("archever", "test", "2016-01-20")`)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCtx(t *testing.T) {
	s := Stmt{}
	assert.Nil(t, s.ctx)
	ctx := context.TODO()
	s.Ctx(ctx)
	assert.Equal(t, ctx, s.ctx)
}

func TestSQL(t *testing.T) {
	a := Stmt{
		sql: "select *",
	}
	sql, args, err := a.SQL()
	assert.Equal(t, "select *", sql)
	assert.Nil(t, args)
	assert.NoError(t, err)
}

func TestFilter(t *testing.T) {
	s1 := Stmt{}
	s1.Filter(FilterS("a=?", 1))
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " where a=?", sql)
	assert.Equal(t, []interface{}{1}, args)

	s2 := Stmt{}
	s2.Filter(FilterS("a=?", 1))
	s2.Filter(FilterS("b=?", 2))
	sql, args, err = s2.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " where (a=? and b=?)", sql)
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestWhere(t *testing.T) {
	s1 := Stmt{}
	s1.Where("a=?", 1)
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " where a=?", sql)
	assert.Equal(t, []interface{}{1}, args)

	s2 := Stmt{}
	s2.Where("a=?", 1)
	s2.Where("b=?", 2)
	sql, args, err = s2.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " where (a=? and b=?)", sql)
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestOrderBy(t *testing.T) {
	s1 := Stmt{}
	s1.OrderBy("a", false)
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " order by `a`", sql)
	assert.Nil(t, args)

	s2 := Stmt{}
	s2.OrderBy("a", true)
	sql, args, err = s2.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " order by `a` desc", sql)
	assert.Nil(t, args)

	s3 := Stmt{}
	s3.OrderBy("a", true)
	s3.OrderBy("b", false)
	sql, args, err = s3.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " order by `a` desc, `b`", sql)
	assert.Nil(t, args)
}

func TestGroupBy(t *testing.T) {
	s1 := Stmt{}
	s1.GroupBy("a")
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " group by `a`", sql)
	assert.Nil(t, args)

	s2 := Stmt{}
	s2.GroupBy("a", "b")
	sql, args, err = s2.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " group by `a`, `b`", sql)
	assert.Nil(t, args)

	s3 := Stmt{}
	s3.GroupBy("a")
	s3.GroupBy("b")
	sql, args, err = s3.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " group by `a`, `b`", sql)
	assert.Nil(t, args)
}

func TestLimit(t *testing.T) {
	s1 := Stmt{}
	s1.Limit(1)
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " limit ?", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestOffset(t *testing.T) {
	s1 := Stmt{}
	s1.Offset(1)
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " offset ?", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestPage(t *testing.T) {
	s1 := Stmt{}
	s1.Page(2, 10)
	sql, args, err := s1.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " limit ? offset ?", sql)
	assert.Equal(t, []interface{}{10, 10}, args)

	s2 := Stmt{}
	s2.Page(1, 10)
	sql, args, err = s2.SQL()
	assert.NoError(t, err)
	assert.Equal(t, " limit ?", sql)
	assert.Equal(t, []interface{}{10}, args)
}

func TestDo(t *testing.T) {
	s, err := Open("mysql", dbURI)
	assert.NoError(t, err)
	i, c, e := s.Exec(`select 1 as d`).Do()
	assert.NoError(t, e)
	assert.Equal(t, int64(0), i)
	assert.Equal(t, int64(0), c)
}

func TestCount(t *testing.T) {
	s, err := Open("mysql", dbURI)
	assert.NoError(t, err)
	c, e := s.Table("test").Select().Count()
	assert.NoError(t, e)
	assert.Equal(t, int64(1), c)
}

func TestCountEmpty(t *testing.T) {
	s, err := Open("mysql", dbURI)
	assert.NoError(t, err)
	c, e := s.Table("test").Select().Where("id=?", 999).Count()
	assert.NoError(t, e)
	assert.Equal(t, int64(0), c)
}

func TestGet(t *testing.T) {
	s, err := Open("mysql", dbURI)
	assert.NoError(t, err)
	var res []M
	e := s.Table("test").Select().Get(&res)
	assert.NoError(t, e)
	assert.Equal(t, 1, len(res))

	var res2 []M
	e = s.Table("test").Select().Where("id=?", 9999).Get(&res2)
	assert.NoError(t, e)
	assert.Nil(t, res2)

	var res3 M
	e = s.Table("test").Select().Get(&res3)
	assert.Equal(t, errors.New("dest type not match: *orm.M"), e)
}

func TestOne(t *testing.T) {
	s, err := Open("mysql", dbURI)
	assert.NoError(t, err)

	var res M
	e := s.Table("test").Select().One(&res)
	assert.NoError(t, e)
	assert.Equal(t, "archever", res["name"])

	var res2 []M
	e = s.Table("test").Select().One(&res2)
	assert.Equal(t, errors.New("dest type not match: *[]orm.M"), e)
}
