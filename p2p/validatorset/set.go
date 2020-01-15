// Package validatorset offers a way to combine a set of libp2p.Validators into
// a single validator. The combined validator set only passes if *all* of its
// constituent validators pass.
package validatorset

import (
	"context"
	"sync"

	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

// Set is a set of libp2p.Validators.
type Set struct {
	mu         sync.RWMutex
	validators []*namedValidator
}

type namedValidator struct {
	name      string
	validator pubsub.Validator
}

// New creates a new validator set
func New() *Set {
	return &Set{}
}

// Add adds a new validator to the set with the given name. The name is used
// in error messages.
func (s *Set) Add(name string, validator pubsub.Validator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	named := &namedValidator{
		name:      name,
		validator: validator,
	}
	s.validators = append(s.validators, named)
}

// Validate validates the message. It returns true if all of the constituent
// validators in the set also return true. If one or more of them return false,
// Validate returns false.
func (s *Set) Validate(ctx context.Context, sender peer.ID, msg *pubsub.Message) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, validator := range s.validators {
		// If the context is done return immediately.
		select {
		case <-ctx.Done():
			return false
		default:
		}

		// Otherwise continue by running this validator.
		isValid := validator.validator(ctx, sender, msg)
		if !isValid {
			// TODO(albrow): Should we reduce a peer's score as a penalty for invalid
			//               messages?
			log.WithField("validatorName", validator.name).Trace("pubsub message validation failed")
			return false
		}
	}
	return true
}
