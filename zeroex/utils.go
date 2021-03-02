package zeroex

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}

// Bytes32 represents the Solidity `bytes32` type.
// It is mostly copied from go-ethereum's Hash type.
// See <https://github.com/ethereum/go-ethereum/blob/053ed9cc847647a9b3ef707d0efe7104c4ab2a4c/common/types.go#L47>
type Bytes32 [32]byte

// BytesToBytes32 converts []byte to Bytes32
func BytesToBytes32(b []byte) Bytes32 {
	var result Bytes32
	result.SetBytes(b)
	return result
}

// HexToBytes32 creates a Bytes32 from a hex string.
func HexToBytes32(s string) Bytes32 {
	return BytesToBytes32(common.FromHex(s))
}

// BigToBytes32 creates a Bytes32 from a big.Int.
func BigToBytes32(i *big.Int) Bytes32 {
	return BytesToBytes32(i.Bytes())
}

// HashToBytes32 creates a Bytes32 from a Hash.
func HashToBytes32(h common.Hash) Bytes32 {
	return BytesToBytes32(h.Bytes())
}

// SetBytes sets the Bytes32 to the value of bytes.
// If bytes is larger than 32, bytes will be cropped from the left.
func (b *Bytes32) SetBytes(bytes []byte) {
	if len(bytes) > 32 {
		bytes = bytes[len(bytes)-32:]
	}

	copy(b[32-len(bytes):], bytes)
}

// Hex converts a Bytes32 to a hex string.
func (b Bytes32) Hex() string {
	return hexutil.Encode(b[:])
}

// Bytes converts a Bytes32 to a []bytes.
func (b Bytes32) Bytes() []byte {
	return b[:]
}

// Raw converts a Bytes32 to a [32]bytes.
func (b Bytes32) Raw() [32]byte {
	return b
}

// Raw converts a Bytes32 to a big.Int
func (b Bytes32) Big() *big.Int {
	return new(big.Int).SetBytes(b[:])
}

// String prints the value in hex
func (b Bytes32) String() string {
	return b.Hex()
}
