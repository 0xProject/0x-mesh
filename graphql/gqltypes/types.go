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

// MarshalGQL implements the graphql.Marshaler interface
func (h Hash) MarshalGQL(w io.Writer) {
	quotedHexString := strconv.Quote(common.Hash(h).Hex())
	_, _ = w.Write([]byte(quotedHexString))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (h *Hash) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Hash must be a hex-encoded string")
	}
	hash := common.HexToHash(unquoteIfNeeded(s))
	if (hash == common.Hash{}) {
		return fmt.Errorf("invalid Hash value: %q", s)
	}
	(*h) = Hash(hash)

	return nil
}

func (h Hash) MarshalJSON() ([]byte, error) {
	quotedHexString := strconv.Quote(common.Hash(h).Hex())
	return []byte(quotedHexString), nil
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("Hash should be a hex-encoded string: %s", err)
	}
	hash := common.HexToHash(unquoted)
	if (hash == common.Hash{}) {
		return fmt.Errorf("invalid Hash value: %q", data)
	}
	(*h) = Hash(hash)

	return nil
}

// Address is an Ethereum address encoded as a hexadecimal string.
type Address common.Address

// MarshalGQL implements the graphql.Marshaler interface
func (a Address) MarshalGQL(w io.Writer) {
	hexString := strings.ToLower(common.Address(a).Hex())
	quotedHexString := strconv.Quote(hexString)
	_, _ = w.Write([]byte(quotedHexString))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (a *Address) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Address must be a hex-encoded string")
	}
	// TODO(albrow): Check if valid hex.
	address := common.HexToAddress(unquoteIfNeeded(s))
	(*a) = Address(address)

	return nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	hexString := strings.ToLower(common.Address(a).Hex())
	quotedHexString := strconv.Quote(hexString)
	return []byte(quotedHexString), nil
}

func (a *Address) UnmarshalJSON(data []byte) error {
	// TODO(albrow): Check if valid hex.
	address := common.HexToAddress(unquoteIfNeeded(string(data)))
	(*a) = Address(address)
	return nil
}

// BigNumber is a uint256 value encoded as a numerical string.
type BigNumber big.Int

// MarshalGQL implements the graphql.Marshaler interface
func (b BigNumber) MarshalGQL(w io.Writer) {
	bigInt := big.Int(b)
	quotedString := strconv.Quote(bigInt.String())
	_, _ = w.Write([]byte(quotedString))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (b *BigNumber) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("BigNumber must be a numerical string but got %T %v", v, v)
	}
	bigInt, ok := math.ParseBig256(unquoteIfNeeded(s))
	if !ok {
		return fmt.Errorf("invalid BigNumber value: %s", s)
	}
	(*b) = BigNumber(*bigInt)

	return nil
}

func (b BigNumber) MarshalJSON() ([]byte, error) {
	bigInt := big.Int(b)
	quotedString := strconv.Quote(bigInt.String())
	return []byte(quotedString), nil
}

func (b *BigNumber) UnmarshalJSON(data []byte) error {
	dataString := unquoteIfNeeded(string(data))
	bigInt, ok := math.ParseBig256(dataString)
	if !ok {
		return fmt.Errorf("invalid BigNumber value: %s", data)
	}
	(*b) = BigNumber(*bigInt)
	return nil
}

// Bytes is an array of arbitrary bytes encoded as a hexadecimal string.
type Bytes []byte

// MarshalGQL implements the graphql.Marshaler interface
func (b Bytes) MarshalGQL(w io.Writer) {
	if len(b) == 0 {
		_, _ = w.Write([]byte(`"0x"`))
		return
	}
	quotedHexString := strconv.Quote(common.ToHex([]byte(b)))
	_, _ = w.Write([]byte(quotedHexString))
}

func (b *Bytes) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("Bytes must be a hex-encoded string")
	}
	// TODO(albrow): Check if valid hex.
	bytes := common.FromHex(unquoteIfNeeded(s))
	(*b) = Bytes(bytes)

	return nil
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	if len(b) == 0 {
		return []byte(`"0x"`), nil
	}
	quotedHexString := strconv.Quote(common.ToHex([]byte(b)))
	return []byte(quotedHexString), nil
}

func (b *Bytes) UnmarshalJSON(data []byte) error {
	// TODO(albrow): Check if valid hex.
	bytes := common.FromHex(unquoteIfNeeded(string(data)))
	(*b) = Bytes(bytes)

	return nil
}

func lowerCaseAndQuote(s string) string {
	return strconv.Quote(strings.ToLower(s))
}

// unquoteIfNeeded removes surrounding quotes (if present) from s.
// Note(albrow): This is needed because generated GraphQL code appears
// to sometimes strip out quotes.
func unquoteIfNeeded(s string) string {
	unquotedString, err := strconv.Unquote(s)
	if err == nil {
		return unquotedString
	}
	return s
}
