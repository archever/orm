package orm

import (
	"fmt"
	"log"
	"reflect"
)

type FieldIfc interface {
	ColName(withEscape bool) string
	DBColName(withEscape bool) string
	RefVal() any
	AutoIncrement() bool
	Val() any
	Dirty() bool

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

func (f Field[T]) Eq(val T) Cond {
	return Cond{
		left:     &f,
		Op:       "=",
		rightVal: val,
	}
}

func (f Field[T]) EqCol(col Field[T]) Cond {
	return Cond{
		left:       &f,
		Op:         "=",
		rightField: &col,
	}
}
