package orm

type PayloadIfc interface {
	Bind()
	Fields() []FieldIfc
}

type PayloadBase struct {
	binds []FieldIfc
}

func (p *PayloadBase) BindField(f FieldIfc) {
	// TODO: check if field is already bound
	p.binds = append(p.binds, f)
}

func (p *PayloadBase) Fields() []FieldIfc {
	return p.binds
}
