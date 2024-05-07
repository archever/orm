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
	err = cli.Table(user).Select().Where(user.ID.Eq(10)).TakePayload(ctx, &payload)
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
	err = cli.Table(user).InsertPayload(&payload1, &payload2).Do(ctx)
	assert.NoError(t, err)
	t.Logf("%v", payload1)
	t.Logf("%v", payload2)
}
