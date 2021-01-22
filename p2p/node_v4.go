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

func (n *Node) receiveV4(ctx context.Context) (*Message, error) {
	if n.subV4 == nil {
		var err error
		n.subV4, err = n.pubsub.Subscribe(gossipSubOrdersV4Topic)
		if err != nil {
			return nil, err
		}
	}
	msg, err := n.subV4.Next(ctx)
	if err != nil {
		return nil, err
	}
	return &Message{From: msg.GetFrom(), Data: msg.Data}, nil
}

// receiveBatch returns up to maxReceiveBatch messages which are received from
// peers. There is no guarantee that the messages are unique.
func (n *Node) receiveBatchV4(ctx context.Context) ([]*Message, error) {
	messages := []*Message{}
	for {
		if len(messages) >= maxReceiveBatch {
			return messages, nil
		}
		select {
		// If the parent context was canceled, return.
		case <-ctx.Done():
			return messages, nil
		default:
		}
		receiveCtx, receiveCancel := context.WithTimeout(n.ctx, receiveTimeout)
		msg, err := n.receiveV4(receiveCtx)
		receiveCancel()
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				return messages, nil
			}
			return nil, err
		}
		if msg.From == n.host.ID() {
			continue
		}
		messages = append(messages, msg)
	}
}

func (n *Node) receiveAndHandleMessagesV4(ctx context.Context) error {
	// Receive up to maxReceiveBatch messages.
	incoming, err := n.receiveBatchV4(ctx)
	if err != nil {
		return err
	}
	if len(incoming) == 0 {
		return nil
	}
	if err := n.messageHandler.HandleMessagesV4(ctx, incoming); err != nil {
		return fmt.Errorf("could not validate or store messages: %s", err.Error())
	}
	return nil
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
