package db

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

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

// TestTransactionExclusion is designed to test whether a collection-based
// transaction has exclusive write access for the collection while open.
func TestTransactionExclusion(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col0, err := db.NewCollection("people0", &testModel{})
	require.NoError(t, err)
	col1, err := db.NewCollection("people1", &testModel{})
	require.NoError(t, err)

	txn := col0.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// discardSignal is fired right before the original transaction on col0 is
	// discarded.
	discardSignal := make(chan struct{}, 1)
	// col0TxnOpenSignal is fired when a transaction on col0 is opened.
	col0TxnOpenSignal := make(chan struct{}, 1)
	// col1TxnOpenSignal is fired when a transaction on col1 is opened.
	col1TxnOpenSignal := make(chan struct{}, 1)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		col1TxnWasOpened := false
		for {
			select {
			case <-discardSignal:
				assert.True(t, col1TxnWasOpened, "expected col1 transaction to open before col0 transaction was committed/discarded")
				return
			case <-col0TxnOpenSignal:
				t.Error("a new transaction was opened on col0 before the first transaction was committed/discarded")
			case <-col1TxnOpenSignal:
				// col1 transactions should be independent of col0 transactions and do
				// not need to wait for the col0 transaction to be committed/discarded.
				col1TxnWasOpened = true
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		txn := col0.OpenTransaction()
		col0TxnOpenSignal <- struct{}{}
		defer func() {
			_ = txn.Discard()
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		txn := col1.OpenTransaction()
		col1TxnOpenSignal <- struct{}{}
		defer func() {
			_ = txn.Discard()
		}()
	}()

	// A short sleep is necessary to ensure that the goroutines have time to run.
	time.Sleep(1 * time.Millisecond)
	discardSignal <- struct{}{}
	require.NoError(t, txn.Discard())

	wg.Wait()
}
