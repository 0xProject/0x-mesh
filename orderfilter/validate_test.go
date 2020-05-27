// +build !js

package orderfilter

import (
	"github.com/0xProject/0x-mesh/ethereum"
)

var (
	testFilterValidateOrderExpectedErrors = [][]string{
		{},
		{
			"exchangeAddress must be one of the following",
		},
		{
			"chainId does not match",
		},
		{},
		{
			"senderAddress does not match",
		},
	}
	testFilterValidateOrderJSONExpectedErrors = [][]string{
		{},
		{
			"makerAddress is required",
		},
		{
			"makerAddress is required",
		},
		{
			"takerAddress: Does not match pattern '^0x[0-9a-fA-F]{40}$",
		},
		{
			"exchangeAddress must be one of the following",
		},
		{
			"chainId does not match",
		},
	}
)

func setupTestCase(chainID int, contractAddresses ethereum.ContractAddresses, customOrderSchema string) {
	// NOTE(jalextowle): Setup is only required in the WebAssembly tests
}
