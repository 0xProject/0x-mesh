// +build js,wasm

package types

import (
	"encoding/json"
	"syscall/js"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
)

func (r GetOrdersResponse) JSValue() js.Value {
	// TODO(albrow): Optimize this. Remove other uses of the JSON
	// encoding/decoding hack.
	encodedResponse, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	responseJS := js.Global().Get("JSON").Call("parse", string(encodedResponse))
	return responseJS
}

func (l LatestBlock) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"number": l.Number,
		"hash":   l.Hash.String(),
	})
}

func (s Stats) JSValue() js.Value {
	secondaryRendezvous := make([]interface{}, len(s.SecondaryRendezvous))
	for i, rendezvousPoint := range s.SecondaryRendezvous {
		secondaryRendezvous[i] = rendezvousPoint
	}
	return js.ValueOf(map[string]interface{}{
		"version":                           s.Version,
		"pubSubTopic":                       s.PubSubTopic,
		"rendezvous":                        s.Rendezvous,
		"secondaryRendezvous":               secondaryRendezvous,
		"peerID":                            s.PeerID,
		"ethereumChainID":                   s.EthereumChainID,
		"latestBlock":                       s.LatestBlock.JSValue(),
		"numPeers":                          s.NumPeers,
		"numOrders":                         s.NumOrders,
		"numOrdersIncludingRemoved":         s.NumOrdersIncludingRemoved,
		"numPinnedOrders":                   s.NumPinnedOrders,
		"maxExpirationTime":                 s.MaxExpirationTime,
		"startOfCurrentUTCDay":              s.StartOfCurrentUTCDay.String(),
		"ethRPCRequestsSentInCurrentUTCDay": s.EthRPCRequestsSentInCurrentUTCDay,
		"ethRPCRateLimitExpiredRequests":    s.EthRPCRateLimitExpiredRequests,
	})
}

func (o OrderWithMetadata) JSValue() js.Value {
	value, _ := jsutil.InefficientlyConvertToJS(o)
	return value
}
