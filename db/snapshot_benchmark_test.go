package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetSnapshot(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	for i := 0; i < 1000; i++ {
		model := &testModel{
			Name: fmt.Sprintf("person_%d", i),
			Age:  i,
		}
		require.NoError(b, col.Insert(model))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		snapshot, err := col.GetSnapshot()
		b.StopTimer()
		require.NoError(b, err)
		snapshot.Release()
		b.StartTimer()
	}
}

func BenchmarkSnapshotFindByID(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(b, col.Insert(model))
	snapshot, err := col.GetSnapshot()
	require.NoError(b, err)
	defer snapshot.Release()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var found testModel
		err := snapshot.FindByID(model.ID(), &found)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}
