// +build !js

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

var (
	validSignedOrderJSON           = `{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x4f833a24e1f95d70f028921e27040ca56e09ab0b","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`
	misspelledFieldSignedOrderJSON = `{"makerAdress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x4f833a24e1f95d70f028921e27040ca56e09ab0b","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`
	wrongTypeSignedOrderJSON       = `{"makerAddress":"0xa3ece5d5b6319fa785efc10d3112769a46c6e149","takerAddress":"hi","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x4f833a24e1f95d70f028921e27040ca56e09ab0b","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`
	missingFieldSignedOrderJSON    = `{"takerAddress":"0x0000000000000000000000000000000000000000","makerAssetAmount":"100000000000000000000","takerAssetAmount":"100000000000000000000000","expirationTimeSeconds":"1559856615025","makerFee":"0","takerFee":"0","feeRecipientAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","salt":"46108882540880341679561755865076495033942060608820537332859096815711589201849","makerAssetData":"0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498","takerAssetData":"0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","exchangeAddress":"0x4f833a24e1f95d70f028921e27040ca56e09ab0b","signature":"0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603"}`
)

func TestValidateOrderSchema(t *testing.T) {
	schema, err := setupSchemaValidator()
	require.NoError(t, err)

	orderLoader := gojsonschema.NewStringLoader(validSignedOrderJSON)
	result, err := schema.Validate(orderLoader)
	require.NoError(t, err)
	assert.Equal(t, true, result.Valid())

	orderLoader = gojsonschema.NewStringLoader(misspelledFieldSignedOrderJSON)
	result, err = schema.Validate(orderLoader)
	require.NoError(t, err)
	assert.Equal(t, false, result.Valid())
	errs := result.Errors()
	require.Len(t, errs, 2)
	assert.Equal(t, "makerAddress is required", errs[0].Description())

	orderLoader = gojsonschema.NewStringLoader(missingFieldSignedOrderJSON)
	result, err = schema.Validate(orderLoader)
	require.NoError(t, err)
	assert.Equal(t, false, result.Valid())
	errs = result.Errors()
	require.Len(t, errs, 2)
	assert.Equal(t, "makerAddress is required", errs[0].Description())

	orderLoader = gojsonschema.NewStringLoader(wrongTypeSignedOrderJSON)
	result, err = schema.Validate(orderLoader)
	require.NoError(t, err)
	assert.Equal(t, false, result.Valid())
	errs = result.Errors()
	require.Len(t, errs, 2)
	assert.Equal(t, "takerAddress", errs[0].Field())
	assert.Equal(t, "Does not match pattern '^0x[0-9a-fA-F]{40}$'", errs[0].Description())
}
