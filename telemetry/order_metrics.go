package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ValidOrdersAccepted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_valid_orders_seen",
		Help: "The total number of valid orders mesh has accepted via JSON RPC",
	})

	InvalidOrdersRejected = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_invalid_orders_seen",
		Help: "The total number of invalid orders mesh has rejected via JSON RPC",
	})
)
