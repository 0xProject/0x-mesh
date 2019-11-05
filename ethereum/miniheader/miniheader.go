package miniheader

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MiniHeader is a representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash      common.Hash
	Parent    common.Hash
	Number    *big.Int
	Timestamp time.Time
	Logs      []types.Log
}

// ID returns the MiniHeader's ID
// HACK(fabio): This method is only used by DBStack and not SimpleStack
// Ideally this would live in the `meshdb` package but it adds the need
// to cast back-and-forth between two almost identical types so we keep
// it here for convenience sake.
func (m *MiniHeader) ID() []byte {
	return m.Hash.Bytes()
}
