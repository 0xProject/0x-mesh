package rpc

import (
	"context"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// Client is a JSON RPC 2.0 client implementation over WebSockets. It can be
// used to communicate with a 0x Mesh node and add orders.
type Client struct {
	rpcClient *rpc.Client
}

// NewClient creates and returns a new client. addr is the address of the server
// (i.e. a 0x Mesh node) to dial.
func NewClient(addr string) (*Client, error) {
	rpcClient, err := rpc.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		rpcClient: rpcClient,
	}, nil
}

// AddOrdersOpts is a set of options for the AddOrders RPC method.
type AddOrdersOpts struct {
	// Pinned determines whether or not the added orders should be pinned. Pinned
	// orders will not be affected by any DDoS prevention or incentive mechanisms
	// and will always stay in storage until they are no longer fillable. Defaults
	// to true.
	Pinned bool `json:"pinned"`
}

// AddOrders adds orders to the 0x Mesh node and broadcasts them throughout the
// 0x Mesh network.
func (c *Client) AddOrders(orders []*zeroex.SignedOrder, opts ...AddOrdersOpts) (*ordervalidator.ValidationResults, error) {
	var validationResults ordervalidator.ValidationResults
	if len(opts) > 0 {
		if err := c.rpcClient.Call(&validationResults, "mesh_addOrders", orders, opts[0]); err != nil {
			return nil, err
		}
	}
	if err := c.rpcClient.Call(&validationResults, "mesh_addOrders", orders); err != nil {
		return nil, err
	}
	return &validationResults, nil
}

// GetOrdersResponse is the response returned for an RPC request to mesh_getOrders
type GetOrdersResponse struct {
	SnapshotID  string       `json:"snapshotID"`
	OrdersInfos []*OrderInfo `json:"ordersInfos"`
}

// GetOrders gets all orders stored on the Mesh node at a particular point in time in a paginated fashion
func (c *Client) GetOrders(page, perPage int, snapshotID string) (*GetOrdersResponse, error) {
	var getOrdersResponse GetOrdersResponse
	if err := c.rpcClient.Call(&getOrdersResponse, "mesh_getOrders", page, perPage, snapshotID); err != nil {
		return nil, err
	}
	return &getOrdersResponse, nil
}

// AddPeer adds the peer to the node's list of peers. The node will attempt to
// connect to this new peer and return an error if it cannot.
func (c *Client) AddPeer(peerInfo peerstore.PeerInfo) error {
	peerIDString := peer.IDB58Encode(peerInfo.ID)
	multiAddrStrings := make([]string, len(peerInfo.Addrs))
	for i, addr := range peerInfo.Addrs {
		multiAddrStrings[i] = addr.String()
	}
	if err := c.rpcClient.Call(nil, "mesh_addPeer", peerIDString, multiAddrStrings); err != nil {
		return err
	}
	return nil
}

// LatestBlock is the latest block processed by the Mesh node
type LatestBlock struct {
	Number int         `json:"number"`
	Hash   common.Hash `json:"hash"`
}

// GetStatsResponse is the response returned for an RPC request to mesh_getStats
type GetStatsResponse struct {
	Version                   string      `json:"version"`
	PubSubTopic               string      `json:"pubSubTopic"`
	Rendezvous                string      `json:"rendervous"`
	PeerID                    string      `json:"peerID"`
	EthereumNetworkID         int         `json:"ethereumNetworkID"`
	LatestBlock               LatestBlock `json:"latestBlock"`
	NumPeers                  int         `json:"numPeers"`
	NumOrders                 int         `json:"numOrders"`
	NumOrdersIncludingRemoved int         `json:"numOrdersIncludingRemoved"`
	MaxExpirationTime         string      `json:"maxExpirationTime"`
}

// GetStats retrieves stats about the Mesh node
func (c *Client) GetStats() (*GetStatsResponse, error) {
	var getStatsResponse *GetStatsResponse
	if err := c.rpcClient.Call(&getStatsResponse, "mesh_getStats"); err != nil {
		return nil, err
	}
	return getStatsResponse, nil
}

// SubscribeToOrders subscribes a stream of order events
// Note copied from `go-ethereum` codebase: Slow subscribers will be dropped eventually. Client
// buffers up to 8000 notifications before considering the subscriber dead. The subscription Err
// channel will receive ErrSubscriptionQueueOverflow. Use a sufficiently large buffer on the channel
// or ensure that the channel usually has at least one reader to prevent this issue.
func (c *Client) SubscribeToOrders(ctx context.Context, ch chan<- []*zeroex.OrderEvent) (*rpc.ClientSubscription, error) {
	return c.rpcClient.Subscribe(ctx, "mesh", ch, "orders")
}

// SubscribeToHeartbeat subscribes a stream of heartbeats in order to have certainty that the WS
// connection is still alive.
// Note copied from `go-ethereum` codebase: Slow subscribers will be dropped eventually. Client
// buffers up to 8000 notifications before considering the subscriber dead. The subscription Err
// channel will receive ErrSubscriptionQueueOverflow. Use a sufficiently large buffer on the channel
// or ensure that the channel usually has at least one reader to prevent this issue.
func (c *Client) SubscribeToHeartbeat(ctx context.Context, ch chan<- string) (*rpc.ClientSubscription, error) {
	return c.rpcClient.Subscribe(ctx, "mesh", ch, "heartbeat")
}
