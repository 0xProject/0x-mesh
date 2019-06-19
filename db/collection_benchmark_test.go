package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkInsert(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		model := &testModel{
			Name: fmt.Sprintf("person_%d", i),
			Age:  i,
		}
		b.StartTimer()
		err := col.Insert(model)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkFindByIDHot(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(b, col.Insert(model))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var found testModel
		err := col.FindByID(model.ID(), &found)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkFindByIDCold(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		model := &testModel{
			Name: fmt.Sprintf("person_%d", i),
			Age:  i,
		}
		require.NoError(b, col.Insert(model))
		b.StartTimer()
		var found testModel
		err := col.FindByID(model.ID(), &found)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkUpdateHot(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	original := &testModel{
		Name: "person_0",
		Age:  0,
	}
	require.NoError(b, col.Insert(original))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		updated := &testModel{
			Name: original.Name,
			Age:  i + 1,
		}
		b.StartTimer()
		err := col.Update(updated)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkUpdateCold(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		original := &testModel{
			Name: fmt.Sprintf("person_%d", i),
			Age:  i,
		}
		require.NoError(b, col.Insert(original))
		updated := &testModel{
			Name: original.Name,
			Age:  i + 1,
		}
		b.StartTimer()
		err := col.Update(updated)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkFindAll100(b *testing.B) {
	benchmarkFindAll(b, 100)
}

func BenchmarkFindAll1000(b *testing.B) {
	benchmarkFindAll(b, 1000)
}

func benchmarkFindAll(b *testing.B, count int) {
	b.Helper()
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	expected := []*testModel{}
	for i := 0; i < count; i++ {
		model := &testModel{
			Name: "person_%d" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(b, col.Insert(model))
		expected = append(expected, model)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var actual []*testModel
		err := col.FindAll(&actual)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkDelete(b *testing.B) {
	db := newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		model := &testModel{
			Name: "person_%d" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(b, col.Insert(model))
		b.StartTimer()
		err := col.Delete(model.ID())
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}
