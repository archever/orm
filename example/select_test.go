package example

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
	"github.com/stretchr/testify/assert"
)

var user = &userSchema{
	ID:   orm.Field[int64]{Name: "id"},
	Name: orm.Field[string]{Name: "name"},
}

type userSchema struct {
	ID   orm.Field[int64]
	Name orm.Field[string]
}

type userPayload struct {
	orm.PayloadBase
	ID   int64
	Name string
}

func (p *userPayload) Bind() {
	p.PayloadBase.BindField(user.ID.WithRef(&p.ID))
	p.PayloadBase.BindField(user.Name.WithRef(&p.Name))
}

func (s *userSchema) TableName() string {
	return "user"
}

func TestSelect(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: db,
	}
	mockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `id`=? LIMIT ?").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
	)
	var payload userPayload
	err := s.Table(user).Select().Where(user.ID.Eq(10)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: db,
	}
	mockDB.ExpectQuery("SELECT `id`, `name` FROM `user` WHERE `id`=? LIMIT ?").
		WithArgs(10, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
		)
	mockDB.ExpectExec("UPDATE `user` SET `name`=? WHERE `id`=?").
		WithArgs("name2", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	var payload userPayload
	err := s.Table(user).Select().Where(user.ID.Eq(10)).TakePayload(ctx, &payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)

	payload.Name = "name2"
	err = s.Table(user).UpdatePayload(&payload).Where(user.ID.Eq(10)).Do(ctx)
	assert.NoError(t, err)
}
