// test custom struct marshal unmarshal
package orm

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MyTime struct {
	Time time.Time
}

var _ Marshaler = &MyTime{}
var _ UnMarshaler = &MyTime{}

func (t *MyTime) UnMarshalSQL(field []byte) error {
	v, err := time.ParseInLocation("2006-01-02 15:04:05", string(field), time.Now().Location())
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

func (t MyTime) MarshalSQL() (string, error) {
	return t.Time.Format("2006-01-02 15:04:05"), nil
}

type destCT struct {
	ID         int64
	Name       string
	Now        MyTime
	CreateTime time.Time
}

func TestMarshal(t *testing.T) {
	data := &destCT{
		300, "test300", MyTime{time.Date(2018, time.September, 13, 12, 0, 0, 0, time.Now().Location())},
		time.Date(2018, time.September, 13, 12, 1, 0, 0, time.Now().Location()),
	}
	sql, args, err := o.Table("test").Insert(data).SQL()
	assert.NoError(t, err)
	// the generate field order seems random
	assert.Contains(t, sql, "id")
	assert.Contains(t, sql, "name")
	assert.Contains(t, sql, "now")
	assert.Contains(t, sql, "createtime")
	assert.Contains(t, args, int64(300))
	assert.Contains(t, args, "test300")
	assert.Contains(t, args, "2018-09-13 12:00:00")
	assert.Contains(t, args, time.Date(2018, time.September, 13, 12, 1, 0, 0, time.Now().Location()))
}

func TestUnMarshal(t *testing.T) {
	var dest destCT
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.IsType(t, MyTime{}, dest.Now)
	assert.IsType(t, time.Time{}, dest.CreateTime)
}
