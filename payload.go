package orm

import (
	"log"
	"reflect"

	"github.com/elliotchance/orderedmap/v2"
)

// type fieldBind[T any] struct {
// 	field Field[T]
// 	// 是否被缓存值
// 	scanned  bool
// 	preVal   T
// 	ref      *T
// 	extraRef []*T
// }

type fieldBind struct {
	field FieldIfc
	// 是否被缓存值
	scanned bool
	preVal  any
	ref     any
	// extraRef []any
}

type PayloadIfc interface {
	Bind()
	BoundFields() []*fieldBind
	Fields() []FieldIfc
}

type PayloadBase struct {
	binds *orderedmap.OrderedMap[string, *fieldBind]
}

func (p *PayloadBase) BindField(ref any, f FieldIfc) {
	key := f.key()
	if p.binds == nil {
		p.binds = orderedmap.NewOrderedMap[string, *fieldBind]()
	}
	if _, ok := p.binds.Get(key); !ok {
		p.binds.Set(key, &fieldBind{
			ref:   ref,
			field: f,
		})
	}
}

func (p *PayloadBase) BoundFields() []*fieldBind {
	dst := []*fieldBind{}
	for _, key := range p.binds.Keys() {
		value, _ := p.binds.Get(key)
		dst = append(dst, value)
	}
	return dst
}

func (p *PayloadBase) Fields() []FieldIfc {
	dst := []FieldIfc{}
	for _, key := range p.binds.Keys() {
		value, _ := p.binds.Get(key)
		dst = append(dst, value.field)
	}
	return dst
}

func BindField[T any](ref *T, f Field[T], base *PayloadBase) {
	base.BindField(ref, &f)
}

func BindFieldIfc(ref any, f FieldIfc, base *PayloadBase) {
	base.BindField(ref, f)
}

func boundFields(p PayloadIfc) []*fieldBind {
	p.Bind()
	// TODO: 查找嵌套的结构中的 PayloadIfc
	return p.BoundFields()
}

func (f *fieldBind) RefVal() any {
	return f.ref
}

func (f *fieldBind) Val() any {
	return reflect.ValueOf(f.ref).Elem().Interface()
}

func (f *fieldBind) Set(v any) {
	reflect.ValueOf(f.ref).Elem().Set(reflect.ValueOf(v))
}

func (f *fieldBind) setPreVal(val any) {
	f.scanned = true
	f.preVal = val
}

func (f *fieldBind) Dirty() bool {
	pre, cur := f.preVal, f.Val()
	log.Printf("scanned: %v, preVal: %v, curVal: %v", f.scanned, pre, cur)
	if f.scanned {
		return !reflect.DeepEqual(pre, cur)
	}
	return true
}
