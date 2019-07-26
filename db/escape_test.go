package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var trickyByteValues = [][]byte{
	[]byte(":"),
	[]byte(`\`),
	[]byte("::"),
	[]byte(`\\`),
	[]byte(`\:`),
	[]byte(`:\`),
	[]byte(`\\:`),
	[]byte(`::\`),
	[]byte(`\:\:`),
	[]byte(`:\:\`),
	[]byte(`:\\`),
	[]byte(`\::`),
	[]byte(`::\\`),
	[]byte(`\\::`),
}

func TestEscapeUnescape(t *testing.T) {
	t.Parallel()
	for _, expected := range trickyByteValues {
		actual := unescape(escape(expected))
		assert.Equal(t, expected, actual)
	}
}

func TestFindWithValueWithEscape(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	ageIndex := col.AddIndex("age", func(m Model) []byte {
		// Note: We add the ':' to the index value to try and trip up the escaping
		// algorithm.
		return []byte(fmt.Sprintf(":%d:", m.(*testModel).Age))
	})
	models := make([]*testModel, len(trickyByteValues))
	// Use the trickyByteValues as the names for each model.
	for i, name := range trickyByteValues {
		models[i] = &testModel{
			Name: string(name),
			Age:  i,
		}
	}
	for i, expected := range models {
		require.NoError(t, col.Insert(expected), "testModel %d", i)
		actual := []*testModel{}
		query := col.NewQuery(ageIndex.ValueFilter([]byte(fmt.Sprintf(":%d:", expected.Age))))
		require.NoError(t, query.Run(&actual), "testModel %d", i)
		require.Len(t, actual, 1, "testModel %d", i)
		assert.Equal(t, expected, actual[0])
	}
}
