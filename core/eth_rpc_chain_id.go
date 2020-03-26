package core

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

func (app *App) getEthRPCChainID(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, ethereumRPCRequestTimeout)
	defer cancel()

	var chainIDRaw string
	err := app.ethRPCClient.CallContext(ctx, &chainIDRaw, "eth_chainId")
	if err != nil {
		return nil, err
	}

	rpcChainID, ok := math.ParseBig256(chainIDRaw)
	if !ok {
		return nil, errors.New("Failed to parse big.Int value from hex-encoded chainID returned from eth_chainId")
	}

	return rpcChainID, nil
}
