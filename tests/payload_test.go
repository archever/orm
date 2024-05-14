package tests

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_TakePayload_Simple(t *testing.T) {
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

func Test_FindPayload_Simple(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	m.MockDB.ExpectQuery("SELECT `id`, `name` FROM `user`").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(10, "name1").
			AddRow(11, "name2"),
		)
	var payload []*userPayload
	err := cli.Table(user).Select().FindPayload(ctx, &payload)
	assert.NoError(t, err)
	assert.EqualValues(t, &userPayload{
		PayloadBase: payload[0].PayloadBase,
		ID:          10,
		Name:        "name1",
	}, payload[0])
	assert.EqualValues(t, &userPayload{
		PayloadBase: payload[1].PayloadBase,
		ID:          11,
		Name:        "name2",
	}, payload[1])
}
