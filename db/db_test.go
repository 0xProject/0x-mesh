package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testModel struct {
	Name      string
	Age       int
	Nicknames []string
}

func (tm *testModel) ID() []byte {
	return []byte(tm.Name)
}

func newTestDB(t require.TestingT) *DB {
	db, err := Open("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	return db
}

func TestOpen(t *testing.T) {
	t.Parallel()
	db, err := Open("/tmp/leveldb_testing")
	require.NoError(t, err)
	require.NoError(t, db.Close())
}
