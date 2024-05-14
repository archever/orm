package orm

import (
	"fmt"
)

var _ FieldIfc = (*Field[any])(nil)

type FieldIfc interface {
	ColName(withEscape bool) string
	DBColName(withEscape bool) string
	IsAutoIncrement() bool
	ExprIfc

	key() string
}

type Field[T any] struct {
	Name          string
	Schema        Schema
	AutoIncrement bool
}

func NewField[T any](name string, schema Schema) *Field[T] {
	return &Field[T]{
		Name:          name,
		Schema:        schema,
		AutoIncrement: false,
	}
}

func (f *Field[T]) SetAutoIncrement(b bool) {
	f.AutoIncrement = b
}

func (f Field[T]) IsAutoIncrement() bool {
	return f.AutoIncrement
}

func (f Field[T]) key() string {
	return fmt.Sprintf("%s.%s", f.Schema.TableName(), f.Name)
}

func (f Field[T]) ColName(withEscape bool) string {
	wrapFn := func(s string) string {
		return s
	}
	if withEscape {
		wrapFn = FieldWrapper
	}
	return wrapFn(f.Name)
}

func (f Field[T]) DBColName(withEscape bool) string {
	wrapFn := func(s string) string {
		return s
	}
	if withEscape {
		wrapFn = FieldWrapper
	}
	return fmt.Sprintf("%s.%s",
		wrapFn(f.Schema.TableName()),
		wrapFn(f.Name),
	)
}

func (f Field[T]) Expr() (string, []any) {
	return f.DBColName(true), []any{}
}

func (f Field[T]) Eq(val T) Cond {
	return Cond{
		left:  &f,
		Op:    "=",
		right: anyVal{val},
	}
}

func (f Field[T]) EqCol(col Field[T]) Cond {
	return Cond{
		left:  &f,
		Op:    "=",
		right: &col,
	}
}

func (f Field[T]) EqQuery(q ExprIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "=",
		right: brackets{q},
	}
}

func (f Field[T]) NotEq(val T) Cond {
	return Cond{
		left:  &f,
		Op:    "<>",
		right: anyVal{val},
	}
}

func (f Field[T]) NotEqCol(col Field[T]) Cond {
	return Cond{
		left:  &f,
		Op:    "<>",
		right: &col,
	}
}

func (f Field[T]) In(val ...T) Cond {
	anyList := []any{}
	for _, v := range val {
		anyList = append(anyList, v)
	}
	return Cond{
		left:  &f,
		Op:    "IN",
		right: brackets{anyValList(anyList)},
	}
}

func (f Field[T]) NotIn(val ...T) Cond {
	anyList := []any{}
	for _, v := range val {
		anyList = append(anyList, v)
	}
	return Cond{
		left:  &f,
		Op:    "NOT IN",
		right: brackets{anyValList(anyList)},
	}
}

func (f Field[T]) InQuery(q ExprIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "IN",
		right: brackets{q},
	}
}

func (f Field[T]) NotInQuery(q ExprIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "NOT IN",
		right: brackets{q},
	}
}

func (f Field[T]) IsNull(isNull bool) Cond {
	if isNull {
		return Cond{
			left: &f,
			Op:   "IS NULL",
		}
	}
	return Cond{
		left: &f,
		Op:   "IS NOT NULL",
	}
}

func (f Field[T]) Gt(val T) Cond {
	return Cond{
		left:  &f,
		Op:    ">",
		right: anyVal{val},
	}
}

func (f Field[T]) Gte(val T) Cond {
	return Cond{
		left:  &f,
		Op:    ">=",
		right: anyVal{val},
	}
}

func (f Field[T]) Lt(val T) Cond {
	return Cond{
		left:  &f,
		Op:    "<",
		right: anyVal{val},
	}
}

func (f Field[T]) Lte(val T) Cond {
	return Cond{
		left:  &f,
		Op:    "<=",
		right: anyVal{val},
	}
}

func (f Field[T]) GtCol(col FieldIfc) Cond {
	return Cond{
		left:  &f,
		Op:    ">",
		right: col,
	}
}

func (f Field[T]) GteCol(col FieldIfc) Cond {
	return Cond{
		left:  &f,
		Op:    ">=",
		right: col,
	}
}

func (f Field[T]) LtCol(col FieldIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "<",
		right: col,
	}
}

func (f Field[T]) LteCol(col FieldIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "<=",
		right: col,
	}
}

func (f Field[T]) Desc(desc bool) Order {
	return Order{
		Field: &f,
		Desc:  desc,
	}
}

func (f Field[T]) Asc() Order {
	return Order{
		Field: &f,
		Desc:  false,
	}
}

type FieldGroup []FieldIfc

func Group(field ...FieldIfc) FieldGroup {
	return FieldGroup(field)
}

func (fg FieldGroup) InQuery(q ExprIfc) Cond {
	left := brackets{
		ExprIfc: fields{
			fields:        fg,
			withTableName: true,
		},
	}
	return Cond{
		left:  left,
		Op:    "IN",
		right: brackets{q},
	}
}

func (fg FieldGroup) NotInQuery(q ExprIfc) Cond {
	left := brackets{
		ExprIfc: fields{
			fields:        fg,
			withTableName: true,
		},
	}
	return Cond{
		left:  left,
		Op:    "NOT IN",
		right: brackets{q},
	}
}
