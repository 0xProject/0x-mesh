// +build js, wasm

package core

import (
	"errors"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/orderfilter"
)

// convertConfig converts a JavaScript config object into a Config. It also
// adds default values for any that are missing in the JavaScript object.
func ConvertConfig(jsConfig js.Value) (Config, error) {
	if types.IsNullOrUndefined(jsConfig) {
		return Config{}, errors.New("config is required")
	}

	// Default config options. Some might be overridden.
	config := Config{
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
	}

	// Required config options
	if ethereumRPCURL := jsConfig.Get("ethereumRPCURL"); types.IsNullOrUndefined(ethereumRPCURL) || ethereumRPCURL.String() == "" {
		return Config{}, errors.New("ethereumRPCURL is required")
	} else {
		config.EthereumRPCURL = ethereumRPCURL.String()
	}
	if ethereumChainID := jsConfig.Get("ethereumChainID"); types.IsNullOrUndefined(ethereumChainID) {
		return Config{}, errors.New("ethereumChainID is required")
	} else {
		config.EthereumChainID = ethereumChainID.Int()
	}

	// Optional config options
	if verbosity := jsConfig.Get("verbosity"); !types.IsNullOrUndefined(verbosity) {
		config.Verbosity = verbosity.Int()
	}
	if useBootstrapList := jsConfig.Get("useBootstrapList"); !types.IsNullOrUndefined(useBootstrapList) {
		config.UseBootstrapList = useBootstrapList.Bool()
	}
	if bootstrapList := jsConfig.Get("bootstrapList"); !types.IsNullOrUndefined(bootstrapList) {
		config.BootstrapList = bootstrapList.String()
	}
	if blockPollingIntervalSeconds := jsConfig.Get("blockPollingIntervalSeconds"); !types.IsNullOrUndefined(blockPollingIntervalSeconds) {
		config.BlockPollingInterval = time.Duration(blockPollingIntervalSeconds.Int()) * time.Second
	}
	if ethereumRPCMaxContentLength := jsConfig.Get("ethereumRPCMaxContentLength"); !types.IsNullOrUndefined(ethereumRPCMaxContentLength) {
		config.EthereumRPCMaxContentLength = ethereumRPCMaxContentLength.Int()
	}
	if ethereumRPCMaxRequestsPer24HrUTC := jsConfig.Get("ethereumRPCMaxRequestsPer24HrUTC"); !types.IsNullOrUndefined(ethereumRPCMaxRequestsPer24HrUTC) {
		config.EthereumRPCMaxRequestsPer24HrUTC = ethereumRPCMaxRequestsPer24HrUTC.Int()
	}
	if ethereumRPCMaxRequestsPerSecond := jsConfig.Get("ethereumRPCMaxRequestsPerSecond"); !types.IsNullOrUndefined(ethereumRPCMaxRequestsPerSecond) {
		config.EthereumRPCMaxRequestsPerSecond = ethereumRPCMaxRequestsPerSecond.Float()
	}
	if enableEthereumRPCRateLimiting := jsConfig.Get("enableEthereumRPCRateLimiting"); !types.IsNullOrUndefined(enableEthereumRPCRateLimiting) {
		config.EnableEthereumRPCRateLimiting = enableEthereumRPCRateLimiting.Bool()
	}
	if customContractAddresses := jsConfig.Get("customContractAddresses"); !types.IsNullOrUndefined(customContractAddresses) {
		config.CustomContractAddresses = customContractAddresses.String()
	}
	if maxOrdersInStorage := jsConfig.Get("maxOrdersInStorage"); !types.IsNullOrUndefined(maxOrdersInStorage) {
		config.MaxOrdersInStorage = maxOrdersInStorage.Int()
	}
	if customOrderFilter := jsConfig.Get("customOrderFilter"); !types.IsNullOrUndefined(customOrderFilter) {
		config.CustomOrderFilter = customOrderFilter.String()
	}

	return config, nil
}
