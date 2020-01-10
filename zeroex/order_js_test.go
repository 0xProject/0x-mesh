// +build js,wasm

package zeroex

import (
	"fmt"
	"math/big"
	"syscall/js"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestContractEvent(t *testing.T) {
	contracts, err := ethereum.GetContractAddressesForChainID(1337)
	require.NoError(t, err)

	for _, event := range []ContractEvent{
		ContractEvent{
			BlockHash: common.HexToHash("0x1"),
			TxHash:    common.HexToHash("0x2"),
			TxIndex:   1,
			LogIndex:  2,
			IsRemoved: true,
			Address:   contracts.WETH9,
			Kind:      "ERC20TransferEvent",
			Parameters: decoder.ERC20TransferEvent{
				From:  constants.GanacheAccount0,
				To:    constants.GanacheAccount1,
				Value: big.NewInt(3),
			},
		},
		ContractEvent{
			BlockHash: constants.GanacheAccount1.Hash(),
			TxHash:    constants.GanacheAccount2.Hash(),
			TxIndex:   342424,
			LogIndex:  1000,
			IsRemoved: false,
			// NOTE(jalextowle): The ERC1155Proxy doesn't actually emit this event,
			// but that isn't important in the context of this test.
			Address: contracts.ERC1155Proxy,
			Kind:    "ERC1155TransferSingleEvent",
			Parameters: decoder.ERC1155TransferSingleEvent{
				Operator: constants.GanacheAccount1,
				From:     constants.GanacheAccount2,
				To:       constants.GanacheAccount3,
				Id:       big.NewInt(32423),
				Value:    big.NewInt(10000),
			},
		},
	} {
		jsEvent := event.JSValue()
		require.Equal(t, jsEvent.Get("address").String(), event.Address.Hex())
		require.Equal(t, jsEvent.Get("blockHash").String(), event.BlockHash.Hex())
		require.Equal(t, jsEvent.Get("txHash").String(), event.TxHash.Hex())
		require.Equal(t, jsEvent.Get("txIndex").Int(), int(event.TxIndex))
		require.Equal(t, jsEvent.Get("logIndex").Int(), int(event.LogIndex))
		require.Equal(t, jsEvent.Get("isRemoved").Bool(), event.IsRemoved)
		require.Equal(t, jsEvent.Get("kind").String(), event.Kind)
		assertDeepEqual(t, jsEvent.Get("parameters"), event.Parameters.(js.Wrapper).JSValue())
	}
}

func assertDeepEqual(t *testing.T, value js.Value, other js.Value) bool {
	keys := js.Global().Get("Object").Get("keys").Invoke(value)
	otherKeys := js.Global().Get("Object").Get("keys").Invoke(other)
	length := keys.Get("length").Int()
	otherLength := otherKeys.Get("length").Int()
	require.Equal(t, length, otherLength)
	for i := 0; i < length; i++ {
		k := keys.Index(i)
		isInBothArrays := false
		fmt.Println(js.Global().Get("JSON").Get("stringify").Invoke(k).String())
		for j := 0; j < otherLength; j++ {
			oK := otherKeys.Index(j)
			if k.String() == oK.String() {
				isInBothArrays = true
				require.True(t)
				break
			}
		}
		require.True(t, isInBothArrays)
	}
	t.Fail()
	return true
}
