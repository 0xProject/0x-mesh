// +build js,wasm

package main

import (
	"syscall/js"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/stretchr/testify/require"
)

func TestConvertConfig(t *testing.T) {
	config, err := convertConfig(js.ValueOf(map[string]interface{}{
		"verbosity":                        js.Undefined(),
		"ethereumRPCURL":                   js.ValueOf("http://localhost:8545"),
		"ethereumChainID":                  js.ValueOf(1337),
		"useBootstrapList":                 js.Undefined(),
		"bootstrapList":                    js.Undefined(),
		"blockPollingIntervalSeconds":      js.Undefined(),
		"ethereumRPCMaxContentLength":      js.Undefined(),
		"ethereumRPCMaxRequestsPer24HrUTC": js.Undefined(),
		"ethereumRPCMaxRequestsPerSecond":  js.Undefined(),
		"enableEthereumRPCRateLimiting":    js.Undefined(),
		"customContractAddresses":          js.Undefined(),
		"maxOrdersInStorage":               js.Undefined(),
	}))
	require.NoError(t, err)
	require.Equal(t, config, core.Config{
		EthereumRPCURL:                   "http://localhost:8545",
		EthereumChainID:                  1337,
		Verbosity:                        2,
		DataDir:                          "0x-mesh",
		P2PTCPPort:                       0,
		P2PWebSocketsPort:                0,
		UseBootstrapList:                 true,
		BlockPollingInterval:             5 * time.Second,
		EthereumRPCMaxContentLength:      524288,
		EthereumRPCMaxRequestsPer24HrUTC: 100000,
		EthereumRPCMaxRequestsPerSecond:  30,
		EnableEthereumRPCRateLimiting:    true,
		MaxOrdersInStorage:               100000,
	})
}
