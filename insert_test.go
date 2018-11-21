package orm

import (
	"log"
	"testing"

	"github.com/archever/orm/f"

	"github.com/stretchr/testify/assert"
)

// TODO null detect bug
func TestInsertNil(t *testing.T) {
	row := &destT{
		ID:       100,
		Name:     "archever2",
		UserType: Male,
	}
	rowID, _ := s.Table("test").Insert(row).MustDo()
	assert.Equal(t, int64(100), rowID)
	var ret *destT
	s.Table("test").Select().Filter(f.Equal("id", 100)).MustOne(&ret)
	log.Printf("ret: %+v", ret.Datetime)
}
