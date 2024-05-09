package orm

import (
	"fmt"
	"log"
	"reflect"
)

var _ FieldIfc = (*Field[any])(nil)

type FieldIfc interface {
	ColName(withEscape bool) string
	DBColName(withEscape bool) string
	RefVal() any
	Val() any
	AutoIncrement() bool
	Dirty() bool
	ExprIfc

	set(any)
	WithRef(val any) FieldIfc
	setPreVal(val any)
}

type Field[T any] struct {
	Name            string
	Schema          Schema
	IsAutoIncrement bool
	// 是否被缓存值
	scanned bool
	preVal  T
	refVal  *T
}

func (f *Field[T]) AutoIncrement() bool {
	return f.IsAutoIncrement
}

func (f *Field[T]) set(v any) {
	*f.refVal = v.(T)
}

func (f *Field[T]) ColName(withEscape bool) string {
	wrapFn := func(s string) string {
		return s
	}
	if withEscape {
		wrapFn = FieldWrapper
	}
	return wrapFn(f.Name)
}

func (f *Field[T]) DBColName(withEscape bool) string {
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

func (f *Field[T]) WithRef(ref any) FieldIfc {
	new := &Field[T]{}
	*new = *f
	new.refVal = ref.(*T)
	return new
}

func (f *Field[T]) RefVal() any {
	return f.refVal
}

func (f *Field[T]) Val() any {
	return *f.refVal
}

func (f *Field[T]) setPreVal(val any) {
	f.scanned = true
	// t := val.(*T)
	f.preVal = val.(T)
}

func (f *Field[T]) Dirty() bool {
	pre, cur := any(f.preVal), any(*f.refVal)
	log.Printf("scanned: %v, preVal: %v, curVal: %v", f.scanned, pre, cur)
	if f.scanned {
		return !reflect.DeepEqual(pre, cur)
	}
	return true
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

func (f Field[T]) InQuery(q ExprIfc) Cond {
	return Cond{
		left:  &f,
		Op:    "IN",
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

type FieldGroup []FieldIfc

func Group(field ...FieldIfc) FieldGroup {
	return FieldGroup(field)
}

func (fg FieldGroup) InQuery(ExprIfc) Cond {
	// TODO: 实现
	return Cond{}
}
