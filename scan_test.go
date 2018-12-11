package orm

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var table = `
create table test (
	id int auto_increment,
	name varchar(32) default null,
	data blob,
	content text,
	value float default null,
	birth date default null,
	create_at datetime default current_timestamp,
	update_at datetime default null,
	primary key (id)
);`

var drop = `drop table if exists test;`
var db *sql.DB
var dbURI = "root:zxcvbnm@tcp(127.0.0.1:3306)/test"

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(drop)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(table)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`insert into test(name, data, birth) values("archever", "test", "2016-01-20")`)
	if err != nil {
		log.Fatal(err)
	}
}

func TestScan(t *testing.T) {
	rows, err := db.Query("select * from test")
	assert.NoError(t, err)
	res, err := Scan(rows)
	assert.NoError(t, err)
	for _, row := range res {
		for k, v := range row {
			switch k {
			case "create_at":
				assert.NotEmpty(t, v.Value)
			case "update_at":
				assert.Equal(t, false, v.Valid)
			case "id":
				assert.Equal(t, "1", v.Value)
			case "name":
				assert.Equal(t, "archever", v.Value)
			case "data":
				assert.Equal(t, true, v.Valid)
			case "content":
				assert.Equal(t, false, v.Valid)
			case "value":
				assert.Equal(t, false, v.Valid)
			case "birth":
				assert.Equal(t, "2016-01-20", v.Value)
			}
		}
	}
}
