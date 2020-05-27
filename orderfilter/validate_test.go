// +build !js

package orderfilter

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
