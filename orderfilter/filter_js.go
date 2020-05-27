// +build js, wasm

package orderfilter

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/ethereum"
)

type Filter struct {
	validatorLoaded      bool
	encodedSchema        string
	chainID              int
	rawCustomOrderSchema string
}

func New(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	chainIDSchema := fmt.Sprintf(`{"$id": "/chainId", "const":%d}`, chainID)
	exchangeAddressSchema := fmt.Sprintf(`{"$id": "/exchangeAddress", "enum":[%q,%q]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))
	// NOTE(jalextowle): The order of the schemas within the two arrays
	// defines their order of compilation.
	js.Global().Call(
		"setSchemaValidator",
		customOrderSchema,
		[]interface{}{
			addressSchema,
			wholeNumberSchema,
			hexSchema,
			chainIDSchema,
			exchangeAddressSchema,
			orderSchema,
			signedOrderSchema,
		},
		[]interface{}{
			rootOrderSchema,
			rootOrderMessageSchema,
		})
	return &Filter{
		validatorLoaded:      true,
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
	}, nil
}
