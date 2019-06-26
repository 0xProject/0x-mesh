package blockwatch

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/ethereum/go-ethereum/common"
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
	if err != nil {
		t.Fatal(err.Error())
	}

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
	if err != nil {
		t.Fatal(err.Error())
	}

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

type testCase struct {
	from                int64
	to                  int64
	expectedBlockRanges []*blockRange
}

func TestGetBlockRangeChunks(t *testing.T) {
	chunkSize := int64(6)
	testCases := []testCase{
		testCase{
			from: 10,
			to:   10,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   10,
				},
			},
		},
		testCase{
			from: 10,
			to:   16,
			expectedBlockRanges: []*blockRange{
				&blockRange{
					FromBlock: 10,
					ToBlock:   16,
				},
			},
		},
		testCase{
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
		testCase{
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
	if err != nil {
		t.Fatal(err.Error())
	}
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
	if err != nil {
		t.Fatal(err.Error())
	}

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
	if err != nil {
		t.Fatal(err.Error())
	}

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
