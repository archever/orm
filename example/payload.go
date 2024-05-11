package example

import "github.com/archever/orm"

type userPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
}

func (p *userPayload) Bind() {
	p.PayloadBase.BindField(&p.ID, user.ID)
	p.PayloadBase.BindField(&p.Name, user.Name)
}

type teamPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
}

func (p *teamPayload) Bind() {
	p.PayloadBase.BindField(&p.ID, team.ID)
	p.PayloadBase.BindField(&p.Name, team.Name)
}

type userWithTeamPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
	Team *teamPayload
}

func (p *userWithTeamPayload) Bind() {
	p.PayloadBase.BindField(&p.ID, user.ID)
	p.PayloadBase.BindField(&p.Name, user.Name)
	if p.Team == nil {
		p.Team = &teamPayload{}
	}
	p.Team.Bind()
}
