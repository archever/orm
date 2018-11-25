package orm

import (
	"log"
	"testing"
	"time"

	"github.com/archever/orm/f"
	_ "github.com/go-sql-driver/mysql"
)

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

const (
	Male   userT = 1
	FeMale userT = 2
)

type destT struct {
	ID       int64      `column:"omitempty"`
	Name     string     `column:"name"`
	Datetime *time.Time `column:"datetime"`
	UserType userT      `column:"type"`
}

func initdata() {
	data1 := f.M{
		"name":     "archever",
		"type":     Male,
		"datetime": "2018-09-13 12:11:00",
	}
	_, _, err := s.Table("test").Insert(data1).Do()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	s, err = Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/unittest")
	if err != nil {
		panic(err)
	}
	_, _, err = s.Exec("drop table if exists test").Do()
	if err != nil {
		panic(err)
	}
	_, _, err = s.Exec(table).Do()
	if err != nil {
		panic(err)
	}
	initdata()
	m.Run()
	defer s.Exec("drop table if exists test").Do()
}
