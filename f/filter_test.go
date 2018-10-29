package f

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS(t *testing.T) {
	ret := S("a=?", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
	ret = S("a=? and b>=?", "1", "2")
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "a=? and b>=?", ret.Where)
}

func TestEqual(t *testing.T) {
	ret := Equal("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
	ret = Equal("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
}

func TestNotEqual(t *testing.T) {
	ret := NotEqual("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a!=?", ret.Where)
	ret = NotEqual("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a!=?", ret.Where)
}

func TestLte(t *testing.T) {
	ret := Lte("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a<=?", ret.Where)
	ret = Lte("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a<=?", ret.Where)
}

func TestLt(t *testing.T) {
	ret := Lt("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a<?", ret.Where)
	ret = Lt("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a<?", ret.Where)
}

func TestGte(t *testing.T) {
	ret := Gte("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a>=?", ret.Where)
	ret = Gte("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a>=?", ret.Where)
}

func TestGt(t *testing.T) {
	ret := Gt("a", "1")
	assert.Equal(t, []interface{}{"1"}, ret.Args)
	assert.Equal(t, "a>?", ret.Where)
	ret = Gt("a", 1)
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a>?", ret.Where)
}

func TestIn(t *testing.T) {
	ret := In("a", []string{"1", "2"})
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "a in (?,?)", ret.Where)
	ret = In("a", []int{1, 2})
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "a in (?,?)", ret.Where)
}

func TestNotIn(t *testing.T) {
	ret := NotIn("a", []string{"1", "2"})
	assert.Equal(t, []interface{}{"1", "2"}, ret.Args)
	assert.Equal(t, "a not in (?,?)", ret.Where)
	ret = NotIn("a", []int{1, 2})
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "a not in (?,?)", ret.Where)
}

func TestLike(t *testing.T) {
	ret := Like("a", "%a%")
	assert.Equal(t, []interface{}{"%a%"}, ret.Args)
	assert.Equal(t, "a like ?", ret.Where)
}

func TestNotLike(t *testing.T) {
	ret := NotLike("a", "%a%")
	assert.Equal(t, []interface{}{"%a%"}, ret.Args)
	assert.Equal(t, "a not like ?", ret.Where)
}

func TestAnd(t *testing.T) {
	ret := And(S("a=?", 1), S("b=?", 2))
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "(a=? and b=?)", ret.Where)
	ret = And(S("a=?", 1))
	assert.Equal(t, []interface{}{1}, ret.Args)
	assert.Equal(t, "a=?", ret.Where)
}

func TestOr(t *testing.T) {
	ret := Or(S("a=?", 1), S("b=?", 2))
	assert.Equal(t, []interface{}{1, 2}, ret.Args)
	assert.Equal(t, "(a=? or b=?)", ret.Where)
}
