// +build js, wasm

// NOTE(jalextowle): This file contains utilities used by browser based mesh nodes
// that need to be tested and would cause cyclic dependencies if they were moved
// to jsutil.
package browserutil

import (
	"errors"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/jsutil"
	"github.com/0xProject/0x-mesh/packages/mesh-browser/go/providerwrapper"
)

// ConvertConfig converts a JavaScript config object into a core.Config. It also
// adds default values for any that are missing in the JavaScript object.
func ConvertConfig(jsConfig js.Value) (core.Config, error) {
	if jsutil.IsNullOrUndefined(jsConfig) {
		return core.Config{}, errors.New("config is required")
	}

	// Default config options. Some might be overridden.
	config := core.Config{
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
		CustomOrderFilter:                orderfilter.DefaultCustomOrderSchema,
		MaxBytesPerSecond:                5242880, // 5 MiB
	}

	// Required config options
	if ethereumChainID := jsConfig.Get("ethereumChainID"); jsutil.IsNullOrUndefined(ethereumChainID) {
		return core.Config{}, errors.New("ethereumChainID is required")
	} else {
		config.EthereumChainID = ethereumChainID.Int()
	}

	// Optional config options
	if verbosity := jsConfig.Get("verbosity"); !jsutil.IsNullOrUndefined(verbosity) {
		config.Verbosity = verbosity.Int()
	}
	if useBootstrapList := jsConfig.Get("useBootstrapList"); !jsutil.IsNullOrUndefined(useBootstrapList) {
		config.UseBootstrapList = useBootstrapList.Bool()
	}
	if bootstrapList := jsConfig.Get("bootstrapList"); !jsutil.IsNullOrUndefined(bootstrapList) {
		config.BootstrapList = bootstrapList.String()
	}
	if blockPollingIntervalSeconds := jsConfig.Get("blockPollingIntervalSeconds"); !jsutil.IsNullOrUndefined(blockPollingIntervalSeconds) {
		config.BlockPollingInterval = time.Duration(blockPollingIntervalSeconds.Int()) * time.Second
	}
	if ethereumRPCMaxContentLength := jsConfig.Get("ethereumRPCMaxContentLength"); !jsutil.IsNullOrUndefined(ethereumRPCMaxContentLength) {
		config.EthereumRPCMaxContentLength = ethereumRPCMaxContentLength.Int()
	}
	if ethereumRPCMaxRequestsPer24HrUTC := jsConfig.Get("ethereumRPCMaxRequestsPer24HrUTC"); !jsutil.IsNullOrUndefined(ethereumRPCMaxRequestsPer24HrUTC) {
		config.EthereumRPCMaxRequestsPer24HrUTC = ethereumRPCMaxRequestsPer24HrUTC.Int()
	}
	if ethereumRPCMaxRequestsPerSecond := jsConfig.Get("ethereumRPCMaxRequestsPerSecond"); !jsutil.IsNullOrUndefined(ethereumRPCMaxRequestsPerSecond) {
		config.EthereumRPCMaxRequestsPerSecond = ethereumRPCMaxRequestsPerSecond.Float()
	}
	if enableEthereumRPCRateLimiting := jsConfig.Get("enableEthereumRPCRateLimiting"); !jsutil.IsNullOrUndefined(enableEthereumRPCRateLimiting) {
		config.EnableEthereumRPCRateLimiting = enableEthereumRPCRateLimiting.Bool()
	}
	if customContractAddresses := jsConfig.Get("customContractAddresses"); !jsutil.IsNullOrUndefined(customContractAddresses) {
		config.CustomContractAddresses = customContractAddresses.String()
	}
	if maxOrdersInStorage := jsConfig.Get("maxOrdersInStorage"); !jsutil.IsNullOrUndefined(maxOrdersInStorage) {
		config.MaxOrdersInStorage = maxOrdersInStorage.Int()
	}
	if customOrderFilter := jsConfig.Get("customOrderFilter"); !jsutil.IsNullOrUndefined(customOrderFilter) {
		config.CustomOrderFilter = customOrderFilter.String()
	}
	if ethereumRPCURL := jsConfig.Get("ethereumRPCURL"); !jsutil.IsNullOrUndefined(ethereumRPCURL) && ethereumRPCURL.String() != "" {
		config.EthereumRPCURL = ethereumRPCURL.String()
	}
	if web3Provider := jsConfig.Get("web3Provider"); !jsutil.IsNullOrUndefined(web3Provider) {
		config.EthereumRPCClient = providerwrapper.NewRPCClient(web3Provider)
	}
	if maxBytesPerSecond := jsConfig.Get("maxBytesPerSecond"); !jsutil.IsNullOrUndefined(maxBytesPerSecond) {
		config.MaxBytesPerSecond = maxBytesPerSecond.Float()
	}

	return config, nil
}
