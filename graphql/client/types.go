package client

import (
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
)

type AcceptedOrderResult struct {
	// The order that was accepted, including metadata.
	Order *OrderWithMetadata `json:"order"`
	// Whether or not the order is new. Set to true if this is the first time this Mesh node has accepted the order
	// and false otherwise.
	IsNew bool `json:"isNew"`
}

// The results of the addOrders mutation. Includes which orders were accepted and which orders where rejected.
type AddOrdersResults struct {
	// The set of orders that were accepted. Accepted orders will be watched and order events will be emitted if
	// their status changes.
	Accepted []*AcceptedOrderResult `json:"accepted"`
	// The set of orders that were rejected, including the reason they were rejected. Rejected orders will not be
	// watched.
	Rejected []*RejectedOrderResult `json:"rejected"`
}

// An on-chain contract event.
type ContractEvent struct {
	// The hash of the block where the event was generated.
	BlockHash common.Hash `json:"blockHash"`
	// The hash of the transaction where the event was generated.
	TxHash common.Hash `json:"txHash"`
	// The index of the transaction where the event was generated.
	TxIndex int `json:"txIndex"`
	// The index of the event log.
	LogIndex int `json:"logIndex"`
	// True when this was an event that was removed due to a block-reorg. False otherwise.
	IsRemoved bool `json:"isRemoved"`
	// The address of the contract that generated the event.
	Address common.Address `json:"address"`
	// The kind of event (e.g. "ERC20TransferEvent").
	Kind string `json:"kind"`
	// The parameters for the event. The parameters are different for each event kind, but will always
	// be a set of key-value pairs.
	Parameters interface{} `json:"parameters"`
}

// The block number and block hash for the latest block that has been processed by Mesh.
type LatestBlock struct {
	Number *big.Int    `json:"number"`
	Hash   common.Hash `json:"hash"`
}

// A signed 0x order according to the [protocol specification](https://github.com/0xProject/0x-protocol-specification/blob/master/v3/v3-specification.md#order-message-format.)
type Order struct {
	ChainID               *big.Int       `json:"chainId"`
	ExchangeAddress       common.Address `json:"exchangeAddress"`
	MakerAddress          common.Address `json:"makerAddress"`
	MakerAssetData        []byte         `json:"makerAssetData"`
	MakerAssetAmount      *big.Int       `json:"makerAssetAmount"`
	MakerFeeAssetData     []byte         `json:"makerFeeAssetData"`
	MakerFee              *big.Int       `json:"makerFee"`
	TakerAddress          common.Address `json:"takerAddress"`
	TakerAssetData        []byte         `json:"takerAssetData"`
	TakerAssetAmount      *big.Int       `json:"takerAssetAmount"`
	TakerFeeAssetData     []byte         `json:"takerFeeAssetData"`
	TakerFee              *big.Int       `json:"takerFee"`
	SenderAddress         common.Address `json:"senderAddress"`
	FeeRecipientAddress   common.Address `json:"feeRecipientAddress"`
	ExpirationTimeSeconds *big.Int       `json:"expirationTimeSeconds"`
	Salt                  *big.Int       `json:"salt"`
	Signature             []byte         `json:"signature"`
}

type OrderEvent struct {
	// The order that was affected.
	Order *OrderWithMetadata `json:"order"`
	// A way of classifying the effect that the order event had on the order. You can
	// think of different end states as different "types" of order events.
	EndState OrderEndState `json:"endState"`
	// The timestamp for the order event, which can be used for bookkeeping purposes.
	// If the order event was generated as a direct result of on-chain events (e.g., FILLED,
	// UNFUNDED, CANCELED), then it is set to the latest block timestamp at which the order
	// was re-validated. Otherwise (e.g., for ADDED, STOPPED_WATCHING), the timestamp corresponds
	// when the event was generated on the server side.
	Timestamp time.Time `json:"timestamp"`
	// Contains all the contract events that triggered the order to be re-validated.
	// All events that _may_ have affected the state of the order are included here.
	// It is guaranteed that at least one of the events included here will have affected
	// the order's state, but there may also be some false positives.
	ContractEvents []*ContractEvent `json:"contractEvents"`
}

// A filter on orders. Can be used in queries to only return orders that meet certain criteria.
type OrderFilter = gqltypes.OrderFilter

// A sort ordering for orders. Can be used in queries to control the order in which results are returned.
type OrderSort = gqltypes.OrderSort

// A signed 0x order along with some additional metadata about the order which is not part of the 0x protocol specification.
type OrderWithMetadata struct {
	Order
	// The hash, which can be used to uniquely identify an order.
	Hash common.Hash `json:"hash"`
	// The remaining amount of the maker asset which has not yet been filled.
	FillableTakerAssetAmount *big.Int `json:"fillableTakerAssetAmount"`
}

type RejectedOrderResult struct {
	// The hash of the order. May be null if the hash could not be computed.
	Hash *common.Hash `json:"hash"`
	// The order that was rejected.
	Order *Order `json:"order"`
	// A machine-readable code indicating why the order was rejected. This code is designed to
	// be used by programs and applications and will never change without breaking backwards-compatibility.
	Code RejectedOrderCode `json:"code"`
	// A human-readable message indicating why the order was rejected. This message may change
	// in future releases and is not covered by backwards-compatibility guarantees.
	Message string `json:"message"`
}

// Contains configuration options and various stats for Mesh.
type Stats struct {
	Version                           string       `json:"version"`
	PubSubTopic                       string       `json:"pubSubTopic"`
	Rendezvous                        string       `json:"rendezvous"`
	PeerID                            string       `json:"peerID"`
	EthereumChainID                   int          `json:"ethereumChainID"`
	LatestBlock                       *LatestBlock `json:"latestBlock"`
	NumPeers                          int          `json:"numPeers"`
	NumOrders                         int          `json:"numOrders"`
	NumOrdersIncludingRemoved         int          `json:"numOrdersIncludingRemoved"`
	StartOfCurrentUTCDay              time.Time    `json:"startOfCurrentUTCDay"`
	EthRPCRequestsSentInCurrentUTCDay int          `json:"ethRPCRequestsSentInCurrentUTCDay"`
	EthRPCRateLimitExpiredRequests    int          `json:"ethRPCRateLimitExpiredRequests"`
	MaxExpirationTime                 *big.Int     `json:"maxExpirationTime"`
}

// The kind of comparison to be used in a filter.
type FilterKind = gqltypes.FilterKind

const (
	FilterKindEqual          FilterKind = "EQUAL"
	FilterKindNotEqual       FilterKind = "NOT_EQUAL"
	FilterKindGreater        FilterKind = "GREATER"
	FilterKindGreaterOrEqual FilterKind = "GREATER_OR_EQUAL"
	FilterKindLess           FilterKind = "LESS"
	FilterKindLessOrEqual    FilterKind = "LESS_OR_EQUAL"
)

var AllFilterKind = gqltypes.AllFilterKind

type OrderEndState = gqltypes.OrderEndState

const (
	// The order was successfully validated and added to the Mesh node. The order is now being watched and any changes to
	// the fillability will result in subsequent order events.
	OrderEndStateAdded OrderEndState = "ADDED"
	// The order was filled for a partial amount. The order is still fillable up to the fillableTakerAssetAmount.
	OrderEndStateFilled OrderEndState = "FILLED"
	// The order was fully filled and its remaining fillableTakerAssetAmount is 0. The order is no longer fillable.
	OrderEndStateFullyFilled OrderEndState = "FULLY_FILLED"
	// The order was cancelled and is no longer fillable.
	OrderEndStateCancelled OrderEndState = "CANCELLED"
	// The order expired and is no longer fillable.
	OrderEndStateExpired OrderEndState = "EXPIRED"
	// The order was previously expired, but due to a block re-org it is no longer considered expired (should be rare).
	OrderEndStateUnexpired OrderEndState = "UNEXPIRED"
	// The order has become unfunded and is no longer fillable. This can happen if the maker makes a transfer or changes their allowance.
	OrderEndStateUnfunded OrderEndState = "UNFUNDED"
	// The fillability of the order has increased. This can happen if a previously processed fill event gets reverted due to a block re-org,
	// or if a maker makes a transfer or changes their allowance.
	OrderEndStateFillabilityIncreased OrderEndState = "FILLABILITY_INCREASED"
	// The order is potentially still valid but was removed for a different reason (e.g.
	// the database is full or the peer that sent the order was misbehaving). The order will no longer be watched
	// and no further events for this order will be emitted. In some cases, the order may be re-added in the
	// future.
	OrderEndStateStoppedWatching OrderEndState = "STOPPED_WATCHING"
)

var AllOrderEndState = gqltypes.AllOrderEndState

// An enum containing all the order fields for which filters and/or sorting is supported.
type OrderField = gqltypes.OrderField

const (
	OrderFieldHash                     OrderField = "hash"
	OrderFieldChainID                  OrderField = "chainId"
	OrderFieldExchangeAddress          OrderField = "exchangeAddress"
	OrderFieldMakerAddress             OrderField = "makerAddress"
	OrderFieldMakerAssetData           OrderField = "makerAssetData"
	OrderFieldMakerAssetAmount         OrderField = "makerAssetAmount"
	OrderFieldMakerFeeAssetData        OrderField = "makerFeeAssetData"
	OrderFieldMakerFee                 OrderField = "makerFee"
	OrderFieldTakerAddress             OrderField = "takerAddress"
	OrderFieldTakerAssetData           OrderField = "takerAssetData"
	OrderFieldTakerAssetAmount         OrderField = "takerAssetAmount"
	OrderFieldTakerFeeAssetData        OrderField = "takerFeeAssetData"
	OrderFieldTakerFee                 OrderField = "takerFee"
	OrderFieldSenderAddress            OrderField = "senderAddress"
	OrderFieldFeeRecipientAddress      OrderField = "feeRecipientAddress"
	OrderFieldExpirationTimeSeconds    OrderField = "expirationTimeSeconds"
	OrderFieldSalt                     OrderField = "salt"
	OrderFieldFillableTakerAssetAmount OrderField = "fillableTakerAssetAmount"
)

var AllOrderField = gqltypes.AllOrderField

// A set of all possible codes included in RejectedOrderResult.
type RejectedOrderCode = gqltypes.RejectedOrderCode

const (
	RejectedOrderCodeEthRPCRequestFailed              RejectedOrderCode = "ETH_RPC_REQUEST_FAILED"
	RejectedOrderCodeOrderHasInvalidMakerAssetAmount  RejectedOrderCode = "ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT"
	RejectedOrderCodeOrderHasInvalidTakerAssetAmount  RejectedOrderCode = "ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT"
	RejectedOrderCodeOrderExpired                     RejectedOrderCode = "ORDER_EXPIRED"
	RejectedOrderCodeOrderFullyFilled                 RejectedOrderCode = "ORDER_FULLY_FILLED"
	RejectedOrderCodeOrderCancelled                   RejectedOrderCode = "ORDER_CANCELLED"
	RejectedOrderCodeOrderUnfunded                    RejectedOrderCode = "ORDER_UNFUNDED"
	RejectedOrderCodeOrderHasInvalidMakerAssetData    RejectedOrderCode = "ORDER_HAS_INVALID_MAKER_ASSET_DATA"
	RejectedOrderCodeOrderHasInvalidMakerFeeAssetData RejectedOrderCode = "ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA"
	RejectedOrderCodeOrderHasInvalidTakerAssetData    RejectedOrderCode = "ORDER_HAS_INVALID_TAKER_ASSET_DATA"
	RejectedOrderCodeOrderHasInvalidTakerFeeAssetData RejectedOrderCode = "ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA"
	RejectedOrderCodeOrderHasInvalidSignature         RejectedOrderCode = "ORDER_HAS_INVALID_SIGNATURE"
	RejectedOrderCodeOrderMaxExpirationExceeded       RejectedOrderCode = "ORDER_MAX_EXPIRATION_EXCEEDED"
	RejectedOrderCodeInternalError                    RejectedOrderCode = "INTERNAL_ERROR"
	RejectedOrderCodeMaxOrderSizeExceeded             RejectedOrderCode = "MAX_ORDER_SIZE_EXCEEDED"
	RejectedOrderCodeOrderAlreadyStoredAndUnfillable  RejectedOrderCode = "ORDER_ALREADY_STORED_AND_UNFILLABLE"
	RejectedOrderCodeOrderForIncorrectChain           RejectedOrderCode = "ORDER_FOR_INCORRECT_CHAIN"
	RejectedOrderCodeIncorrectExchangeAddress         RejectedOrderCode = "INCORRECT_EXCHANGE_ADDRESS"
	RejectedOrderCodeSenderAddressNotAllowed          RejectedOrderCode = "SENDER_ADDRESS_NOT_ALLOWED"
	RejectedOrderCodeDatabaseFullOfOrders             RejectedOrderCode = "DATABASE_FULL_OF_ORDERS"
	RejectedOrderCodeInvalidSchema                    RejectedOrderCode = "INVALID_SCHEMA"
)

var AllRejectedOrderCode = gqltypes.AllRejectedOrderCode

// The direction to sort in. Ascending means lowest to highest. Descending means highest to lowest.
type SortDirection = gqltypes.SortDirection

var AllSortDirection = gqltypes.AllSortDirection

const (
	SortDirectionAsc  = gqltypes.SortDirectionAsc
	SortDirectionDesc = gqltypes.SortDirectionDesc
)
