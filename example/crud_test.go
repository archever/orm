package example

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	// SELECT `id`, `name` FROM `user` WHERE `id`=? LIMIT ?
	mockDB.ExpectQuery("SELECT `user`.`id`, `user`.`name` FROM `user` WHERE `user`.`id`=? LIMIT ?").WillReturnRows(
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
		DB: orm.NewDefaultExecutor(db),
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
	_, err = s.Table(user).UpdatePayload(&payload).Where(user.ID.Eq(10)).Do(ctx)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	mockDB.ExpectExec("DELETE FROM `user` WHERE `id`=?").
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// var payload userPayload
	_, err := s.Table(user).Delete().Where(user.ID.Eq(10)).Do(ctx)
	assert.NoError(t, err)
}

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	row1 := &userPayload{
		ID:   10,
		Name: "archever",
	}
	mockDB.ExpectExec("INSERT INTO `user` (`id`,`name`) VALUES(?,?)").
		WithArgs(10, "archever").
		WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := s.Table(user).InsertPayload(row1).Do(ctx)
	assert.NoError(t, err)
}

func TestInsertMany(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	row1 := &userPayload{
		ID:   10,
		Name: "archever",
	}
	row2 := &userPayload{
		ID:   12,
		Name: "archever2",
	}
	mockDB.ExpectExec("INSERT INTO `user` (`id`,`name`) VALUES(?,?),(?,?)").
		WithArgs(10, "archever", 12, "archever2").
		WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := s.Table(user).InsertPayload(row1, row2).Do(ctx)
	assert.NoError(t, err)
}
