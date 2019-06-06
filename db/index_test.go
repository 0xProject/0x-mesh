package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertWithIndex(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.True(t, exists, "Index not stored in database at the expected key")
}

func TestUpdateWithIndex(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	oldKeyExists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.False(t, oldKeyExists, "Old index was still stored after update")
	updatedKeyExists, err := db.ldb.Has([]byte("index:people:age:43:foo"), nil)
	require.NoError(t, err)
	assert.True(t, updatedKeyExists, "Index not stored in database at the updated key")
}
