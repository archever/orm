package example

import (
	"github.com/archever/orm"
)

var user = &userSchema{
	ID:     orm.Field[int64]{Name: "id", Schema: &userSchema{}},
	Name:   orm.Field[string]{Name: "name", Schema: &userSchema{}},
	TeamID: orm.Field[int64]{Name: "team_id", Schema: &userSchema{}},
}

type userSchema struct {
	ID     orm.Field[int64]
	Name   orm.Field[string]
	TeamID orm.Field[int64]
}

func (s *userSchema) TableName() string {
	return "user"
}

var team = &teamSchema{
	ID:   orm.Field[int64]{Name: "id", Schema: &teamSchema{}},
	Name: orm.Field[string]{Name: "name", Schema: &teamSchema{}},
}

type teamSchema struct {
	ID   orm.Field[int64]
	Name orm.Field[string]
}

func (s *teamSchema) TableName() string {
	return "team"
}
