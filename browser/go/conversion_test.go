// +build js,wasm

package main

import (
	"encoding/json"
	"syscall/js"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/stretchr/testify/require"
)

func TestConvertConfigEmpty(t *testing.T) {
	data, err := json.Marshal(struct {
		Exchange     string `json:"exchange"`
		DevUtils     string `json:"devUtils"`
		ERC20Proxy   string `json:"erc20Proxy"`
		ERC721Proxy  string `json:"erc721Proxy"`
		ERC1155Proxy string `json:"erc1155Proxy"`
	}{
		Exchange:     "0x48bacb9266a570d521063ef5dd96e61686dbe788",
		DevUtils:     "0x38ef19fdf8e8415f18c307ed71967e19aac28ba1",
		ERC20Proxy:   "0x1dc4c1cefef38a777b15aa20260a54e584b16c48",
		ERC721Proxy:  "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401",
		ERC1155Proxy: "0x2d7022f5b17d2f8b695918fb48fa1089c9f85401",
	})
	require.NoError(t, err)
	for _, testCase := range []struct {
		jsConfig       map[string]interface{}
		expectedConfig core.Config
	}{
		{
			jsConfig: map[string]interface{}{
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
			},
			expectedConfig: core.Config{
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
			},
		},
		{
			jsConfig: map[string]interface{}{
				"verbosity":        js.ValueOf(6),
				"ethereumRPCURL":   js.ValueOf("http://localhost:8545"),
				"ethereumChainID":  js.ValueOf(1337),
				"useBootstrapList": js.ValueOf(true),
				"bootstrapList": js.ValueOf(
					"/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF," +
						"/ip4/3.214.190.68/tcp/60559/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumG",
				),
				"blockPollingIntervalSeconds":      js.ValueOf(50),
				"ethereumRPCMaxContentLength":      js.ValueOf(10000),
				"ethereumRPCMaxRequestsPer24HrUTC": js.ValueOf(10000),
				"ethereumRPCMaxRequestsPerSecond":  js.ValueOf(50),
				"enableEthereumRPCRateLimiting":    js.ValueOf(true),
				"customContractAddresses":          string(data),
				"maxOrdersInStorage":               js.ValueOf(10000),
			},
			expectedConfig: core.Config{
				BootstrapList: "/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF," +
					"/ip4/3.214.190.68/tcp/60559/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumG",
				CustomContractAddresses:          string(data),
				EthereumRPCURL:                   "http://localhost:8545",
				EthereumChainID:                  1337,
				Verbosity:                        6,
				DataDir:                          "0x-mesh",
				P2PTCPPort:                       0,
				P2PWebSocketsPort:                0,
				UseBootstrapList:                 true,
				BlockPollingInterval:             50 * time.Second,
				EthereumRPCMaxContentLength:      10000,
				EthereumRPCMaxRequestsPer24HrUTC: 10000,
				EthereumRPCMaxRequestsPerSecond:  50,
				EnableEthereumRPCRateLimiting:    true,
				MaxOrdersInStorage:               10000,
			},
		},
	} {
		config, err := convertConfig(js.ValueOf(testCase.jsConfig))
		require.NoError(t, err)
		require.Equal(t, testCase.expectedConfig, config)
	}
}
