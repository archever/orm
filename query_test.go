// test dest

package orm

import (
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

type destT struct {
	ID   int64
	Name string
}

func TestSelectRowsStrcut(t *testing.T) {
	dest := []destT{}
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
}

func TestSelectRowsStrcutPointer(t *testing.T) {
	dest := []*destT{}
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
}

func TestSelectRowsStrcutInterface(t *testing.T) {
	var dest []destT
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
}

func TestSelectRowsStrcutInterfacePointer(t *testing.T) {
	var dest []*destT
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
}

func TestSelectRowsStrcutOne(t *testing.T) {
	dest := destT{}
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest.ID)
	assert.Equal(t, "archever", dest.Name)
}

func TestSelectRowsStrcutOneInterface(t *testing.T) {
	var dest destT
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest.ID)
	assert.Equal(t, "archever", dest.Name)
}

func TestSelectRowsMap(t *testing.T) {
	dest := []M{}
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0]["id"])
	assert.Equal(t, "archever", dest[0]["name"])
}

func TestSelectRowsMapNil(t *testing.T) {
	var dest []M
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0]["id"])
	assert.Equal(t, "archever", dest[0]["name"])
}

func TestSelectRowsMapOne(t *testing.T) {
	dest := M{}
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest["id"])
	assert.Equal(t, "archever", dest["name"])
}

func TestSelectRowsMapOneNil(t *testing.T) {
	var dest M
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest["id"])
	assert.Equal(t, "archever", dest["name"])
}

func TestSelectRowsInterface(t *testing.T) {
	var dest interface{}
	err := o.Table("test").Select().Get(&dest)
	log.Printf("res: %v", dest)
	destv, ok := dest.([]map[string]interface{})
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), destv[0]["id"])
	assert.Equal(t, "archever", destv[0]["name"])
}

func TestSelectRowsOneInterface(t *testing.T) {
	var dest interface{}
	err := o.Table("test").Select().One(&dest)
	log.Printf("res: %v", dest)
	destv, ok := dest.(map[string]interface{})
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), destv["id"])
	assert.Equal(t, "archever", destv["name"])
}
