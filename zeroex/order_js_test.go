// +build js,wasm

package zeroex

import (
	"fmt"
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
		TxIndex:   0,
		LogIndex:  0,
		IsRemoved: true,
		Address:   contracts.WETH9,
		Kind:      "ERC20TransferEvent",
		Parameters: decoder.ERC20TransferEvent{
			From:  constants.GanacheAccount0,
			To:    constants.GanacheAccount0,
			Value: big.NewInt(1),
		},
	}.JSValue()
	fmt.Printf("%+v", jsContractEvent)
	t.Logf("%+v", jsContractEvent)
	t.Fail()
}
