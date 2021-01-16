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

func TestGanacheOrderHashV4(t *testing.T) {
	// See <https://0xproject.slack.com/archives/CAU8U19LJ/p1610762234026600?thread_ts=1610761930.023900&cid=CAU8U19LJ>
	expectedOrderHash := common.HexToHash("0xef61679248399669a4dd10de335d0c151a5c42568618abace01f7a8ec1e693e1")

	order := testOrderV4
	order.ChainID = big.NewInt(1337)
	order.ExchangeAddress = common.HexToAddress("0x5315e44798395d4a952530d131249fE00f554565")
	order.ResetHash()

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
	order.ResetHash()
	actualOrderHash, err := order.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, common.HexToHash("0xddee5b1b08f4df161cafa12fd347e374779573773161906999a19a1cf5f692cc"), actualOrderHash)

	// Sign order
	signedOrder, err := SignOrderV4(localSigner, order)
	require.NoError(t, err)
	// See <https://github.com/0xProject/protocol/blob/edda1edc507fbfceb6dcb02ef212ee4bdcb123a6/packages/protocol-utils/test/orders_test.ts#L67>
	assert.Equal(t, EthSignSignatureV4, signedOrder.SignatureTypeV4)
	assert.Equal(t, uint8(28), signedOrder.V)
	assert.Equal(t, HexToBytes32("0x5d4fe9b4c8f94efc46ef9e7e3f996c238f9c930fd5c03014ec6db6d4d18a34e5"), signedOrder.R)
	assert.Equal(t, HexToBytes32("0x0949269d29524aec1ba5b19236c392a3d1866ca39bb8c7b6345e90a3fbf404fc"), signedOrder.S)
}
