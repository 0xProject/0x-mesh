// +build js, wasm

package orderfilter

import (
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/ethereum"
)

var (
	testFilterValidateOrderExpectedErrors = [][]string{
		{},
		{
			`{"keyword":"enum","dataPath":".exchangeAddress","schemaPath":"/exchangeAddress/enum","params":{"allowedValues":["0x48bacb9266a570d521063ef5dd96e61686dbe788"]},"message":"should be equal to one of the allowed values"}`,
		},
		{
			`{"keyword":"const","dataPath":".chainId","schemaPath":"/chainId/const","params":{"allowedValue":1337},"message":"should be equal to constant"}`,
		},
		{},
		{
			`{"keyword":"const","dataPath":".senderAddress","schemaPath":"/customOrder/properties/senderAddress/const","params":{"allowedValue":"0x00000000000000000000000000000000ba5eba11"},"message":"should be equal to constant"}`,
		},
	}
	testFilterValidateOrderJSONExpectedErrors = [][]string{
		{},
		{
			`{"keyword":"required","dataPath":"","schemaPath":"#/required","params":{"missingProperty":"makerAddress"},"message":"should have required property 'makerAddress'"}`,
		},
		{

			`{"keyword":"required","dataPath":"","schemaPath":"#/required","params":{"missingProperty":"makerAddress"},"message":"should have required property 'makerAddress'"}`,
		},
		{
			`{"keyword":"pattern","dataPath":".takerAddress","schemaPath":"/address/pattern","params":{"pattern":"^0x[0-9a-fA-F]{40}$"},"message":"should match pattern \"^0x[0-9a-fA-F]{40}$\""}`,
		},
		{
			`{"keyword":"enum","dataPath":".exchangeAddress","schemaPath":"/exchangeAddress/enum","params":{"allowedValues":["0x48bacb9266a570d521063ef5dd96e61686dbe788"]},"message":"should be equal to one of the allowed values"}`,
		},
		{
			`{"keyword":"const","dataPath":".chainId","schemaPath":"/chainId/const","params":{"allowedValue":1337},"message":"should be equal to constant"}`,
		},
	}
)

func setupTestCase(chainID int, contractAddresses ethereum.ContractAddresses, customOrderSchema string) {
	normalizedExchangeAddress := strings.ToLower(contractAddresses.Exchange.String())
	js.Global().Call("setSchemaValidator", js.ValueOf(chainID), js.ValueOf(normalizedExchangeAddress), js.ValueOf(customOrderSchema))
}
