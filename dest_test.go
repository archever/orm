package orm

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectRowsStrcut(t *testing.T) {
	dest := []destT{}
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
	assert.IsType(t, Male, dest[0].UserType)
}

func TestSelectRowsStrcutPointer(t *testing.T) {
	dest := []*destT{}
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
	assert.IsType(t, Male, dest[0].UserType)
}

func TestSelectRowsStrcutInterface(t *testing.T) {
	var dest []destT
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
	assert.IsType(t, Male, dest[0].UserType)
}

func TestSelectRowsStrcutInterfacePointer(t *testing.T) {
	var dest []*destT
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0].ID)
	assert.Equal(t, "archever", dest[0].Name)
	assert.IsType(t, Male, dest[0].UserType)
}

func TestSelectRowsStrcutOne(t *testing.T) {
	dest := destT{}
	err := s.Table("test").Select().One(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest.ID)
	assert.Equal(t, "archever", dest.Name)
	assert.IsType(t, Male, dest.UserType)
}

func TestSelectRowsStrcutOneInterface(t *testing.T) {
	var dest destT
	err := s.Table("test").Select().One(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest.ID)
	assert.Equal(t, "archever", dest.Name)
	assert.IsType(t, Male, dest.UserType)
}

func TestSelectRowsMap(t *testing.T) {
	dest := []M{}
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0]["id"])
	assert.Equal(t, "archever", dest[0]["name"])
}

func TestSelectRowsMapNil(t *testing.T) {
	var dest []M
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest[0]["id"])
	assert.Equal(t, "archever", dest[0]["name"])
}

func TestSelectRowsMapOne(t *testing.T) {
	dest := M{}
	err := s.Table("test").Select().One(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest["id"])
	assert.Equal(t, "archever", dest["name"])
}

func TestSelectRowsMapOneNil(t *testing.T) {
	var dest M
	err := s.Table("test").Select().One(&dest)
	log.Printf("res: %#v", dest)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), dest["id"])
	assert.Equal(t, "archever", dest["name"])
}

func TestSelectRowsInterface(t *testing.T) {
	var dest interface{}
	err := s.Table("test").Select().Get(&dest)
	log.Printf("res: %#v", dest)
	destv, ok := dest.([]map[string]interface{})
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), destv[0]["id"])
	assert.Equal(t, "archever", destv[0]["name"])
}

func TestSelectRowsOneInterface(t *testing.T) {
	var dest interface{}
	err := s.Table("test").Select().One(&dest)
	log.Printf("res: %#v", dest)
	destv, ok := dest.(map[string]interface{})
	assert.NoError(t, err)
	assert.Equal(t, true, ok)
	assert.Equal(t, int64(1), destv["id"])
	assert.Equal(t, "archever", destv["name"])
}
