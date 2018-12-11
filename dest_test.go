package orm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DestT struct {
	name      string
	ID        int64  `column:"omitempty"`
	TestName  string `column:"test_name"`
	Ignore    bool   `column:"-"`
	TestName2 string `column:"test_name2,omitempty"`
}

func TestGetFieldName(t *testing.T) {
	var dt DestT
	rv := reflect.TypeOf(dt)
	f1, _ := rv.FieldByName("name")
	n, o, i := GetFieldName(f1)
	assert.Equal(t, "name", n)
	assert.Equal(t, false, o)
	assert.Equal(t, true, i)

	f2, _ := rv.FieldByName("ID")
	n, o, i = GetFieldName(f2)
	assert.Equal(t, "ID", n)
	assert.Equal(t, true, o)
	assert.Equal(t, false, i)

	f3, _ := rv.FieldByName("TestName")
	n, o, i = GetFieldName(f3)
	assert.Equal(t, "test_name", n)
	assert.Equal(t, false, o)
	assert.Equal(t, false, i)

	f4, _ := rv.FieldByName("Ignore")
	n, o, i = GetFieldName(f4)
	assert.Equal(t, "Ignore", n)
	assert.Equal(t, false, o)
	assert.Equal(t, true, i)

	f5, _ := rv.FieldByName("TestName2")
	n, o, i = GetFieldName(f5)
	assert.Equal(t, "test_name2", n)
	assert.Equal(t, true, o)
	assert.Equal(t, false, i)
}

type aT struct {
	Msg string
}

func (t *aT) MarshalSQL() (string, error) { return t.Msg, nil }
func (t *aT) UnmarshalSQL(*ScanRow) error {
	t.Msg = "aT dest"
	return nil
}

type bT []string

func (t *bT) MarshalSQL() (string, error) {
	return strings.Join(*t, ","), nil
}

func (t *bT) UnmarshalSQL(*ScanRow) error {
	*t = []string{"t", "e", "s", "t"}
	return nil
}

func TestITOMarshaler1(t *testing.T) {
	a := aT{
		Msg: "test a",
	}
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&a))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "test a", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&a))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, "aT dest", a.Msg)
	}

	a2 := &aT{
		Msg: "test a2",
	}
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&a2))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "test a2", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&a2))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, "aT dest", a2.Msg)
	}

	a3_ := &aT{
		Msg: "test a3",
	}
	a3 := &a3_
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&a3))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "test a3", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&a3))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, "aT dest", (*a3).Msg)
	}
}

func TestITOMarshaler2(t *testing.T) {
	b := bT{"b", "t"}
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&b))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "b,t", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&b))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, bT{"t", "e", "s", "t"}, b)
	}

	b2 := &bT{"b", "t"}
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&b2))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "b,t", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&b2))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, &bT{"t", "e", "s", "t"}, b2)
	}

	b3_ := &bT{"b", "t"}
	b3 := &b3_
	{
		m, ok := ITOMarshaler(reflect.ValueOf(&b3))
		assert.Equal(t, true, ok)
		data, err := m.MarshalSQL()
		assert.NoError(t, err)
		assert.Equal(t, "b,t", data)
	}
	{
		m, ok := ITOUnmarshaler(reflect.ValueOf(&b3))
		assert.Equal(t, true, ok)
		err := m.UnmarshalSQL(nil)
		assert.NoError(t, err)
		assert.Equal(t, &bT{"t", "e", "s", "t"}, *b3)
	}
}
