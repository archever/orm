package tests

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Insert_Simple(t *testing.T) {
	ctx := context.Background()
	m := (&mockInc{}).MustBuild()
	cli := getClient(m)
	m.MockDB.ExpectExec("INSERT INTO `user` (`name`) VALUES(?),(?)").
		WithArgs("name1", "name2").
		WillReturnResult(sqlmock.NewResult(1, 2))
	payload1 := userPayload{Name: "name1"}
	payload2 := userPayload{Name: "name2"}
	cnt, err := cli.Table(user).InsertPayload(&payload1, &payload2).Do(ctx)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, cnt)
	assert.EqualValues(t, 1, payload1.ID)
	assert.EqualValues(t, 2, payload2.ID)
}
