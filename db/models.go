package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

// Order is the database representation a 0x order along with some relevant metadata
type Order struct {
	Hash                  common.Hash    `db:"hash"`
	ChainID               *Uint256       `db:"chainID"`
	ExchangeAddress       common.Address `db:"exchangeAddress"`
	MakerAddress          common.Address `db:"makerAddress"`
	MakerAssetData        []byte         `db:"makerAssetData"`
	MakerFeeAssetData     []byte         `db:"makerFeeAssetData"`
	MakerAssetAmount      *Uint256       `db:"makerAssetAmount"`
	MakerFee              *Uint256       `db:"makerFee"`
	TakerAddress          common.Address `db:"takerAddress"`
	TakerAssetData        []byte         `db:"takerAssetData"`
	TakerFeeAssetData     []byte         `db:"takerFeeAssetData"`
	TakerAssetAmount      *Uint256       `db:"takerAssetAmount"`
	TakerFee              *Uint256       `db:"takerFee"`
	SenderAddress         common.Address `db:"senderAddress"`
	FeeRecipientAddress   common.Address `db:"feeRecipientAddress"`
	ExpirationTimeSeconds *Uint256       `db:"expirationTimeSeconds"`
	Salt                  *Uint256       `db:"salt"`
	Signature             []byte         `db:"signature"`
	// When was this order last validated
	LastUpdated time.Time `db:"lastUpdated"`
	// How much of this order can still be filled
	FillableTakerAssetAmount *Uint256 `db:"fillableTakerAssetAmount"`
	// Was this order flagged for removal? Due to the possibility of block-reorgs, instead
	// of immediately removing an order when FillableTakerAssetAmount becomes 0, we instead
	// flag it for removal. After this order isn't updated for X time and has IsRemoved = true,
	// the order can be permanently deleted.
	IsRemoved bool `db:"isRemoved"`
	// IsPinned indicates whether or not the order is pinned. Pinned orders are
	// not removed from the database unless they become unfillable.
	IsPinned bool `db:"isPinned"`
	// JSON-encoded list of assetdatas contained in MakerAssetData. For non-MAP
	// orders, the list contains only one element which is equal to MakerAssetData.
	// For MAP orders, it contains each component assetdata.
	ParsedMakerAssetData *ParsedAssetData `db:"parsedMakerAssetData"`
	// Same as ParsedMakerAssetData but for MakerFeeAssetData instead of MakerAssetData.
	ParsedMakerFeeAssetData *ParsedAssetData `db:"parsedMakerFeeAssetData"`
}

// Metadata is the database representation of MeshDB instance metadata
type Metadata struct {
	EthereumChainID                   int
	MaxExpirationTime                 *big.Int
	EthRPCRequestsSentInCurrentUTCDay int
	StartOfCurrentUTCDay              time.Time
}

// MiniHeader is a representation of a succinct Ethereum block headers
type MiniHeader struct {
	Hash      common.Hash `db:"hash"`
	Parent    common.Hash `db:"parent"`
	Number    *Uint256    `db:"number"`
	Timestamp time.Time   `db:"timestamp"`
	Logs      *EventLogs  `db:"logs"`
}

type Uint256 struct {
	*big.Int
}

func NewUint256(v *big.Int) *Uint256 {
	return &Uint256{
		Int: v,
	}
}

func (u *Uint256) Value() (driver.Value, error) {
	if u == nil || u.Int == nil {
		return nil, nil
	}
	return u.String(), nil
}

func (u *Uint256) Scan(value interface{}) error {
	if value == nil {
		u = nil
		return nil
	}
	switch v := value.(type) {
	case int64:
		u.Int = big.NewInt(v)
	case string:
		parsed, ok := math.ParseBig256(v)
		if !ok {
			return fmt.Errorf("could not scan string value %q into Uint256", v)
		}
		u.Int = parsed
	default:
		return fmt.Errorf("could not scan type %T into Uint256", value)
	}

	return nil
}

func (u *Uint256) MarshalJSON() ([]byte, error) {
	if u == nil || u.Int == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(u.Int.String())
}

func (u *Uint256) UnmarshalJSON(data []byte) error {
	unqouted, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON data into Uint256: %s", string(data))
	}
	bigInt, ok := math.ParseBig256(unqouted)
	if !ok {
		return fmt.Errorf("could not unmarshal JSON data into Uint256: %s", string(data))
	}
	u.Int = bigInt
	return nil
}

type EventLogs struct {
	Logs []types.Log
}

func NewEventLogs(logs []types.Log) *EventLogs {
	eventLogs := EventLogs{Logs: logs}
	return &eventLogs
}

func (e *EventLogs) Value() (driver.Value, error) {
	if e == nil {
		return nil, nil
	}
	logsJSON, err := json.Marshal(e.Logs)
	if err != nil {
		return nil, err
	}
	return logsJSON, err
}

func (e *EventLogs) Scan(value interface{}) error {
	if value == nil {
		e = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &e.Logs)
	case string:
		return json.Unmarshal([]byte(v), &e.Logs)
	default:
		return fmt.Errorf("could not scan type %T into EventLogs", value)
	}
}

type ParsedAssetData []SingleAssetData

type SingleAssetData struct {
	Address common.Address `json:"address"`
	TokenID *Uint256       `json:"tokenID"`
}

func (s *ParsedAssetData) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *ParsedAssetData) Scan(value interface{}) error {
	if value == nil {
		s = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, s)
	case string:
		return json.Unmarshal([]byte(v), s)
	default:
		return fmt.Errorf("could not scan type %T into EventLogs", value)
	}
}

func ParseContractAddressesAndTokenIdsFromAssetData(assetData []byte, contractAddresses ethereum.ContractAddresses) (ParsedAssetData, error) {
	singleAssetDatas := []SingleAssetData{}
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
		a := SingleAssetData{
			Address: decodedAssetData.Address,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		a := SingleAssetData{
			Address: decodedAssetData.Address,
			TokenID: NewUint256(decodedAssetData.TokenId),
		}
		singleAssetDatas = append(singleAssetDatas, a)
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return nil, err
		}
		for _, id := range decodedAssetData.Ids {
			a := SingleAssetData{
				Address: decodedAssetData.Address,
				TokenID: NewUint256(id),
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
		// HACK(fabio): Despite Chai ERC20Bridge orders encoding the Dai address as
		// the tokenAddress, we actually want to react to the Chai token's contract
		// events, so we actually return it instead.
		if decodedAssetData.BridgeAddress == contractAddresses.ChaiBridge {
			tokenAddress = contractAddresses.ChaiToken
		}
		a := SingleAssetData{
			Address: tokenAddress,
		}
		singleAssetDatas = append(singleAssetDatas, a)
	default:
		return nil, fmt.Errorf("unrecognized assetData type name found: %s", assetDataName)
	}
	return singleAssetDatas, nil
}
