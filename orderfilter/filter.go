// +build !js

package orderfilter

import (
	"fmt"
	"strings"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/ethereum/go-ethereum/common"
	jsonschema "github.com/xeipuuv/gojsonschema"
)

var (
	// Built-in schemas
	addressSchemaLoader     = jsonschema.NewStringLoader(addressSchema)
	wholeNumberSchemaLoader = jsonschema.NewStringLoader(wholeNumberSchema)
	hexSchemaLoader         = jsonschema.NewStringLoader(hexSchema)
	orderSchemaLoader       = jsonschema.NewStringLoader(orderSchema)
	signedOrderSchemaLoader = jsonschema.NewStringLoader(signedOrderSchema)

	// Root schemas
	rootOrderSchemaLoader        = jsonschema.NewStringLoader(rootOrderSchema)
	rootOrderMessageSchemaLoader = jsonschema.NewStringLoader(rootOrderMessageSchema)
)

var builtInSchemas = []jsonschema.JSONLoader{
	addressSchemaLoader,
	wholeNumberSchemaLoader,
	hexSchemaLoader,
	orderSchemaLoader,
	signedOrderSchemaLoader,
}

type Filter struct {
	encodedSchema        string
	chainID              int
	rawCustomOrderSchema string
	orderSchema          *jsonschema.Schema
	messageSchema        *jsonschema.Schema
	exchangeAddress      common.Address
}

// TODO(jalextowle): We do not need `contractAddresses` since we only use `contractAddresses.Exchange`.
// In a future refactor, we should update this interface.
func New(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	orderLoader, err := newLoader(chainID, customOrderSchema, contractAddresses)
	if err != nil {
		return nil, err
	}
	compiledRootOrderSchema, err := orderLoader.Compile(rootOrderSchemaLoader)
	if err != nil {
		return nil, err
	}

	messageLoader, err := newLoader(chainID, customOrderSchema, contractAddresses)
	if err != nil {
		return nil, err
	}
	if err := messageLoader.AddSchemas(rootOrderSchemaLoader); err != nil {
		return nil, err
	}
	compiledRootOrderMessageSchema, err := messageLoader.Compile(rootOrderMessageSchemaLoader)
	if err != nil {
		return nil, err
	}
	return &Filter{
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
		orderSchema:          compiledRootOrderSchema,
		messageSchema:        compiledRootOrderMessageSchema,
		exchangeAddress:      contractAddresses.Exchange,
	}, nil
}

func loadExchangeAddress(loader *jsonschema.SchemaLoader, contractAddresses ethereum.ContractAddresses) error {
	// Note that exchangeAddressSchema accepts both checksummed and
	// non-checksummed (i.e. all lowercase) addresses.
	exchangeAddressSchema := fmt.Sprintf(`{"enum":[%q,%q]}`, contractAddresses.Exchange.Hex(), strings.ToLower(contractAddresses.Exchange.Hex()))
	return loader.AddSchema("/exchangeAddress", jsonschema.NewStringLoader(exchangeAddressSchema))
}

func loadChainID(loader *jsonschema.SchemaLoader, chainID int) error {
	chainIDSchema := fmt.Sprintf(`{"const":%d}`, chainID)
	return loader.AddSchema("/chainId", jsonschema.NewStringLoader(chainIDSchema))
}

func newLoader(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadChainID(loader, chainID); err != nil {
		return nil, err
	}
	if err := loadExchangeAddress(loader, contractAddresses); err != nil {
		return nil, err
	}
	if err := loader.AddSchemas(builtInSchemas...); err != nil {
		return nil, err
	}
	if err := loader.AddSchema("/customOrder", jsonschema.NewStringLoader(customOrderSchema)); err != nil {
		return nil, err
	}
	return loader, nil
}
