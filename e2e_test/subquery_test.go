package e2etest

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSubQuery(t *testing.T) {
	ctx := context.Background()
	cli := getClient()
	var payload userPayload
	subQuery := cli.Table(user).Select(&user.ID).Where(user.Name.Eq("name")).SubQuery()
	err := cli.Table(user).Select().
		Where(user.ID.EqQuery(subQuery)).
		Limit(1).
		TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

// TODO: 支持 from subquery, columns in subquery
