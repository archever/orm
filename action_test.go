package orm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var aWithT = &Action{
	table: "T",
}

var a = &Action{}

func TestSelect(t *testing.T) {
	s1 := aWithT.Select("a", "b", "c")
	assert.Equal(t, "select `a`, `b`, `c` from `T`", s1.sql)
	s2 := aWithT.Select("a")
	assert.Equal(t, "select `a` from `T`", s2.sql)

	s3 := a.Select("a", "b", "c")
	assert.Equal(t, "select `a`, `b`, `c`", s3.sql)
	s4 := a.Select("a")
	assert.Equal(t, "select `a`", s4.sql)
}

func TestSelectS(t *testing.T) {
	s1 := aWithT.SelectS("a", "b", "c")
	assert.Equal(t, "select a, b, c from `T`", s1.sql)
	assert.Nil(t, s1.args)
	s2 := aWithT.SelectS("sum(`a`) as s")
	assert.Equal(t, "select sum(`a`) as s from `T`", s2.sql)
	assert.Nil(t, s1.args)

	s3 := a.SelectS("a", "b", "c")
	assert.Equal(t, "select a, b, c", s3.sql)
	assert.Nil(t, s1.args)
	s4 := a.SelectS("sum(`a`) as s")
	assert.Equal(t, "select sum(`a`) as s", s4.sql)
	assert.Nil(t, s1.args)
}

func TestUpdate(t *testing.T) {
	s1 := aWithT.Update(M{
		"a": "a",
	})
	assert.Equal(t, "update `T` set `a`=?", s1.sql)
	assert.Equal(t, []interface{}{"a"}, s1.args)
	s2 := aWithT.Update(M{
		"a": "a",
		"b": 1,
	})
	assert.Equal(t, "update `T` set `a`=?, `b`=?", s2.sql)
	assert.Equal(t, []interface{}{"a", 1}, s2.args)

	s3 := a.Update(M{
		"a": "a",
	})
	assert.Equal(t, errors.New("table not set"), s3.err)
	s4 := a.Update(M{})
	assert.Equal(t, errors.New("update empty data"), s4.err)
}

type InsertT struct {
	Name string
}

func TestInsert(t *testing.T) {
	s1 := aWithT.Insert(M{"a": "a"})
	assert.Equal(t, "insert into `T`(`a`) values (?)", s1.sql)
	assert.Equal(t, []interface{}{"a"}, s1.args)

	s2 := aWithT.Insert(M{})
	assert.Equal(t, errors.New("save empty data"), s2.err)

	s3 := aWithT.Insert(InsertT{Name: "name"})
	assert.Equal(t, "insert into `T`(`Name`) values (?)", s3.sql)
	assert.Equal(t, []interface{}{"name"}, s3.args)

	s4 := a.Insert(InsertT{Name: "name"})
	assert.Equal(t, errors.New("table not set"), s4.err)

	s5 := aWithT.Insert([]InsertT{InsertT{Name: "name"}})
	assert.Equal(t, errors.New("invalid data type: []orm.InsertT"), s5.err)

	s6 := aWithT.Insert(1)
	assert.Equal(t, errors.New("invalid data type: int"), s6.err)
}

func TestInsertMany(t *testing.T) {
	s1 := aWithT.InsertMany([]M{M{"a": "a"}})
	assert.Equal(t, "insert into `T`(`a`) values (?)", s1.sql)
	assert.Equal(t, []interface{}{"a"}, s1.args)

	s11 := aWithT.InsertMany([]M{M{"a": "a"}, M{"a": "b"}})
	assert.Equal(t, "insert into `T`(`a`) values (?), (?)", s11.sql)
	assert.Equal(t, []interface{}{"a", "b"}, s11.args)

	s12 := aWithT.InsertMany([]M{M{"a": "a"}, M{"b": "b"}})
	assert.Equal(t, errors.New("can not save many with different data field"), s12.err)

	s13 := aWithT.InsertMany([]M{M{"a": "a", "b": "b"}, M{"b": "b"}})
	assert.Equal(t, errors.New("can not save many with different data field"), s13.err)

	s2 := aWithT.InsertMany([]M{})
	assert.Equal(t, errors.New("save empty data"), s2.err)

	s3 := aWithT.InsertMany([]InsertT{InsertT{Name: "name"}})
	assert.Equal(t, "insert into `T`(`Name`) values (?)", s3.sql)
	assert.Equal(t, []interface{}{"name"}, s3.args)

	s4 := a.InsertMany(InsertT{Name: "name"})
	assert.Equal(t, errors.New("table not set"), s4.err)

	s5 := aWithT.InsertMany(InsertT{Name: "name"})
	assert.Equal(t, errors.New("invalid data type: orm.InsertT"), s5.err)

	s6 := aWithT.InsertMany(1)
	assert.Equal(t, errors.New("invalid data type: int"), s6.err)
}

func TestReplace(t *testing.T) {
	s1 := aWithT.Replace(M{"a": "a"})
	assert.Equal(t, "replace into `T`(`a`) values (?)", s1.sql)
	assert.Equal(t, []interface{}{"a"}, s1.args)

	s2 := aWithT.Replace(M{})
	assert.Equal(t, errors.New("save empty data"), s2.err)

	s3 := aWithT.Replace(InsertT{Name: "name"})
	assert.Equal(t, "replace into `T`(`Name`) values (?)", s3.sql)
	assert.Equal(t, []interface{}{"name"}, s3.args)

	s4 := a.Replace(InsertT{Name: "name"})
	assert.Equal(t, errors.New("table not set"), s4.err)

	s5 := aWithT.Replace([]InsertT{InsertT{Name: "name"}})
	assert.Equal(t, errors.New("invalid data type: []orm.InsertT"), s5.err)

	s6 := aWithT.Replace(1)
	assert.Equal(t, errors.New("invalid data type: int"), s6.err)
}

func TestReplaceMany(t *testing.T) {
	s1 := aWithT.ReplaceMany([]M{M{"a": "a"}})
	assert.Equal(t, "replace into `T`(`a`) values (?)", s1.sql)
	assert.Equal(t, []interface{}{"a"}, s1.args)

	s11 := aWithT.ReplaceMany([]M{M{"a": "a"}, M{"a": "b"}})
	assert.Equal(t, "replace into `T`(`a`) values (?), (?)", s11.sql)
	assert.Equal(t, []interface{}{"a", "b"}, s11.args)

	s12 := aWithT.ReplaceMany([]M{M{"a": "a"}, M{"b": "b"}})
	assert.Equal(t, errors.New("can not save many with different data field"), s12.err)

	s13 := aWithT.ReplaceMany([]M{M{"a": "a", "b": "b"}, M{"b": "b"}})
	assert.Equal(t, errors.New("can not save many with different data field"), s13.err)

	s2 := aWithT.ReplaceMany([]M{})
	assert.Equal(t, errors.New("save empty data"), s2.err)

	s3 := aWithT.ReplaceMany([]InsertT{InsertT{Name: "name"}})
	assert.Equal(t, "replace into `T`(`Name`) values (?)", s3.sql)
	assert.Equal(t, []interface{}{"name"}, s3.args)

	s4 := a.ReplaceMany(InsertT{Name: "name"})
	assert.Equal(t, errors.New("table not set"), s4.err)

	s5 := aWithT.ReplaceMany(InsertT{Name: "name"})
	assert.Equal(t, errors.New("invalid data type: orm.InsertT"), s5.err)

	s6 := aWithT.ReplaceMany(1)
	assert.Equal(t, errors.New("invalid data type: int"), s6.err)
}

func TestDelete(t *testing.T) {
	s1 := aWithT.Delete()
	assert.Equal(t, "delete from `T`", s1.sql)
	assert.Nil(t, s1.args)

	s2 := a.Delete()
	assert.Equal(t, errors.New("table not set"), s2.err)
}
