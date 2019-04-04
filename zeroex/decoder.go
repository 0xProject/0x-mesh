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

type AssetDataDecoder struct {
	idToAbi  map[string]abi.ABI
	idToName map[string]string
}

func NewAssetDataDecoder() (*AssetDataDecoder, error) {
	idToName := map[string]string{
		ERC20_ASSET_DATA_ID:  "ERC20Token",
		ERC721_ASSET_DATA_ID: "ERC721Token",
		MULTI_ASSET_DATA_ID:  "MultiAsset",
	}
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
	idToAbi := map[string]abi.ABI{
		ERC20_ASSET_DATA_ID:  erc20AssetDataABI,
		ERC721_ASSET_DATA_ID: erc721AssetDataABI,
		MULTI_ASSET_DATA_ID:  multiAssetDataABI,
	}
	decoder := &AssetDataDecoder{
		idToAbi:  idToAbi,
		idToName: idToName,
	}
	return decoder, nil
}

func (d *AssetDataDecoder) Decode(assetData []byte) (interface{}, error) {
	id := assetData[:4]
	idHex := common.Bytes2Hex(id)
	abi, ok := d.idToAbi[idHex]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unrecognized assetData with prefix: %s", idHex))
	}
	name, ok := d.idToName[idHex]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Found assetData prefix without corresponding name: %s", idHex))
	}

	switch name {
	case "ERC20Token":
		var decodedAssetData ERC20AssetData
		err := abi.Methods[name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil
	case "ERC721Token":
		var decodedAssetData ERC721AssetData
		err := abi.Methods[name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil

	case "MultiAsset":
		var decodedAssetData MultiAssetData
		err := abi.Methods[name].Inputs.Unpack(&decodedAssetData, assetData[4:])
		if err != nil {
			return nil, err
		}
		return decodedAssetData, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unsupported AssetData with name %s", name))
	}
}
