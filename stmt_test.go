package orm

import (
	"testing"

	"github.com/archever/orm/f"

	"github.com/stretchr/testify/assert"
)

func TestSQL(t *testing.T) {
	sql, args, err := s.Exec("select now()").SQL()
	assert.Equal(t, "select now()", sql)
	assert.Equal(t, []interface{}{}, args)
	assert.NoError(t, err)
}

func TestDo(t *testing.T) {
	rowID, rowCNT, err := s.Table("test").Insert(M{
		"name":     "archever",
		"type":     Male,
		"datetime": "2018-09-13 12:11:00",
	}).Do()
	assert.NotEqual(t, int64(0), rowID)
	assert.Equal(t, int64(1), rowCNT)
	assert.NoError(t, err)
}

func TestCount(t *testing.T) {
	var dest []M
	stmt := s.Table("test").Select("name").Where("name=?", "archever")
	cnt, err1 := stmt.Count()
	err2 := stmt.Get(&dest)
	assert.NotEqual(t, int64(0), cnt)
	assert.Equal(t, "archever", dest[0]["name"])
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestFilter(t *testing.T) {
	sql, args, err := s.Table("test").Select().Filter(f.Equal("a", 1), f.Gt("b", 2)).SQL()
	assert.Equal(t, "select * from `test` where (`a`=? and `b`>?)", sql)
	assert.Equal(t, []interface{}{1, 2}, args)
	assert.NoError(t, err)
}

type A struct{}

func (a *A) MarshalSQL() (string, error) {
	return "cmt", nil
}

var _ Marshaler = &A{}

func TestFilterMarshal(t *testing.T) {
	sql, args, err := s.Table("test").Select().Filter(f.Or(f.Equal("a", &A{}), f.Gt("b", 2))).SQL()
	assert.Equal(t, "select * from `test` where (`a`=? or `b`>?)", sql)
	assert.Equal(t, []interface{}{"cmt", 2}, args)
	assert.NoError(t, err)
}
