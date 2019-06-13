package constants

import (
	"github.com/ethereum/go-ethereum/common"
)

/**
 * General
 */

// TestNetworkID is the test (Ganache) networkId used for testing
const TestNetworkID = 50

// GanacheEndpoint specifies the Ganache test Ethereum node JSON RPC endpoint used in tests
const GanacheEndpoint = "http://localhost:8545"

// ContractNameToAddress maps a contract's name to it's Ethereum address
type ContractNameToAddress struct {
	ERC20Proxy           common.Address
	ERC721Proxy          common.Address
	OrderValidationUtils common.Address
	Exchange             common.Address
	EthBalanceChecker    common.Address
}

// NetworkIDToContractAddresses maps networkId to a mapping of contract name to Ethereum address
// on that given network
var NetworkIDToContractAddresses = map[int]ContractNameToAddress{
	// Mainnet
	1: ContractNameToAddress{
		ERC20Proxy:           common.HexToAddress("0x2240dab907db71e64d3e0dba4800c83b5c502d4e"),
		ERC721Proxy:          common.HexToAddress("0x208e41fb445f1bb1b6780d58356e81405f3e6127"),
		Exchange:             common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b"),
		OrderValidationUtils: common.HexToAddress("0xaccc650ae2669b9a48c525c40f28c497ff88d839"),
		EthBalanceChecker:    common.HexToAddress("0x9bc2c6ae8b1a8e3c375b6ccb55eb4273b2c3fbde"),
	},
	// Ropsten
	3: ContractNameToAddress{
		ERC20Proxy:           common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
		ERC721Proxy:          common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
		Exchange:             common.HexToAddress("0x4530c0483a1633c7a1c97d2c53721caff2caaaaf"),
		OrderValidationUtils: common.HexToAddress("0x078b18930c5c16f73b3d9e8e4911d1720c98e292"),
		EthBalanceChecker:    common.HexToAddress("0xd5d960219af544b6f2f3e14a8bfd03dec12292fa"),
	},
	// Rinkeby
	4: ContractNameToAddress{
		ERC20Proxy:           common.HexToAddress("0x3e809c563c15a295e832e37053798ddc8d6c8dab"),
		ERC721Proxy:          common.HexToAddress("0x8e1ff02637cb5e39f2fa36c14706aa348b065b09"),
		Exchange:             common.HexToAddress("0x22ebc052f43a88efa06379426120718170f2204e"),
		OrderValidationUtils: common.HexToAddress("0x2d645bce387b25bdd807be48e57f374226eadff9"),
		EthBalanceChecker:    common.HexToAddress("0x08b71282431009022eda2dda8af0fbee535e1507"),
	},
	// Kovan
	42: ContractNameToAddress{
		ERC20Proxy:           common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
		ERC721Proxy:          common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
		Exchange:             common.HexToAddress("0x35dd2932454449b14cee11a94d3674a936d5d7b2"),
		OrderValidationUtils: common.HexToAddress("0x98cafab5518485d84fe3c2c495f0fcc807f1a7c0"),
		EthBalanceChecker:    common.HexToAddress("0x505aa534485bf80ee919339717cff90eb2e3364c"),
	},
	// Ganache snapshot
	50: ContractNameToAddress{
		ERC20Proxy:           common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:          common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		Exchange:             common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		OrderValidationUtils: common.HexToAddress("0x8d42e38980ce74736c21c059b2240df09958d3c8"),
		EthBalanceChecker:    common.HexToAddress("0xa31e64ea55b9b6bbb9d6a676738e9a5b23149f84"),
	},
}

// NullAddress is an Ethereum address with all zeroes.
var NullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

var (
	// GanacheAccount0 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount0           = common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	ganacheAccount0PrivateKey = common.Hex2Bytes("f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d")
	// GanacheAccount1 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount1           = common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84")
	ganacheAccount1PrivateKey = common.Hex2Bytes("df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1")
	// GanacheAccount2 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount2           = common.HexToAddress("0xe834ec434daba538cd1b9fe1582052b880bd7e63")
	ganacheAccount2PrivateKey = common.Hex2Bytes("ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0")
	// GanacheAccount3 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount3           = common.HexToAddress("0x78dc5d2d739606d31509c31d654056a45185ecb6")
	ganacheAccount3PrivateKey = common.Hex2Bytes("752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249")
)

// GanacheAccountToPrivateKey maps Ganache test Ethereum node accounts to their private key
var GanacheAccountToPrivateKey = map[common.Address][]byte{
	GanacheAccount0: ganacheAccount0PrivateKey,
	GanacheAccount1: ganacheAccount1PrivateKey,
	GanacheAccount2: ganacheAccount2PrivateKey,
	GanacheAccount3: ganacheAccount3PrivateKey,
}
