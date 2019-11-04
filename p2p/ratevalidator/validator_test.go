package ratevalidator

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

var peerIDStrings = []string{
	"16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7",
	"16Uiu2HAmVqV4kepwSiNRmvKiBxwpt4EQJi3pAe9auSMyGjzA1eBZ",
	"16Uiu2HAmAmmoyR4M492Aq8vWFh4gyVr9Gz2uEGAWjdpGPfKpcw5F",
}

var peerIDs []peer.ID

func init() {
	for _, peerIDString := range peerIDStrings {
		peerID, _ := peer.IDB58Decode(peerIDString)
		peerIDs = append(peerIDs, peerID)
	}
}

func TestValidatorPerPeer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	validator, err := New(Config{
		MyPeerID:     peerIDs[0],
		GlobalLimit:  rate.Inf,
		PerPeerLimit: 1,
		PerPeerBurst: 5,
	})
	require.NoError(t, err)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, validator.Start(ctx))
	}()

	// Wait for validator to start.
	require.NoError(t, validator.waitForStart(ctx))

	for _, peerID := range peerIDs[1:] {
		// All messages should be valid until we hit GlobalBurst.
		for i := 0; i < validator.config.PerPeerBurst; i++ {
			valid := validator.Validate(ctx, peerID, &pubsub.Message{})
			assert.True(t, valid, "message should be valid")
		}
		// Next message should be invalid.
		valid := validator.Validate(ctx, peerID, &pubsub.Message{})
		assert.False(t, valid, "message should be invalid")
	}

	// Wait one second. Limiter should now allow each peer to send one additional
	// message.
	time.Sleep(1 * time.Second)

	for _, peerID := range peerIDs[1:] {
		// First message should be valid.
		valid := validator.Validate(ctx, peerID, &pubsub.Message{})
		assert.True(t, valid, "message should be valid")

		// Next message should be invalid.
		valid = validator.Validate(ctx, peerID, &pubsub.Message{})
		assert.False(t, valid, "message should be invalid")
	}

	cancel()
	wg.Wait()
}

func TestValidatorWithOwnPeerID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	validator, err := New(Config{
		MyPeerID:     peerIDs[0],
		GlobalLimit:  1,
		GlobalBurst:  1,
		PerPeerLimit: 1,
		PerPeerBurst: 1,
	})
	require.NoError(t, err)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		require.NoError(t, validator.Start(ctx))
	}()

	// Wait for validator to start.
	require.NoError(t, validator.waitForStart(ctx))

	// All messages should sent by us should be valid.
	messagesToSend := validator.config.PerPeerBurst + validator.config.GlobalBurst + 5
	for i := 0; i < messagesToSend; i++ {
		valid := validator.Validate(ctx, peerIDs[0], &pubsub.Message{})
		assert.True(t, valid, "message should be valid")
	}

	cancel()
	wg.Wait()
}
