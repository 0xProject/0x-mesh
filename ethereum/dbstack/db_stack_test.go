package dbstack

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const limit = 10

var (
	miniHeaderOne = &miniheader.MiniHeader{
		Number:    big.NewInt(1),
		Hash:      common.Hash{},
		Parent:    common.Hash{},
		Timestamp: time.Now().UTC(),
	}
)

func TestDBStackPushPeekPop(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack := New(meshDB, 10)
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	expectedLen := 1
	miniHeaders, err := stack.PeekAll()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, expectedLen)

	miniHeader, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, miniHeaders[0], miniHeader)

	expectedLen = 1
	miniHeaders, err = stack.PeekAll()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, expectedLen)

	miniHeader, err = stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaders[0], miniHeader)

	expectedLen = 0
	miniHeaders, err = stack.PeekAll()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, expectedLen)
}

func TestDBStackReset(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack := New(meshDB, 10)
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	err = stack.Reset()
	require.NoError(t, err)

	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Nil(t, miniHeader)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)
}

func TestDBStackCheckpoint(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack := New(meshDB, 10)
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 1)
	assert.Equal(t, miniHeaderOne, miniHeaders[0])
}
