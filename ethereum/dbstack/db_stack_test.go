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
		Hash:      common.HexToHash("0x1"),
		Parent:    common.HexToHash("0x0"),
		Timestamp: time.Now().UTC(),
	}
	miniHeaderTwo = &miniheader.MiniHeader{
		Number:    big.NewInt(2),
		Hash:      common.HexToHash("0x2"),
		Parent:    common.HexToHash("0x1"),
		Timestamp: time.Now().UTC(),
	}
)

func TestDBStackPushPeekPop(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)
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

func TestDBStackErrorIfPushTwoHeadersWithSameNumber(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)
	// Push miniHeaderOne
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	assert.Error(t, err)
}

func TestDBStackErrorIfResetWithoutCheckpointFirst(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)

	checkpointID := 123
	err = stack.Reset(checkpointID)
	require.Error(t, err)
}

func TestDBStackClear(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)

	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	miniHeader, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	err = stack.Clear()
	require.NoError(t, err)

	miniHeader, err = stack.Peek()
	require.NoError(t, err)
	require.Nil(t, miniHeader)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)
}

func TestDBStackErrorIfResetWithOldCheckpoint(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)

	checkpointIDOne, err := stack.Checkpoint()
	require.NoError(t, err)

	checkpointIDTwo, err := stack.Checkpoint()
	require.NoError(t, err)

	err = stack.Reset(checkpointIDOne)
	require.Error(t, err)

	err = stack.Reset(checkpointIDTwo)
	require.NoError(t, err)
}

func TestDBStackCheckpoint(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)
	err = stack.Push(miniHeaderTwo)
	require.NoError(t, err)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)

	_, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 2)
	assert.Equal(t, miniHeaderOne, miniHeaders[0])
	assert.Equal(t, miniHeaderTwo, miniHeaders[1])

	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderTwo, miniHeader)

	miniHeader, err = stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 2)
	assert.Equal(t, miniHeaderOne, miniHeaders[0])
	assert.Equal(t, miniHeaderTwo, miniHeaders[1])

	_, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)
}

func TestDBStackCheckpointAfterSameHeaderPushedAndPopped(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)
	// Push miniHeaderOne
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Pop miniHeaderOne
	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	_, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)
}

func TestDBStackCheckpointAfterSameHeaderPushedThenPoppedThenPushed(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)
	// Push miniHeaderOne
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Pop miniHeaderOne
	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	_, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 1)
	assert.Equal(t, miniHeaderOne, miniHeaders[0])
}

func TestDBStackCheckpointThenReset(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack, err := New(meshDB, 10)
	require.NoError(t, err)

	checkpointID, err := stack.Checkpoint()
	require.NoError(t, err)

	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)

	err = stack.Reset(checkpointID)
	require.NoError(t, err)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)

	err = stack.Push(miniHeaderTwo)
	require.NoError(t, err)

	checkpointID, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 1)
	assert.Equal(t, miniHeaderTwo, miniHeaders[0])

	miniHeader, err := stack.Pop()
	require.NoError(t, err)
	assert.Equal(t, miniHeader, miniHeaderTwo)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 1)
	assert.Equal(t, miniHeaderTwo, miniHeaders[0])

	checkpointID, err = stack.Checkpoint()
	require.NoError(t, err)

	miniHeaders, err = meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, 0)
}
