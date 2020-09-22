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

// StaticCallAssetDataID is the assetDataId for staticcalls
const StaticCallAssetDataID = "c339d10a"

// CheckGasDefaultID is the function selector for the `checkGas` function that does
// not accept a gasPrice.
const CheckGasPriceDefaultID = "d728f5b7"

// CheckGasID is the function selector for the `checkGas` function that accepts a gasPrice.
const CheckGasPriceID = "da5b166a"

// MultiAssetDataID is the assetDataId for multiAsset tokens
const MultiAssetDataID = "94cfcdd7"

// ERC20BridgeAssetDataID is the assetDataId for ERC20Bridge assets
const ERC20BridgeAssetDataID = "dc1600f3"

const erc20AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"}],\"name\":\"ERC20Token\",\"type\":\"function\"}]"
const erc721AssetDataAbi = "[{\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721Token\",\"type\":\"function\"}]"
const erc1155AssetDataAbi = "[{\"constant\":false,\"inputs\":[{\"name\":\"address\",\"type\":\"address\"},{\"name\":\"ids\",\"type\":\"uint256[]\"},{\"name\":\"values\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"ERC1155Assets\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
const staticCallAssetDataAbi = "[{\"inputs\":[{\"name\":\"staticCallTargetAddress\",\"type\":\"address\"},{\"name\":\"staticCallData\",\"type\":\"bytes\"},{\"name\":\"expectedReturnHashData\", \"type\":\"bytes32\"}],\"name\":\"StaticCall\",\"type\":\"function\"}]"
const checkGasPriceDefaultStaticCallDataAbi = "[{\"inputs\":[],\"name\":\"checkGasPrice\",\"type\":\"function\"}]"
const checkGasPriceStaticCallDataAbi = "[{\"inputs\":[{\"name\":\"maxGasPrice\",\"type\":\"uint256\"}],\"name\":\"checkGasPrice\",\"type\":\"function\"}]"
const multiAssetDataAbi = "[{\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"MultiAsset\",\"type\":\"function\"}]"
const erc20BridgeAssetDataAbi = "[{\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"name\":\"ERC20Bridge\",\"type\":\"function\"}]"

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

// ERC20BridgeAssetData represents an ERC20 Bridge assetData
type ERC20BridgeAssetData struct {
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}

type StaticCallAssetData struct {
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnHashData  [32]byte
}

type CheckGasPriceStaticCallData struct {
	MaxGasPrice *big.Int
}

// MultiAssetData represents a MultiAssetData
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
	staticCallAssetDataABI, err := abi.JSON(strings.NewReader(staticCallAssetDataAbi))
	if err != nil {
		log.WithField("staticCallAssetDataAbi", staticCallAssetDataAbi).Panic("staticCallAssetDataAbi should be ABI parsable")
	}
	checkGasPriceDefaultStaticCallDataABI, err := abi.JSON(strings.NewReader(checkGasPriceDefaultStaticCallDataAbi))
	if err != nil {
		log.WithField("checkGasPriceDefaultStaticCallDataAbi", checkGasPriceDefaultStaticCallDataAbi).Panic("checkGasPriceDefaultStaticCallDataAbi should be ABI parsable")
	}
	checkGasPriceStaticCallDataABI, err := abi.JSON(strings.NewReader(checkGasPriceStaticCallDataAbi))
	if err != nil {
		log.WithField("checkGasPriceStaticCallDataAbi", checkGasPriceStaticCallDataAbi).Panic("checkGasStaticCallDataAbi should be ABI parsable")
	}
	multiAssetDataABI, err := abi.JSON(strings.NewReader(multiAssetDataAbi))
	if err != nil {
		log.WithField("multiAssetDataAbi", multiAssetDataAbi).Panic("multiAssetDataAbi should be ABI parsable")
	}
	erc20BridgeAssetDataABI, err := abi.JSON(strings.NewReader(erc20BridgeAssetDataAbi))
	if err != nil {
		log.WithField("erc20BridgeAssetDataABI", erc20BridgeAssetDataAbi).Panic("erc20BridgeAssetDataABI should be ABI parsable")
	}
	idToAssetDataInfo := map[string]assetDataInfo{
		ERC20AssetDataID: {
			name: "ERC20Token",
			abi:  erc20AssetDataABI,
		},
		ERC721AssetDataID: {
			name: "ERC721Token",
			abi:  erc721AssetDataABI,
		},
		ERC1155AssetDataID: {
			name: "ERC1155Assets",
			abi:  erc1155AssetDataABI,
		},
		StaticCallAssetDataID: {
			name: "StaticCall",
			abi:  staticCallAssetDataABI,
		},
		CheckGasPriceDefaultID: {
			name: "checkGasPrice",
			abi:  checkGasPriceDefaultStaticCallDataABI,
		},
		CheckGasPriceID: {
			name: "checkGasPrice",
			abi:  checkGasPriceStaticCallDataABI,
		},
		MultiAssetDataID: {
			name: "MultiAsset",
			abi:  multiAssetDataABI,
		},
		ERC20BridgeAssetDataID: {
			name: "ERC20Bridge",
			abi:  erc20BridgeAssetDataABI,
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
		return "", fmt.Errorf("Unrecognized assetData with prefix: %s", idHex)
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
		return fmt.Errorf("Unrecognized assetData with prefix: %s", idHex)
	}

	// This is necessary to prevent a nil pointer exception for ABIs with no inputs
	if len(info.abi.Methods[info.name].Inputs.NonIndexed()) == 0 {
		return nil
	}
	err := info.abi.Methods[info.name].Inputs.Unpack(decodedAssetData, assetData[4:])
	if err != nil {
		return err
	}

	return nil
}
