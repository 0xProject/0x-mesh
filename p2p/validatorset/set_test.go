package validatorset

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"testing"
	"time"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE(jalextowle): We must ignore this flag to prevent the flag package from
// panicking when this flag is provided to `wasmbrowsertest` in the browser tests.
func init() {
	_ = flag.String("initFile", "", "")
}

// alwaysFalseValidator is a pubsub.Validator that always returns false.
func alwaysFalseValidator(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	return false
}

// alwaysTrueValidator is a pubsub.Validator that always returns true.
func alwaysTrueValidator(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	return true
}

func TestValidatorSet(t *testing.T) {
	t.Parallel()

	// For each test case, we construct a new validator set with the given
	// validators. Then we check that the actual results of set.Validate matches
	// the expected result.
	testCases := []struct {
		validators     []pubsub.Validator
		expectedResult bool
	}{
		{
			validators:     []pubsub.Validator{alwaysTrueValidator},
			expectedResult: true,
		},
		{
			validators:     []pubsub.Validator{alwaysFalseValidator},
			expectedResult: false,
		},
		{
			validators:     []pubsub.Validator{alwaysTrueValidator, alwaysTrueValidator},
			expectedResult: true,
		},
		{
			validators:     []pubsub.Validator{alwaysFalseValidator, alwaysFalseValidator},
			expectedResult: false,
		},
		{
			validators:     []pubsub.Validator{alwaysTrueValidator, alwaysFalseValidator},
			expectedResult: false,
		},
		{
			validators:     []pubsub.Validator{alwaysTrueValidator, alwaysFalseValidator, alwaysTrueValidator},
			expectedResult: false,
		},
	}

	for i, testCase := range testCases {
		testCaseName := fmt.Sprintf("validatorset test case %d", i)
		t.Run(testCaseName, func(t *testing.T) {
			set := New()
			for j, validator := range testCase.validators {
				set.Add(fmt.Sprintf("validator %d", j), validator)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			actualResult := set.Validate(ctx, getRandomPeerID(t), &pubsub.Message{})
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}

func getRandomPeerID(t *testing.T) peer.ID {
	privKey, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
	require.NoError(t, err)
	id, err := peer.IDFromPrivateKey(privKey)
	require.NoError(t, err)
	return id
}
