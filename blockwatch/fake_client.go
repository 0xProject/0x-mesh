package blockwatch

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// FixtureTimestep holds the JSON-RPC data available at every timestep of the simulation.
type FixtureTimestep struct {
	GetLatestBlock   MiniBlockHeader                 `json:"getLatestBlock"  gencodec:"required"`
	GetBlockByNumber map[uint64]MiniBlockHeader      `json:"getBlockByNumber"  gencodec:"required"`
	GetBlockByHash   map[common.Hash]MiniBlockHeader `json:"getBlockByHash"  gencodec:"required"`
	GetCorrectChain  []*MiniBlockHeader              `json:"getCorrectChain" gencodec:"required"`
	Events           []*Event                        `json:"Events" gencodec:"required"`
	ScenarioLabel    string                          `json:"scenarioLabel" gencodec:"required"`
}

// fakeClient is a fake Client for testing purposes.
type fakeClient struct {
	currentTimestep uint
	fixtureData     []FixtureTimestep
}

// newFakeClient instantiates a fakeClient for testing purposes.
func newFakeClient() *fakeClient {
	blob, err := ioutil.ReadFile("testdata/fake_client_fixtures.json")
	if err != nil {
		panic("Failed to read blockwatch fixture file")
	}

	var fixtureData []FixtureTimestep
	_ = json.Unmarshal(blob, &fixtureData)

	var startTimestep uint = 0
	return &fakeClient{startTimestep, fixtureData}
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (fc *fakeClient) HeaderByNumber(ctx context.Context, number *big.Int) (*MiniBlockHeader, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	var miniBlockHeader MiniBlockHeader
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
func (fc *fakeClient) HeaderByHash(ctx context.Context, hash common.Hash) (*MiniBlockHeader, error) {
	timestep := fc.fixtureData[fc.currentTimestep]
	miniBlockHeader, ok := timestep.GetBlockByHash[hash]
	if !ok {
		return nil, ethereum.NotFound
	}
	return &miniBlockHeader, nil
}

// IncrementTimestep increments the timestep of the simulation.
func (fc *fakeClient) IncrementTimestep() {
	fc.currentTimestep++
}

// NumberOfTimesteps returns the number of timesteps in the simulation
func (fc *fakeClient) NumberOfTimesteps() int {
	return len(fc.fixtureData)
}

// ExpectedRetainedBlocks returns the expected retained blocks at the current timestep.
func (fc *fakeClient) ExpectedRetainedBlocks() []*MiniBlockHeader {
	return fc.fixtureData[fc.currentTimestep].GetCorrectChain
}

// GetScenarioLabel returns a label describing the test case being tested by the current timestep
// of the simulation.
func (fc *fakeClient) GetScenarioLabel() string {
	return fc.fixtureData[fc.currentTimestep].ScenarioLabel
}

// GetEvents returns the events in the order they should have been emitted by Watcher for
// the current timestep of the simulation.
func (fc *fakeClient) GetEvents() []*Event {
	return fc.fixtureData[fc.currentTimestep].Events
}
