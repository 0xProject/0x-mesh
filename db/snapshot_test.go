package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)

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

	// Take a snapshot.
	snapshot, err := col.GetSnapshot()
	require.NoError(t, err)
	defer snapshot.Release()

	// Any models we add after taking the snapshot should not affect the query.
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "OtherPerson_" + strconv.Itoa(i),
			Age:  42,
		}
		require.NoError(t, col.Insert(model))
	}

	// Any models we delete after taking the snapshot should not affect the query.
	for _, model := range expected {
		require.NoError(t, col.Delete(model.ID()))
	}

	// Any new indexes we add should not affect indexes in the snapshot.
	col.AddIndex("name", func(m Model) []byte {
		return []byte(m.(*testModel).Name)
	})
	assert.Equal(t, []*Index{ageIndex}, snapshot.colInfo.indexes)

	// Make sure that the query only return results that match the state at the
	// time the snapshot was taken.
	filter := ageIndex.ValueFilter([]byte("42"))
	query := snapshot.NewQuery(filter)
	var actual []*testModel
	require.NoError(t, query.Run(&actual))
	assert.Equal(t, expected, actual)
	actualCount, err := snapshot.Count()
	require.NoError(t, err)
	assert.Equal(t, len(expected), actualCount)
}
