package orm

import (
	"log"
	"reflect"
)

type FieldIfc interface {
	ColName() string
	RefVal() any
	Val() any
	Dirty() bool
	WithRef(val any) FieldIfc
	setPreVal(val any)
}

type Field[T any] struct {
	Name string
	// 是否被缓存值
	scanned bool
	preVal  T
	refVal  *T
}

func (f *Field[T]) ColName() string {
	return f.Name
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
