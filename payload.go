package orm

type PayloadIfc interface {
	Bind()
	Fields() []FieldIfc
}

type PayloadBase struct {
	bindMap map[string]bool
	binds   []FieldIfc
}

func (p *PayloadBase) BindField(f FieldIfc) {
	key := f.DBColName(true)
	if p.bindMap == nil {
		p.bindMap = make(map[string]bool)
	}
	if p.bindMap[key] {
		return
	}
	p.bindMap[key] = true
	p.binds = append(p.binds, f)
}

func (p *PayloadBase) Fields() []FieldIfc {
	return p.binds
}
