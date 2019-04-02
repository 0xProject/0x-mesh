package blockwatch

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
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
	watcher := New(pollingInterval, startBlockDepth, blockRetentionLimit, fakeClient)

	// HACK(fabio): By default `blockwatch.Watcher` starts polling for blocks as soon as the first
	// subscription is established. Since we want to control the polling interval, we hijack the
	// internal ticker channel
	watcher.isWatching = true // When isWatching = true, `Subscribe` doesn't start the polling loop
	// Having a buffer of 1 unblocks the below for-loop without resorting to a goroutine
	events := make(chan []*Event, 1)
	sub := watcher.Subscribe(events)
	fakeTickerChan := make(chan time.Time, 1)
	fakeTicker := &time.Ticker{
		C: fakeTickerChan,
	}
	watcher.ticker = fakeTicker // Replace default ticker with our own custom ticker
	go watcher.startPolling()   // Start polling. Blocks waiting to receive from ticker channel

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
