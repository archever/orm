package tests

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Update_Simple(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	m.MockDB.ExpectExec("UPDATE `user` SET `user`.`name` = ? WHERE `user`.`id` = ?").
		WithArgs("name1", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	cnt, err := cli.Table(user).Update(user.Name.Eq("name1")).Where(user.ID.Eq(10)).Do(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

func Test_UpdatePayload_Simple(t *testing.T) {
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

	payload.Name = "name2"
	m.MockDB.ExpectExec("UPDATE `user` SET `user`.`name` = ? WHERE `user`.`id` = ?").
		WithArgs("name2", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	cnt, err := cli.Table(user).UpdatePayload(&payload).Where(user.ID.Eq(payload.ID)).Do(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}
