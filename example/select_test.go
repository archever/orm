package example

import (
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
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: db,
	}
	mockDB.ExpectQuery("select `id`, `name` from `user` where `id`=? limit ?").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
	)
	var payload userPayload
	err := s.Table(user).Select().Where(user.ID.Eq(10)).Get(&payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)
}

func TestUpdate(t *testing.T) {
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: db,
	}
	mockDB.ExpectQuery("select `id`, `name` from `user` where `id`=? limit ?").
		WithArgs(10, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(10, "archever"),
		)
	mockDB.ExpectExec("update `user` set `name`=? where `id`=?").
		WithArgs("name2", 10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	var payload userPayload
	err := s.Table(user).Select().Where(user.ID.Eq(10)).Get(&payload)
	assert.NoError(t, err)
	t.Logf("%v", payload)

	payload.Name = "name2"
	_, _, err = s.Table(user).Update(&payload).Where(user.ID.Eq(10)).Do()
	assert.NoError(t, err)
}
