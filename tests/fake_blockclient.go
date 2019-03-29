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
	GetLatestBlock   blockwatch.SuccinctBlock                 `json:"getLatestBlock"  gencodec:"required"`
	GetBlockByNumber map[uint64]blockwatch.SuccinctBlock      `json:"getBlockByNumber"  gencodec:"required"`
	GetBlockByHash   map[common.Hash]blockwatch.SuccinctBlock `json:"getBlockByHash"  gencodec:"required"`
	GetCorrectChain  []*blockwatch.SuccinctBlock              `json:"getCorrectChain" gencodec:"required"`
	BlockEvents      []*blockwatch.BlockEvent                 `json:"blockEvents" gencodec:"required"`
	ScenarioLabel    string                                   `json:"scenarioLabel" gencodec:"required"`
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

// BlockByNumber fetches a block by its number. If no `number` is supplied, it will return the latest block.
// If not block exists with this number it will return a `ethereum.NotFound` error.
func (fc *FakeBlockClient) BlockByNumber(ctx context.Context, number *big.Int) (*blockwatch.SuccinctBlock, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	var succinctBlock blockwatch.SuccinctBlock
	var ok bool
	if number == nil {
		succinctBlock = timestep.GetLatestBlock
	} else {
		succinctBlock, ok = timestep.GetBlockByNumber[number.Uint64()]
		if !ok {
			return nil, ethereum.NotFound
		}
	}
	return &succinctBlock, nil
}

// BlockByHash fetches a block by its block hash. If not block exists with this number it will return a `ethereum.NotFound` error.
func (fc *FakeBlockClient) BlockByHash(ctx context.Context, hash common.Hash) (*blockwatch.SuccinctBlock, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	succinctBlock, ok := timestep.GetBlockByHash[hash]
	if !ok {
		return nil, ethereum.NotFound
	}
	return &succinctBlock, nil
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
func (fc *FakeBlockClient) ExpectedRetainedBlocks() []*blockwatch.SuccinctBlock {
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
