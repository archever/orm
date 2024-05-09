package e2etest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func Test_initData(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	team1 := teamPayload{
		Name: "team1",
	}
	team2 := teamPayload{
		Name: "team2",
	}
	cnt, err := cli.Table(team).InsertPayload(&team1, &team2).Do(ctx)
	assert.NoError(t, err)
	t.Logf("team1: %v", team1)
	t.Logf("team2: %v", team2)
	t.Logf("cnt: %v", cnt)
}

func TestJoin(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userPayload
	// select user.* from user join team on user.team_id=team.id
	err := cli.Table(user).Select().
		Join(team, user.TeamID.EqCol(team.ID)).
		Limit(1).
		TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestJoinSelectAll(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userAndTeamPayload
	// select * from user join team on user.team_id=team.id
	err := cli.Table(user).Select().
		Join(team, user.TeamID.EqCol(team.ID)).
		Limit(1).
		TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%#v", payload)
}

func TestJoinSelectNest(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userWithTeamPayload
	// select * from user join team on user.team_id=team.id
	err := cli.Table(user).Select().
		Join(team, user.TeamID.EqCol(team.ID)).
		Limit(1).
		// TakePayload(ctx, &payload, &payload.TeamPtr, &payload.Team)
		TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%#v", payload)
	// t.Logf("%#v", payload.Team)
	t.Logf("%#v", *payload.TeamPtr)
}
