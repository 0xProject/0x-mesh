package zeroex

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const ERC20_ASSET_DATA_ID = "f47261b0"
const ERC721_ASSET_DATA_ID = "02571792"
const MULTI_ASSET_DATA_ID = "94cfcdd7"

const ERC20_ASSET_DATA_ABI = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"}],\"name\":\"ERC20Token\",\"type\":\"function\"}]"
const ERC721_ASSET_DATA_ABI = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721Token\",\"type\":\"function\"}]"
const MULTI_ASSET_DATA_ABI = "[{\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"MultiAsset\",\"type\":\"function\"}]"

// ERC20AssetData represents an ERC20 assetData
type ERC20AssetData struct {
	Address common.Address
}

// ERC721AssetData represents an ERC721 assetData
type ERC721AssetData struct {
	Address common.Address
	TokenId *big.Int
}

// MultiAssetData represents an MultiAssetData
type MultiAssetData struct {
	Amounts         []*big.Int
	NestedAssetData [][]byte
}

type assetDataInfo struct {
	name string
	abi  abi.ABI
}

type AssetDataDecoder struct {
	idToAssetDataInfo map[string]assetDataInfo
}

func NewAssetDataDecoder() (*AssetDataDecoder, error) {
	erc20AssetDataABI, err := abi.JSON(strings.NewReader(ERC20_ASSET_DATA_ABI))
	if err != nil {
		return nil, err
	}
	erc721AssetDataABI, err := abi.JSON(strings.NewReader(ERC721_ASSET_DATA_ABI))
	if err != nil {
		return nil, err
	}
	multiAssetDataABI, err := abi.JSON(strings.NewReader(MULTI_ASSET_DATA_ABI))
	if err != nil {
		return nil, err
	}
	idToAssetDataInfo := map[string]assetDataInfo{
		ERC20_ASSET_DATA_ID: assetDataInfo{
			name: "ERC20Token",
			abi:  erc20AssetDataABI,
		},
		ERC721_ASSET_DATA_ID: assetDataInfo{
			name: "ERC721Token",
			abi:  erc721AssetDataABI,
		},
		MULTI_ASSET_DATA_ID: assetDataInfo{
			name: "MultiAsset",
			abi:  multiAssetDataABI,
		},
	}
	decoder := &AssetDataDecoder{
		idToAssetDataInfo: idToAssetDataInfo,
	}
	return decoder, nil
}

func (d *AssetDataDecoder) Decode(assetData []byte) (interface{}, error) {
	id := assetData[:4]
	idHex := common.Bytes2Hex(id)
	info, ok := d.idToAssetDataInfo[idHex]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unrecognized assetData with prefix: %s", idHex))
	}

	switch info.name {
	case "ERC20Token":
		var decodedAssetData ERC20AssetData
		err := info.abi.Methods[info.name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil
	case "ERC721Token":
		var decodedAssetData ERC721AssetData
		err := info.abi.Methods[info.name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil

	case "MultiAsset":
		var decodedAssetData MultiAssetData
		err := info.abi.Methods[info.name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unsupported AssetData with name %s", info.name))
	}
}
