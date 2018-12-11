package orm_test

import (
	"log"
	orm "ormV2"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var s *orm.Session
var dbURI = "root:zxcvbnm@tcp(127.0.0.1:3306)/test"
var dropSQL = `drop table if exists T`
var tSQL = `create table T (
	id int auto_increment,
	name varchar(32) default null,
	namenil varchar(32) default null,
	type int default 0,
	data blob,
	content text,
	array varchar(256) default "",
	value float default null,
	amount decimal(20, 4) default 0,
	birth date default null,
	create_at datetime default current_timestamp,
	update_at datetime default null,
	primary key (id)
)`

type MyType int
type MyArray []string

func (a *MyArray) MarshalSQL() (string, error) {
	if a == nil {
		return "", orm.ErrNull
	}
	return strings.Join(*a, ","), nil
}

func (a *MyArray) UnmarshalSQL(row *orm.ScanRow) error {
	if row == nil || !row.Valid {
		*a = []string{}
		return nil
	}
	*a = strings.Split(row.String(), ",")
	return nil
}

type ModelT struct {
	ID       int64      `column:"id,omitempty"`
	Name     *string    `column:"name"`
	Type     MyType     `column:"type"`
	Data     []byte     `column:"data"`
	Content  string     `column:"content"`
	Value    float64    `column:"value"`
	Amount   float64    `column:"amount"`
	Array    *MyArray   `column:"array"`
	Array2   MyArray    `column:"array"`
	Birth    *time.Time `column:"birth"`
	CreateAt *time.Time `column:"create_at,omitempty"`
	UpdateAt *time.Time `column:"update_at,omitempty"`
}

func init() {
	var err error
	s, err = orm.Open("mysql", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	s.Exec(dropSQL).MustDo()
	s.Exec(tSQL).MustDo()
	s.Table("T").Insert(orm.M{
		"type":   1,
		"array":  "",
		"amount": 1.23,
		"birth":  "2000-01-01",
	}).MustDo()
	s.Table("T").Insert(orm.M{
		"name":   "name2",
		"type":   2,
		"data":   []byte("test"),
		"array":  "a,b,c,d,e",
		"amount": 4,
		"birth":  "2000-01-01",
	}).MustDo()
}

func AssertSelect(t *testing.T, res []orm.M) {
	birth, _ := time.ParseInLocation("2006-01-02", "2000-01-01", time.Now().Location())
	assert.Equal(t, 1, res[0]["id"])
	assert.Equal(t, 1, res[0]["type"])
	assert.Equal(t, 1.23, res[0]["amount"])
	assert.Equal(t, &birth, res[0]["birth"])
	assert.Nil(t, res[0]["name"])
	assert.Nil(t, res[0]["content"])
	assert.Nil(t, res[0]["data"])
	assert.Nil(t, res[0]["value"])
	assert.NotNil(t, res[0]["create_at"])
	assert.Nil(t, res[0]["update_at"])

	assert.Equal(t, 2, res[1]["id"])
	assert.Equal(t, 2, res[1]["type"])
	assert.Equal(t, 4.0, res[1]["amount"])
	assert.Equal(t, "a,b,c,d,e", res[1]["array"])
	assert.Equal(t, &birth, res[1]["birth"])
	assert.Equal(t, "name2", res[1]["name"])
	assert.Nil(t, res[1]["content"])
	assert.Equal(t, []byte("test"), res[1]["data"])
}

func AssertSelectModel(t *testing.T, res []*ModelT) {
	birth, _ := time.ParseInLocation("2006-01-02", "2000-01-01", time.Now().Location())
	assert.Equal(t, int64(1), res[0].ID)
	assert.Equal(t, MyType(1), res[0].Type)
	assert.Equal(t, 1.23, res[0].Amount)
	assert.Equal(t, &birth, res[0].Birth)
	assert.Nil(t, res[0].Name)
	assert.Equal(t, "", res[0].Content)
	assert.Nil(t, res[0].Data)
	assert.Equal(t, 0.0, res[0].Value)
	assert.NotNil(t, res[0].CreateAt)
	assert.Nil(t, res[0].UpdateAt)

	assert.Equal(t, int64(2), res[1].ID)
	assert.Equal(t, MyType(2), res[1].Type)
	assert.Equal(t, 4.0, res[1].Amount)
	assert.Equal(t, &MyArray{"a", "b", "c", "d", "e"}, res[1].Array)
	assert.Equal(t, MyArray{"a", "b", "c", "d", "e"}, res[1].Array2)
	assert.Equal(t, &birth, res[1].Birth)
	assert.Equal(t, "name2", *res[1].Name)
	assert.Equal(t, "", res[1].Content)
	assert.Equal(t, []byte("test"), res[1].Data)
}

func TestSelectMulti(t *testing.T) {
	var res1 []orm.M
	s.Table("T").Select().MustGet(&res1)
	AssertSelect(t, res1)

	var res2 interface{}
	s.Table("T").Select().MustGet(&res2)
	res := res2.([]orm.M)
	AssertSelect(t, res)

	var res3 [4]orm.M
	s.Table("T").Select().MustGet(&res3)
	assert.Nil(t, res3[2])
	assert.Nil(t, res3[3])
	AssertSelect(t, res3[:])

	var res5 []*ModelT
	s.Table("T").Select().MustGet(&res5)
	AssertSelectModel(t, res5)

	var res6 [4]*ModelT
	s.Table("T").Select().MustGet(&res6)
	AssertSelectModel(t, res6[:])

	var res7 []ModelT
	s.Table("T").Select().MustGet(&res7)
	var res7_ []*ModelT
	for i := range res7 {
		res7_ = append(res7_, &res7[i])
	}
	AssertSelectModel(t, res7_)

	var res8 [4]ModelT
	s.Table("T").Select().MustGet(&res8)
	var res8_ []*ModelT
	for i := range res8 {
		res8_ = append(res8_, &res8[i])
	}
	AssertSelectModel(t, res8_)
}

func AssertSelectOne(t *testing.T, res orm.M) {
	birth, _ := time.ParseInLocation("2006-01-02", "2000-01-01", time.Now().Location())
	assert.Equal(t, 2, res["id"])
	assert.Equal(t, 2, res["type"])
	assert.Equal(t, 4.0, res["amount"])
	assert.Equal(t, "a,b,c,d,e", res["array"])
	assert.Equal(t, &birth, res["birth"])
	assert.Equal(t, "name2", res["name"])
	assert.Nil(t, res["content"])
	assert.Equal(t, []byte("test"), res["data"])
	assert.Nil(t, res["value"])
	assert.NotNil(t, res["create_at"])
	assert.Nil(t, res["update_at"])
}

func AssertSelectModelOne(t *testing.T, res *ModelT) {
	birth, _ := time.ParseInLocation("2006-01-02", "2000-01-01", time.Now().Location())
	assert.Equal(t, int64(2), res.ID)
	assert.Equal(t, MyType(2), res.Type)
	assert.Equal(t, 4.0, res.Amount)
	assert.Equal(t, &MyArray{"a", "b", "c", "d", "e"}, res.Array)
	assert.Equal(t, MyArray{"a", "b", "c", "d", "e"}, res.Array2)
	assert.Equal(t, &birth, res.Birth)
	assert.Equal(t, "name2", *res.Name)
	assert.Equal(t, "", res.Content)
	assert.Equal(t, []byte("test"), res.Data)
	assert.Equal(t, 0.0, res.Value)
	assert.NotNil(t, res.CreateAt)
	assert.Nil(t, res.UpdateAt)
}

func TestSelectOne(t *testing.T) {
	var res1 orm.M
	s.Table("T").Select().Where("id=?", 2).MustOne(&res1)
	AssertSelectOne(t, res1)

	var res2 interface{}
	s.Table("T").Select().Where("id=?", 2).MustOne(&res2)
	AssertSelectOne(t, res2.(orm.M))

	var res3 *orm.M
	s.Table("T").Select().Where("id=?", 2).MustOne(&res3)
	AssertSelectOne(t, *res3)

	var res4 *ModelT
	s.Table("T").Select().Where("id=?", 2).MustOne(&res4)
	AssertSelectModelOne(t, res4)

	var res5 ModelT
	s.Table("T").Select().Where("id=?", 2).MustOne(&res5)
	AssertSelectModelOne(t, &res5)
}

func TestInsertOne(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")
	var name *string
	name2 := "name3"
	name = &name2
	var row = &ModelT{
		ID:      3,
		Name:    name,
		Type:    MyType(23),
		Data:    []byte("abc"),
		Content: "content",
		Value:   12.12,
		Array:   &MyArray{"a", "b"},
		Array2:  MyArray{"a", "b"},
		Birth:   &birth,
	}
	i, c := s.Table("T").Insert(row).MustDo()
	assert.Equal(t, int64(3), i)
	assert.Equal(t, int64(1), c)
	var res *ModelT
	s.Table("T").Select().Where("id=?", 3).MustOne(&res)
	assert.Equal(t, row.ID, res.ID)
	assert.Equal(t, row.Name, res.Name)
	assert.Equal(t, row.Type, res.Type)
	assert.Equal(t, row.Data, res.Data)
	assert.Equal(t, row.Content, res.Content)
	assert.Equal(t, row.Value, res.Value)
	assert.Equal(t, row.Amount, res.Amount)
	assert.Equal(t, row.Array, res.Array)
	assert.Equal(t, row.Array2, res.Array2)
	assert.Equal(t, row.Birth.Format("2006-01-02"), res.Birth.Format("2006-01-02"))
	assert.NotNil(t, res.CreateAt)
	assert.Nil(t, res.UpdateAt)
}

func TestInsertOneMap(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")
	birth2, _ := time.ParseInLocation("2006-01-02", "2000-01-01", time.Now().Location())
	var row = orm.M{
		"id":      4,
		"name":    "name4",
		"type":    MyType(23),
		"data":    []byte("abc"),
		"content": "content",
		"value":   12.12,
		"array":   &MyArray{"a", "b"},
		"birth":   &birth,
	}
	i, c := s.Table("T").Insert(row).MustDo()
	assert.Equal(t, int64(4), i)
	assert.Equal(t, int64(1), c)
	var res orm.M
	s.Table("T").Select().Where("id=?", 4).MustOne(&res)
	assert.Equal(t, 4, res["id"])
	assert.Equal(t, "name4", res["name"])
	assert.Equal(t, 23, res["type"])
	assert.Equal(t, []byte("abc"), res["data"])
	assert.Equal(t, "content", res["content"])
	assert.Equal(t, 12.12, res["value"])
	assert.Equal(t, 0.0, res["amount"])
	assert.Equal(t, "a,b", res["array"])
	assert.Equal(t, &birth2, res["birth"])
	assert.NotNil(t, res["create_at"])
	assert.Nil(t, res["update_at"])
}

func TestInsertManyMap(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")
	var row = orm.M{
		"name":    "name4",
		"type":    MyType(23),
		"data":    []byte("abc"),
		"content": "content",
		"value":   12.12,
		"array":   &MyArray{"a", "b"},
		"birth":   &birth,
	}
	_, c := s.Table("T").InsertMany([]orm.M{row, row}).MustDo()
	assert.Equal(t, int64(2), c)
}

func TestInsertManyStruct(t *testing.T) {
	birth, _ := time.Parse("2006-01-02", "2000-01-01")
	var name *string
	name2 := "name3"
	name = &name2
	var row = &ModelT{
		Name:    name,
		Type:    MyType(23),
		Data:    []byte("abc"),
		Content: "content",
		Value:   12.12,
		Array:   &MyArray{"a", "b"},
		Array2:  MyArray{"a", "b"},
		Birth:   &birth,
	}
	_, c := s.Table("T").InsertMany([]*ModelT{row, row}).MustDo()
	assert.Equal(t, int64(2), c)
}

func TestUpdate(t *testing.T) {
	_, c := s.Table("T").Update(orm.M{
		"array": &MyArray{"x", "y"},
		"name":  interface{}(nil),
	}).Where("id=?", 2).MustDo()
	assert.Equal(t, int64(1), c)
	var res orm.M
	s.Table("T").Select("array", "name").Where("id=?", 2).MustOne(&res)
	assert.Equal(t, "x,y", res["array"])
	assert.Nil(t, res["name"])
}

func TestDelete(t *testing.T) {
	_, c := s.Table("T").Delete().Where("id=?", 2).MustDo()
	assert.Equal(t, int64(1), c)
	cnt, _ := s.Table("T").Select().Where("id=?", 2).Count()
	assert.Equal(t, int64(0), cnt)
}
