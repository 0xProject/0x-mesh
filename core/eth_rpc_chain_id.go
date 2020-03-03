package core

import (
	"context"
	"strconv"
	"strings"
)

func (app *App) getEthRPCChainID(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, ethereumRPCRequestTimeout)
	defer cancel()

	var chainIDRaw string

	err := app.ethRPCClient.CallContext(ctx, &chainIDRaw, "eth_chainId")
	if err != nil {
		return -1, err
	}

	// Value in RPC response is a string in hexadecimal form e.g. "0x539"
	// Need to remove the "0x" prefix before parsing
	rpcChainID, err := strconv.ParseInt(strings.Replace(chainIDRaw, "0x", "", -1), 16, 64)
	if err != nil {
		return -1, err
	}

	return rpcChainID, nil
}
