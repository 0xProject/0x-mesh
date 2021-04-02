package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ProtocolVersionLabel  = "protocol_version"
	ProtocolV3            = "v3"
	ProtocolV4            = "v4"
	ValidationStatusLabel = "validation_status"
	ValidationAccepted    = "accepted"
	ValidationRejected    = "rejected"
	QueryLabel            = "query"
	OrdersyncStatusLabel  = "status"
	OrdersyncSuccess      = "success"
	OrdersyncFailure      = "failure"
)

var (
	OrdersShared = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_orders_shared_total",
		Help: "The total number of shared orders",
	},
		[]string{
			ProtocolVersionLabel,
		})

	OrdersStored = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mesh_orders_total",
		Help: "Current total number of stored orders",
	},
		[]string{
			ProtocolVersionLabel,
		})
	PeersConnected = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mesh_peers_connected_total",
		Help: "Current total number of connected peers",
	})

	OrdersAddedViaGraphQl = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_graphql_orders_added",
		Help: "Total number of orders added / rejected via graphql",
	},
		[]string{
			ProtocolVersionLabel,
			ValidationStatusLabel,
		})

	P2POrdersReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_p2p_orders_received",
		Help: "Number of orders received",
	}, []string{
		ProtocolVersionLabel,
		ValidationStatusLabel,
	})

	GraphqlQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_graphql_queries_total",
		Help: "Total number of GraphQL endpoint queries handled",
	}, []string{
		QueryLabel,
	})

	OrdersyncRequestsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_ordersync_requests_received",
		Help: "Number of ordersync requests received",
	}, []string{
		ProtocolVersionLabel,
	})

	OrdersyncRequestsSent = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mesh_ordersync_requests_sent",
		Help: "Number of ordersync requests sent",
	}, []string{
		ProtocolVersionLabel,
		OrdersyncStatusLabel,
	})
)

func ServeMetrics(ctx context.Context, serveAddr string) error {
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(serveAddr, nil)
}
