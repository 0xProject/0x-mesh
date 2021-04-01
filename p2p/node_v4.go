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
	for _, topic := range n.config.PublishTopicsV4 {
		err := n.pubsub.Publish(topic, data) //nolint:staticcheck
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (n *Node) receiveAndHandleMessagesV4(ctx context.Context) error {
	// Subscribe to topic if we haven't already
	if n.subV4 == nil {
		var err error
		n.subV4, err = n.pubsub.Subscribe(gossipSubOrdersV4Topic)
		if err != nil {
			return err
		}
	}

	// Receive up to maxReceiveBatch messages.
	incoming, err := n.receiveBatchV4(ctx)
	if err != nil {
		return err
	}
	if len(incoming) == 0 {
		return nil
	}
	if err := n.messageHandler.HandleMessagesV4(ctx, incoming); err != nil {
		return fmt.Errorf("could not validate or store v4 messages: %s", err.Error())
	}

	return nil
}

func (n *Node) receiveV4(ctx context.Context) (*Message, error) {
	msg, err := n.subV4.Next(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{From: msg.GetFrom(), Data: msg.Data}, nil
}

func (n *Node) receiveBatchV4(ctx context.Context) ([]*Message, error) {
	// Receive a batch of messages so we can handle them in bulk
	messages := []*Message{}
	for {
		if len(messages) >= maxReceiveBatch {
			return messages, nil
		}
		// If the parent context was canceled, break.
		select {
		case <-ctx.Done():
			return messages, nil
		default:
		}

		// Receive message
		receiveCtx, receiveCancel := context.WithTimeout(n.ctx, receiveTimeout)
		msg, err := n.receiveV4(receiveCtx)
		receiveCancel()
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return messages, nil
			} else {
				return nil, err
			}
		}
		// Skip self messages
		if msg.From == n.host.ID() {
			continue
		}

		// Add to batch
		messages = append(messages, msg)
	}
}

func (n *Node) startMessageHandlerV4(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := n.receiveAndHandleMessagesV4(ctx); err != nil {
			return err
		}

		// Check bandwidth usage non-deterministically
		if mathrand.Float64() <= chanceToCheckBandwidthUsage {
			n.banner.CheckBandwidthUsage()
		}
	}
}
