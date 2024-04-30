package orm

type PayloadIfc interface {
	Bind()
	Fields() []FieldIfc
}

type PayloadBase struct {
	binds []FieldIfc
}

func (p *PayloadBase) BindField(f FieldIfc) {
	p.binds = append(p.binds, f)
}

func (p *PayloadBase) Fields() []FieldIfc {
	return p.binds
}
