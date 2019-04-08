package blockwatch

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// fixtureTimestep holds the JSON-RPC data available at every timestep of the simulation.
type fixtureTimestep struct {
	GetLatestBlock   MiniHeader                 `json:"getLatestBlock"  gencodec:"required"`
	GetBlockByNumber map[uint64]MiniHeader      `json:"getBlockByNumber"  gencodec:"required"`
	GetBlockByHash   map[common.Hash]MiniHeader `json:"getBlockByHash"  gencodec:"required"`
	GetCorrectChain  []*MiniHeader              `json:"getCorrectChain" gencodec:"required"`
	BlockEvents      []*Event                   `json:"blockEvents" gencodec:"required"`
	ScenarioLabel    string                     `json:"scenarioLabel" gencodec:"required"`
}

// fakeClient is a fake Client for testing purposes.
type fakeClient struct {
	currentTimestep uint
	fixtureData     []fixtureTimestep
	fixtureMut      sync.Mutex
}

// newFakeClient instantiates a fakeClient for testing purposes.
func newFakeClient() (*fakeClient, error) {
	blob, err := ioutil.ReadFile("testdata/fake_client_fixtures.json")
	if err != nil {
		return nil, errors.New("Failed to read blockwatch fixture file")
	}

	var fixtureData []fixtureTimestep
	_ = json.Unmarshal(blob, &fixtureData)

	var startTimestep uint = 0
	return &fakeClient{currentTimestep: startTimestep, fixtureData: fixtureData}, nil
}

// HeaderByNumber fetches a block header by its number. If no `number` is supplied, it will return the latest
// block header. If no block exists with this number it will return a `ethereum.NotFound` error.
func (fc *fakeClient) HeaderByNumber(number *big.Int) (*MiniHeader, error) {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	timestep := fc.fixtureData[fc.currentTimestep]
	var miniHeader MiniHeader
	var ok bool
	if number == nil {
		miniHeader = timestep.GetLatestBlock
	} else {
		miniHeader, ok = timestep.GetBlockByNumber[number.Uint64()]
		if !ok {
			return nil, ethereum.NotFound
		}
	}
	return &miniHeader, nil
}

// HeaderByHash fetches a block header by its block hash. If no block exists with this number it will return
// a `ethereum.NotFound` error.
func (fc *fakeClient) HeaderByHash(hash common.Hash) (*MiniHeader, error) {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	timestep := fc.fixtureData[fc.currentTimestep]
	miniHeader, ok := timestep.GetBlockByHash[hash]
	if !ok {
		return nil, ethereum.NotFound
	}
	return &miniHeader, nil
}

// FilterLogs returns the logs that satisfy the supplied filter query.
func (fc *fakeClient) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	// NOT IMPLEMENTED SINCE NOT USED IN TESTS
	return []types.Log{}, nil
}

// IncrementTimestep increments the timestep of the simulation.
func (fc *fakeClient) IncrementTimestep() {
	fc.currentTimestep++
}

// NumberOfTimesteps returns the number of timesteps in the simulation
func (fc *fakeClient) NumberOfTimesteps() int {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	return len(fc.fixtureData)
}

// ExpectedRetainedBlocks returns the expected retained blocks at the current timestep.
func (fc *fakeClient) ExpectedRetainedBlocks() []*MiniHeader {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	return fc.fixtureData[fc.currentTimestep].GetCorrectChain
}

// GetScenarioLabel returns a label describing the test case being tested by the current timestep
// of the simulation.
func (fc *fakeClient) GetScenarioLabel() string {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	return fc.fixtureData[fc.currentTimestep].ScenarioLabel
}

// GetEvents returns the events in the order they should have been emitted by Watcher for
// the current timestep of the simulation.
func (fc *fakeClient) GetEvents() []*Event {
	fc.fixtureMut.Lock()
	defer fc.fixtureMut.Unlock()
	return fc.fixtureData[fc.currentTimestep].BlockEvents
}
