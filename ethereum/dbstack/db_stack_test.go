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
		Timestamp: time.Now(),
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
