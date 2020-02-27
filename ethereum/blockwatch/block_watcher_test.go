// +build !browser

package blockwatch

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/ethereum/simplestack"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var config = Config{
	PollingInterval: 1 * time.Second,
	WithLogs:        false,
	Topics:          []common.Hash{},
}

var (
	basicFakeClientFixture = "testdata/fake_client_basic_fixture.json"
	blockRetentionLimit    = 10
	startMiniHeaders       = []*miniheader.MiniHeader{}
)

func TestWatcher(t *testing.T) {
	fakeClient, err := newFakeClient("testdata/fake_client_block_poller_fixtures.json")
	require.NoError(t, err)

	// Polling interval unused because we hijack the ticker for this test
	require.NoError(t, err)
	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)
	config.Client = fakeClient
	watcher := New(config)

	// Having a buffer of 1 unblocks the below for-loop without resorting to a goroutine
	events := make(chan []*Event, 1)
	sub := watcher.Subscribe(events)

	for i := 0; i < fakeClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeClient.GetScenarioLabel()

		err = watcher.SyncToLatestBlock()
		if strings.HasPrefix(scenarioLabel, "ERROR") {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		retainedBlocks, err := watcher.getAllRetainedBlocks()
		require.NoError(t, err)
		expectedRetainedBlocks := fakeClient.ExpectedRetainedBlocks()
		assert.Equal(t, expectedRetainedBlocks, retainedBlocks, scenarioLabel)

		expectedEvents := fakeClient.GetEvents()
		if len(expectedEvents) != 0 {
			select {
			case gotEvents := <-events:
				assert.Equal(t, expectedEvents, gotEvents, scenarioLabel)

			case <-time.After(3 * time.Second):
				t.Fatal("Timed out waiting for Events channel to deliver expected events")
			}
		}

		fakeClient.IncrementTimestep()

		if i == fakeClient.NumberOfTimesteps()-1 {
			sub.Unsubscribe()
		}
	}
}

func TestWatcherStartStop(t *testing.T) {
	fakeClient, err := newFakeClient(basicFakeClientFixture)
	require.NoError(t, err)

	require.NoError(t, err)
	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)
	config.Client = fakeClient
	watcher := New(config)

	// Start the watcher in a goroutine. We use the done channel to signal when
	// watcher.Watch returns.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	defer cancel()
	go func() {
		require.NoError(t, watcher.Watch(ctx))
		done <- struct{}{}
	}()

	// Wait a bit and then stop the watcher by calling cancel.
	time.Sleep(100 * time.Millisecond)
	cancel()

	// Make sure that the watcher actually stops.
	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("timed out waiting for watcher to stop")
	case <-done:
		break
	}
}

type blockRangeChunksTestCase struct {
	from                int
	to                  int
	expectedBlockRanges []*blockRange
}

func TestGetSubBlockRanges(t *testing.T) {
	rangeSize := 6
	testCases := []blockRangeChunksTestCase{
		blockRangeChunksTestCase{
			from: 10,
			to:   10,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   10,
				},
			},
		},
		blockRangeChunksTestCase{
			from: 10,
			to:   16,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   15,
				},
				&blockRange{
					FromBlock: 16,
					ToBlock:   16,
				},
			},
		},
		blockRangeChunksTestCase{
			from: 10,
			to:   22,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   15,
				},
				&blockRange{
					FromBlock: 16,
					ToBlock:   21,
				},
				&blockRange{
					FromBlock: 22,
					ToBlock:   22,
				},
			},
		},
		blockRangeChunksTestCase{
			from: 10,
			to:   24,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   15,
				},
				&blockRange{
					FromBlock: 16,
					ToBlock:   21,
				},
				&blockRange{
					FromBlock: 22,
					ToBlock:   24,
				},
			},
		},
	}

	fakeClient, err := newFakeClient(basicFakeClientFixture)
	require.NoError(t, err)
	require.NoError(t, err)
	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)
	config.Client = fakeClient
	watcher := New(config)

	for _, testCase := range testCases {
		blockRanges := watcher.getSubBlockRanges(testCase.from, testCase.to, rangeSize)
		assert.Equal(t, testCase.expectedBlockRanges, blockRanges)
	}
}

func TestFastSyncToLatestBlockLessThan128Missed(t *testing.T) {
	// Fixture will return block 132 as the tip of the chain (127 blocks from block 5)
	fakeClient, err := newFakeClient("testdata/fake_client_fast_sync_fixture.json")
	require.NoError(t, err)

	require.NoError(t, err)
	// Add block number 5 as the last block seen by BlockWatcher
	lastBlockSeen := &miniheader.MiniHeader{
		Number:    big.NewInt(5),
		Hash:      common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent:    common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
		Timestamp: time.Now(),
	}

	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)

	err = config.Stack.Push(lastBlockSeen)
	require.NoError(t, err)

	config.Client = fakeClient
	watcher := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	blocksElapsed, err := watcher.FastSyncToLatestBlock(ctx)
	require.NoError(t, err)
	assert.Equal(t, 127, blocksElapsed)

	// Check that block 132 is now in the DB, and block 5 was removed.
	headers, err := config.Stack.PeekAll()
	require.NoError(t, err)
	require.Len(t, headers, 1)
	assert.Equal(t, big.NewInt(132), headers[0].Number)
}

func TestFastSyncToLatestBlockMoreThanOrExactly128Missed(t *testing.T) {
	// Fixture will return block 133 as the tip of the chain (128 blocks from block 5)
	fakeClient, err := newFakeClient("testdata/fake_client_reset_fixture.json")
	require.NoError(t, err)

	require.NoError(t, err)
	// Add block number 5 as the last block seen by BlockWatcher
	lastBlockSeen := &miniheader.MiniHeader{
		Number:    big.NewInt(5),
		Hash:      common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent:    common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
		Timestamp: time.Now(),
	}

	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)

	err = config.Stack.Push(lastBlockSeen)
	require.NoError(t, err)

	config.Client = fakeClient
	watcher := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	blocksElapsed, err := watcher.FastSyncToLatestBlock(ctx)
	require.NoError(t, err)
	assert.Equal(t, 128, blocksElapsed)

	// Check that all blocks have been removed from BlockWatcher
	headers, err := config.Stack.PeekAll()
	require.NoError(t, err)
	require.Len(t, headers, 0)
}

func TestFastSyncToLatestBlockNoneMissed(t *testing.T) {
	// Fixture will return block 5 as the tip of the chain
	fakeClient, err := newFakeClient("testdata/fake_client_basic_fixture.json")
	require.NoError(t, err)

	require.NoError(t, err)
	// Add block number 5 as the last block seen by BlockWatcher
	lastBlockSeen := &miniheader.MiniHeader{
		Number:    big.NewInt(5),
		Hash:      common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent:    common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
		Timestamp: time.Now(),
	}

	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)

	err = config.Stack.Push(lastBlockSeen)
	require.NoError(t, err)

	config.Client = fakeClient
	watcher := New(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	blocksElapsed, err := watcher.FastSyncToLatestBlock(ctx)
	require.NoError(t, err)
	assert.Equal(t, blocksElapsed, 0)

	// Check that block 5 is still in the DB
	headers, err := config.Stack.PeekAll()
	require.NoError(t, err)
	require.Len(t, headers, 1)
	assert.Equal(t, big.NewInt(5), headers[0].Number)
}

var logStub = types.Log{
	Address: common.HexToAddress("0x21ab6c9fac80c59d401b37cb43f81ea9dde7fe34"),
	Topics: []common.Hash{
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x0000000000000000000000004d8a4aa1f304f9632cf3877473445d85c577fe5d"),
		common.HexToHash("0x0000000000000000000000004bdd0d16cfa18e33860470fc4d65c6f5cee60959"),
	},
	Data:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000337ad34c0"),
	BlockNumber: uint64(30),
	TxHash:      common.HexToHash("0xd9bb5f9e888ee6f74bedcda811c2461230f247c205849d6f83cb6c3925e54586"),
	TxIndex:     uint(0),
	BlockHash:   common.HexToHash("0x6bbf9b6e836207ab25379c20e517a89090cbbaf8877746f6ed7fb6820770816b"),
	Index:       uint(0),
	Removed:     false,
}

var errUnexpected = errors.New("Something unexpected")

type filterLogsRecusivelyTestCase struct {
	Label                     string
	rangeToFilterLogsResponse map[string]filterLogsResponse
	Err                       error
	Logs                      []types.Log
}

func TestFilterLogsRecursively(t *testing.T) {
	from := 10
	to := 20
	testCases := []filterLogsRecusivelyTestCase{
		filterLogsRecusivelyTestCase{
			Label: "HAPPY_PATH",
			rangeToFilterLogsResponse: map[string]filterLogsResponse{
				"10-20": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs: []types.Log{logStub},
		},
		filterLogsRecusivelyTestCase{
			Label: "TOO_MANY_RESULTS_INFURA_ERROR",
			rangeToFilterLogsResponse: map[string]filterLogsResponse{
				"10-20": filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				"10-15": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				"16-20": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs: []types.Log{logStub, logStub},
		},
		filterLogsRecusivelyTestCase{
			Label: "TOO_MANY_RESULTS_INFURA_ERROR_DEEPER_RECURSION",
			rangeToFilterLogsResponse: map[string]filterLogsResponse{
				"10-20": filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				"10-15": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				"16-20": filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				"16-18": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				"19-20": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs: []types.Log{logStub, logStub, logStub},
		},
		filterLogsRecusivelyTestCase{
			Label: "TOO_MANY_RESULTS_INFURA_ERROR_DEEPER_RECURSION_FAILURE",
			rangeToFilterLogsResponse: map[string]filterLogsResponse{
				"10-20": filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				"10-15": filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				"16-20": filterLogsResponse{
					Err: errUnexpected,
				},
			},
			Err: errUnexpected,
		},
		filterLogsRecusivelyTestCase{
			Label: "UNEXPECTED_ERROR",
			rangeToFilterLogsResponse: map[string]filterLogsResponse{
				"10-20": filterLogsResponse{
					Err: errUnexpected,
				},
			},
			Err: errUnexpected,
		},
	}

	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)

	for _, testCase := range testCases {
		fakeLogClient, err := newFakeLogClient(testCase.rangeToFilterLogsResponse)
		require.NoError(t, err)
		config.Client = fakeLogClient
		watcher := New(config)

		logs, err := watcher.filterLogsRecurisively(from, to, []types.Log{})
		require.Equal(t, testCase.Err, err, testCase.Label)
		require.Equal(t, testCase.Logs, logs, testCase.Label)
		assert.Equal(t, len(testCase.rangeToFilterLogsResponse), fakeLogClient.Count())
	}
}

type logsInBlockRangeTestCase struct {
	Label                     string
	From                      int
	To                        int
	RangeToFilterLogsResponse map[string]filterLogsResponse
	Logs                      []types.Log
	FurthestBlockProcessed    int
}

func TestGetLogsInBlockRange(t *testing.T) {
	from := 10
	to := 20
	testCases := []logsInBlockRangeTestCase{
		logsInBlockRangeTestCase{
			Label: "HAPPY_PATH",
			From:  from,
			To:    to,
			RangeToFilterLogsResponse: map[string]filterLogsResponse{
				aRange(from, to): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs:                   []types.Log{logStub},
			FurthestBlockProcessed: to,
		},
		logsInBlockRangeTestCase{
			Label: "SPLIT_REQUEST_BY_MAX_BLOCKS_IN_QUERY",
			From:  from,
			To:    from + maxBlocksInGetLogsQuery + 10,
			RangeToFilterLogsResponse: map[string]filterLogsResponse{
				aRange(from, from+maxBlocksInGetLogsQuery-1): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				aRange(from+maxBlocksInGetLogsQuery, from+maxBlocksInGetLogsQuery+10): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs:                   []types.Log{logStub, logStub},
			FurthestBlockProcessed: from + maxBlocksInGetLogsQuery + 10,
		},
		logsInBlockRangeTestCase{
			Label: "SHORT_CIRCUIT_SEMAPHORE_BLOCKED_REQUESTS_ON_ERROR",
			From:  from,
			To:    from + (maxBlocksInGetLogsQuery * (getLogsRequestChunkSize + 1)),
			RangeToFilterLogsResponse: map[string]filterLogsResponse{
				// Same number of responses as the getLogsRequestChunkSize since the
				// error response will stop any further requests.
				aRange(from, from+maxBlocksInGetLogsQuery-1): filterLogsResponse{
					Err: errUnexpected,
				},
				aRange(from+maxBlocksInGetLogsQuery, from+(maxBlocksInGetLogsQuery*2)-1): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				aRange(from+(maxBlocksInGetLogsQuery*2), from+(maxBlocksInGetLogsQuery*3)-1): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
			},
			Logs:                   []types.Log{},
			FurthestBlockProcessed: from - 1,
		},
		logsInBlockRangeTestCase{
			Label: "CORRECT_FURTHEST_BLOCK_PROCESSED_ON_ERROR",
			From:  from,
			To:    from + maxBlocksInGetLogsQuery + 10,
			RangeToFilterLogsResponse: map[string]filterLogsResponse{
				aRange(from, from+maxBlocksInGetLogsQuery-1): filterLogsResponse{
					Logs: []types.Log{
						logStub,
					},
				},
				aRange(from+maxBlocksInGetLogsQuery, from+maxBlocksInGetLogsQuery+10): filterLogsResponse{
					Err: errUnexpected,
				}},
			Logs:                   []types.Log{logStub},
			FurthestBlockProcessed: from + maxBlocksInGetLogsQuery - 1,
		},
	}

	config.Stack = simplestack.New(blockRetentionLimit, startMiniHeaders)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, testCase := range testCases {
		fakeLogClient, err := newFakeLogClient(testCase.RangeToFilterLogsResponse)
		require.NoError(t, err)
		config.Client = fakeLogClient
		watcher := New(config)

		logs, furthestBlockProcessed := watcher.getLogsInBlockRange(ctx, testCase.From, testCase.To)
		require.Equal(t, testCase.FurthestBlockProcessed, furthestBlockProcessed, testCase.Label)
		require.Equal(t, testCase.Logs, logs, testCase.Label)
		assert.Equal(t, len(testCase.RangeToFilterLogsResponse), fakeLogClient.Count())
	}
}

func TestIsWarning(t *testing.T) {
	errs := map[error]bool{
		errors.New("not found"):     true,
		errors.New("unknown block"): true,
		errors.New("Post https://eth-mainnet.alchemyapi.io/jsonrpc: context deadline exceeded"): true,
		errors.New("couldn't parse parameters: blockhash"):                                      false,
	}

	for err, expected := range errs {
		assert.Equal(t, expected, isWarning(err))
	}
}

func aRange(from, to int) string {
	r := fmt.Sprintf("%d-%d", from, to)
	return r
}
