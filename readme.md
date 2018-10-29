# orm

## usage

```golang
// 1. init a session
db := sql.Open(...)
s := orm.NewSession(db)

// 2. make a varible
var dest orm.M

// 3. query
s.Exec("select now()").One(&dest)
```

### select

```golang
var dest orm.M
var dests []orm.M
rowID, rowCount, err := s.Exec("select 1").Do()
// select 1;

err := s.Table("t").Select().One(&dest)
// select * from t limit 1;

err := s.Table("t").Select("name", "id").Where("a", 1).Limit(1).Get(&dests)
// select name id from t where a=? limit ?, [1, 1]

err := s.Table("t").Select().Filter(f.Equel("a", 1), f.Gte("b", 2)).One(&dest)
// select * from t where a=? and b>=? limit ?, [1,2,1]

sql, args := s.Table("t").Select().SQL()
// it won't execute sql but return the sql string and sql arguements
```

### use a transction
```golang
tx, err := s.Begin()
// use the tx as well as s

tx.Table("t").Select().Do()

tx.Commit()

tx, err := s.Begin()
// the next time`s.Begin()` get a new tx
// if you forget to commit, the next time get a tx, it will rollback the priviouse one
```

### Insert

```golang
// add a `omitempty` tag to avoid insert the column while value is empty
type TestTable struct {
    ID int64    `column:"id,omitempty"`
    Name string `column:"name"`
}
row1 := &TestTable{
    Name: "archever"
}
row2 := orm.M{
    "name": "archever",
}
s.Table("t").Insert(row1).Do()
s.Table("t").Insert(row2).Do()
```

### costom sql marshaler
```golang
type MyType struct {
    ...
}

// or no pointer reciver
func (t *MyType) MarshalSQL () (string, error) {
    ...
}

func (t *MyType) UnMarshalSQL (*orm.ScanRow) (error) {
    ...
}

// use myType in select or fiter or insert sqls
s.Table("t").Insert(orm.M{
    "column": MyType
}).Do()
s.Table("t").Select().Where("a=?", MyType{...})
```

### some filters
```golang
cond1 := f.Equel("a", 1)
cond2 := f.Lt("b", 1)
cond := f.Or(cond1, cond2)

// or 
conds := []*f.FilterItem{}
conds = append(conds, cond1, cond2)
s.Table(...).Filter(...conds)
```

### more
more usage see [example](./example/)
