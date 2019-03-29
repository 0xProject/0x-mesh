package blockwatch

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

// TestWatcher tests that Watcher properly stores up to blockRetentionLimit blocks in the correct
// order and also emits block events in the proper order.
func TestWatcher(t *testing.T) {
	fakeClient, err := newFakeClient()
	if err != nil {
		t.Fatal(err.Error())
	}

	var blockRetentionLimit uint = 15
	startBlockDepth := rpc.LatestBlockNumber
	bs := New(startBlockDepth, blockRetentionLimit, fakeClient)

	for i := 0; i < fakeClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeClient.GetScenarioLabel()

		bs.PollNextBlock(context.Background())
		retainedBlocks := bs.InspectRetainedBlocks()
		expectedRetainedBlocks := fakeClient.ExpectedRetainedBlocks()
		assert.Equal(t, expectedRetainedBlocks, retainedBlocks, scenarioLabel)

		expectedEvents := fakeClient.GetEvents()
		// If we expect events to be emitted, check them
		if len(expectedEvents) != 0 {
			select {
			case gotEvents := <-bs.Events:
				assert.Equal(t, expectedEvents, gotEvents, scenarioLabel)

			case <-time.After(3 * time.Second):
				t.Fatal("Timed out waiting for Events channel to deliver expected events")
			}
		}

		fakeClient.IncrementTimestep()
	}
}
