package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const defaultTargetNickname = "target"

func setupQueryBenchmark(b *testing.B) (db *DB, col *Collection, nicknameIndex *Index) {
	b.Helper()
	db = newTestDB(b)
	defer db.Close()
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(b, err)
	nicknameIndex = col.AddMultiIndex("nicknames", func(m Model) [][]byte {
		person := m.(*testModel)
		indexValues := make([][]byte, len(person.Nicknames))
		for i, nickname := range person.Nicknames {
			indexValues[i] = []byte(nickname)
		}
		return indexValues
	})
	return db, col, nicknameIndex
}

func insertModelsForQueryBenchmark(b *testing.B, col *Collection, targetNickname string, targetCount int, otherCount int) {
	txn := col.OpenTransaction()

	defer func() {
		_ = txn.Discard()
	}()
	// Insert targetCount models with nickname = targetNickname
	for i := 0; i < targetCount; i++ {
		model := &testModel{
			Name:      fmt.Sprintf("person_%d", i),
			Age:       i,
			Nicknames: []string{targetNickname},
		}
		require.NoError(b, txn.Insert(model))
	}
	// Insert otherCount with nickname != targetnickName
	for i := 0; i < otherCount; i++ {
		model := &testModel{
			Name:      fmt.Sprintf("person_%d", i),
			Age:       i,
			Nicknames: []string{fmt.Sprintf("not_%s_%d", targetNickname, i)},
		}
		require.NoError(b, txn.Insert(model))
	}
	require.NoError(b, txn.Commit())
}

func benchmarkQueryFind(b *testing.B, targetCount int, total int) {
	benchmarkQueryFindWithMaxAndOffset(b, targetCount, total, 0, 0)
}

func BenchmarkQueryFind100OutOf100(b *testing.B) {
	benchmarkQueryFind(b, 100, 100)
}

func BenchmarkQueryFind100OutOf1000(b *testing.B) {
	benchmarkQueryFind(b, 100, 1000)
}

func BenchmarkQueryFind100OutOf10000(b *testing.B) {
	benchmarkQueryFind(b, 100, 10000)
}

func BenchmarkQueryFind1000OutOf1000(b *testing.B) {
	benchmarkQueryFind(b, 1000, 1000)
}

func BenchmarkQueryFind1000OutOf10000(b *testing.B) {
	benchmarkQueryFind(b, 1000, 10000)
}

func BenchmarkQueryFind10000OutOf10000(b *testing.B) {
	benchmarkQueryFind(b, 10000, 10000)
}

func benchmarkQueryCount(b *testing.B, targetCount int, total int) {
	db, col, nicknameIndex := setupQueryBenchmark(b)
	defer db.Close()
	insertModelsForQueryBenchmark(b, col, defaultTargetNickname, targetCount, total-targetCount)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := col.NewQuery(nicknameIndex.ValueFilter([]byte(defaultTargetNickname)))
		_, err := query.Count()
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkQueryCount100OutOf100(b *testing.B) {
	benchmarkQueryCount(b, 100, 100)
}

func BenchmarkQueryCount100OutOf1000(b *testing.B) {
	benchmarkQueryCount(b, 100, 1000)
}

func BenchmarkQueryCount100OutOf10000(b *testing.B) {
	benchmarkQueryCount(b, 100, 10000)
}

func BenchmarkQueryCount1000OutOf1000(b *testing.B) {
	benchmarkQueryCount(b, 1000, 1000)
}

func BenchmarkQueryCount1000OutOf10000(b *testing.B) {
	benchmarkQueryCount(b, 1000, 10000)
}

func BenchmarkQueryCount10000OutOf10000(b *testing.B) {
	benchmarkQueryCount(b, 10000, 10000)
}

func benchmarkQueryFindWithMaxAndOffset(b *testing.B, targetCount int, total int, max int, offset int) {
	db, col, nicknameIndex := setupQueryBenchmark(b)
	defer db.Close()
	insertModelsForQueryBenchmark(b, col, defaultTargetNickname, targetCount, total-targetCount)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := col.NewQuery(nicknameIndex.ValueFilter([]byte(defaultTargetNickname)))
		var actual []*testModel
		err := query.Max(max).Offset(offset).Run(&actual)
		b.StopTimer()
		require.NoError(b, err)
		b.StartTimer()
	}
}

func BenchmarkQueryFind1000OutOf10000Max100Offset0(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 1000, 10000, 100, 0)
}

func BenchmarkQueryFind1000OutOf10000Max100Offset100(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 1000, 10000, 100, 100)
}

func BenchmarkQueryFind1000OutOf10000Max100Offset900(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 1000, 10000, 100, 900)
}

func BenchmarkQueryFind1000OutOf10000Max100Offset1000(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 1000, 10000, 100, 1000)
}

func BenchmarkQueryFind10000OutOf10000Max1000Offset0(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 10000, 10000, 1000, 0)
}

func BenchmarkQueryFind10000OutOf10000Max1000Offset1000(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 10000, 10000, 1000, 1000)
}

func BenchmarkQueryFind10000OutOf10000Max1000Offset9000(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 10000, 10000, 1000, 9000)
}

func BenchmarkQueryFind10000OutOf10000Max1000Offset10000(b *testing.B) {
	benchmarkQueryFindWithMaxAndOffset(b, 10000, 10000, 1000, 10000)
}
