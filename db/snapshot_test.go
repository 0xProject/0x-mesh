package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshot(t *testing.T) {
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

	// Take a snapshot.
	snapshot, err := db.GetSnapshot(col)
	require.NoError(t, err)
	defer snapshot.Release()

	// Any models we add after taking the snapshot should not be included in the
	// query.
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "OtherPerson_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
	}

	// Make sure that only models we inserted before taking the snapshot are
	// included in the query results.
	filter := ageIndex.ValueFilter([]byte("42"))
	query := snapshot.NewQuery(filter)
	var actual []*testModel
	require.NoError(t, query.Run(&actual))
	assert.Equal(t, expected, actual)
}
