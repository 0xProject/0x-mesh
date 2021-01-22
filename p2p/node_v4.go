// Implements the GossipSub / PubSub sharing of v4 orders
package p2p

import (
	"context"
	"fmt"
	mathrand "math/rand"
)

const (
	// GossipSub topic for V4 orders
	gossipSubOrdersV4Topic = "0x-orders-v4"
)

func (n *Node) SendV4(data []byte) error {
	var firstErr error
	err := n.pubsub.Publish(gossipSubOrdersV4Topic, data) //nolint:staticcheck
	if err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

func (n *Node) startMessageHandlerV4(ctx context.Context) error {
	// Subscribe to topic if we haven't already
	if n.subV4 == nil {
		var err error
		n.subV4, err = n.pubsub.Subscribe(gossipSubOrdersV4Topic)
		if err != nil {
			return err
		}
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Receive a batch of messages so we can handle them in bulk
		messages := []*Message{}
		for {
			// If the parent context was canceled, break.
			select {
			case <-ctx.Done():
				break
			default:
			}

			// Receive message
			receiveCtx, receiveCancel := context.WithTimeout(n.ctx, receiveTimeout)
			msg, err := n.subV4.Next(receiveCtx)
			receiveCancel()
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					break
				} else {
					return err
				}
			}
			// Skip self messages
			if msg.GetFrom() == n.host.ID() {
				continue
			}

			// Add to batch
			messages = append(messages, &Message{From: msg.GetFrom(), Data: msg.Data})
			if len(messages) >= maxReceiveBatch {
				break
			}
		}

		// Handle messages
		if err := n.messageHandler.HandleMessagesV4(ctx, messages); err != nil {
			return fmt.Errorf("could not validate or store V4 order messages: %s", err.Error())
		}

		// Check bandwidth usage non-deterministically
		if mathrand.Float64() <= chanceToCheckBandwidthUsage {
			n.banner.CheckBandwidthUsage()
		}
	}
}
