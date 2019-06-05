package telemetry

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheusServer(port int) error {
	http.Handle("/metrics", promhttp.Handler())
	listenString := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(listenString, nil)
}
