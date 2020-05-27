// +build js, wasm

package orderfilter

import (
	"github.com/0xProject/0x-mesh/ethereum"
)

type Filter struct {
	validatorLoaded      bool
	encodedSchema        string
	chainID              int
	rawCustomOrderSchema string
}

func New(chainID int, customOrderSchema string, contractAddresses ethereum.ContractAddresses) (*Filter, error) {
	return &Filter{
		chainID:              chainID,
		rawCustomOrderSchema: customOrderSchema,
	}, nil
}
