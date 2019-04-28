package meshdb

import (
	"math/big"
	"sort"

	"github.com/0xProject/0x-mesh/db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MiniHeader is the database representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash   common.Hash `json:"hash"   gencodec:"required"`
	Parent common.Hash `json:"parent" gencodec:"required"`
	Number *big.Int    `json:"number" gencodec:"required"`
	Logs   []types.Log `json:"logs" gencodec:"required"`
}

// NewMiniHeader returns a new MiniHeader.
func NewMiniHeader(hash common.Hash, parent common.Hash, number *big.Int) *MiniHeader {
	miniHeader := MiniHeader{Hash: hash, Parent: parent, Number: number, Logs: []types.Log{}}
	return &miniHeader
}

// ID returns the MiniHeader's ID
func (m *MiniHeader) ID() []byte {
	return m.Hash.Bytes()
}

// MeshDB instantiates the DB connection and creates all the collections used by the application
type MeshDB struct {
	database    *db.DB
	MiniHeaders *db.Collection
}

// NewMeshDB instantiates a new MeshDB instance
func NewMeshDB(path string) (*MeshDB, error) {
	database, err := db.Open(path)
	if err != nil {
		return nil, err
	}
	miniHeaders := database.NewCollection("miniHeader", &MiniHeader{})

	return &MeshDB{
		database:    database,
		MiniHeaders: miniHeaders,
	}, nil
}

// Close closes the database connection
func (m *MeshDB) Close() {
	m.database.Close()
}

// FindAllMiniHeadersSortedByNumber returns all MiniHeaders sorted by block number
func (m *MeshDB) FindAllMiniHeadersSortedByNumber() []*MiniHeader {
	miniHeaders := []*MiniHeader{}
	m.MiniHeaders.FindAll(&miniHeaders)
	sort.Slice(miniHeaders, func(i, j int) bool {
		return miniHeaders[i].Number.Cmp(miniHeaders[j].Number) == -1
	})
	return miniHeaders
}
