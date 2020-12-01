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

	// V3 Order Schemas
	orderSchemaV3Loader       = jsonschema.NewStringLoader(orderSchemaV3)
	signedOrderSchemaV3Loader = jsonschema.NewStringLoader(signedOrderSchemaV3)

	// V4 Order Schemas
	orderSchemaV4Loader       = jsonschema.NewStringLoader(orderSchemaV4)
	signedOrderSchemaV4Loader = jsonschema.NewStringLoader(signedOrderSchemaV4)

	// V3 Root schemas
	rootOrderSchemaV3Loader        = jsonschema.NewStringLoader(rootOrderSchemaV3)
	rootOrderMessageSchemaV3Loader = jsonschema.NewStringLoader(rootOrderMessageSchemaV3)

	// V4 Root schemas
	rootOrderSchemaV4Loader        = jsonschema.NewStringLoader(rootOrderSchemaV4)
	rootOrderMessageSchemaV4Loader = jsonschema.NewStringLoader(rootOrderMessageSchemaV4)
)

var builtInSchemasV3 = []jsonschema.JSONLoader{
	addressSchemaLoader,
	wholeNumberSchemaLoader,
	hexSchemaLoader,
	orderSchemaV3Loader,
	signedOrderSchemaV3Loader,
}

var builtInSchemasV4 = []jsonschema.JSONLoader{
	addressSchemaLoader,
	wholeNumberSchemaLoader,
	hexSchemaLoader,
	orderSchemaV4Loader,
	signedOrderSchemaV4Loader,
}

type Filter struct {
	encodedSchemaV3        string
	encodedSchemaV4        string
	chainID                int
	rawCustomOrderSchemaV3 string
	orderSchemaV3          *jsonschema.Schema
	messageSchemaV3        *jsonschema.Schema
	rawCustomOrderSchemaV4 string
	orderSchemaV4          *jsonschema.Schema
	messageSchemaV4        *jsonschema.Schema
	exchangeAddressV3      common.Address
	exchangeAddressV4      common.Address
}

// FIXME(jalextowle): This will need to be able to handle orderfilters that
// only have v3 information so that v11 and < v11 can continue to communicate
// using ordersync
//
// TODO(jalextowle): We do not need `contractAddresses` since we only use `contractAddresses.Exchange`.
// In a future refactor, we should update this interface.
func New(chainID int, customOrderSchemaV3 string, customOrderSchemaV4 string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	// Create v3 loaders
	orderV3Loader, err := newV3Loader(chainID, customOrderSchemaV3, contractAddresses)
	if err != nil {
		return nil, err
	}
	compiledRootOrderSchemaV3, err := orderV3Loader.Compile(rootOrderSchemaV3Loader)
	if err != nil {
		return nil, err
	}
	messageV3Loader, err := newV3Loader(chainID, customOrderSchemaV3, contractAddresses)
	if err != nil {
		return nil, err
	}
	if err := messageV3Loader.AddSchemas(rootOrderSchemaV3Loader); err != nil {
		return nil, err
	}
	compiledRootOrderMessageSchemaV3, err := messageV3Loader.Compile(rootOrderMessageSchemaV3Loader)
	if err != nil {
		return nil, err
	}

	// Create v4 loaders
	orderV4Loader, err := newV4Loader(chainID, customOrderSchemaV4, contractAddresses)
	if err != nil {
		return nil, err
	}
	compiledRootOrderSchemaV4, err := orderV4Loader.Compile(rootOrderSchemaV4Loader)
	if err != nil {
		return nil, err
	}
	messageV4Loader, err := newV4Loader(chainID, customOrderSchemaV4, contractAddresses)
	if err != nil {
		return nil, err
	}
	if err := messageV4Loader.AddSchemas(rootOrderSchemaV4Loader); err != nil {
		return nil, err
	}
	compiledRootOrderMessageSchemaV4, err := messageV4Loader.Compile(rootOrderMessageSchemaV4Loader)
	if err != nil {
		return nil, err
	}

	// Create the order filter
	return &Filter{
		chainID:                chainID,
		rawCustomOrderSchemaV3: customOrderSchemaV3,
		orderSchemaV3:          compiledRootOrderSchemaV3,
		messageSchemaV3:        compiledRootOrderMessageSchemaV3,
		rawCustomOrderSchemaV4: customOrderSchemaV4,
		orderSchemaV4:          compiledRootOrderSchemaV4,
		messageSchemaV4:        compiledRootOrderMessageSchemaV4,
		exchangeAddressV3:      contractAddresses.ExchangeV3,
		exchangeAddressV4:      contractAddresses.ExchangeV4,
	}, nil
}

func loadExchangeV3Address(loader *jsonschema.SchemaLoader, contractAddresses ethereum.ContractAddresses) error {
	// Note that exchangeAddressSchema accepts both checksummed and
	// non-checksummed (i.e. all lowercase) addresses.
	exchangeAddressSchema := fmt.Sprintf(`{"enum":[%q,%q]}`, contractAddresses.ExchangeV3.Hex(), strings.ToLower(contractAddresses.ExchangeV3.Hex()))
	// fmt.Println(exchangeAddressSchema)
	return loader.AddSchema("/exchangeAddress", jsonschema.NewStringLoader(exchangeAddressSchema))
}

func loadExchangeV4Address(loader *jsonschema.SchemaLoader, contractAddresses ethereum.ContractAddresses) error {
	checksummed_address := contractAddresses.ExchangeV4.Hex()
	lower_case_address := strings.ToLower(checksummed_address)
	enum_variants := []string{checksummed_address}
	if checksummed_address != lower_case_address {
		enum_variants = append(enum_variants, lower_case_address)
	}
	return loader.AddSchema("/exchange", jsonschema.NewStringLoader(fmt.Sprintf(`{"enum":%q}`, enum_variants)))
}

func loadChainID(loader *jsonschema.SchemaLoader, chainID int) error {
	chainIDSchema := fmt.Sprintf(`{"const":%d}`, chainID)
	return loader.AddSchema("/chainId", jsonschema.NewStringLoader(chainIDSchema))
}

func newV3Loader(chainID int, customOrderSchemaV3 string, contractAddresses ethereum.ContractAddresses) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadChainID(loader, chainID); err != nil {
		return nil, err
	}
	if err := loadExchangeV3Address(loader, contractAddresses); err != nil {
		return nil, err
	}
	if err := loader.AddSchemas(builtInSchemasV3...); err != nil {
		return nil, err
	}
	//TODO(mason) some error here!>!>!
	if err := loader.AddSchema("/customOrderV3", jsonschema.NewStringLoader(customOrderSchemaV3)); err != nil {
		return nil, err
	}
	return loader, nil
}

func newV4Loader(chainID int, customOrderSchemaV4 string, contractAddresses ethereum.ContractAddresses) (*jsonschema.SchemaLoader, error) {
	loader := jsonschema.NewSchemaLoader()
	if err := loadChainID(loader, chainID); err != nil {
		return nil, err
	}
	if err := loadExchangeV4Address(loader, contractAddresses); err != nil {
		return nil, err
	}
	if err := loader.AddSchemas(builtInSchemasV4...); err != nil {
		return nil, err
	}
	if err := loader.AddSchema("/customOrderV4", jsonschema.NewStringLoader(customOrderSchemaV4)); err != nil {
		return nil, err
	}
	return loader, nil
}
