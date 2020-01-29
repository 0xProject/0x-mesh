// +build !js

package ordersync

import (
	"context"
	"testing"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type inMemoryOrderProvider struct {
	orderData map[string][]byte
}

func (p inMemoryOrderProvider) ProvideOrders(topic string) ([]byte, error) {
	if p.orderData == nil {
		return nil, nil
	}
	return p.orderData[topic], nil
}

func newTestService(t *testing.T, orderData map[string][]byte) *Service {
	basicHost, err := libp2p.New(context.Background())
	require.NoError(t, err)
	provider := &inMemoryOrderProvider{
		orderData: orderData,
	}
	return New(basicHost, provider)
}

func TestOrderSync(t *testing.T) {
	orderData := map[string][]byte{
		"topic-0": []byte("orders-for-topic-0"),
		"topic-1": []byte("orders-for-topic-1"),
	}
	providingService := newTestService(t, orderData)
	requestingService := newTestService(t, nil)

	err := providingService.host.Connect(context.Background(), peer.AddrInfo{
		ID:    requestingService.host.ID(),
		Addrs: requestingService.host.Addrs(),
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for topic, expectedData := range orderData {
		actualData, err := requestingService.GetOrders(ctx, topic)
		require.NoError(t, err)
		assert.Equal(t, string(expectedData), string(actualData), "incorrect order data for topic: %q", topic)
	}

	actualData, err := requestingService.GetOrders(ctx, "unkown-topic")
	assert.Error(t, err, "expected error when getting orders for an unknown topic")
	assert.Equal(t, ErrNoOrders, err, "wrong error type")
	assert.Nil(t, actualData, "actual data should be nil for unknown topic")
}
