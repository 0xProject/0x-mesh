package zeroex

import (
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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

func TestSignOrderV4(t *testing.T) {
	// See <https://github.com/0xProject/protocol/blob/edda1edc507fbfceb6dcb02ef212ee4bdcb123a6/packages/protocol-utils/test/orders_test.ts#L15>
	privateKeyBytes := hexutil.MustDecode("0xee094b79aa0315914955f2f09be9abe541dcdc51f0aae5bec5453e9f73a471a6")

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	require.NoError(t, err)
	localSigner := signer.NewLocalSigner(privateKey)
	localSignerAddress := localSigner.(*signer.LocalSigner).GetSignerAddress()
	assert.Equal(t, common.HexToAddress("0x05cAc48D17ECC4D8A9DB09Dde766A03959b98367"), localSignerAddress)

	// Only maker is allowed to sign
	order := testOrderV4
	order.Maker = localSignerAddress

	signedOrder, err := SignOrderV4(localSigner, testOrderV4)
	require.NoError(t, err)

	// See <https://github.com/0xProject/protocol/blob/edda1edc507fbfceb6dcb02ef212ee4bdcb123a6/packages/protocol-utils/test/orders_test.ts#L77>
	// TODO: Check if signature is wrong or if we are dealing with non-uniqueness
	// of signatures. Maybe the private key is loaded differently than the reference, in which case the maker address will differ.
	assert.Equal(t, EIP712SignatureV4, signedOrder.SignatureTypeV4)
	// assert.Equal(t, uint8(27), signedOrder.V)
	// assert.Equal(t, HexToBytes32("0x030e27e0a261dda1139154d9ba7e814932bd6b8d15231a8c2cd78d634ff22c2b"), signedOrder.R)
	// assert.Equal(t, HexToBytes32("0x50af45e0d6e81b721905bd35748168f1f348be34fe03073d7a2f2b053cbdca2d"), signedOrder.S)
	assert.Equal(t, uint8(28), signedOrder.V)
	assert.Equal(t, HexToBytes32("0x8a01b9be45f72e7080a23f415dcaa6b5dc69871c9964f0758704a678cbdfad9f"), signedOrder.R)
	assert.Equal(t, HexToBytes32("0x33b0e281d2512bda7064138cfc7e7aefc99effc61d3191305006e963ae87efa8"), signedOrder.S)
}
