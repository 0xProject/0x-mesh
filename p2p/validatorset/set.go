package validatorset

import (
	"context"
	"sync"

	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

type Set struct {
	mu         sync.RWMutex
	validators []*namedValidator
}

type namedValidator struct {
	name      string
	validator pubsub.Validator
}

func New() *Set {
	return &Set{}
}

func (s *Set) Add(name string, validator pubsub.Validator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	named := &namedValidator{
		name:      name,
		validator: validator,
	}
	s.validators = append(s.validators, named)
}

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
			// TODO(albrow): Change the verbosity of this log to Trace.
			// TODO(albrow): Should we reduce a peer's score as a penalty for invalid
			//               messages?
			log.WithField("validatorName", validator.name).Debug("pubsub message validation failed")
			return false
		}
	}
	return true
}
