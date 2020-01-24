// +build js, wasm

package main

import (
	"math/big"
	"syscall/js"

	"github.com/0xProject/0x-mesh/common/types"
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
