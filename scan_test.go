package orm_test

import (
	"log"

	"github.com/archever/orm"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"testing"
)

var table2 = `create table test (data varchar(32) default null);`
var drop = `drop table if exists test;`

func TestScan(t *testing.T) {
	db, _ := sql.Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/test")
	rows, _ := db.Query("select * from test")
	res, err := orm.Scan(rows)
	log.Printf("err: %#v", err)
	for _, row := range res {
		for k, v := range row {
			log.Printf("%s: res: %#v", k, v)
			if v.Value != nil {
				log.Printf("%s: res: %#v", k, v.ToString())
			} else {
				log.Printf("%s: res is null", k)
			}
		}
	}
}

type DestT struct {
	Data *string
}

func TestDest(t *testing.T) {
	var res []*DestT
	s, _ := orm.Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/test")
	s.Table("test").Select().MustGet(&res)
	for _, i := range res {
		if i.Data == nil {
			log.Printf("res is null")
		} else {
			log.Printf("res: %s", *i.Data)
		}
	}
}
