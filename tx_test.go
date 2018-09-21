package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxCommit(t *testing.T) {
	tx, err := NewTx(db)
	assert.NoError(t, err)
	err = tx.Begin()
	assert.NoError(t, err)
	_, _, err = tx.Table("test").Insert(M{"name": "tx_test1"}).Do()
	assert.NoError(t, err)
	tx.Commit()

	var dest map[string]interface{}
	err = o.Table("test").Select().Where(S("name=?", "tx_test1")).One(&dest)
	assert.NoError(t, err)
	assert.Equal(t, "tx_test1", dest["name"])
}
func TestTxRollBack(t *testing.T) {
	tx, err := NewTx(db)
	assert.NoError(t, err)
	err = tx.Begin()
	assert.NoError(t, err)
	_, _, err = tx.Table("test").Insert(M{"name": "tx_test2"}).Do()
	assert.NoError(t, err)
	tx.RollBack()

	var dest map[string]interface{}
	err = o.Table("test").Select().Where(S("name=?", "tx_test2")).One(&dest)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(dest))
}
