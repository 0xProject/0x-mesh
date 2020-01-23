// +build js, wasm

package main

import (
	"math/big"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
)

const (
	loadEventName = "0xmeshtest"
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

func main() {
	/*
		conversionTestCases := map[string]interface{}{
			"data": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				return types.WrapInPromise(func() (interface{}, error) {
					return 1, nil
				})
			}),
		}
		js.Global().Set("conversionTestCases", conversionTestCases)
	*/
	setGlobals()
	triggerLoadEvent()
	select {}
}

func setGlobals() {
	order := zeroex.SignedOrder{
		Order: zeroex.Order{
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
	}
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		panic("BOO")
	}
	order.ResetHash()
	conversionTestCases := map[string]interface{}{
		"orderEvents": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return types.WrapInPromise(func() (interface{}, error) {
				return []interface{}{
					zeroex.OrderEvent{
						Timestamp:                time.Now().UTC(),
						OrderHash:                orderHash,
						SignedOrder:              &order,
						EndState:                 "ADDED",
						FillableTakerAssetAmount: big.NewInt(10000000),
						ContractEvents: []*zeroex.ContractEvent{
							&zeroex.ContractEvent{
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
						},
					},
				}, nil
			})
		}),
		"a": 1,
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
