// +build !js

package ws

import (
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dummyOrderHandler is used for testing purposes. It allows declaring handlers
// for some requests or all of them, depending on testing needs.
type dummyOrderHandler struct {
	addOrderHandler func(order *zeroex.SignedOrder) error
}

func (d *dummyOrderHandler) AddOrder(order *zeroex.SignedOrder) error {
	if d.addOrderHandler == nil {
		return errors.New("dummyOrderHandler: no handler set for AddOrder")
	}
	return d.addOrderHandler(order)
}

// newTestServerAndClient returns a server and client which have been connected
// to one another on the local network. The server will use the given
// orderHandler to handle incoming requests. Useful for testing purposes. Will
// block until both the server and client are running and connected to one
// another.
func newTestServerAndClient(t *testing.T, orderHandler *dummyOrderHandler) (*Server, *Client) {
	// Start a new server.
	server, err := NewServer(":0", orderHandler)
	require.NoError(t, err)
	go func() {
		_ = server.Listen()
	}()

	// We need to wait for the OS to choose an available port and for server.Addr
	// to return a non-nil value.
	for server.Addr() == nil {
		time.Sleep(10 * time.Millisecond)
	}

	// Create a new client which is connected to the server.
	client, err := NewClient("ws://" + server.Addr().String())
	require.NoError(t, err)

	return server, client
}

var testOrder = &zeroex.SignedOrder{
	MakerAddress:          common.HexToAddress("0x6924a03bb710eaf199ab6ac9f2bb148215ae9b5d"),
	TakerAddress:          constants.NullAddress,
	SenderAddress:         constants.NullAddress,
	FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
	MakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064"),
	TakerAssetData:        common.Hex2Bytes("f47261b000000000000000000000000025b8fe1de9daf8ba351890744ff28cf7dfa8f5e3"),
	Salt:                  big.NewInt(1548619145450),
	MakerFee:              big.NewInt(0),
	TakerFee:              big.NewInt(0),
	MakerAssetAmount:      big.NewInt(3551808554499581700),
	TakerAssetAmount:      big.NewInt(300000000000000),
	ExpirationTimeSeconds: big.NewInt(1548619325),
	ExchangeAddress:       constants.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
}

func TestAddOrder(t *testing.T) {
	// Set up the dummy handler with an addOrderHandler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	orderHandler := &dummyOrderHandler{
		addOrderHandler: func(order *zeroex.SignedOrder) error {
			assert.Equal(t, testOrder, order, "AddOrder was called with an unexpected order argument")
			wg.Done()
			return nil
		},
	}

	server, client := newTestServerAndClient(t, orderHandler)
	defer server.Close()

	actualOrderHash, err := client.AddOrder(testOrder)
	require.NoError(t, err)
	expectedOrderHash, err := testOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash.String(), actualOrderHash.String(), "returned orderHash did not match")

	// The WaitGroup signals that AddOrder was called on the server-side.
	wg.Wait()
}
