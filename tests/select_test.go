package tests

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
	"github.com/stretchr/testify/assert"
)

func Test_Select_Simple(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `user`.`id` = ? LIMIT ?").
		WithArgs(10, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
		)
	var payload userPayload
	err := cli.Table(user).Select().Where(user.ID.Eq(10)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	assert.EqualValues(t, userPayload{
		PayloadBase: payload.PayloadBase,
		ID:          10,
		Name:        "archever",
	}, payload)
}

func Test_Select_Where(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	{
		m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE (`user`.`id` = ? AND `user`.`name` IS NULL AND `user`.`name` > ? AND `user`.`name` >= ? AND `user`.`name` < ? AND `user`.`name` <= ? AND `user`.`name` <> ? AND `user`.`name` IN (?,?) AND `user`.`name` NOT IN (?,?)) LIMIT ?").
			WithArgs(10, "gt", "gte", "lt", "lte", "nte", "ina", "inb", "nIna", "nInb", 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload userPayload
		err := cli.Table(user).Select().
			Where(user.ID.Eq(10)).
			Where(user.Name.IsNull(true)).
			Where(user.Name.Gt("gt")).
			Where(user.Name.Gte("gte")).
			Where(user.Name.Lt("lt")).
			Where(user.Name.Lte("lte")).
			Where(user.Name.NotEq("nte")).
			Where(user.Name.In("ina", "inb")).
			Where(user.Name.NotIn("nIna", "nInb")).
			TakePayload(ctx, &payload)
		assert.NoError(t, err)
	}
	{
		m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE (`user`.`name` > ? OR (`user`.`name` < ? AND `user`.`name` IS NULL)) LIMIT ?").
			WithArgs(10, "gt", "gte", "lt", "lte", "nte", "ina", "inb", "nIna", "nInb", 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload userPayload
		err := cli.Table(user).Select().
			Where(orm.Or(
				user.Name.Gt("gt"),
				orm.And(
					user.Name.Lt("lt"),
					user.Name.IsNull(true),
				),
			)).
			TakePayload(ctx, &payload)
		assert.NoError(t, err)
	}
}

func Test_Select_GroupOrder(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `user`.`id` = ? GROUP BY `id`, `name` order by `user`.`id`, `user`.`name` DESC LIMIT ?").
		WithArgs(10, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
		)
	var payload userPayload
	err := cli.Table(user).Select().
		Where(user.ID.Eq(10)).
		GroupBy(user.ID).
		GroupBy(user.Name).
		OrderBy(user.ID.Asc(), user.Name.Desc(true)).
		TakePayload(ctx, &payload)
	assert.NoError(t, err)
	assert.EqualValues(t, userPayload{
		PayloadBase: payload.PayloadBase,
		ID:          10,
		Name:        "archever",
	}, payload)
}

func Test_Select_Page(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	{
		m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` LIMIT ?").
			WithArgs(10).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload []*userPayload
		err := cli.Table(user).Select().
			Page(1, 10).
			FindPayload(ctx, &payload)
		assert.NoError(t, err)
	}
	{
		m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` LIMIT ? OFFSET ?").
			WithArgs(10, 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload []*userPayload
		err := cli.Table(user).Select().
			Offset(1).Limit(10).
			FindPayload(ctx, &payload)
		assert.NoError(t, err)
	}
}

func Test_Select_Join(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	{
		m.MockDB.ExpectQuery("SELECT `user`.`id`, `user`.`name` FROM `user` JOIN `team` ON `user`.`id` = `team`.`id`").
			WithArgs(10).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload []*userPayload
		err := cli.Table(user).Select().
			Join(team, user.ID.EqCol(team.ID)).
			FindPayload(ctx, &payload)
		assert.NoError(t, err)
	}
	{
		m.MockDB.ExpectQuery("SELECT `user`.`id`, `user`.`name` FROM `user` LEFT JOIN `team` ON `user`.`id` = `team`.`id`").
			WithArgs(10, 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
			)
		var payload []*userPayload
		err := cli.Table(user).Select().
			LeftJoin(team, user.ID.EqCol(team.ID)).
			FindPayload(ctx, &payload)
		assert.NoError(t, err)
	}
}

func Test_Select_SubQuery(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)

	m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `user`.`team_id` IN (SELECT `id` FROM `team` WHERE `team`.`id` = ?)").
		WithArgs(10).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
		)
	var payload []*userPayload
	subQuery := cli.Table(team).Select(team.ID).Where(team.ID.Eq(10)).SubQuery()
	err := cli.Table(user).Select().
		Where(user.TeamID.InQuery(subQuery)).
		FindPayload(ctx, &payload)
	assert.NoError(t, err)
}
