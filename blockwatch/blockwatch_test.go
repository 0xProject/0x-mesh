package blockwatch

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

// TestBlockWatch tests that BlockWatch properly stores up to blockRetentionLimit blocks in the correct
// order and also emits block events in the proper order.
func TestBlockWatch(t *testing.T) {
	fakeBlockClient := NewFakeBlockClient()

	var blockRetentionLimit uint = 15
	startBlockDepth := rpc.LatestBlockNumber
	bs := NewBlockWatch(startBlockDepth, blockRetentionLimit, fakeBlockClient)

	for i := 0; i < fakeBlockClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeBlockClient.GetScenarioLabel()

		bs.PollNextBlock(context.Background())
		retainedBlocks := bs.GetRetainedBlocks()
		expectedRetainedBlocks := fakeBlockClient.ExpectedRetainedBlocks()
		assert.Equal(t, expectedRetainedBlocks, retainedBlocks, scenarioLabel)

		expectedBlockEvents := fakeBlockClient.GetBlockEvents()
		// If we expect events to be emitted, check them
		if len(expectedBlockEvents) != 0 {
			gotBlockEvents := <-bs.Events
			assert.Equal(t, expectedBlockEvents, gotBlockEvents, scenarioLabel)
		}

		fakeBlockClient.IncrementTimestep()
	}
}
