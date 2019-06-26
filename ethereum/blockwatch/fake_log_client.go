package blockwatch

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type filterLogsResponse struct {
	Logs []types.Log
	Err  error
}

// fakeLogClient is a fake Client for testing purposes.
type fakeLogClient struct {
	count     int
	responses []filterLogsResponse
	Mu        sync.Mutex
}

// newFakeLogClient instantiates a fakeLogClient for testing log fetching
func newFakeLogClient(responses []filterLogsResponse) (*fakeLogClient, error) {
	return &fakeLogClient{count: 0, responses: responses}, nil
}

// HeaderByNumber fetches a block header by its number
func (fc *fakeLogClient) HeaderByNumber(number *big.Int) (*meshdb.MiniHeader, error) {
	return nil, errors.New("NOT_IMPLEMENTED")
}

// HeaderByHash fetches a block header by its block hash
func (fc *fakeLogClient) HeaderByHash(hash common.Hash) (*meshdb.MiniHeader, error) {
	return nil, errors.New("NOT_IMPLEMENTED")
}

// FilterLogs returns the logs that satisfy the supplied filter query
func (fc *fakeLogClient) FilterLogs(q ethereum.FilterQuery) ([]types.Log, error) {
	fc.Mu.Lock()
	// Add a slight delay to simulate an actual network request. This also gives
	// BlockWatcher.getLogsInBlockRange multi-requests to hit the concurrent request
	// limit semaphore and simulate more realistic conditions.
	<-time.Tick(5 * time.Millisecond)
	defer fc.Mu.Unlock()
	res := fc.responses[fc.count]
	fc.count = fc.count + 1
	return res.Logs, res.Err
}
