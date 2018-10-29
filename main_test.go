package orm

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var s *Session

var table = `
create table test (
	id int unsigned auto_increment,
	name varchar(64),
	type int,
	datetime datetime,
	primary key (id)
);`

type userT int64

func (u *userT) UnMarshalSQL(raw *ScanRow) error {
	var err error
	d, err := raw.ToInt64()
	if err != nil {
		return err
	}
	*u = userT(d)
	return nil
}

func (u *userT) MarshalSQL() (string, error) {
	return fmt.Sprintf("%#v", u), nil
}

const (
	Male   userT = 1
	FeMale userT = 2
)

type destT struct {
	ID       int64 `column:"omitempty"`
	Name     string
	Datetime time.Time
	UserType userT `column:"type"`
}

func initdata() {
	data1 := M{
		"name":     "archever",
		"type":     Male,
		"datetime": "2018-09-13 12:11:00",
	}
	data2 := destT{
		Name:     "archever2",
		UserType: FeMale,
		Datetime: time.Now(),
	}
	_, _, err := s.Table("test").Insert(data1).Do()
	_, _, err = s.Table("test").Insert(data2).Do()
	if err != nil {
		log.Panic(err)
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	s, err = Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/unittest")
	if err != nil {
		log.Panic(err)
	}
	_, _, err = s.Exec("drop table if exists test").Do()
	if err != nil {
		log.Panic(err)
	}
	_, _, err = s.Exec(table).Do()
	if err != nil {
		log.Panic(err)
	}
	initdata()
	m.Run()
	defer s.Exec("drop table if exists test").Do()
}
