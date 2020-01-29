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
	topic0ActualOrders, err := requestingService.GetOrders(ctx, "topic-0")
	require.NoError(t, err)
	assert.Equal(t, string(orderData["topic-0"]), string(topic0ActualOrders), "incorrect order data for topic 0")

	// TODO(albrow): Same test for topic-1
}
