// +build js, wasm

package orderfilter

import (
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	"github.com/ethereum/go-ethereum/common"
)

// TODO(mason) have two filters instead of this craziness!!
type Filter struct {
	orderValidatorV3       js.Value
	orderValidatorV4       js.Value
	messageValidatorV3     js.Value
	messageValidatorV4     js.Value
	encodedSchemaV3        string
	encodedSchemaV4        string
	chainID                int
	rawCustomOrderSchemaV3 string
	rawCustomOrderSchemaV4 string
	exchangeAddressV3      common.Address
	exchangeAddressV4      common.Address
}

func New(chainID int, customOrderSchemaV3 string, customOrderSchemaV4 string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	chainIDSchema := fmt.Sprintf(`{"$id": "/chainId", "const":%d}`, chainID)
	exchangeAddressSchema := fmt.Sprintf(`{"$id": "/exchangeAddress", "enum":[%q,%q]}`, contractAddresses.ExchangeV3.Hex(), strings.ToLower(contractAddresses.ExchangeV3.Hex()))

	if jsutil.IsNullOrUndefined(js.Global().Get("__mesh_createSchemaValidator__")) {
		return nil, errors.New(`"__mesh_createSchemaValidator__" has not been set on the Javascript "global" object`)
	}
	// NOTE(jalextowle): The order of the schemas within the two arrays
	// defines their order of compilation.
	schemaValidatorV3 := js.Global().Call(
		"__mesh_createSchemaValidator__",
		customOrderSchemaV3,
		[]interface{}{
			addressSchema,
			wholeNumberSchema,
			hexSchema,
			chainIDSchema,
			exchangeAddressSchema,
			orderSchemaV3,
			signedOrderSchemaV3,
		},
		[]interface{}{
			rootOrderSchemaV3,
			rootOrderMessageSchemaV3,
		})

	schemaValidatorV4 := js.Global().Call(
		"__mesh_createSchemaValidator__",
		customOrderSchemaV4,
		[]interface{}{
			addressSchema,
			wholeNumberSchema,
			hexSchema,
			chainIDSchema,
			exchangeAddressSchema,
			orderSchemaV4,
			signedOrderSchemaV4,
		},
		[]interface{}{
			rootOrderSchemaV4,
			rootOrderMessageSchemaV4,
		})

	orderValidatorV3 := schemaValidatorV3.Get("orderValidator")
	orderValidatorV4 := schemaValidatorV4.Get("orderValidator")

	if jsutil.IsNullOrUndefined(orderValidatorV3) {
		return nil, errors.New(`"orderValidator" has not been set on the provided "schemaValidatorV3"`)
	}
	if jsutil.IsNullOrUndefined(orderValidatorV4) {
		return nil, errors.New(`"orderValidator" has not been set on the provided "schemaValidatorV4"`)
	}

	messageValidatorV3 := schemaValidatorV3.Get("messageValidator")
	if jsutil.IsNullOrUndefined(messageValidatorV3) {
		return nil, errors.New(`"messageValidator" has not been set on the provided "schemaValidatorV3"`)
	}
	messageValidatorV4 := schemaValidatorV4.Get("messageValidator")
	if jsutil.IsNullOrUndefined(messageValidatorV4) {
		return nil, errors.New(`"messageValidator" has not been set on the provided "schemaValidatorV3"`)
	}

	return &Filter{
		orderValidatorV3:       orderValidatorV3,
		messageValidatorV3:     messageValidatorV3,
		chainID:                chainID,
		rawCustomOrderSchemaV3: customOrderSchemaV3,
		exchangeAddressV3:      contractAddresses.ExchangeV3,
		orderValidatorV4:       orderValidatorV4,
		messageValidatorV4:     messageValidatorV4,
		rawCustomOrderSchemaV4: customOrderSchemaV4,
		exchangeAddressV4:      contractAddresses.ExchangeV4,
	}, nil
}
