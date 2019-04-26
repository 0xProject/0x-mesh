package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *DB {
	db, err := Open("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	return db
}

func TestOpen(t *testing.T) {
	db, err := Open("/tmp/leveldb_testing")
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

type testModel struct {
	Name string
	Age  int
}

func (tm *testModel) ID() []byte {
	return []byte(tm.Name)
}

func TestInsert(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
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
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
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
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
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
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
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

func TestInsertWithIndex(t *testing.T) {
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

func TestDelete(t *testing.T) {
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

func TestFindWithValue(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	// expected is a set of testModels with Age = 42
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "ExpectedPerson_" + strconv.Itoa(i),
			Age:  42,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}

	// We also insert some other models with a different age.
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "OtherPerson_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithValue(ageIndex, []byte(fmt.Sprint(42)), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithRange(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	all := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		all = append(all, model)
	}
	// expected is the set of people with 1 <= age < 4
	expected := all[1:4]

	var actual []*testModel
	require.NoError(t, col.FindWithRange(ageIndex, []byte(fmt.Sprint(1)), []byte(fmt.Sprint(4)), &actual))
	assert.Equal(t, expected, actual)
}
