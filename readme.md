
# orm for mysql

## main principles
* similar to `encoding/json` usage
* `nil` equal mysql `NULL`
* dest `nil` equal `NOT FUND`
* automatic dest to `map[string]interface` or `struct`
* custom marshaler types
* convenient to generate where cases

## install
`go get github.com/archever/orm`

## test
```shell
docker run --rm -e MYSQL_ROOT_PASSWORD=zxcvbnm -e MYSQL_DATABASE=test -p 3306:3306 mysql
go test ./...
```

## usage
```golang
// init a instance
s := orm.Open("mysql", "root:zxcvbnm@tcp(127.0.0.1:3306)/test")
```
you can store your query result into a `map[string]interface` or `struct`

```golang
// first define a variable to store
var res []orm.M

// make a sql statement
stmt := s.Table(<table name>)
.Select(<select fields>) // action update or delete or insert ...
.Filter(<filter item>)
.OrderBy(<order field>)
.GroupBy(<group fields>)
.Limit(<limit>)
.Offest(<limit>)
.Page(<page>, <limit>)

// execute it and fetch
stmt.MustGet(&res)
// count the result num without limit
cnt, err := stmt.Count()
```

### costom sql marshaler
```golang
type MyType struct {
    ...
}

func (t *MyType) MarshalSQL () (string, error) {
    ...
}

func (t *MyType) UnmarshalSQL (*orm.ScanRow) (error) {
    ...
}

// use MyType in select or fiter or insert sqls
s.Table("t").Insert(orm.M{
    "column": MyType{}
}).Do()
s.Table("t").Select().Where("a=?", MyType{...})
```

### use a transaction
```golang
// start a transaction
tx := s.MustBegin()

// use the tx as well as s
tx.Table("t").Select().Do()

// commit or rollback
tx.Commit()
tx.RollBack()
```

### set logger
here are two loggers, you can modify them, if you need
* orm.Log
* orm.Debug

```golang
orm.Debug.SetOutput(ioutil.Discard)
```
