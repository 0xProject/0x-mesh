package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	P2POrdersAlreadyStored = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_p2p_valid_orders_already_stored",
		Help: "The total number of valid orders mesh has already stored and rejected via p2p",
	})
	P2PValidOrdersSeen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_p2p_valid_orders_seen",
		Help: "The total number of valid orders mesh has seen",
	})

	P2PInvalidOrdersSeen = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mesh_p2p_invalid_orders_seen",
		Help: "The total number of invalid orders mesh has seen",
	})
)
