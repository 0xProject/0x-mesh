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

// transactionExclusionTestTimeout is used in transaction exclusion tests to
// timeout while waiting for one transaction to open.
const transactionExclusionTestTimeout = 500 * time.Millisecond

func TestTransaction(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	defer db.Close()
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

func TestTransactionCount(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	defer db.Close()
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

	// Open a transaction.
	txn := col.OpenTransaction()
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
		require.NoError(t, txn.Insert(model))
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
		require.NoError(t, txn.Delete(id))
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

// TestTransactionExclusion is designed to test whether a collection-based
// transaction has exclusive write access for the collection while open.
func TestTransactionExclusion(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	defer db.Close()
	col0, err := db.NewCollection("people0", &testModel{})
	require.NoError(t, err)
	col1, err := db.NewCollection("people1", &testModel{})
	require.NoError(t, err)

	txn := col0.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// col0TxnOpenSignal is fired when a transaction on col0 is opened.
	col0TxnOpenSignal := make(chan struct{})
	// col1TxnOpenSignal is fired when a transaction on col1 is opened.
	col1TxnOpenSignal := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		txn := col0.OpenTransaction()
		close(col0TxnOpenSignal)
		defer func() {
			_ = txn.Discard()
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		txn := col1.OpenTransaction()
		close(col1TxnOpenSignal)
		defer func() {
			_ = txn.Discard()
		}()
	}()

	select {
	case <-col1TxnOpenSignal:
		// This is expected. Continue the test.
		break
	case <-time.After(transactionExclusionTestTimeout):
		t.Fatal("timed out waiting for col1 transaction to open")
	case <-col0TxnOpenSignal:
		t.Error("a new transaction was opened on col0 before the first transaction was committed/discarded")
	}

	require.NoError(t, txn.Discard())

	select {
	case <-col0TxnOpenSignal:
		// This is expected. Continue the test.
		break
	case <-time.After(transactionExclusionTestTimeout):
		t.Fatal("timed out waiting for second col0 transaction to open")
	}

	wg.Wait()
}
