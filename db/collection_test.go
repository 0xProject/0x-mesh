package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCollection(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	_, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	_, err = db.NewCollection("people", &testModel{})
	require.Error(t, err, "Expected an error when creating new collection with the same name")
}

func TestInsert(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	expected := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(expected))
	exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
	require.NoError(t, err)
	assert.True(t, exists, "Model not stored in database at the expected key")
}

func TestFindByID(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	expected := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(expected))
	actual := &testModel{}
	require.NoError(t, col.FindByID(expected.ID(), actual))
	assert.Equal(t, expected, actual)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	original := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(original))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	actual := &testModel{}
	require.NoError(t, col.FindByID(original.ID(), actual))
	assert.Equal(t, updated, actual)
}

func TestFindAll(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}
	var actual []*testModel
	require.NoError(t, col.FindAll(&actual))
	assert.Equal(t, expected, actual)
}

func TestCount(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)

	// Insert some test models and make sure Count is equal to the number of
	// models inserted.
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}
	{
		actualCount, err := col.Count()
		require.NoError(t, err)
		assert.Equal(t, len(expected), actualCount, "Count returned wrong results")
	}

	// Delete a model and then check that the count decremented by 1.
	require.NoError(t, col.Delete(expected[0].ID()))
	{
		actualCount, err := col.Count()
		require.NoError(t, err)
		assert.Equal(t, len(expected)-1, actualCount, "Count returned wrong results")
	}

	// Delete all remaining models and check that the countKey was deleted.
	for _, model := range expected[1:] {
		require.NoError(t, col.Delete(model.ID()))
	}
	{
		actualCount, err := col.Count()
		require.NoError(t, err)
		assert.Equal(t, 0, actualCount, "Count returned wrong results")
		countKeyExists, err := col.ldb.Has(col.info.countKey(), nil)
		require.NoError(t, err)
		require.False(t, countKeyExists, "expected countKey to be deleted but it was not")
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	require.NoError(t, col.Delete(model.ID()))
	{
		exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Primary key should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Index should not be stored in database after calling Delete")
	}
}

func TestDeleteAfterUpdate(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
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
	require.NoError(t, col.Delete(model.ID()))
	{
		exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Primary key should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Old index should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:43:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Updated index should not be stored in database after calling Delete")
	}
}
