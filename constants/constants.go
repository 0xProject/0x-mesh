package constants

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
)

/**
 * General
 */

// TestNetworkID is the test (Ganache) networkId used for testing
const TestNetworkID = 50

// GanacheEndpoint specifies the Ganache test Ethereum node JSON RPC endpoint used in tests
const GanacheEndpoint = "http://localhost:8545"

// NullAddress is an Ethereum address with all zeroes.
var NullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

var (
	// GanacheAccount0 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount0           = common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631")
	ganacheAccount0PrivateKey = common.Hex2Bytes("f2f48ee19680706196e2e339e5da3491186e0c4c5030670656b0e0164837257d")
	// GanacheAccount1 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount1           = common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb")
	ganacheAccount1PrivateKey = common.Hex2Bytes("5d862464fe9303452126c8bc94274b8c5f9874cbd219789b3eb2128075a76f72")
	// GanacheAccount2 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount2           = common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84")
	ganacheAccount2PrivateKey = common.Hex2Bytes("df02719c4df8b9b8ac7f551fcb5d9ef48fa27eef7a66453879f4d8fdc6e78fb1")
	// GanacheAccount3 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount3           = common.HexToAddress("0xe834ec434daba538cd1b9fe1582052b880bd7e63")
	ganacheAccount3PrivateKey = common.Hex2Bytes("ff12e391b79415e941a94de3bf3a9aee577aed0731e297d5cfa0b8a1e02fa1d0")
	// GanacheAccount4 is the first account exposed on the Ganache test Ethereum node
	GanacheAccount4           = common.HexToAddress("0x78dc5d2d739606d31509c31d654056a45185ecb6")
	ganacheAccount4PrivateKey = common.Hex2Bytes("752dd9cf65e68cfaba7d60225cbdbc1f4729dd5e5507def72815ed0d8abc6249")
)

// GanacheAccountToPrivateKey maps Ganache test Ethereum node accounts to their private key
var GanacheAccountToPrivateKey = map[common.Address][]byte{
	GanacheAccount0: ganacheAccount0PrivateKey,
	GanacheAccount1: ganacheAccount1PrivateKey,
	GanacheAccount2: ganacheAccount2PrivateKey,
	GanacheAccount3: ganacheAccount3PrivateKey,
	GanacheAccount4: ganacheAccount4PrivateKey,
}

// ErrInternal is used whenever we don't wish to expose internal errors to a client
var ErrInternal = errors.New("internal error")
