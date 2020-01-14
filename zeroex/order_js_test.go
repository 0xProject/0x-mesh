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
	"github.com/ethereum/go-ethereum/common"
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
		eventObject := testCase.event.JSValue()
		eventString := stringify(eventObject)
		var eventJSON contractEventJSON
		err := json.Unmarshal([]byte(eventString), &eventJSON)
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

func TestSignedOrderConversion(t *testing.T) {
	for _, order := range []SignedOrder{
		SignedOrder{
			Order: Order{
				MakerAddress:          constants.NullAddress,
				TakerAddress:          constants.NullAddress,
				SenderAddress:         constants.NullAddress,
				FeeRecipientAddress:   constants.NullAddress,
				MakerAssetData:        common.FromHex(""),
				MakerAssetAmount:      big.NewInt(0),
				MakerFeeAssetData:     common.FromHex(""),
				MakerFee:              big.NewInt(0),
				TakerAssetData:        common.FromHex(""),
				TakerAssetAmount:      big.NewInt(0),
				TakerFeeAssetData:     common.FromHex(""),
				TakerFee:              big.NewInt(0),
				ChainID:               big.NewInt(1337),
				ExpirationTimeSeconds: big.NewInt(0),
				Salt:                  big.NewInt(0),
			},
			Signature: common.FromHex(""),
		},
		SignedOrder{
			Order: Order{
				MakerAddress:          constants.GanacheAccount0,
				TakerAddress:          constants.NullAddress,
				SenderAddress:         constants.NullAddress,
				FeeRecipientAddress:   constants.GanacheAccount4,
				MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
				MakerAssetAmount:      big.NewInt(10000000),
				MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
				MakerFee:              big.NewInt(10000000),
				TakerAssetData:        common.FromHex("0xf47261b000000000000000000000000081228eA33D680B0F51271aBAb1105886eCd01C2c"),
				TakerAssetAmount:      big.NewInt(10000000),
				TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
				TakerFee:              big.NewInt(10000000),
				ChainID:               big.NewInt(1337),
				ExpirationTimeSeconds: big.NewInt(0),
				Salt:                  big.NewInt(0),
			},
			Signature: common.Hex2Bytes("0x1befcf4b6b1da4d207067a4b06e9bfbf21f85e2b6644f3ecf3a15f009e484756f251e3e00e909447ce45a16c620d14920a9acf516d9f4fe45bc36c914be6c9ec2703"),
		},
	} {
		orderObject := order.JSValue()
		orderString := stringify(orderObject)
		recoveredOrder := &SignedOrder{}
		err := recoveredOrder.UnmarshalJSON([]byte(orderString))
		require.NoError(t, err)

		if len(order.MakerAssetData) == 0 {
			require.Equal(t, "0x", orderObject.Get("makerAssetData").String())
		}
		if len(order.MakerFeeAssetData) == 0 {
			require.Equal(t, "0x", orderObject.Get("makerFeeAssetData").String())
		}
		if len(order.TakerAssetData) == 0 {
			require.Equal(t, "0x", orderObject.Get("takerAssetData").String())
		}
		if len(order.TakerFeeAssetData) == 0 {
			require.Equal(t, "0x", orderObject.Get("takerFeeAssetData").String())
		}
		if len(order.Signature) == 0 {
			require.Equal(t, "0x", orderObject.Get("signature").String())
		}

		require.Equal(t, order, *recoveredOrder)
	}
}

func stringify(value js.Value) string {
	return js.Global().Get("JSON").Call("stringify", value).String()
}
