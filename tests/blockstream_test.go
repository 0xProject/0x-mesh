package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/0xProject/0x-mesh/blockstream"
	"github.com/ethereum/go-ethereum/rpc"
)

// TestBlockstream tests that BlockStream properly stores up to blockRetentionLimit blocks in the correct
// order and also emits block events in the proper order.
func TestBlockstream(t *testing.T) {
	fakeBlockClient := NewFakeBlockClient()

	var blockRetentionLimit uint = 15
	startBlockDepth := rpc.LatestBlockNumber
	bs := blockstream.NewBlockStream(startBlockDepth, blockRetentionLimit, fakeBlockClient)

	for i := 0; i < fakeBlockClient.NumberOfTimesteps(); i++ {
		scenarioLabel := fakeBlockClient.GetScenarioLabel()

		expectedBlockEvents := fakeBlockClient.GetBlockEvents()
		// If we expect events to be emitted, check them
		if len(expectedBlockEvents) != 0 {
			go func() {
				var gotBlockEvents []*blockstream.BlockEvent
				for i = 0; i <= len(expectedBlockEvents)-1; i++ {
					gotBlockEvent := <-bs.Stream
					gotBlockEvents = append(gotBlockEvents, gotBlockEvent)
				}
				if !reflect.DeepEqual(expectedBlockEvents, gotBlockEvents) {
					gotBlockEventsJson, err := json.MarshalIndent(gotBlockEvents, "", "	")
					if err != nil {
						panic(err)
					}
					expectedBlockEventsJson, err := json.MarshalIndent(expectedBlockEvents, "", "	")
					if err != nil {
						panic(err)
					}
					fmt.Printf("GOT EVENTS: %v\n", string(gotBlockEventsJson))
					fmt.Printf("EXPECTED EVENTS: %v\n", string(expectedBlockEventsJson))
					t.Fatal("Failed emitted event test:", scenarioLabel)
				}
			}()
		}

		bs.PollNextBlock(context.Background())
		retainedBlocks := bs.GetRetainedBlocks()
		if uint(len(retainedBlocks)) > blockRetentionLimit {
			t.Fatal("Blockstream retained more blocks then specified in block retention limit")
		}
		expectedRetainedBlocks := fakeBlockClient.ExpectedRetainedBlocks()
		if !reflect.DeepEqual(retainedBlocks, expectedRetainedBlocks) {
			retainedBlocksJson, err := json.MarshalIndent(retainedBlocks, "", "	")
			if err != nil {
				panic(err)
			}
			expectedRetainedBlocksJson, err := json.MarshalIndent(expectedRetainedBlocks, "", "	")
			if err != nil {
				panic(err)
			}
			fmt.Printf("GOT BLOCKS: %v\n", string(retainedBlocksJson))
			fmt.Printf("EXPECTED BLOCK: %v\n", string(expectedRetainedBlocksJson))
			t.Fatal("Failed retained block test: ", scenarioLabel)
		}

		fakeBlockClient.IncrementTimestep()
	}
}
