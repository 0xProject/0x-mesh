package blockwatch

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	miniHeaderOne = &types.MiniHeader{
		Number:    big.NewInt(1),
		Hash:      common.HexToHash("0x1"),
		Parent:    common.HexToHash("0x0"),
		Timestamp: time.Now().UTC(),
	}
	miniHeaderTwo = &types.MiniHeader{
		Number:    big.NewInt(2),
		Hash:      common.HexToHash("0x2"),
		Parent:    common.HexToHash("0x1"),
		Timestamp: time.Now().UTC(),
	}
)

func newTestStack(t *testing.T, ctx context.Context) *Stack {
	database, err := db.New(ctx, &db.Options{Path: "/tmp/orderwatcher_testing/" + uuid.New().String()})
	require.NoError(t, err)
	return NewStack(database)
}

func TestSimpleStackPushPeekPop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	expectedLen := 1
	miniHeaders, err := stack.PeekAll()
	require.NoError(t, err)
	require.Len(t, miniHeaders, expectedLen)

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

func TestSimpleStackErrorIfPushTwoHeadersWithSameNumber(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)
	// Push miniHeaderOne
	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)
	// Push miniHeaderOne again
	err = stack.Push(miniHeaderOne)
	assert.Error(t, err)
}

func TestSimpleStackClear(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := newTestStack(t, ctx)

	err := stack.Push(miniHeaderOne)
	require.NoError(t, err)

	miniHeader, err := stack.Peek()
	require.NoError(t, err)
	assert.Equal(t, miniHeaderOne, miniHeader)

	err = stack.Clear()
	require.NoError(t, err)

	miniHeader, err = stack.Peek()
	require.NoError(t, err)
	require.Nil(t, miniHeader)
}
