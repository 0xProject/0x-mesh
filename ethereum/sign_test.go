package ethereum

import (
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestECSign(t *testing.T) {
	// Test parameters lifted from @0x/order-utils' `signature_utils_test.ts`
	signerAddress := common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	message := common.Hex2Bytes("6927e990021d23b1eb7b8789f6a6feaf98fe104bb0cf8259421b79f9a34222b0")
	expectedSignature := &ECSignature{
		V: byte(27),
		R: common.Hex2Bytes("61a3ed31b43c8780e905a260a35faefcc527be7516aa11c0256729b5b351bc33"),
		S: common.Hex2Bytes("40349190569279751135161d22529dc25add4f6069af05be04cacbda2ace2254"),
	}

	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	actualSignature, err := ECSign(message, signerAddress, rpcClient)
	require.NoError(t, err)

	assert.Equal(t, expectedSignature, actualSignature)
}
