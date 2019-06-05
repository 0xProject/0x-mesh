package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AddOrdersRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_jrpc_request_add_orders_total",
		Help: "The total number of mesh_addOrders JRPC requests",
	})

	AddPeerRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_jrpc_request_add_peer_total",
		Help: "The total number of mesh_addPeer JRPC requests",
	})
)
