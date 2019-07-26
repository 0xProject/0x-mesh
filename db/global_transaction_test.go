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

	// This short sleep prevents false positives. Without this, the test might
	// pass simply because the other goroutines did not have time do do anything
	// before we committed the transaction. We want to rule this out and make sure
	// that the mutexes are the thing enforcing that no new collections are
	// created and no new writes are made to the db state until after the global
	// transaction is committed.
	time.Sleep(transactionTestSleepDuration)

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

func TestGlobalTransactionCount(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)

	// insertedBeforeTransaction is a set of testModels inserted before the
	// transaction is opened.
	insertedBeforeTransaction := []*testModel{}
	for i := 0; i < 10; i++ {
		model := &testModel{
			Name: "Before_Transaction_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		insertedBeforeTransaction = append(insertedBeforeTransaction, model)
	}

	// Open a global transaction.
	txn := db.OpenGlobalTransaction()
	defer func() {
		err := txn.Discard()
		if err != nil && err != ErrCommitted {
			t.Error(err)
		}
	}()

	// Insert some models inside the transaction.
	for i := 0; i < 7; i++ {
		model := &testModel{
			Name: "Inside_Transaction_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, txn.Insert(col, model))
	}

	// The WaitGroup will be used to wait for all goroutines to finish.
	wg := &sync.WaitGroup{}

	// Insert some models outside the transaction.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 4; i++ {
			model := &testModel{
				Name: "Outside_Transaction_" + strconv.Itoa(i),
				Age:  42,
			}
			require.NoError(t, col.Insert(model))
		}
	}()

	// Delete some models inside of the transaction.
	idsToDeleteInside := [][]byte{
		insertedBeforeTransaction[0].ID(),
		insertedBeforeTransaction[1].ID(),
		insertedBeforeTransaction[2].ID(),
	}
	for _, id := range idsToDeleteInside {
		require.NoError(t, txn.Delete(col, id))
	}

	// Delete some models outside of the transaction.
	idsToDeleteOutside := [][]byte{
		insertedBeforeTransaction[3].ID(),
		insertedBeforeTransaction[4].ID(),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, id := range idsToDeleteOutside {
			require.NoError(t, col.Delete(id))
		}
	}()

	// Make sure that prior to commiting the transaction, Count only includes the
	// models inserted/deleted before the transaction was open.
	expectedPreCommitCount := 10
	actualPreCommitCount, err := col.Count()
	require.NoError(t, err)
	assert.Equal(t, expectedPreCommitCount, actualPreCommitCount)

	// Commit the transaction.
	require.NoError(t, txn.Commit())

	// Wait for any goroutines to finish.
	wg.Wait()

	// Make sure that after commiting the transaction, Count includes the models
	// inserted/deleted in the transaction and outside of the transaction.
	//   10 before transaction.
	//   +7 inserted inside transaction
	//   +4 inserted outside transaction
	//   -3 deleted inside transaction
	//   -2 deleted outside transaction
	// = 16 total
	expectedPostCommitCount := 16
	actualPostCommitCount, err := col.Count()
	require.NoError(t, err)
	assert.Equal(t, expectedPostCommitCount, actualPostCommitCount)
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

	time.Sleep(transactionTestSleepDuration)
	discardSignal <- struct{}{}
	require.NoError(t, txn.Discard())

	wg.Wait()
}
