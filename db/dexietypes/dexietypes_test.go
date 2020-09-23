// +build js,wasm

package dexietypes

import (
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestMiniHeadersConversion(t *testing.T) {
	originalMiniHeader := &types.MiniHeader{
		Hash:      common.HexToHash("0x00a3ce0e9cbcb5c4d79c1c19df276a0db954a487b895dca1d4deb35e39859eb8"),
		Parent:    common.HexToHash("0x302febf685d86eaa2339e6f9b226e36d69ebf48b1bfd10b44fc51fcaaefbf148"),
		Number:    ethmath.MaxBig256,
		Timestamp: time.Date(1992, time.September, 29, 8, 45, 15, 1230, time.UTC),
		Logs: []ethtypes.Log{
			{
				Address: common.HexToAddress("0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"),
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
					common.HexToHash("0x0000000000000000000000004d8a4aa1f304f9632cf3877473445d85c577fe5d"),
					common.HexToHash("0x0000000000000000000000004bdd0d16cfa18e33860470fc4d65c6f5cee60959"),
				},
				Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000337ad34c0"),
				BlockNumber: 30,
				TxHash:      common.HexToHash("0xd9bb5f9e888ee6f74bedcda811c2461230f247c205849d6f83cb6c3925e54586"),
				TxIndex:     0,
				BlockHash:   common.HexToHash("0x6bbf9b6e836207ab25379c20e517a89090cbbaf8877746f6ed7fb6820770816b"),
				Index:       0,
				Removed:     false,
			},
			{
				Address: common.HexToAddress("0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"),
				Topics: []common.Hash{
					common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
					common.HexToHash("0x0000000000000000000000004d8a4aa1f304f9632cf3877473445d85c577fe5d"),
					common.HexToHash("0x0000000000000000000000004bdd0d16cfa18e33860470fc4d65c6f5cee60959"),
				},
				Data:        common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000deadbeef"),
				BlockNumber: 31,
				TxHash:      common.HexToHash("0xd9bb5f9e888ee6f74bedcda811c2461230f247c205849d6f83cb6c3925e54586"),
				TxIndex:     1,
				BlockHash:   common.HexToHash("0x6bbf9b6e836207ab25379c20e517a89090cbbaf8877746f6ed7fb6820770816b"),
				Index:       2,
				Removed:     true,
			},
		},
	}

	// Convert to JS/Dexie type and back. Make sure we get back the same values that we started with.
	convertedMiniHeader := MiniHeaderToCommonType(MiniHeaderFromCommonType(originalMiniHeader))
	assert.Equal(t, originalMiniHeader, convertedMiniHeader)
}
