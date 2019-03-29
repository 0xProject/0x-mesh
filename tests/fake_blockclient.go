package tests

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// FixtureTimestep holds the JSON-RPC data available at every timestep of the simulation.
type FixtureTimestep struct {
	GetLatestBlock   blockwatch.MiniBlockHeader                 `json:"getLatestBlock"  gencodec:"required"`
	GetBlockByNumber map[uint64]blockwatch.MiniBlockHeader      `json:"getBlockByNumber"  gencodec:"required"`
	GetBlockByHash   map[common.Hash]blockwatch.MiniBlockHeader `json:"getBlockByHash"  gencodec:"required"`
	GetCorrectChain  []*blockwatch.MiniBlockHeader              `json:"getCorrectChain" gencodec:"required"`
	BlockEvents      []*blockwatch.BlockEvent                   `json:"blockEvents" gencodec:"required"`
	ScenarioLabel    string                                     `json:"scenarioLabel" gencodec:"required"`
}

// FakeBlockClient is a fake BlockClient for testing purposes.
type FakeBlockClient struct {
	currentTimestep uint
	fixtureData     []FixtureTimestep
}

// NewFakeBlockClient instantiates a FakeBlockClient for testing purposes.
func NewFakeBlockClient() *FakeBlockClient {
	blob, err := ioutil.ReadFile("testdata/blockwatch_generated_test.json")
	if err != nil {
		panic("Failed to read blockwatch fixture file")
	}

	var fixtureData []FixtureTimestep
	_ = json.Unmarshal(blob, &fixtureData)

	var startTimestep uint = 0
	return &FakeBlockClient{startTimestep, fixtureData}
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (fc *FakeBlockClient) HeaderByNumber(ctx context.Context, number *big.Int) (*blockwatch.MiniBlockHeader, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	var miniBlockHeader blockwatch.MiniBlockHeader
	var ok bool
	if number == nil {
		miniBlockHeader = timestep.GetLatestBlock
	} else {
		miniBlockHeader, ok = timestep.GetBlockByNumber[number.Uint64()]
		if !ok {
			return nil, ethereum.NotFound
		}
	}
	return &miniBlockHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (fc *FakeBlockClient) HeaderByHash(ctx context.Context, hash common.Hash) (*blockwatch.MiniBlockHeader, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	miniBlockHeader, ok := timestep.GetBlockByHash[hash]
	if !ok {
		return nil, ethereum.NotFound
	}
	return &miniBlockHeader, nil
}

// IncrementTimestep increments the timestep of the simulation.
func (fc *FakeBlockClient) IncrementTimestep() {
	fc.currentTimestep++
}

// NumberOfTimesteps returns the number of timesteps in the simulation
func (fc *FakeBlockClient) NumberOfTimesteps() int {
	return len(fc.fixtureData)
}

// ExpectedRetainedBlocks returns the expected retained blocks at the current timestep.
func (fc *FakeBlockClient) ExpectedRetainedBlocks() []*blockwatch.MiniBlockHeader {
	return fc.fixtureData[fc.currentTimestep].GetCorrectChain
}

// GetScenarioLabel returns a label describing the test case being tested by the current timestep
// of the simulation.
func (fc *FakeBlockClient) GetScenarioLabel() string {
	return fc.fixtureData[fc.currentTimestep].ScenarioLabel
}

// GetBlockEvents returns the events in the order they should have been emitted by BlockWatch for
// the current timestep of the simulation.
func (fc *FakeBlockClient) GetBlockEvents() []*blockwatch.BlockEvent {
	return fc.fixtureData[fc.currentTimestep].BlockEvents
}
