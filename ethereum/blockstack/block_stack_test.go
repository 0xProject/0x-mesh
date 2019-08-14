package blockstack

import (
	"math/big"
	"testing"

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
		Number: big.NewInt(1),
		Hash:   common.Hash{},
		Parent: common.Hash{},
	}
)

func TestBlockStackPush(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	stack := New(meshDB, 10)
	err = stack.Push(miniHeaderOne)
	require.NoError(t, err)

	expectedLen := 1
	miniHeaders, err := stack.Inspect()
	require.NoError(t, err)
	assert.Len(t, miniHeaders, expectedLen)
}
