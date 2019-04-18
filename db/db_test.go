package db

import (
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

func (tm testModel) ID() []byte {
	return []byte(tm.Name)
}

func TestInsert(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people")
	expected := testModel{
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
	col := db.NewCollection("people")
	expected := testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(expected))
	var got testModel
	require.NoError(t, col.FindByID(expected.ID(), &got))
	assert.Equal(t, expected, got)
}

func TestFindAll(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people")
	expected := []testModel{}
	for i := 0; i < 5; i++ {
		model := testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}
	var got []testModel
	require.NoError(t, col.FindAll(&got))
	assert.Equal(t, expected, got)
}
