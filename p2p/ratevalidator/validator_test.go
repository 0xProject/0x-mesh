package ratevalidator

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
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

	validator, err := New(ctx, Config{
		MyPeerID:       peerIDs[0],
		GlobalLimit:    rate.Inf,
		PerPeerLimit:   1,
		PerPeerBurst:   5,
		MaxMessageSize: 1024,
	})
	require.NoError(t, err)

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
}

func TestValidatorWithOwnPeerID(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	validator, err := New(ctx, Config{
		MyPeerID:       peerIDs[0],
		GlobalLimit:    1,
		GlobalBurst:    1,
		PerPeerLimit:   1,
		PerPeerBurst:   1,
		MaxMessageSize: 1024,
	})
	require.NoError(t, err)

	// All messages should sent by us should be valid.
	messagesToSend := validator.config.PerPeerBurst + validator.config.GlobalBurst + 5
	for i := 0; i < messagesToSend; i++ {
		valid := validator.Validate(ctx, peerIDs[0], &pubsub.Message{})
		assert.True(t, valid, "message should be valid")
	}
}

func TestValidatorMaxMessageSize(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	maxSize := 48
	validator, err := New(ctx, Config{
		MyPeerID:       peerIDs[0],
		GlobalLimit:    rate.Inf,
		PerPeerLimit:   rate.Inf,
		MaxMessageSize: maxSize,
	})
	require.NoError(t, err)

	valid := validator.Validate(ctx, peerIDs[1], &pubsub.Message{
		Message: &pb.Message{
			Data: make([]byte, maxSize+1),
		},
	})
	assert.False(t, valid, "message should be valid")
}
