// +build js, wasm

package main

import (
	"math/big"
	"syscall/js"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
)

const (
	loadEventName = "0xmeshtest"
)

func main() {
	setGlobals()
	triggerLoadEvent()
	select {}
}

func setGlobals() {
	conversionTestCases := map[string]interface{}{
		"contractEventsAsync": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return types.WrapInPromise(func() (interface{}, error) {
				return []interface{}{
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC20ApprovalEvent",
						Parameters: decoder.ERC20ApprovalEvent{
							Owner:   common.HexToAddress("0x4"),
							Spender: common.HexToAddress("0x5"),
							Value:   big.NewInt(1000),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC20TransferEvent",
						Parameters: decoder.ERC20TransferEvent{
							From:  common.HexToAddress("0x4"),
							To:    common.HexToAddress("0x5"),
							Value: big.NewInt(1000),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC721ApprovalEvent",
						Parameters: decoder.ERC721ApprovalEvent{
							Owner:    common.HexToAddress("0x4"),
							Approved: common.HexToAddress("0x5"),
							TokenId:  big.NewInt(1),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC721ApprovalForAllEvent",
						Parameters: decoder.ERC721ApprovalForAllEvent{
							Owner:    common.HexToAddress("0x4"),
							Operator: common.HexToAddress("0x5"),
							Approved: true,
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC721TransferEvent",
						Parameters: decoder.ERC721TransferEvent{
							From:    common.HexToAddress("0x4"),
							To:      common.HexToAddress("0x5"),
							TokenId: big.NewInt(1),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC1155ApprovalForAllEvent",
						Parameters: decoder.ERC1155ApprovalForAllEvent{
							Owner:    common.HexToAddress("0x4"),
							Operator: common.HexToAddress("0x5"),
							Approved: false,
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC1155TransferSingleEvent",
						Parameters: decoder.ERC1155TransferSingleEvent{
							Operator: common.HexToAddress("0x4"),
							From:     common.HexToAddress("0x5"),
							To:       common.HexToAddress("0x6"),
							Id:       big.NewInt(1),
							Value:    big.NewInt(100),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ERC1155TransferBatchEvent",
						Parameters: decoder.ERC1155TransferBatchEvent{
							Operator: common.HexToAddress("0x4"),
							From:     common.HexToAddress("0x5"),
							To:       common.HexToAddress("0x6"),
							Ids:      []*big.Int{big.NewInt(1)},
							Values:   []*big.Int{big.NewInt(100)},
						},
					},
					// FIXME(jalextowle): Should I include another event with non-null asset data?
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ExchangeFillEvent",
						Parameters: decoder.ExchangeFillEvent{
							MakerAddress:           common.HexToAddress("0x4"),
							TakerAddress:           constants.NullAddress,
							SenderAddress:          common.HexToAddress("0x5"),
							FeeRecipientAddress:    common.HexToAddress("0x6"),
							MakerAssetFilledAmount: big.NewInt(456),
							TakerAssetFilledAmount: big.NewInt(654),
							MakerFeePaid:           big.NewInt(12),
							TakerFeePaid:           big.NewInt(21),
							ProtocolFeePaid:        big.NewInt(150000),
							OrderHash:              common.HexToHash("0x7"),
							MakerAssetData:         constants.NullBytes,
							TakerAssetData:         constants.NullBytes,
							MakerFeeAssetData:      constants.NullBytes,
							TakerFeeAssetData:      constants.NullBytes,
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ExchangeCancelEvent",
						Parameters: decoder.ExchangeCancelEvent{
							MakerAddress:        common.HexToAddress("0x4"),
							SenderAddress:       common.HexToAddress("0x5"),
							FeeRecipientAddress: common.HexToAddress("0x6"),
							OrderHash:           common.HexToHash("0x7"),
							MakerAssetData:      constants.NullBytes,
							TakerAssetData:      constants.NullBytes,
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "ExchangeCancelUpToEvent",
						Parameters: decoder.ExchangeCancelUpToEvent{
							MakerAddress:       common.HexToAddress("0x4"),
							OrderSenderAddress: common.HexToAddress("0x5"),
							OrderEpoch:         big.NewInt(50),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "WethDepositEvent",
						Parameters: decoder.WethDepositEvent{
							Owner: common.HexToAddress("0x4"),
							Value: big.NewInt(150000),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "WethWithdrawalEvent",
						Parameters: decoder.WethWithdrawalEvent{
							Owner: common.HexToAddress("0x4"),
							Value: big.NewInt(150000),
						},
					},
					zeroex.ContractEvent{
						BlockHash: common.HexToHash("0x1"),
						TxHash:    common.HexToHash("0x2"),
						TxIndex:   123,
						LogIndex:  321,
						IsRemoved: false,
						Address:   common.HexToAddress("0x3"),
						Kind:      "FooBarBazEvent",
						// NOTE(jalextowle): We have to use something non-empty
						// that implements `js.Wrapper` or else we'll experience
						// a runtime panic.
						Parameters: decoder.ERC20ApprovalEvent{
							Owner:   common.HexToAddress("0x4"),
							Spender: common.HexToAddress("0x5"),
							Value:   big.NewInt(1),
						},
					},
				}, nil
			})
		}),
	}
	js.Global().Set("conversionTestCases", conversionTestCases)
}

// triggerLoadEvent triggers the global load event to indicate that the Wasm is
// done loading.
func triggerLoadEvent() {
	event := js.Global().Get("document").Call("createEvent", "Event")
	event.Call("initEvent", loadEventName, true, true)
	js.Global().Call("dispatchEvent", event)
}
