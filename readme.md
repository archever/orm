
# simple sql genarater

## TODO
* [ ] auto get table name
* [ ] support dest array type
* [ ] support mysql array/json type
* [ ] support pg

## usage

### init
install

```bash
go get "github.com/archever/orm/orm"
```

```golang
import "github.com/archever/orm/orm"

var err error
db, err = sql.Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/demo")
if err != nil {
    log.Panic(err)
}
// init a instance
o = orm.New(db)

// fetch data
var dest interface{}
err = o.Select("now() as now").One(&dest)
if err != nil {
    log.Panic(err)
}
log.Printf("now: %v", dest["now"].(time.Time))
```

### execute

```golang
var table = ` 
create table test (
	id int unsigned auto_increment,
	name varchar(64),
	primary key (id)
);`

_, _, err = o.Exec("drop table if exists test").Do()
_, _, err = o.Exec(table).Do()
```

### insert/replace

```golang
// insert via map
_, _, err = o.Table("test").Insert(orm.M{
	"name": "arhever",
}).Do()

// insert via struct
type testT struct {
	ID int64
	Name string
}

dataS1 := &testT{2, "Archever"}
dataS2 := &testT{3, "data2"}
dataS3 := &testT{4, "data3"}
dataSlice := []*testT{
	{10, "data10"},
	{11, "data11"},
	{12, "data12"},
	{13, "data13"},
}

_, _, err = o.Table("test").Insert(dataS1).Do()
_, _, err = o.Table("test").Insert(dataS2, dataS3).Do()
_, _, err = o.Table("test").Insert(dataSlice...).Do()
```

### Update

```golang
o.Table("test").Update(orm.M{"name": "archever"}).Where(orm.Equel("id", 10)).Do()
o.Table("test").Update(orm.M{"name": "archever"}).WhereS("id=?", 10).Do()
```

### Delete

```golang
o.Table("test").Delete().Where(orm.Equel("id", 10)).Do()
```

### select

```golang
// for get all, dest must be a multiply value or interface like
var dests interface{}
// or
var dests []interface{}
var dests []MyStruct
var dests []*MyStruct
var dests []map[string]interface{}
var dests []map[string]int64
// or
dests := []interface{}{}
dests := []MyStruct{}
dests := []*MyStruct{}
dests := []map[string]interface{}{}
dests := []map[string]int64{}

o.Table("test").Select().Get(&dests)

// for get one, dest must be a single value or interface like
// and it will limit 1 automaticly
var dest interface{}
// or
var dests MyStruct
var dests *MyStruct
var dests map[string]interface{}
var dests map[string]int64
// or
dests := MyStruct{}
dests := &MyStruct{}
dests := map[string]interface{}{}
dests := map[string]int64{}

o.Table("test").Select().One(&dest)
```

### custom struct serialize
there are tow interfaces to handler serialize, similar to encoding/json

* (*)UnMarshaler
* Marshaler

if the interface implemented, orm will use the func to handler sql data to struct and struct to sql data

```golang
type Date struct {
	value Time.time
}

// UnMarshalSQL sql field to struct
func (d *date) UnMarshalSQL(field []byte) error {
	v, err := time.ParseInLocation("2006-01-02 15:04:05", string(field), time.Now().Location())
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

// MarshalSQL struct to sql field
func (d date) MarshalSQL() (string, error) {
	return t.Time.Format("2006-01-02 15:04:05"), nil
}

type MyTable struct {
	ID int64
	CreateDate Date `column:"create_date"` 
}

var dest MyTale
data := &MyTable{
	CreateDate: time.Now(),
} 
o.Table("mytable").Insert(data).Do()
o.Table("mytable").Select().One(&dest)
```


## developer

### init a dev mysql	

`docker-compose up -d`

### test

`go test`

* init test data in main_test.go
* use `.Sql()` to check the generated sql and args
