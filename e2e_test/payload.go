package e2etest

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
	ID      int64
	Name    string
	Team    teamPayload
	TeamPtr *teamPayload
}

func (p *userWithTeamPayload) Bind() {
	// p.PayloadBase.BindField(user.ID.WithRef(&p.ID))
	orm.BindField(&p.ID, user.ID, &p.PayloadBase)
	p.PayloadBase.BindField(&p.Name, user.Name)
}

type userAndTeamPayload struct {
	orm.PayloadBase
	UserID   int64
	Name     string
	TeamName string
}

func (p *userAndTeamPayload) Bind() {
	p.PayloadBase.BindField(&p.UserID, user.ID)
	p.PayloadBase.BindField(&p.Name, user.Name)
	p.PayloadBase.BindField(&p.TeamName, team.Name)
}
