package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *DB {
	db, err := Open("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	return db
}

func TestOpen(t *testing.T) {
	db, err := Open("/tmp/leveldb_testing")
	require.NoError(t, err)
	require.NoError(t, db.Close())
}

type testModel struct {
	Name      string
	Age       int
	Nicknames []string
}

func (tm *testModel) ID() []byte {
	return []byte(tm.Name)
}

func TestInsert(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	expected := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(expected))
	exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
	require.NoError(t, err)
	assert.True(t, exists, "Model not stored in database at the expected key")
}

func TestFindByID(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	expected := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(expected))
	actual := &testModel{}
	require.NoError(t, col.FindByID(expected.ID(), actual))
	assert.Equal(t, expected, actual)
}

func TestUpdate(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	original := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(original))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	actual := &testModel{}
	require.NoError(t, col.FindByID(original.ID(), actual))
	assert.Equal(t, updated, actual)
}

func TestFindAll(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	expected := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		expected = append(expected, model)
	}
	var actual []*testModel
	require.NoError(t, col.FindAll(&actual))
	assert.Equal(t, expected, actual)
}

func TestInsertWithIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.True(t, exists, "Index not stored in database at the expected key")
}

func TestUpdateWithIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	oldKeyExists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
	require.NoError(t, err)
	assert.False(t, oldKeyExists, "Old index was still stored after update")
	updatedKeyExists, err := db.ldb.Has([]byte("index:people:age:43:foo"), nil)
	require.NoError(t, err)
	assert.True(t, updatedKeyExists, "Index not stored in database at the updated key")
}

func TestDelete(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	require.NoError(t, col.Delete(model.ID()))
	{
		exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Primary key should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Index should not be stored in database after calling Delete")
	}
}

func TestDeleteAfterUpdate(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})
	model := &testModel{
		Name: "foo",
		Age:  42,
	}
	require.NoError(t, col.Insert(model))
	updated := &testModel{
		Name: "foo",
		Age:  43,
	}
	require.NoError(t, col.Update(updated))
	require.NoError(t, col.Delete(model.ID()))
	{
		exists, err := db.ldb.Has([]byte("model:people:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Primary key should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:42:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Old index should not be stored in database after calling Delete")
	}
	{
		exists, err := db.ldb.Has([]byte("index:people:age:43:foo"), nil)
		require.NoError(t, err)
		assert.False(t, exists, "Updated index should not be stored in database after calling Delete")
	}
}

func TestFindWithValue(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

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

	// We also insert some other models with a different age.
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "OtherPerson_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithValue(ageIndex, []byte(fmt.Sprint(42)), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithRange(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})

	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	all := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		all = append(all, model)
	}
	// expected is the set of people with 1 <= age < 4
	expected := all[1:4]

	var actual []*testModel
	require.NoError(t, col.FindWithRange(ageIndex, []byte(fmt.Sprint(1)), []byte(fmt.Sprint(4)), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithValueWithMultiIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	nicknameIndex := col.AddMultiIndex("nicknames", func(m Model) [][]byte {
		person := m.(*testModel)
		indexValues := make([][]byte, len(person.Nicknames))
		for i, nickname := range person.Nicknames {
			indexValues[i] = []byte(nickname)
		}
		return indexValues
	})

	// expected is a set of testModels that include the nickname "Bob"
	expected := []*testModel{
		{
			Name:      "ExpectedPerson_0",
			Age:       42,
			Nicknames: []string{"Bob", "Jim", "John"},
		},
		{
			Name:      "ExpectedPerson_1",
			Age:       43,
			Nicknames: []string{"Alice", "Bob", "Emily"},
		},
		{
			Name:      "ExpectedPerson_2",
			Age:       44,
			Nicknames: []string{"Bob", "No one"},
		},
	}
	for _, model := range expected {
		require.NoError(t, col.Insert(model))
	}

	// We also insert some other models with different nicknames.
	excluded := []*testModel{
		{
			Name:      "ExcludedPerson_0",
			Age:       42,
			Nicknames: []string{"Bill", "Jim", "John"},
		},
		{
			Name:      "ExcludedPerson_1",
			Age:       43,
			Nicknames: []string{"Alice", "Jane", "Emily"},
		},
		{
			Name:      "ExcludedPerson_2",
			Age:       44,
			Nicknames: []string{"Nemo", "No one"},
		},
	}
	for _, model := range excluded {
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithValue(nicknameIndex, []byte("Bob"), &actual))
	assert.Equal(t, expected, actual)
}

func TestFindWithRangeWithMultiIndex(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	nicknameIndex := col.AddMultiIndex("nicknames", func(m Model) [][]byte {
		person := m.(*testModel)
		indexValues := make([][]byte, len(person.Nicknames))
		for i, nickname := range person.Nicknames {
			indexValues[i] = []byte(nickname)
		}
		return indexValues
	})

	// expected is a set of testModels that include at least one nickname that
	// satisfies "B" <= nickname < "E"
	expected := []*testModel{
		{
			Name:      "ExpectedPerson_0",
			Age:       42,
			Nicknames: []string{"Alice", "Beth", "Emily"},
		},
		{
			Name:      "ExpectedPerson_1",
			Age:       43,
			Nicknames: []string{"Bob", "Charles", "Dan"},
		},
		{
			Name:      "ExpectedPerson_2",
			Age:       44,
			Nicknames: []string{"James", "Darell"},
		},
	}
	for _, model := range expected {
		require.NoError(t, col.Insert(model))
	}

	// We also insert some other models with different nicknames.
	excluded := []*testModel{
		{
			Name:      "ExcludedPerson_0",
			Age:       42,
			Nicknames: []string{"Allen", "Jim", "John"},
		},
		{
			Name:      "ExcludedPerson_1",
			Age:       43,
			Nicknames: []string{"Sophia", "Jane", "Emily"},
		},
		{
			Name:      "ExcludedPerson_2",
			Age:       44,
			Nicknames: []string{"Nemo", "No one"},
		},
	}
	for _, model := range excluded {
		require.NoError(t, col.Insert(model))
	}

	var actual []*testModel
	require.NoError(t, col.FindWithRange(nicknameIndex, []byte("B"), []byte("E"), &actual))
	assert.Equal(t, expected, actual)
}

var trickyByteValues = [][]byte{
	[]byte(":"),
	[]byte(`\`),
	[]byte("::"),
	[]byte(`\\`),
	[]byte(`\:`),
	[]byte(`:\`),
	[]byte(`\\:`),
	[]byte(`::\`),
	[]byte(`\:\:`),
	[]byte(`:\:\`),
	[]byte(`:\\`),
	[]byte(`\::`),
	[]byte(`::\\`),
	[]byte(`\\::`),
}

func TestEscapeUnescape(t *testing.T) {
	for _, expected := range trickyByteValues {
		actual := unescape(escape(expected))
		assert.Equal(t, expected, actual)
	}
}

func TestFindWithValueWithEscape(t *testing.T) {
	db := newTestDB(t)
	col := db.NewCollection("people", &testModel{})
	ageIndex := col.AddIndex("age", func(m Model) []byte {
		// Note: We add the ':' to the index value to try and trip up the escaping
		// algorithm.
		return []byte(fmt.Sprintf(":%d:", m.(*testModel).Age))
	})
	models := make([]*testModel, len(trickyByteValues))
	// Use the trickyByteValues as the names for each model.
	for i, name := range trickyByteValues {
		models[i] = &testModel{
			Name: string(name),
			Age:  i,
		}
	}
	for i, expected := range models {
		require.NoError(t, col.Insert(expected), "testModel %d", i)
		actual := []*testModel{}
		err := col.FindWithValue(ageIndex, []byte(fmt.Sprintf(":%d:", expected.Age)), &actual)
		require.NoError(t, err, "testModel %d", i)
		require.Len(t, actual, 1, "testModel %d", i)
		assert.Equal(t, expected, actual[0])
	}
}
