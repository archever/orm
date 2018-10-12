package orm

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var o OrmI

var table = `
create table test (
	id int unsigned auto_increment,
	name varchar(64),
	now datetime,
	createtime datetime,
	primary key (id)
);`

func initdata() {
	data1 := M{
		"name":       "archever",
		"now":        "2018-09-13 12:01:00",
		"createtime": "2018-09-13 12:11:00",
	}
	data2 := M{
		"name":       "archever2",
		"now":        "2018-09-13 12:02:00",
		"createtime": "2018-09-13 12:12:00",
	}
	_, _, err := o.Table("test").Insert(data1).Do()
	if err != nil {
		log.Panic(err)
	}
	_, _, err = o.Table("test").Insert(data2).Do()
	if err != nil {
		log.Panic(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/unittest")
	if err != nil {
		log.Panic(err)
	}
	o = New(db)
	_, _, err = o.Exec("drop table if exists test").Do()
	if err != nil {
		log.Panic(err)
	}
	_, _, err = o.Exec(table).Do()
	if err != nil {
		log.Panic(err)
	}
	initdata()
	m.Run()
	defer o.Exec("drop table if exists test").Do()
}
