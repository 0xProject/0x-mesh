package gqltypes

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type LatestBlock struct {
	Number string
	Hash   string
}

func LatestBlockFromCommonType(latestBlock types.LatestBlock) LatestBlock {
	return LatestBlock{
		// TODO(albrow): Change this.
		Number: fmt.Sprint(latestBlock.Number),
		Hash:   latestBlock.Hash.Hex(),
	}
}

type Stats struct {
	Version                           string
	PubSubTopic                       string
	Rendezvous                        string
	PeerID                            string
	EthereumChainID                   int32
	LatestBlock                       LatestBlock
	NumPeers                          int32
	NumOrders                         int32
	NumOrdersIncludingRemoved         int32
	StartOfCurrentUTCDay              string
	EthRPCRequestsSentInCurrentUTCDay int32
	EthRPCRateLimitExpiredRequests    int32
	MaxExpirationTime                 string
}

type SignedOrder struct {
	ChainID               string
	ExchangeAddress       string
	MakerAddress          string
	MakerAssetData        string
	MakerFeeAssetData     string
	MakerAssetAmount      string
	MakerFee              string
	TakerAddress          string
	TakerAssetData        string
	TakerFeeAssetData     string
	TakerAssetAmount      string
	TakerFee              string
	SenderAddress         string
	FeeRecipientAddress   string
	ExpirationTimeSeconds string
	Salt                  string
	Signature             string
}

// TODO(albrow): Can we use the SignedOrder type instead?
type NewOrder struct {
	ChainID               string
	ExchangeAddress       string
	MakerAddress          string
	MakerAssetData        string
	MakerFeeAssetData     string
	MakerAssetAmount      string
	MakerFee              string
	TakerAddress          string
	TakerAssetData        string
	TakerFeeAssetData     string
	TakerAssetAmount      string
	TakerFee              string
	SenderAddress         string
	FeeRecipientAddress   string
	ExpirationTimeSeconds string
	Salt                  string
	Signature             string
}

type OrderWithMetadata struct {
	Hash                     string
	ChainID                  string
	ExchangeAddress          string
	MakerAddress             string
	MakerAssetData           string
	MakerFeeAssetData        string
	MakerAssetAmount         string
	MakerFee                 string
	TakerAddress             string
	TakerAssetData           string
	TakerFeeAssetData        string
	TakerAssetAmount         string
	TakerFee                 string
	SenderAddress            string
	FeeRecipientAddress      string
	ExpirationTimeSeconds    string
	Salt                     string
	Signature                string
	FillableTakerAssetAmount string
	LastUpdated              string
}

type AddOrdersResults struct {
	Accepted []AcceptedOrderResult
	Rejected []RejectedOrderResult
}

type AcceptedOrderResult struct {
	Order *OrderWithMetadata
	IsNew bool
}

type RejectedOrderResult struct {
	Hash    *string
	Order   SignedOrder
	Code    string
	Message string
}

type OrderEvent struct {
	Order          OrderWithMetadata
	EndState       string
	Timestamp      string
	ContractEvents []ContractEvent
}

type ContractEvent struct {
	BlockHash  string
	TxHash     string
	TxIndex    int32
	LogIndex   int32
	Isremoved  bool
	Address    string
	Kind       string
	Parameters ContractEventParams
}

// ContractEventParams corresponds to the ContractEventParams scalar type in the GraphQL Schema.
// We need this custom type because GraphQL doesn't ship with an "any" type.
type ContractEventParams struct {
	value interface{}
}

func (c *ContractEventParams) ImplementsGraphQLType(name string) bool {
	return name == "ContractEventParams"
}

func (c *ContractEventParams) UnmarshalGraphQL(input interface{}) error {
	c.value = input
	return nil
}

func (c *ContractEventParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.value)
}

func StatsFromCommonType(stats *types.Stats) *Stats {
	return &Stats{
		Version:                           stats.Version,
		PubSubTopic:                       stats.PubSubTopic,
		Rendezvous:                        stats.Rendezvous,
		PeerID:                            stats.PeerID,
		EthereumChainID:                   int32(stats.EthereumChainID),
		LatestBlock:                       LatestBlockFromCommonType(stats.LatestBlock),
		NumPeers:                          int32(stats.NumPeers),
		NumOrders:                         int32(stats.NumOrders),
		NumOrdersIncludingRemoved:         int32(stats.NumOrdersIncludingRemoved),
		StartOfCurrentUTCDay:              stats.StartOfCurrentUTCDay.Format(time.RFC3339),
		EthRPCRequestsSentInCurrentUTCDay: int32(stats.EthRPCRequestsSentInCurrentUTCDay),
		EthRPCRateLimitExpiredRequests:    int32(stats.EthRPCRateLimitExpiredRequests),
		MaxExpirationTime:                 stats.MaxExpirationTime.String(),
	}
}

func NewOrderToCommonType(newOrder *NewOrder) (*zeroex.SignedOrder, error) {
	chainID, ok := math.ParseBig256(newOrder.ChainID)
	if !ok {
		return nil, fmt.Errorf("could not parse chainId as big.Int: %q", newOrder.ChainID)
	}
	makerAssetAmount, ok := math.ParseBig256(newOrder.MakerAssetAmount)
	if !ok {
		return nil, fmt.Errorf("could not parse makerAssetAmount as big.Int: %q", newOrder.MakerAssetAmount)
	}
	makerFee, ok := math.ParseBig256(newOrder.MakerFee)
	if !ok {
		return nil, fmt.Errorf("could not parse makerFee as big.Int: %q", newOrder.MakerFee)
	}
	takerAssetAmount, ok := math.ParseBig256(newOrder.TakerAssetAmount)
	if !ok {
		return nil, fmt.Errorf("could not parse takerAssetAmount as big.Int: %q", newOrder.TakerAssetAmount)
	}
	takerFee, ok := math.ParseBig256(newOrder.TakerFee)
	if !ok {
		return nil, fmt.Errorf("could not parse takerFee as big.Int: %q", newOrder.TakerFee)
	}
	expirationTimeSeconds, ok := math.ParseBig256(newOrder.ExpirationTimeSeconds)
	if !ok {
		return nil, fmt.Errorf("could not parse expirationTimeSeconds as big.Int: %q", newOrder.ExpirationTimeSeconds)
	}
	salt, ok := math.ParseBig256(newOrder.Salt)
	if !ok {
		return nil, fmt.Errorf("could not parse salt as big.Int: %q", newOrder.Salt)
	}
	return &zeroex.SignedOrder{
		Order: zeroex.Order{
			ChainID:               chainID,
			ExchangeAddress:       common.HexToAddress(newOrder.ExchangeAddress),
			MakerAddress:          common.HexToAddress(newOrder.MakerAddress),
			MakerAssetData:        common.FromHex(newOrder.MakerAssetData),
			MakerFeeAssetData:     common.FromHex(newOrder.MakerFeeAssetData),
			MakerAssetAmount:      makerAssetAmount,
			MakerFee:              makerFee,
			TakerAddress:          common.HexToAddress(newOrder.TakerAddress),
			TakerAssetData:        common.FromHex(newOrder.TakerAssetData),
			TakerFeeAssetData:     common.FromHex(newOrder.TakerFeeAssetData),
			TakerAssetAmount:      takerAssetAmount,
			TakerFee:              takerFee,
			SenderAddress:         common.HexToAddress(newOrder.SenderAddress),
			FeeRecipientAddress:   common.HexToAddress(newOrder.FeeRecipientAddress),
			ExpirationTimeSeconds: expirationTimeSeconds,
			Salt:                  salt,
		},
		Signature: common.FromHex(newOrder.Signature),
	}, nil
}

func OrderWithMetadataFromCommonType(order *types.OrderWithMetadata) *OrderWithMetadata {
	return &OrderWithMetadata{
		Hash:                     order.Hash.Hex(),
		ChainID:                  order.ChainID.String(),
		ExchangeAddress:          strings.ToLower(order.ExchangeAddress.Hex()),
		MakerAddress:             strings.ToLower(order.MakerAddress.Hex()),
		MakerAssetData:           common.Bytes2Hex(order.MakerAssetData),
		MakerFeeAssetData:        common.Bytes2Hex(order.MakerFeeAssetData),
		MakerAssetAmount:         order.MakerAssetAmount.String(),
		MakerFee:                 order.MakerFee.String(),
		TakerAddress:             strings.ToLower(order.TakerAddress.Hex()),
		TakerAssetData:           common.Bytes2Hex(order.TakerAssetData),
		TakerFeeAssetData:        common.Bytes2Hex(order.TakerFeeAssetData),
		TakerAssetAmount:         order.TakerAssetAmount.String(),
		TakerFee:                 order.TakerFee.String(),
		SenderAddress:            strings.ToLower(order.SenderAddress.Hex()),
		FeeRecipientAddress:      strings.ToLower(order.FeeRecipientAddress.Hex()),
		ExpirationTimeSeconds:    order.ExpirationTimeSeconds.String(),
		Salt:                     order.Salt.String(),
		Signature:                common.Bytes2Hex(order.Signature),
		FillableTakerAssetAmount: order.FillableTakerAssetAmount.String(),
		LastUpdated:              order.LastUpdated.Format(time.RFC3339),
	}
}

func OrdersWithMetadataFromCommonType(orders []*types.OrderWithMetadata) []*OrderWithMetadata {
	result := make([]*OrderWithMetadata, len(orders))
	for i, order := range orders {
		result[i] = OrderWithMetadataFromCommonType(order)
	}
	return result
}

func FilterKindToDBType(kind string) (db.FilterKind, error) {
	switch kind {
	case "EQUAL":
		return db.Equal, nil
	case "NOT_EQUAL":
		return db.NotEqual, nil
	case "GREATER":
		return db.Greater, nil
	case "GREATER_OR_EQUAL":
		return db.GreaterOrEqual, nil
	case "LESS":
		return db.Less, nil
	case "LESS_OR_EQUAL":
		return db.LessOrEqual, nil
	default:
		return "", fmt.Errorf("invalid filter kind: %q", kind)
	}
}

func SortDirectionToDBType(direction string) (db.SortDirection, error) {
	switch direction {
	case "ASC":
		return db.Ascending, nil
	case "DESC":
		return db.Descending, nil
	default:
		return "", fmt.Errorf("invalid sort direction: %q", direction)
	}
}
