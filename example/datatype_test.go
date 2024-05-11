package example

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/archever/orm"
	"github.com/stretchr/testify/assert"
)

type jsonStruct struct {
	Key   string
	Array []string
}

type allTypeSchema struct {
	ID      orm.Field[int64]
	Str     orm.Field[string]
	NilStr  orm.Field[*string]
	Time    orm.Field[time.Time]
	NilTime orm.Field[*time.Time]
	Json    orm.Field[*jsonStruct]
}

func (t *allTypeSchema) TableName() string {
	return "all_type"
}

type allTypePayload struct {
	orm.PayloadBase
	ID      int64
	Str     string
	NilStr  *string
	Time    time.Time
	NilTime *time.Time
	Json    *jsonStruct
}

func (t *allTypePayload) Bind() {
	t.PayloadBase.BindField(&t.ID, allType.ID)
	t.PayloadBase.BindField(&t.Str, allType.Str)
	t.PayloadBase.BindField(&t.NilStr, allType.NilStr)
	t.PayloadBase.BindField(&t.Time, allType.Time)
	t.PayloadBase.BindField(&t.NilTime, allType.NilTime)
	t.PayloadBase.BindField(&t.Json, allType.Json)
}

var allType = &allTypeSchema{
	ID:      orm.Field[int64]{Name: "id"},
	Str:     orm.Field[string]{Name: "str"},
	NilStr:  orm.Field[*string]{Name: "nil_str"},
	Time:    orm.Field[time.Time]{Name: "time"},
	NilTime: orm.Field[*time.Time]{Name: "nil_time"},
	Json:    orm.Field[*jsonStruct]{Name: "json"},
}

func TestSelect_AllType(t *testing.T) {
	ctx := context.Background()
	db, mockDB, _ := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	s := &orm.Client{
		DB: orm.NewDefaultExecutor(db),
	}
	mockDB.ExpectQuery("SELECT `id`, `str`, `nil_str`, `time`, `nil_time`, `json` FROM `all_type` WHERE `id`=? LIMIT ?").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "str", "nil_str", "time", "nil_time", "json"}).
				AddRow(10, "archever", sql.NullString{}, time.Now(), sql.NullTime{}, sql.NullString{String: `{"Key": "val"}`}),
		)

	var dst allTypePayload
	err := s.Table(user).Select().Where(allType.ID.Eq(1)).TakePayload(ctx, &dst)
	assert.NoError(t, err)
	t.Logf("%v", dst.ID)
	t.Logf("%v", dst.Str)
	t.Logf("%v", dst.NilStr)
	t.Logf("%v", dst.Time)
	t.Logf("%v", dst.NilTime)
	t.Logf("%v", dst.Json)
}
