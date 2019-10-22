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

// ERC20AssetDataID is the assetDataId for ERC20 tokens
const ERC20AssetDataID = "f47261b0"

// ERC721AssetDataID is the assetDataId for ERC721 tokens
const ERC721AssetDataID = "02571792"

// ERC1155AssetDataID is the assetDataId for ERC721 tokens
const ERC1155AssetDataID = "a7cb5fb7"

// MultiAssetDataID is the assetDataId for multiAsset tokens
const MultiAssetDataID = "94cfcdd7"

const erc20AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"}],\"name\":\"ERC20Token\",\"type\":\"function\"}]"
const erc721AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721Token\",\"type\":\"function\"}]"
const erc1155AssetDataAbi = "[{\"constant\":false,\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"ERC1155Assets\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
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

// ERC1155AssetData represents an ERC1155 assetData
type ERC1155AssetData struct {
	Address      common.Address
	Ids          []*big.Int
	Values       []*big.Int
	CallbackData []byte
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
		log.WithField("erc721AssetDataAbi", erc721AssetDataAbi).Panic("erc721AssetDataAbi should be ABI parsable")
	}
	erc1155AssetDataABI, err := abi.JSON(strings.NewReader(erc1155AssetDataAbi))
	if err != nil {
		log.WithField("erc1155AssetDataAbi", erc1155AssetDataAbi).Panic("erc1155AssetDataAbi should be ABI parsable")
	}
	multiAssetDataABI, err := abi.JSON(strings.NewReader(multiAssetDataAbi))
	if err != nil {
		log.WithField("erc20AssetDataAbi", erc20AssetDataAbi).Panic("erc20AssetDataAbi should be ABI parsable")
	}
	idToAssetDataInfo := map[string]assetDataInfo{
		ERC20AssetDataID: assetDataInfo{
			name: "ERC20Token",
			abi:  erc20AssetDataABI,
		},
		ERC721AssetDataID: assetDataInfo{
			name: "ERC721Token",
			abi:  erc721AssetDataABI,
		},
		ERC1155AssetDataID: assetDataInfo{
			name: "ERC1155Assets",
			abi:  erc1155AssetDataABI,
		},
		MultiAssetDataID: assetDataInfo{
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
	if len(assetData) < 4 {
		return "", errors.New("assetData must be at least 4 bytes long")
	}
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
	if len(assetData) < 4 {
		return errors.New("assetData must be at least 4 bytes long")
	}
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
