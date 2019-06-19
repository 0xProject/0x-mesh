package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkTransactionInsert100(b *testing.B) {
	benchmarkTransactionInsert(b, 100)
}

func BenchmarkTransactionInsert1000(b *testing.B) {
	benchmarkTransactionInsert(b, 1000)
}

func BenchmarkTransactionInsert10000(b *testing.B) {
	benchmarkTransactionInsert(b, 10000)
}

func benchmarkTransactionInsert(b *testing.B, count int) {
	b.Helper()
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txn := col.OpenTransaction()
		defer func() {
			_ = txn.Discard()
		}()
		for j := 0; j < count; j++ {
			model := &testModel{
				Name: fmt.Sprintf("person_%d_%d", i, j),
				Age:  j,
			}
			err := txn.Insert(model)
			b.StopTimer()
			require.NoError(b, err)
			b.StartTimer()
		}
		err := txn.Commit()
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}
