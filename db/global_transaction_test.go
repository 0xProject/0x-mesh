package db

import (
	"strconv"
	"sync"
	"testing"
	"time"

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

	// beforeTxnOpen is a set of testModels inserted before the transaction is opened.
	beforeTxnOpen := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "ExpectedPerson_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col0.Insert(model))
		beforeTxnOpen = append(beforeTxnOpen, model)
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
	// the database state until after is committed.
	outsideTransaction := []*testModel{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			model := &testModel{
				Name: "OutsideTransaction_" + strconv.Itoa(i),
				Age:  i,
			}
			require.NoError(t, col0.Insert(model))
			outsideTransaction = append(outsideTransaction, model)
		}
	}()

	// Any models we add to col0 within the transaction should not affect
	// the database state until after it is committed.
	insideTransaction := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "InsideTransaction_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, txn.Insert(col0, model))
		insideTransaction = append(insideTransaction, model)
	}

	// Any models we add to col1 after opening the transaction should not affect
	// the database state until after it is committed.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			model := &testModel{
				Name: "OtherPerson_" + strconv.Itoa(i),
				Age:  i,
			}
			require.NoError(t, col1.Insert(model))
		}
	}()

	// Any models we delete after opening the transaction should not affect
	// the database state until after it is committed.
	idToDelete := beforeTxnOpen[2].ID()
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
	assert.Equal(t, beforeTxnOpen, actual)

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

	// Check that all the models are now written.
	// TODO(albrow): Fix bug with Count and transactions, then we can use Count
	// instead of FindAll here.
	var existingModels []*testModel
	require.NoError(t, col0.FindAll(&existingModels))
	assert.Len(t, existingModels, 14)

	col1PostTxnCount, err := col1.Count()
	require.NoError(t, err)
	assert.Equal(t, 5, col1PostTxnCount)
}

// TestGlobalTransactionExclusion is designed to test whether a global
// transaction has exclusive write access for all collections while open.
func TestGlobalTransactionExclusion(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col0, err := db.NewCollection("people0", &testModel{})
	require.NoError(t, err)
	col1, err := db.NewCollection("people1", &testModel{})
	require.NoError(t, err)

	txn := db.OpenGlobalTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// discardSignal is fired right before the original global transaction is
	// discarded.
	discardSignal := make(chan struct{}, 1)
	// newGlobalTxnOpenSignal is fired when a new global transaction is opened.
	newGlobalTxnOpenSignal := make(chan struct{}, 1)
	// col0TxnOpenSignal is fired when a transaction on col0 is opened.
	col0TxnOpenSignal := make(chan struct{}, 1)
	// col1TxnOpenSignal is fired when a transaction on col1 is opened.
	col1TxnOpenSignal := make(chan struct{}, 1)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-discardSignal:
				// Expected outcome. Return from the goroutine.
				return
			case <-newGlobalTxnOpenSignal:
				t.Error("a new global transaction was opened before the first was committed/discarded")
			case <-col0TxnOpenSignal:
				t.Error("a new transaction was opened on col0 before the global transaction was committed/discarded")
			case <-col1TxnOpenSignal:
				t.Error("a new transaction was opened on col1 before the global transaction was committed/discarded")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		txn := db.OpenGlobalTransaction()
		newGlobalTxnOpenSignal <- struct{}{}
		defer func() {
			_ = txn.Discard()
		}()
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

	time.Sleep(1 * time.Millisecond)
	discardSignal <- struct{}{}
	require.NoError(t, txn.Discard())

	wg.Wait()
}
