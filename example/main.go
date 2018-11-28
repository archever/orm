package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/archever/orm"
	_ "github.com/go-sql-driver/mysql"
)

var s *orm.Session

type gender int64

type hobby []string

func (h hobby) MarshalSQL() (string, error) {
	return strings.Join(h, ","), nil
}

func (h *hobby) UnmarshalSQL(raw *orm.ScanRow) error {
	*h = strings.Split(raw.ToString(), ",")
	return nil
}

type Date struct {
	year  int64
	month time.Month
	day   int64
}

func (d Date) MarshalSQL() (string, error) {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day), nil
}
func (d *Date) UnmarshalSQL(raw *orm.ScanRow) error {
	var month int
	_, err := fmt.Sscanf(raw.ToString(), "%04d-%02d-%02d", &d.year, &month, &d.day)
	d.month = time.Month(month)
	return err
}

type TestT struct {
	ID         int64 `column:"id,omitempty"` // use omitempty to avode insert zero field
	Name       string
	Gender     gender
	Birth      Date
	Hobby      hobby
	CreateTime *time.Time `column:"create_time,omitempty"`
}

const (
	Male   gender = 1
	FeMale gender = 2
)

const t = `
create table t(
	id int auto_increment,
	name varchar(64),
	gender int,
	birth date,
	create_time datetime default CURRENT_TIMESTAMP,
	hobby varchar(256),
	primary key (id)
)
`

func init() {
	// run a mysql instance
	// docker run --rm -e MYSQL_ROOT_PASSWORD=zxcvbnm -e MYSQL_DATABASE=unittest -d mysql
	s, _ = orm.Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/unittest")
}

func tearDown() {
	s.Exec("drop table if exists t").Do()
}

func createTable() {
	var err error
	_, _, err = s.Exec("drop table if exists t").Do()
	if err != nil {
		log.Panic(err)
	}
	_, _, err = s.Exec(t).Do()
	if err != nil {
		log.Panic(err)
	}
}

func insertData() {
	row1 := &TestT{
		Name:   "archever1",
		Gender: Male,
		Birth:  Date{1992, time.March, 13},
		Hobby:  hobby{"coding,reading,hiking"},
	}
	row2 := orm.f.M{
		"name":   "archever2",
		"gender": Male,
		"birth":  Date{1992, time.March, 14},
		"hobby":  hobby{"coding,reading"},
	}
	rows := []TestT{
		{Name: "archever3", Gender: Male, Birth: Date{1992, time.March, 15}, Hobby: hobby{"music"}},
		{Name: "archever4", Gender: FeMale, Birth: Date{1992, time.March, 16}},
		{Name: "archever5", Gender: FeMale, Birth: Date{1992, time.March, 17}},
	}
	id, c, err := s.Table("t").Insert(row1).Do()
	log.Printf("%v, %v, %v", id, c, err)
	id, c, err = s.Table("t").Insert(row2).Do()
	log.Printf("%v, %v, %v", id, c, err)
	id, c, err = s.Table("t").InsertMany(rows).Do()
	log.Printf("%v, %v, %v", id, c, err)
}

func selectData() {
	var res TestT
	// err := s.Table("t").Select().Filter(f.Equal("id", 2)).One(&res)
	err := s.Table("t").Select().Where("id=?", 2).One(&res)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res: %#v", res)
	fmt.Printf("res: %s", res.CreateTime)
}

func main() {
	defer tearDown()
	createTable()
	insertData()
	selectData()
}
