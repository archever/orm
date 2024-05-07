package example

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	mockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `id`=? LIMIT ?").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
	)
	var payload userWithTeamPayload
	err := s.Table(user).Select().Where(user.ID.Eq(10)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}
