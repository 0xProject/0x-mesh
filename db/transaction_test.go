package db

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
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

	// Open a transaction.
	txn := col.OpenTransaction()
	defer func() {
		err := txn.Discard()
		if err != nil && err != ErrCommitted {
			t.Error(err)
		}
	}()

	// The WaitGroup will be used to wait for all goroutines to finish.
	wg := &sync.WaitGroup{}

	// Any models we add after opening the transaction should not affect the query.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			model := &testModel{
				Name: "OtherPerson_" + strconv.Itoa(i),
				Age:  42,
			}
			require.NoError(t, col.Insert(model))
		}
	}()

	// Any models we delete after opening the transaction should not affect the query.
	idToDelete := expected[2].ID()
	wg.Add(1)
	go func(idToDelete []byte) {
		defer wg.Done()
		require.NoError(t, col.Delete(idToDelete))
	}(idToDelete)

	// Any new indexes we add should not affect indexes in the transaction.
	col.AddIndex("name", func(m Model) []byte {
		return []byte(m.(*testModel).Name)
	})
	assert.Equal(t, []*Index{ageIndex}, txn.colInfo.indexes)

	// Make sure that the query only return results that match the state inside
	// the transaction.
	filter := ageIndex.ValueFilter([]byte("42"))
	query := col.NewQuery(filter)
	var actual []*testModel
	require.NoError(t, query.Run(&actual))
	assert.Equal(t, expected, actual)

	// Commit the transaction.
	require.NoError(t, txn.Commit())

	// Wait for any goroutines to finish.
	wg.Wait()
}
