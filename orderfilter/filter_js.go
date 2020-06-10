// +build js, wasm

package orderfilter

import (
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
)

type Filter struct {
	orderValidator       js.Value
	messageValidator     js.Value
	encodedSchema        string
	chainID              int
	rawCustomOrderSchema string
}

func New(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	chainIDSchema := fmt.Sprintf(`{"$id": "/chainId", "const":%d}`, chainID)
	exchangeAddressSchema := fmt.Sprintf(`{"$id": "/exchangeAddress", "enum":[%q,%q]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))

	if jsutil.IsNullOrUndefined(js.Global().Get("__mesh_createSchemaValidator__")) {
		return nil, errors.New(`"__mesh_createSchemaValidator__" has not been set on the Javascript "global" object`)
	}
	// NOTE(jalextowle): The order of the schemas within the two arrays
	// defines their order of compilation.
	schemaValidator := js.Global().Call(
		"__mesh_createSchemaValidator__",
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
	orderValidator := schemaValidator.Get("orderValidator")
	if jsutil.IsNullOrUndefined(orderValidator) {
		return nil, errors.New(`"orderValidator" has not been set on the provided "schemaValidator"`)
	}
	messageValidator := schemaValidator.Get("messageValidator")
	if jsutil.IsNullOrUndefined(messageValidator) {
		return nil, errors.New(`"messageValidator" has not been set on the provided "schemaValidator"`)
	}
	return &Filter{
		orderValidator:       orderValidator,
		messageValidator:     messageValidator,
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
	}, nil
}
