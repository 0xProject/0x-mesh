package zeroex

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// See <https://github.com/0xProject/protocol/blob/edda1edc507fbfceb6dcb02ef212ee4bdcb123a6/packages/protocol-utils/test/orders_test.ts#L24>
var testOrderV4 = &OrderV4{
	ChainID:         big.NewInt(8008),
	ExchangeAddress: common.HexToAddress("0x6701704d2421c64ee9aa93ec7f96ede81c4be77d"),

	MakerToken:          common.HexToAddress("0x349e8d89e8b37214d9ce3949fc5754152c525bc3"),
	TakerToken:          common.HexToAddress("0x83c62b2e67dea0df2a27be0def7a22bd7102642c"),
	MakerAmount:         big.NewInt(1234),
	TakerAmount:         big.NewInt(5678),
	TakerTokenFeeAmount: big.NewInt(9101112),
	Maker:               common.HexToAddress("0x8d5e5b5b5d187bdce2e0143eb6b3cc44eef3c0cb"),
	Taker:               common.HexToAddress("0x615312fb74c31303eab07dea520019bb23f4c6c2"),
	Sender:              common.HexToAddress("0x70f2d6c7acd257a6700d745b76c602ceefeb8e20"),
	FeeRecipient:        common.HexToAddress("0xcc3c7ea403427154ec908203ba6c418bd699f7ce"),
	Pool:                HexToBytes32("0x0bbff69b85a87da39511aefc3211cb9aff00e1a1779dc35b8f3635d8b5ea2680"),
	Expiry:              big.NewInt(1001),
	Salt:                big.NewInt(2001),
}

func TestGenerateOrderHashV4(t *testing.T) {
	// See <https://github.com/0xProject/protocol/blob/edda1edc507fbfceb6dcb02ef212ee4bdcb123a6/packages/protocol-utils/test/orders_test.ts#L49>
	expectedOrderHash := common.HexToHash("0x8bb1f6e880b3b4f91a901897c4b914ec606dc3b8b59f64983e1638a45bdf3116")
	actualOrderHash, err := testOrderV4.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}
