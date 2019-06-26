package blockwatch

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var config = Config{
	PollingInterval:     1 * time.Second,
	BlockRetentionLimit: 10,
	StartBlockDepth:     rpc.LatestBlockNumber,
	WithLogs:            false,
	Topics:              []common.Hash{},
}

var basicFakeClientFixture = "testdata/fake_client_basic_fixture.json"

func TestWatcher(t *testing.T) {
	fakeClient, err := newFakeClient("testdata/fake_client_block_poller_fixtures.json")
	require.NoError(t, err)

	// Polling interval unused because we hijack the ticker for this test
	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	config.MeshDB = meshDB
	config.Client = fakeClient
	watcher := New(config)

	// Having a buffer of 1 unblocks the below for-loop without resorting to a goroutine
	events := make(chan []*Event, 1)
	sub := watcher.Subscribe(events)

	for i := 0; i < fakeClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeClient.GetScenarioLabel()

		err := watcher.pollNextBlock()
		require.NoError(t, err)

		retainedBlocks, err := watcher.InspectRetainedBlocks()
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

	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	config.MeshDB = meshDB
	config.Client = fakeClient
	watcher := New(config)
	require.NoError(t, watcher.Start())
	watcher.stopPolling()
	require.NoError(t, watcher.Start())
	watcher.Stop()
}

type blockRangeChunksTestCase struct {
	from                int64
	to                  int64
	expectedBlockRanges []*blockRange
}

func TestGetBlockRangeChunks(t *testing.T) {
	chunkSize := int64(6)
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
	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	config.MeshDB = meshDB
	config.Client = fakeClient
	watcher := New(config)

	for _, testCase := range testCases {
		blockRanges := watcher.getBlockRangeChunks(testCase.from, testCase.to, chunkSize)
		assert.Equal(t, testCase.expectedBlockRanges, blockRanges)
	}
}

func TestGetMissedEventsToBackfillSomeMissed(t *testing.T) {
	// Fixture will return block 30 as the tip of the chain
	fakeClient, err := newFakeClient("testdata/fake_client_fast_sync_fixture.json")
	require.NoError(t, err)

	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	// Add block number 5 as the last block seen by BlockWatcher
	lastBlockSeen := &meshdb.MiniHeader{
		Number: big.NewInt(5),
		Hash:   common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent: common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
	}
	meshDB.MiniHeaders.Insert(lastBlockSeen)

	config.MeshDB = meshDB
	config.Client = fakeClient
	watcher := New(config)

	events, err := watcher.getMissedEventsToBackfill()
	require.NoError(t, err)
	assert.Len(t, events, 1)

	// Check that block 30 is now in the DB, and block 5 was removed.
	headers, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	require.Len(t, headers, 1)
	assert.Equal(t, big.NewInt(30), headers[0].Number)
}

func TestGetMissedEventsToBackfillNoneMissed(t *testing.T) {
	// Fixture will return block 5 as the tip of the chain
	fakeClient, err := newFakeClient("testdata/fake_client_basic_fixture.json")
	require.NoError(t, err)

	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	// Add block number 5 as the last block seen by BlockWatcher
	lastBlockSeen := &meshdb.MiniHeader{
		Number: big.NewInt(5),
		Hash:   common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent: common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
	}
	meshDB.MiniHeaders.Insert(lastBlockSeen)

	config.MeshDB = meshDB
	config.Client = fakeClient
	watcher := New(config)

	events, err := watcher.getMissedEventsToBackfill()
	require.NoError(t, err)
	assert.Len(t, events, 0)

	// Check that block 5 is still in the DB
	headers, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	require.Len(t, headers, 1)
	assert.Equal(t, big.NewInt(5), headers[0].Number)
}

type filterLogsRecusivelyTestCase struct {
	Label               string
	FilterLogsResponses []filterLogsResponse
	Err                 error
	Logs                []types.Log
}

func TestFilterLogsRecursively(t *testing.T) {
	log := types.Log{
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
	testCases := []filterLogsRecusivelyTestCase{
		filterLogsRecusivelyTestCase{
			Label: "HAPPY_PATH",
			FilterLogsResponses: []filterLogsResponse{
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
			},
			Logs: []types.Log{log},
		},
		filterLogsRecusivelyTestCase{
			Label: "TOO_MANY_RESULTS_INFURA_ERROR",
			FilterLogsResponses: []filterLogsResponse{
				filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
			},
			Logs: []types.Log{log, log},
		},
		filterLogsRecusivelyTestCase{
			Label: "TOO_MANY_RESULTS_INFURA_ERROR_DEEPER_RECURSION",
			FilterLogsResponses: []filterLogsResponse{
				filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
				filterLogsResponse{
					Err: errors.New(infuraTooManyResultsErrMsg),
				},
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
				filterLogsResponse{
					Logs: []types.Log{
						log,
					},
				},
			},
			Logs: []types.Log{log, log, log},
		},
		filterLogsRecusivelyTestCase{
			Label: "UNEXPECTED_ERROR",
			FilterLogsResponses: []filterLogsResponse{
				filterLogsResponse{
					Err: errors.New("Something unexpected"),
				},
			},
			Err: errors.New("Something unexpected"),
		},
	}

	meshDB, err := meshdb.NewMeshDB("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	config.MeshDB = meshDB

	for _, testCase := range testCases {
		fakeLogClient, err := newFakeLogClient(testCase.FilterLogsResponses)
		require.NoError(t, err)
		config.Client = fakeLogClient
		watcher := New(config)

		logs, err := watcher.filterLogsRecurisively(int64(10), int64(20), []types.Log{})
		require.Equal(t, testCase.Err, err, testCase.Label)
		require.Equal(t, testCase.Logs, logs, testCase.Label)
	}
}
