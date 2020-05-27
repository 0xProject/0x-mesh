// +build js, wasm

package orderfilter

var (
	wrongExchangeAddressError    = `{"keyword":"enum","dataPath":".exchangeAddress","schemaPath":"/exchangeAddress/enum","params":{"allowedValues":["0x48BaCB9266a570d521063EF5dD96e61686DbE788","0x48bacb9266a570d521063ef5dd96e61686dbe788"]},"message":"should be equal to one of the allowed values"}`
	mismatchedChainIDError       = `{"keyword":"const","dataPath":".chainId","schemaPath":"/chainId/const","params":{"allowedValue":1337},"message":"should be equal to constant"}`
	mismatchedSenderAddressError = `{"keyword":"const","dataPath":".senderAddress","schemaPath":"/customOrder/properties/senderAddress/const","params":{"allowedValue":"0x00000000000000000000000000000000ba5eba11"},"message":"should be equal to constant"}`
	requiredMakerAddressError    = `{"keyword":"required","dataPath":"","schemaPath":"#/required","params":{"missingProperty":"makerAddress"},"message":"should have required property 'makerAddress'"}`
	invalidTakerAddressError     = `{"keyword":"pattern","dataPath":".takerAddress","schemaPath":"/address/pattern","params":{"pattern":"^0x[0-9a-fA-F]{40}$"},"message":"should match pattern \"^0x[0-9a-fA-F]{40}$\""}`
)
