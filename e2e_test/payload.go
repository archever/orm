package e2etest

import "github.com/archever/orm"

type userPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
}

func (p *userPayload) Bind() {
	p.PayloadBase.BindField(user.ID.WithRef(&p.ID))
	p.PayloadBase.BindField(user.Name.WithRef(&p.Name))
}

type teamPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
}

func (p *teamPayload) Bind() {
	p.PayloadBase.BindField(team.ID.WithRef(&p.ID))
	p.PayloadBase.BindField(team.Name.WithRef(&p.Name))
}

type userWithTeamPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
	Team *teamPayload
}

func (p *userWithTeamPayload) Bind() {
	p.PayloadBase.BindField(user.ID.WithRef(&p.ID))
	p.PayloadBase.BindField(user.Name.WithRef(&p.Name))
	// p.PayloadBase.BindNest(&p.Team)
}
