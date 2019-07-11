package db

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrityCheckPass(t *testing.T) {
	t.Parallel()
	db, _, _, _ := setUpIntegrityCheckTest(t)

	// We didn't break anything so the integrity check should pass
	require.NoError(t, db.CheckIntegrity())
}

func TestIntegrityCheckInvalidModelData(t *testing.T) {
	t.Parallel()
	db, col, models, _ := setUpIntegrityCheckTest(t)

	// Manually break integrity by storing invalid model data.
	keyToChange := col.info.primaryKeyForModel(models[0])
	require.NoError(t, db.ldb.Put(keyToChange, []byte("invalid data"), nil))
	expectedError := "integritiy check failed for collection people: could not unmarshal model data for primary key model:people:Person_0: invalid character 'i' looking for beginning of value"
	require.EqualError(t, db.CheckIntegrity(), expectedError)
}

func TestIntegrityCheckIndexKeyWithoutModelData(t *testing.T) {
	t.Parallel()
	db, col, models, _ := setUpIntegrityCheckTest(t)

	// Manually break integrity by deleting a primary key.
	keyToDelete := col.info.primaryKeyForModel(models[0])
	require.NoError(t, db.ldb.Delete(keyToDelete, nil))
	expectedError := "integritiy check failed for index people.age: key exists in index but could not find corresponding model data for primary key: model:people:Person_0"
	require.EqualError(t, db.CheckIntegrity(), expectedError)
}

func TestIntegrityCheckModelNotIndexed(t *testing.T) {
	t.Parallel()
	db, _, models, ageIndex := setUpIntegrityCheckTest(t)

	// Manually break integrity by deleting an index key.
	keyToDelete := ageIndex.keysForModel(models[0])[0]
	require.NoError(t, db.ldb.Delete(keyToDelete, nil))
	expectedError := "integritiy check failed for index people.age: indexKey index:people:age:0:Person_0 does not exist"
	require.EqualError(t, db.CheckIntegrity(), expectedError)
}

func setUpIntegrityCheckTest(t *testing.T) (*DB, *Collection, []*testModel, *Index) {
	db := newTestDB(t)
	col, err := db.NewCollection("people", &testModel{})
	require.NoError(t, err)
	ageIndex := col.AddIndex("age", func(m Model) []byte {
		return []byte(fmt.Sprint(m.(*testModel).Age))
	})

	// Insert some test models
	models := []*testModel{}
	for i := 0; i < 5; i++ {
		model := &testModel{
			Name: "Person_" + strconv.Itoa(i),
			Age:  i,
		}
		require.NoError(t, col.Insert(model))
		models = append(models, model)
	}

	return db, col, models, ageIndex
}
