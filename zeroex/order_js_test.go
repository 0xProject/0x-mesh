// +build js,wasm

package zeroex

import (
	"encoding/json"
	"errors"
	"math/big"
	"syscall/js"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/stretchr/testify/require"
)

var (
	blockHash        = constants.GanacheAccount0.Hash()
	txHash           = constants.GanacheAccount1.Hash()
	txIndex     uint = 1
	logIndex    uint = 2
	address          = constants.GanacheAccount2
	one              = big.NewInt(1)
	protocolFee      = big.NewInt(150000)
	id          *big.Int
)

func init() {
	var success bool
	id, success = (&big.Int{}).SetString("0xdeadbeef", 0)
	if !success {
		panic("Failed to set id to 0xdeadbeef")
	}
}

func TestContractEventConversion(t *testing.T) {
	for _, testCase := range []struct {
		event ContractEvent
		err   error
	}{
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC20ApprovalEvent",
				Parameters: decoder.ERC20ApprovalEvent{
					Owner:   constants.GanacheAccount3,
					Spender: constants.GanacheAccount4,
					Value:   one,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC20TransferEvent",
				Parameters: decoder.ERC20TransferEvent{
					From:  constants.GanacheAccount3,
					To:    constants.GanacheAccount4,
					Value: one,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC721ApprovalEvent",
				Parameters: decoder.ERC721ApprovalEvent{
					Owner:    constants.GanacheAccount3,
					Approved: constants.GanacheAccount4,
					TokenId:  id,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC721ApprovalForAllEvent",
				Parameters: decoder.ERC721ApprovalForAllEvent{
					Owner:    constants.GanacheAccount3,
					Operator: constants.GanacheAccount4,
					Approved: true,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC721TransferEvent",
				Parameters: decoder.ERC721TransferEvent{
					From:    constants.GanacheAccount3,
					To:      constants.GanacheAccount4,
					TokenId: id,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC1155ApprovalForAllEvent",
				Parameters: decoder.ERC1155ApprovalForAllEvent{
					Owner:    constants.GanacheAccount3,
					Operator: constants.GanacheAccount4,
					Approved: false,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC1155TransferSingleEvent",
				Parameters: decoder.ERC1155TransferSingleEvent{
					Operator: constants.GanacheAccount2,
					From:     constants.GanacheAccount3,
					To:       constants.GanacheAccount4,
					Id:       id,
					Value:    one,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ERC1155TransferBatchEvent",
				Parameters: decoder.ERC1155TransferBatchEvent{
					Operator: constants.GanacheAccount2,
					From:     constants.GanacheAccount3,
					To:       constants.GanacheAccount4,
					Ids:      []*big.Int{id},
					Values:   []*big.Int{one},
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ExchangeFillEvent",
				Parameters: decoder.ExchangeFillEvent{
					MakerAddress:           constants.GanacheAccount0,
					TakerAddress:           constants.NullAddress,
					SenderAddress:          constants.GanacheAccount0,
					FeeRecipientAddress:    constants.GanacheAccount1,
					MakerAssetFilledAmount: one,
					TakerAssetFilledAmount: one,
					MakerFeePaid:           one,
					TakerFeePaid:           one,
					ProtocolFeePaid:        one,
					OrderHash:              constants.GanacheAccount2.Hash(),
					MakerAssetData:         constants.NullBytes,
					TakerAssetData:         constants.NullBytes,
					MakerFeeAssetData:      constants.NullBytes,
					TakerFeeAssetData:      constants.NullBytes,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ExchangeCancelEvent",
				Parameters: decoder.ExchangeCancelEvent{
					MakerAddress:        constants.GanacheAccount0,
					SenderAddress:       constants.GanacheAccount2,
					FeeRecipientAddress: constants.GanacheAccount1,
					OrderHash:           constants.GanacheAccount2.Hash(),
					MakerAssetData:      constants.NullBytes,
					TakerAssetData:      constants.NullBytes,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "ExchangeCancelUpToEvent",
				Parameters: decoder.ExchangeCancelUpToEvent{
					MakerAddress:       constants.GanacheAccount0,
					OrderSenderAddress: constants.GanacheAccount2,
					OrderEpoch:         big.NewInt(50),
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "WethDepositEvent",
				Parameters: decoder.WethDepositEvent{
					Owner: constants.GanacheAccount0,
					Value: protocolFee,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "WethWithdrawalEvent",
				Parameters: decoder.WethWithdrawalEvent{
					Owner: constants.GanacheAccount0,
					Value: protocolFee,
				},
			},
			err: nil,
		},
		{
			event: ContractEvent{
				BlockHash: blockHash,
				TxHash:    txHash,
				TxIndex:   txIndex,
				LogIndex:  logIndex,
				IsRemoved: false,
				Address:   constants.GanacheAccount2,
				Kind:      "FooBarBazEvent",
				// NOTE(jalextowle): We have to use something non-empty
				// that implements `js.Wrapper` or else we'll experience
				// a runtime panic.
				Parameters: decoder.ERC20ApprovalEvent{
					Owner:   constants.GanacheAccount3,
					Spender: constants.GanacheAccount4,
					Value:   one,
				},
			},
			err: errors.New("unknown event kind: FooBarBazEvent"),
		},
	} {
		// Convert the contract event to a JSValue, and then attempt to
		// recover it from the stringified JSON.
		jsEvent := testCase.event.JSValue()
		jsString := js.Global().Get("JSON").Call("stringify", jsEvent).String()
		var eventJSON contractEventJSON
		err := json.Unmarshal([]byte(jsString), &eventJSON)
		require.NoError(t, err)
		decodedEvent, err := unmarshalContractEvent(&eventJSON)
		if testCase.err != nil {
			require.Equal(t, testCase.err, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, testCase.event, *decodedEvent)
		}
	}
}
