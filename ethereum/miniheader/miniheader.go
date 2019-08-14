package miniheader

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MiniHeader is a representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash   common.Hash
	Parent common.Hash
	Number *big.Int
	Logs   []types.Log
}

// ID returns the MiniHeader's ID
func (m *MiniHeader) ID() []byte {
	return m.Hash.Bytes()
}
