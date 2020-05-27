// +build !js

package orderfilter

var (
	wrongExchangeAddressError    = "exchangeAddress must be one of the following"
	mismatchedChainIDError       = "chainId does not match"
	mismatchedSenderAddressError = "senderAddress does not match"
	requiredMakerAddressError    = "makerAddress is required"
	invalidTakerAddressError     = "takerAddress: Does not match pattern '^0x[0-9a-fA-F]{40}$"
)
