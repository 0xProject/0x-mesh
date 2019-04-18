package db

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestOpen(t *testing.T) {
	db, err := leveldb.OpenFile("/tmp/leveldb_testing", nil)
	defer db.Close()
	require.NoError(t, err)
}
