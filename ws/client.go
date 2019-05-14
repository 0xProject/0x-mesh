package ws

import (
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
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

// AddOrder adds the order to the 0x Mesh node and broadcasts it throughout the
// 0x Mesh network.
func (c *Client) AddOrder(order *zeroex.SignedOrder) (common.Hash, error) {
	var orderHashHex string
	if err := c.rpcClient.Call(&orderHashHex, "mesh_addOrder", order); err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(orderHashHex), nil
}
