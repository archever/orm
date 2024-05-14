package tests

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
)

type mockInc struct {
	MockDB sqlmock.Sqlmock

	DB *sql.DB
}

func (m *mockInc) Build() (*mockInc, error) {
	var err error
	if m.DB == nil {
		m.DB, m.MockDB, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *mockInc) MustBuild() *mockInc {
	inc, err := m.Build()
	if err != nil {
		panic(err)
	}
	return inc
}

func getClient(m *mockInc) *orm.Client {
	return &orm.Client{
		DB: orm.NewDefaultExecutor(m.DB),
	}
}
