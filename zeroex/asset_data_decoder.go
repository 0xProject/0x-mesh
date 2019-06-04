package zeroex

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

const erc20AssetDataId = "f47261b0"
const erc721AssetDataId = "02571792"
const multiAssetDataId = "94cfcdd7"

const erc20AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"}],\"name\":\"ERC20Token\",\"type\":\"function\"}]"
const erc721AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721Token\",\"type\":\"function\"}]"
const multiAssetDataAbi = "[{\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"MultiAsset\",\"type\":\"function\"}]"

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

// AssetDataDecoder decodes 0x order asset data
type AssetDataDecoder struct {
	idToAssetDataInfo map[string]assetDataInfo
}

// NewAssetDataDecoder instantiates a new asset data decoder
func NewAssetDataDecoder() *AssetDataDecoder {
	erc20AssetDataABI, err := abi.JSON(strings.NewReader(erc20AssetDataAbi))
	if err != nil {
		log.WithField("erc20AssetDataAbi", erc20AssetDataAbi).Panic("erc20AssetDataAbi should be ABI parsable")
	}
	erc721AssetDataABI, err := abi.JSON(strings.NewReader(erc721AssetDataAbi))
	if err != nil {
		log.WithField("erc20AssetDataAbi", erc20AssetDataAbi).Panic("erc20AssetDataAbi should be ABI parsable")
	}
	multiAssetDataABI, err := abi.JSON(strings.NewReader(multiAssetDataAbi))
	if err != nil {
		log.WithField("erc20AssetDataAbi", erc20AssetDataAbi).Panic("erc20AssetDataAbi should be ABI parsable")
	}
	idToAssetDataInfo := map[string]assetDataInfo{
		erc20AssetDataId: assetDataInfo{
			name: "ERC20Token",
			abi:  erc20AssetDataABI,
		},
		erc721AssetDataId: assetDataInfo{
			name: "ERC721Token",
			abi:  erc721AssetDataABI,
		},
		multiAssetDataId: assetDataInfo{
			name: "MultiAsset",
			abi:  multiAssetDataABI,
		},
	}
	decoder := &AssetDataDecoder{
		idToAssetDataInfo: idToAssetDataInfo,
	}
	return decoder
}

// GetName returns the name of the assetData type
func (a *AssetDataDecoder) GetName(assetData []byte) (string, error) {
	id := assetData[:4]
	idHex := common.Bytes2Hex(id)
	info, ok := a.idToAssetDataInfo[idHex]
	if !ok {
		return "", errors.New(fmt.Sprintf("Unrecognized assetData with prefix: %s", idHex))
	}
	return info.name, nil
}

// Decode decodes an encoded asset data into it's sub-components
func (a *AssetDataDecoder) Decode(assetData []byte, decodedAssetData interface{}) error {
	id := assetData[:4]
	idHex := common.Bytes2Hex(id)
	info, ok := a.idToAssetDataInfo[idHex]
	if !ok {
		return errors.New(fmt.Sprintf("Unrecognized assetData with prefix: %s", idHex))
	}

	err := info.abi.Methods[info.name].Inputs.Unpack(decodedAssetData, assetData[4:])
	if err != nil {
		return err
	}
	return nil
}
