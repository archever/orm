package e2etest

import (
	"context"
	"os"
	"testing"

	"github.com/archever/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()
	cli, err := orm.NewClient("mysql", os.Getenv("MYSQL_DSN"))
	assert.NoError(t, err)
	var payload userPayload
	err = cli.Table(user).Select().Where(user.ID.Eq(7)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	cli, err := orm.NewClient("mysql", os.Getenv("MYSQL_DSN"))
	assert.NoError(t, err)
	payload1 := userPayload{
		Name: "archever1",
	}
	payload2 := userPayload{
		Name: "archever2",
	}
	_, err = cli.Table(user).InsertPayload(&payload1, &payload2).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload1)
	t.Logf("%v", payload2)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	cli, err := orm.NewClient("mysql", os.Getenv("MYSQL_DSN"))
	assert.NoError(t, err)
	var payload userPayload
	err = cli.Table(user).Select().Where(user.ID.Eq(7)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)

	payload.Name = "archever1_1"
	_, err = cli.Table(user).UpdatePayload(&payload).Where(user.ID.Eq(7)).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	cli, err := orm.NewClient("mysql", os.Getenv("MYSQL_DSN"))
	assert.NoError(t, err)

	cnt, err := cli.Table(user).Delete().Where(user.ID.In(9, 10)).Do(ctx)
	assert.NoError(t, err)
	t.Log(cnt)
}
