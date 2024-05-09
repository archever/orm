package e2etest

import (
	"context"
	"os"
	"testing"

	"github.com/archever/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func getClient() *orm.Client {
	cli, err := orm.NewClient("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		panic(err)
	}
	return cli
}

func TestSelect(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userPayload
	err := cli.Table(user).Select().Where(user.ID.Eq(7)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestSelectMany(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	payloads := []*userPayload{}
	err := cli.Table(user).Select().Limit(10).FindPayload(ctx, &payloads)
	assert.NoError(t, err)
	t.Logf("%v", payloads[0])
	t.Logf("%v", payloads[1])
}

func TestInsertMany(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	payload1 := userPayload{
		Name: "archever1",
	}
	payload2 := userPayload{
		Name: "archever2",
	}
	_, err := cli.Table(user).InsertPayload(&payload1, &payload2).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload1)
	t.Logf("%v", payload2)
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	payload := teamPayload{
		Name: "team4",
	}
	_, err := cli.Table(user).InsertPayload(&payload).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestUpdatePayload(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userPayload
	err := cli.Table(user).Select().Where(user.ID.Eq(7)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)

	payload.Name = "archever1_1"
	_, err = cli.Table(user).UpdatePayload(&payload).Where(user.ID.Eq(7)).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	cli := getClient()

	cnt, err := cli.Table(user).Update(user.Name.Eq("name2")).Where(user.ID.Eq(7)).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", cnt)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	cli := getClient()

	cnt, err := cli.Table(user).Delete().Where(user.ID.In(9, 10)).Do(ctx)
	assert.NoError(t, err)
	t.Log(cnt)
}
