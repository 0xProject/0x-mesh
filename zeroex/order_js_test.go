// +build js,wasm

package zeroex

import (
	"math/big"
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
	jsContractEvent := ContractEvent{
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
	}.JSValue()
	require.Equal(t, jsContractEvent.Get("blockHash").String(), common.HexToHash("0x1").Hex())
	require.Equal(t, jsContractEvent.Get("txHash").String(), common.HexToHash("0x2").Hex())
	require.Equal(t, jsContractEvent.Get("txIndex").Int(), 1)
	require.Equal(t, jsContractEvent.Get("logIndex").Int(), 2)
	require.Equal(t, jsContractEvent.Get("isRemoved").Bool(), true)
	require.Equal(t, jsContractEvent.Get("kind").String(), "ERC20TransferEvent")
	require.Equal(t, jsContractEvent.Get("parameters").Get("from").String(), constants.GanacheAccount0.Hex())
	require.Equal(t, jsContractEvent.Get("parameters").Get("to").String(), constants.GanacheAccount1.Hex())
	require.Equal(t, jsContractEvent.Get("parameters").Get("value").String(), "3")
}
