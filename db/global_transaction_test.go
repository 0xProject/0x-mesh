package db

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlobalTransaction(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col0, err := db.NewCollection("people0", &testModel{})
	require.NoError(t, err)
	col1, err := db.NewCollection("people1", &testModel{})
	require.NoError(t, err)

	// expected is a set of testModels with Age = 42
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "ExpectedPerson_" + strconv.Itoa(i),
			Age:  42,
		}
		require.NoError(t, col0.Insert(model))
		expected = append(expected, model)
	}

	// Open a global transaction.
	txn := db.OpenGlobalTransaction()
	defer func() {
		err := txn.Discard()
		if err != nil && err != ErrCommitted {
			t.Error(err)
		}
	}()

	// The WaitGroup will be used to wait for all goroutines to finish.
	wg := &sync.WaitGroup{}

	// Any models we add to col0 after opening the transaction should not affect
	// the database state.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			model := &testModel{
				Name: "OtherPerson_" + strconv.Itoa(i),
				Age:  42,
			}
			require.NoError(t, col0.Insert(model))
		}
	}()

	// Any models we add to col1 after opening the transaction should not affect
	// the database state.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			model := &testModel{
				Name: "OtherPerson_" + strconv.Itoa(i),
				Age:  42,
			}
			require.NoError(t, col1.Insert(model))
		}
	}()

	// Any models we delete after opening the transaction should not affect
	// the database state.
	idToDelete := expected[2].ID()
	wg.Add(1)
	go func(idToDelete []byte) {
		defer wg.Done()
		require.NoError(t, col0.Delete(idToDelete))
	}(idToDelete)

	// Attempting to add a new collection should block until after the transaction
	// is committed/discarded. We use two channels to determine the order in which
	// the two events occurred.
	// commitSignal is fired right before the transaction is committed.
	commitSignal := make(chan struct{}, 1)
	// newCollectionSignal is fired after the new collection has been created.
	newCollectionSignal := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-commitSignal:
			// Expected outcome. Return from goroutine.
			return
		case <-newCollectionSignal:
			// Not the expected outcome. commitSignal should have fired first.
			t.Error("new collection was created before the transaction was committed")
			return
		}
	}()

	// TODO(albrow): Test that opening a collection transaction blocks until after
	// the transaction is committed.

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := db.NewCollection("people2", &testModel{})
		require.NoError(t, err)
		// Signal that the new collection was created.
		newCollectionSignal <- struct{}{}
	}()

	// Make sure that col0 only contains models that were created before the
	// transaction was opened.
	var actual []*testModel
	require.NoError(t, col0.FindAll(&actual))
	assert.Equal(t, expected, actual)

	// Make sure that col1 doesn't contain any models (since they were created
	// after the transaction was opened).
	actualCount, err := col1.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, actualCount)

	// Signal that we are about to commit the transaction, then commit it.
	commitSignal <- struct{}{}
	require.NoError(t, txn.Commit())

	// Wait for any goroutines to finish.
	wg.Wait()
}
