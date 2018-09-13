// test insert replace

package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertOneMap(t *testing.T) {
	data := M{
		"id":   100,
		"name": "insert100",
	}
	id, count, err := o.Table("test").Insert(data).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(100), id)
	assert.Equal(t, int64(1), count)
}

func TestInsertManyMap(t *testing.T) {
	data1 := M{
		"id":   110,
		"name": "insert110",
	}
	data2 := M{
		"id":   111,
		"name": "insert111",
	}
	id, count, err := o.Table("test").Insert(data1, data2).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(111), id)
	assert.Equal(t, int64(2), count)
}

func TestInsertOneStruct(t *testing.T) {
	data := &destT{
		ID:   101,
		Name: "insert101",
	}
	id, count, err := o.Table("test").Insert(data).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(101), id)
	assert.Equal(t, int64(1), count)
}

func TestInsertManyStruct(t *testing.T) {
	data1 := &destT{
		ID:   120,
		Name: "insert120",
	}
	data2 := &destT{
		ID:   121,
		Name: "insert121",
	}
	id, count, err := o.Table("test").Insert(data1, data2).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(121), id)
	assert.Equal(t, int64(2), count)
}
func TestReplaceOneMap(t *testing.T) {
	data := M{
		"id":   200,
		"name": "insert200",
	}
	id, count, err := o.Table("test").Replace(data).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(200), id)
	assert.Equal(t, int64(1), count)
}

func TestReplaceManyMap(t *testing.T) {
	data1 := M{
		"id":   210,
		"name": "insert210",
	}
	data2 := M{
		"id":   211,
		"name": "insert211",
	}
	id, count, err := o.Table("test").Replace(data1, data2).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(211), id)
	assert.Equal(t, int64(2), count)
}

func TestReplaceOneStruct(t *testing.T) {
	data := &destT{
		ID:   201,
		Name: "insert201",
	}
	id, count, err := o.Table("test").Replace(data).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(201), id)
	assert.Equal(t, int64(1), count)
}

func TestReplaceManyStruct(t *testing.T) {
	data1 := &destT{
		ID:   220,
		Name: "insert220",
	}
	data2 := &destT{
		ID:   221,
		Name: "insert221",
	}
	id, count, err := o.Table("test").Replace(data1, data2).Do()
	assert.NoError(t, err)
	assert.Equal(t, int64(221), id)
	assert.Equal(t, int64(2), count)
}
