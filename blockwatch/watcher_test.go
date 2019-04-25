package blockwatch

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	fakeClient, err := newFakeClient()
	if err != nil {
		t.Fatal(err.Error())
	}

	// Polling interval unused because we hijack the ticker for this test
	pollingInterval := 1 * time.Second
	blockRetentionLimit := 10
	startBlockDepth := rpc.LatestBlockNumber
	withLogs := false
	topics := []common.Hash{}
	watcher := New(pollingInterval, startBlockDepth, blockRetentionLimit, withLogs, topics, fakeClient)

	// Having a buffer of 1 unblocks the below for-loop without resorting to a goroutine
	events := make(chan []*Event, 1)
	sub := watcher.Subscribe(events)

	// Replace default ticker with our own custom ticker
	fakeTickerChan := make(chan time.Time, 1)
	fakeTicker := &time.Ticker{
		C: fakeTickerChan,
	}
	watcher.ticker = fakeTicker
	watcher.isWatching = true
	go watcher.startPollingLoop()

	for i := 0; i < fakeClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeClient.GetScenarioLabel()

		fakeTickerChan <- time.Now()
		time.Sleep(10 * time.Millisecond) // Ensure pollNextBlock runs

		retainedBlocks := watcher.InspectRetainedBlocks()
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
	fakeClient, err := newFakeClient()
	if err != nil {
		t.Fatal(err.Error())
	}

	pollingInterval := 1 * time.Second
	blockRetentionLimit := 10
	startBlockDepth := rpc.LatestBlockNumber
	withLogs := false
	topics := []common.Hash{}
	watcher := New(pollingInterval, startBlockDepth, blockRetentionLimit, withLogs, topics, fakeClient)
	require.NoError(t, watcher.StartPolling())
	watcher.StopPolling()
	require.NoError(t, watcher.StartPolling())
	watcher.StopPolling()
}
