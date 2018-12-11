package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterS(t *testing.T) {
	ret := FilterS("a=?", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
	ret = FilterS("a=? and b>=?", "1", "2")
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "a=? and b>=?", ret.Where)
}

func TestIsNull(t *testing.T) {
	ret := IsNull("a")
	assert.Nil(t, ret.Args)
	assert.Equal(t, "`a` is null", ret.Where)
}

func TestIsNotNull(t *testing.T) {
	ret := IsNotNull("a")
	assert.Nil(t, ret.Args)
	assert.Equal(t, "`a` is not null", ret.Where)
}

func TestEqual(t *testing.T) {
	ret := Equal("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`=?", ret.Where)
	ret = Equal("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`=?", ret.Where)
}

func TestNotEqual(t *testing.T) {
	ret := NotEqual("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`!=?", ret.Where)
	ret = NotEqual("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`!=?", ret.Where)
}

func TestLte(t *testing.T) {
	ret := Lte("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`<=?", ret.Where)
	ret = Lte("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`<=?", ret.Where)
}

func TestLt(t *testing.T) {
	ret := Lt("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`<?", ret.Where)
	ret = Lt("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`<?", ret.Where)
}

func TestGte(t *testing.T) {
	ret := Gte("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`>=?", ret.Where)
	ret = Gte("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`>=?", ret.Where)
}

func TestGt(t *testing.T) {
	ret := Gt("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "`a`>?", ret.Where)
	ret = Gt("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "`a`>?", ret.Where)
}

func TestBetween(t *testing.T) {
	ret := Between("a", 1, 2)
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "`a` between ? and ?", ret.Where)
}

func TestIn(t *testing.T) {
	ret := In("a", []string{"1", "2"})
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "`a` in (?,?)", ret.Where)
	ret = In("a", []int{1, 2})
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "`a` in (?,?)", ret.Where)
}

func TestNotIn(t *testing.T) {
	ret := NotIn("a", []string{"1", "2"})
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "`a` not in (?,?)", ret.Where)
	ret = NotIn("a", []int{1, 2})
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "`a` not in (?,?)", ret.Where)
}

func TestLike(t *testing.T) {
	ret := Like("a", "%a%")
	assert.Equal(t, []interface{}{"%a%"}, ret.Args)
	assert.Equal(t, "`a` like ?", ret.Where)
}

func TestNotLike(t *testing.T) {
	ret := NotLike("a", "%a%")
	assert.Equal(t, []interface{}{"%a%"}, ret.Args)
	assert.Equal(t, "`a` not like ?", ret.Where)
}

func TestAnd(t *testing.T) {
	ret := And(FilterS("a=?", 1), FilterS("b=?", 2))
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "(a=? and b=?)", ret.Where)
	ret = And(FilterS("a=?", 1))
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
}

func TestOr(t *testing.T) {
	ret := Or(FilterS("a=?", 1), FilterS("b=?", 2))
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "(a=? or b=?)", ret.Where)
}
