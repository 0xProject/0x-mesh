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
