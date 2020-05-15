package db

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
)

// Metadata is the database representation of MeshDB instance metadata
type Metadata struct {
	EthereumChainID                   int
	MaxExpirationTime                 *big.Int
	EthRPCRequestsSentInCurrentUTCDay int
	StartOfCurrentUTCDay              time.Time
}

func ParseContractAddressesAndTokenIdsFromAssetData(assetData []byte, contractAddresses ethereum.ContractAddresses) ([]*types.SingleAssetData, error) {
	if len(assetData) == 0 {
		return []*types.SingleAssetData{}, nil
	}
	singleAssetDatas := []*types.SingleAssetData{}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

	assetDataName, err := assetDataDecoder.GetName(assetData)
	if err != nil {
		return nil, err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := &types.SingleAssetData{
			Address: decodedAssetData.Address,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := &types.SingleAssetData{
			Address: decodedAssetData.Address,
			TokenID: decodedAssetData.TokenId,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, id := range decodedAssetData.Ids {
			a := &types.SingleAssetData{
				Address: decodedAssetData.Address,
				TokenID: id,
			}
			singleAssetDatas = append(singleAssetDatas, a)
		}
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		// NOTE(jalextowle): As of right now, none of the supported staticcalls
		// have important information in the StaticCallData. We choose not to add
		// `singleAssetData` because it would not be used.
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			as, err := ParseContractAddressesAndTokenIdsFromAssetData(assetData, contractAddresses)
			if err != nil {
				return nil, err
			}
			singleAssetDatas = append(singleAssetDatas, as...)
		}
	case "ERC20Bridge":
		var decodedAssetData zeroex.ERC20BridgeAssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		tokenAddress := decodedAssetData.TokenAddress
		// TODO(albrow): Update orderwatcher to account for this instead of storing
		// it in the database. This would mean we can remove contractAddresses as an
		// argument and simplify the implementation. Maybe even have the db package
		// handle parsing asset data automatically.
		// HACK(fabio): Despite Chai ERC20Bridge orders encoding the Dai address as
		// the tokenAddress, we actually want to react to the Chai token's contract
		// events, so we actually return it instead.
		if decodedAssetData.BridgeAddress == contractAddresses.ChaiBridge {
			tokenAddress = contractAddresses.ChaiToken
		}
		a := &types.SingleAssetData{
			Address: tokenAddress,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	default:
		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return singleAssetDatas, nil
}
