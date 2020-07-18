// Contains non-generated custom model code.

package gqltypes

import (
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// Hash is a 32-byte Keccak256 hash encoded as a hexadecimal string.
type Hash common.Hash

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (h *Hash) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Hash must be a hex-encoded string")
	}
	hash := common.HexToHash(s)
	if (hash == common.Hash{}) {
		return fmt.Errorf("invalid Hash value: %q", s)
	}
	(*h) = Hash(hash)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (h Hash) MarshalGQL(w io.Writer) {
	quotedHexString := strconv.Quote(common.Hash(h).Hex())
	_, _ = w.Write([]byte(quotedHexString))
}

// Address is an Ethereum address encoded as a hexadecimal string.
type Address common.Address

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (a *Address) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Address must be a hex-encoded string")
	}
	// TODO(albrow): Check if valid hex.
	address := common.HexToAddress(s)
	(*a) = Address(address)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (a Address) MarshalGQL(w io.Writer) {
	hexString := strings.ToLower(common.Address(a).Hex())
	quotedHexString := strconv.Quote(hexString)
	_, _ = w.Write([]byte(quotedHexString))
}

// BigNumber is a uint256 value encoded as a numerical string.
type BigNumber big.Int

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (b *BigNumber) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("BigNumber must be a numerical string")
	}
	bigInt, ok := math.ParseBig256(s)
	if !ok {
		return fmt.Errorf("invalid BigNumber value: %q", s)
	}
	(*b) = BigNumber(*bigInt)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (b BigNumber) MarshalGQL(w io.Writer) {
	bigInt := big.Int(b)
	_, _ = w.Write([]byte((&bigInt).String()))
}

// Bytes is an array of arbitrary bytes encoded as a hexadecimal string.
type Bytes []byte

func (b *Bytes) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Bytes must be a hex-encoded string")
	}
	// TODO(albrow): Check if valid hex.
	bytes := common.FromHex(s)
	(*b) = Bytes(bytes)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (b Bytes) MarshalGQL(w io.Writer) {
	quotedHexString := strconv.Quote(common.ToHex([]byte(b)))
	_, _ = w.Write([]byte(quotedHexString))
}

func lowerCaseAndQuote(s string) string {
	return strconv.Quote(strings.ToLower(s))
}
