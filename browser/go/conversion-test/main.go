// +build js, wasm

package main

import (
	"math/big"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
)

const (
	loadEventName = "0xmeshtest"
)

// This file has a very simple role in the browser conversion tests: create exposed
// functions that the typescript component of the test can access. These functions
// should expose data or functions that the typescript bundle cannot effectively test
// (for example values that have been converted to Javascript using `JSValue` methods).
func main() {
	setGlobals()
	triggerLoadEvent()
	select {}
}

func setGlobals() {
	conversionTestCases := map[string]interface{}{
		"contractEvents": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
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
			}
		}),
		"getOrdersResponse": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return []interface{}{
				types.GetOrdersResponse{
					SnapshotID:        "208c81f9-6f8d-44aa-b6ea-0a3276ec7318",
					SnapshotTimestamp: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC),
					OrdersInfos:       []*types.OrderInfo{},
				},
				types.GetOrdersResponse{
					SnapshotID:        "208c81f9-6f8d-44aa-b6ea-0a3276ec7318",
					SnapshotTimestamp: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC),
					OrdersInfos: []*types.OrderInfo{
						&types.OrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
									MakerAssetAmount:      big.NewInt(123456789),
									MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
									MakerFee:              big.NewInt(89),
									TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
									TakerAssetAmount:      big.NewInt(987654321),
									TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
									TakerFee:              big.NewInt(12),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
							},
							FillableTakerAssetAmount: big.NewInt(987654321),
						},
					},
				},
				types.GetOrdersResponse{
					SnapshotID:        "208c81f9-6f8d-44aa-b6ea-0a3276ec7318",
					SnapshotTimestamp: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC),
					OrdersInfos: []*types.OrderInfo{
						&types.OrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0x"),
									MakerAssetAmount:      big.NewInt(0),
									MakerFeeAssetData:     common.FromHex("0x"),
									MakerFee:              big.NewInt(0),
									TakerAssetData:        common.FromHex("0x"),
									TakerAssetAmount:      big.NewInt(0),
									TakerFeeAssetData:     common.FromHex("0x"),
									TakerFee:              big.NewInt(0),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x"),
							},
							FillableTakerAssetAmount: big.NewInt(0),
						},
						&types.OrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
									MakerAssetAmount:      big.NewInt(123456789),
									MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
									MakerFee:              big.NewInt(89),
									TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
									TakerAssetAmount:      big.NewInt(987654321),
									TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
									TakerFee:              big.NewInt(12),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
							},
							FillableTakerAssetAmount: big.NewInt(987654321),
						},
					},
				},
			}
		}),
		"orderEvents": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return []interface{}{
				zeroex.OrderEvent{
					Timestamp: time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC),
					OrderHash: common.HexToHash("0x1"),
					SignedOrder: &zeroex.SignedOrder{
						Order: zeroex.Order{
							ChainID:               big.NewInt(1337),
							MakerAddress:          common.HexToAddress("0x1"),
							TakerAddress:          common.HexToAddress("0x2"),
							SenderAddress:         common.HexToAddress("0x3"),
							FeeRecipientAddress:   common.HexToAddress("0x4"),
							ExchangeAddress:       common.HexToAddress("0x5"),
							MakerAssetData:        common.FromHex("0x"),
							MakerAssetAmount:      big.NewInt(0),
							MakerFeeAssetData:     common.FromHex("0x"),
							MakerFee:              big.NewInt(0),
							TakerAssetData:        common.FromHex("0x"),
							TakerAssetAmount:      big.NewInt(0),
							TakerFeeAssetData:     common.FromHex("0x"),
							TakerFee:              big.NewInt(0),
							ExpirationTimeSeconds: big.NewInt(10000000000),
							Salt:                  big.NewInt(1532559225),
						},
						Signature: common.FromHex("0x"),
					},
					EndState:                 zeroex.ESOrderAdded,
					FillableTakerAssetAmount: big.NewInt(1),
					ContractEvents:           []*zeroex.ContractEvent{},
				},
				zeroex.OrderEvent{
					Timestamp: time.Date(2006, time.January, 1, 1, 1, 1, 1, time.UTC),
					OrderHash: common.HexToHash("0x1"),
					SignedOrder: &zeroex.SignedOrder{
						Order: zeroex.Order{
							ChainID:               big.NewInt(1337),
							MakerAddress:          common.HexToAddress("0x1"),
							TakerAddress:          common.HexToAddress("0x2"),
							SenderAddress:         common.HexToAddress("0x3"),
							FeeRecipientAddress:   common.HexToAddress("0x4"),
							ExchangeAddress:       common.HexToAddress("0x5"),
							MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
							MakerAssetAmount:      big.NewInt(123456789),
							MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
							MakerFee:              big.NewInt(89),
							TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
							TakerAssetAmount:      big.NewInt(987654321),
							TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
							TakerFee:              big.NewInt(12),
							ExpirationTimeSeconds: big.NewInt(10000000000),
							Salt:                  big.NewInt(1532559225),
						},
						Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
					},
					EndState:                 zeroex.ESOrderFilled,
					FillableTakerAssetAmount: big.NewInt(0),
					ContractEvents: []*zeroex.ContractEvent{
						&zeroex.ContractEvent{
							BlockHash: common.HexToHash("0x1"),
							TxHash:    common.HexToHash("0x2"),
							TxIndex:   123,
							LogIndex:  321,
							IsRemoved: false,
							Address:   common.HexToAddress("0x5"),
							Kind:      "ExchangeFillEvent",
							Parameters: decoder.ExchangeFillEvent{
								MakerAddress:           common.HexToAddress("0x1"),
								TakerAddress:           common.HexToAddress("0x2"),
								SenderAddress:          common.HexToAddress("0x3"),
								FeeRecipientAddress:    common.HexToAddress("0x4"),
								MakerAssetFilledAmount: big.NewInt(123456789),
								TakerAssetFilledAmount: big.NewInt(987654321),
								MakerFeePaid:           big.NewInt(89),
								TakerFeePaid:           big.NewInt(12),
								ProtocolFeePaid:        big.NewInt(150000),
								OrderHash:              common.HexToHash("0x1"),
								MakerAssetData:         common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
								TakerAssetData:         common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
								MakerFeeAssetData:      common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
								TakerFeeAssetData:      common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
							},
						},
					},
				},
			}
		}),
		"signedOrders": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return []interface{}{
				zeroex.SignedOrder{
					Order: zeroex.Order{
						ChainID:               big.NewInt(1337),
						MakerAddress:          common.HexToAddress("0x1"),
						TakerAddress:          common.HexToAddress("0x2"),
						SenderAddress:         common.HexToAddress("0x3"),
						FeeRecipientAddress:   common.HexToAddress("0x4"),
						ExchangeAddress:       common.HexToAddress("0x5"),
						MakerAssetData:        common.FromHex("0x"),
						MakerAssetAmount:      big.NewInt(0),
						MakerFeeAssetData:     common.FromHex("0x"),
						MakerFee:              big.NewInt(0),
						TakerAssetData:        common.FromHex("0x"),
						TakerAssetAmount:      big.NewInt(0),
						TakerFeeAssetData:     common.FromHex("0x"),
						TakerFee:              big.NewInt(0),
						ExpirationTimeSeconds: big.NewInt(10000000000),
						Salt:                  big.NewInt(1532559225),
					},
					Signature: common.FromHex("0x"),
				},
				zeroex.SignedOrder{
					Order: zeroex.Order{
						ChainID:               big.NewInt(1337),
						MakerAddress:          common.HexToAddress("0x1"),
						TakerAddress:          common.HexToAddress("0x2"),
						SenderAddress:         common.HexToAddress("0x3"),
						FeeRecipientAddress:   common.HexToAddress("0x4"),
						ExchangeAddress:       common.HexToAddress("0x5"),
						MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
						MakerAssetAmount:      big.NewInt(123456789),
						MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
						MakerFee:              big.NewInt(89),
						TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
						TakerAssetAmount:      big.NewInt(987654321),
						TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
						TakerFee:              big.NewInt(12),
						ExpirationTimeSeconds: big.NewInt(10000000000),
						Salt:                  big.NewInt(1532559225),
					},
					Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
				},
			}
		}),
		"stats": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return []interface{}{
				types.Stats{
					Version:         "development",
					PubSubTopic:     "someTopic",
					Rendezvous:      "/0x-mesh/network/1337/version/2",
					PeerID:          "16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7",
					EthereumChainID: 1337,
					LatestBlock: types.LatestBlock{
						Hash:   common.HexToHash("0x1"),
						Number: 1500,
					},
					NumPeers:                          200,
					NumOrders:                         100000,
					NumOrdersIncludingRemoved:         200000,
					NumPinnedOrders:                   400,
					MaxExpirationTime:                 "115792089237316195423570985008687907853269984665640564039457584007913129639935",
					StartOfCurrentUTCDay:              time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC),
					EthRPCRequestsSentInCurrentUTCDay: 100000,
					EthRPCRateLimitExpiredRequests:    5000,
				},
			}
		}),
		"validationResults": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return []interface{}{
				ordervalidator.ValidationResults{},
				ordervalidator.ValidationResults{
					Accepted: []*ordervalidator.AcceptedOrderInfo{
						&ordervalidator.AcceptedOrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0x"),
									MakerAssetAmount:      big.NewInt(0),
									MakerFeeAssetData:     common.FromHex("0x"),
									MakerFee:              big.NewInt(0),
									TakerAssetData:        common.FromHex("0x"),
									TakerAssetAmount:      big.NewInt(0),
									TakerFeeAssetData:     common.FromHex("0x"),
									TakerFee:              big.NewInt(0),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x"),
							},
							FillableTakerAssetAmount: big.NewInt(0),
							IsNew:                    true,
						},
					},
				},
				ordervalidator.ValidationResults{
					Rejected: []*ordervalidator.RejectedOrderInfo{
						&ordervalidator.RejectedOrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0x"),
									MakerAssetAmount:      big.NewInt(0),
									MakerFeeAssetData:     common.FromHex("0x"),
									MakerFee:              big.NewInt(0),
									TakerAssetData:        common.FromHex("0x"),
									TakerAssetAmount:      big.NewInt(0),
									TakerFeeAssetData:     common.FromHex("0x"),
									TakerFee:              big.NewInt(0),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x"),
							},
							Kind:   ordervalidator.ZeroExValidation,
							Status: ordervalidator.ROInvalidMakerAssetData,
						},
					},
				},
				ordervalidator.ValidationResults{
					Accepted: []*ordervalidator.AcceptedOrderInfo{
						&ordervalidator.AcceptedOrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0x"),
									MakerAssetAmount:      big.NewInt(0),
									MakerFeeAssetData:     common.FromHex("0x"),
									MakerFee:              big.NewInt(0),
									TakerAssetData:        common.FromHex("0x"),
									TakerAssetAmount:      big.NewInt(0),
									TakerFeeAssetData:     common.FromHex("0x"),
									TakerFee:              big.NewInt(0),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x"),
							},
							FillableTakerAssetAmount: big.NewInt(0),
							IsNew:                    true,
						},
						&ordervalidator.AcceptedOrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
									MakerAssetAmount:      big.NewInt(123456789),
									MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
									MakerFee:              big.NewInt(89),
									TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
									TakerAssetAmount:      big.NewInt(987654321),
									TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
									TakerFee:              big.NewInt(12),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
							},
							FillableTakerAssetAmount: big.NewInt(987654321),
							IsNew:                    true,
						},
					},
					Rejected: []*ordervalidator.RejectedOrderInfo{
						&ordervalidator.RejectedOrderInfo{
							OrderHash: common.HexToHash("0x1"),
							SignedOrder: &zeroex.SignedOrder{
								Order: zeroex.Order{
									ChainID:               big.NewInt(1337),
									MakerAddress:          common.HexToAddress("0x1"),
									TakerAddress:          common.HexToAddress("0x2"),
									SenderAddress:         common.HexToAddress("0x3"),
									FeeRecipientAddress:   common.HexToAddress("0x4"),
									ExchangeAddress:       common.HexToAddress("0x5"),
									MakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
									MakerAssetAmount:      big.NewInt(123456789),
									MakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
									MakerFee:              big.NewInt(89),
									TakerAssetData:        common.FromHex("0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
									TakerAssetAmount:      big.NewInt(987654321),
									TakerFeeAssetData:     common.FromHex("0xf47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
									TakerFee:              big.NewInt(12),
									ExpirationTimeSeconds: big.NewInt(10000000000),
									Salt:                  big.NewInt(1532559225),
								},
								Signature: common.FromHex("0x012761a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
							},
							Kind:   ordervalidator.MeshError,
							Status: ordervalidator.ROCoordinatorEndpointNotFound,
						},
					},
				},
			}
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
