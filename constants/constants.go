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
	ERC20Proxy        common.Address
	ERC721Proxy       common.Address
	OrderValidator    common.Address
	Exchange          common.Address
	EthBalanceChecker common.Address
}

// NetworkIDToContractAddresses maps networkId to a mapping of contract name to Ethereum address
// on that given network
var NetworkIDToContractAddresses = map[int]ContractNameToAddress{
	50: ContractNameToAddress{
		ERC20Proxy:        common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:       common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		OrderValidator:    common.HexToAddress("0x32eecaf51dfea9618e9bc94e9fbfddb1bbdcba15"),
		Exchange:          common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		EthBalanceChecker: common.HexToAddress("0xaa86dda78e9434aca114b6676fc742a18d15a1cc"),
	},
	1: ContractNameToAddress{
		ERC20Proxy:        common.HexToAddress("0x2240dab907db71e64d3e0dba4800c83b5c502d4e"),
		ERC721Proxy:       common.HexToAddress("0x208e41fb445f1bb1b6780d58356e81405f3e6127"),
		OrderValidator:    common.HexToAddress("0x9463e518dea6810309563c81d5266c1b1d149138"),
		Exchange:          common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b"),
		EthBalanceChecker: common.HexToAddress("0x9bc2c6ae8b1a8e3c375b6ccb55eb4273b2c3fbde"),
	},
}

// NullAddress is an Ethereum address with all zeroes.
var NullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// GanacheAccount0 is the first account exposed on the Ganache test Ethereum node
var GanacheAccount0 = common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
var ganacheAccount0PrivateKey = common.Hex2Bytes("f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d")

// GanacheAccount1 is the first account exposed on the Ganache test Ethereum node
var GanacheAccount1 = common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84")
var ganacheAccount1PrivateKey = common.Hex2Bytes("df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1")

// GanacheAccount2 is the first account exposed on the Ganache test Ethereum node
var GanacheAccount2 = common.HexToAddress("0xe834ec434daba538cd1b9fe1582052b880bd7e63")
var ganacheAccount2PrivateKey = common.Hex2Bytes("ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0")

// GanacheAccount3 is the first account exposed on the Ganache test Ethereum node
var GanacheAccount3 = common.HexToAddress("0x78dc5d2d739606d31509c31d654056a45185ecb6")
var ganacheAccount3PrivateKey = common.Hex2Bytes("752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249")

// GanacheAccountToPrivateKey maps Ganache test Ethereum node accounts to their private key
var GanacheAccountToPrivateKey = map[common.Address][]byte{
	GanacheAccount0: ganacheAccount0PrivateKey,
	GanacheAccount1: ganacheAccount1PrivateKey,
	GanacheAccount2: ganacheAccount2PrivateKey,
	GanacheAccount3: ganacheAccount3PrivateKey,
}
