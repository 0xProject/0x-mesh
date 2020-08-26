package ordervalidator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

// Specifies the max number of eth_call requests we want to make concurrently.
// Additional requests will block until an ongoing request has completed.
const concurrencyLimit = 5

// RejectedOrderInfo encapsulates all the needed information to understand _why_ a 0x order
// was rejected (i.e. did not pass) order validation. Since there are many potential reasons, some
// Mesh-specific, others 0x-specific and others due to external factors (i.e., network
// disruptions, etc...), we categorize them into `Kind`s and uniquely identify the reasons for
// machines with a `Code`
type RejectedOrderInfo struct {
	OrderHash   common.Hash         `json:"orderHash"`
	SignedOrder *zeroex.SignedOrder `json:"signedOrder"`
	Kind        RejectedOrderKind   `json:"kind"`
	Status      RejectedOrderStatus `json:"status"`
}

// AcceptedOrderInfo represents an fillable order and how much it could be filled for
type AcceptedOrderInfo struct {
	OrderHash                common.Hash         `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount *big.Int            `json:"fillableTakerAssetAmount"`
	IsNew                    bool                `json:"isNew"`
}

type acceptedOrderInfoJSON struct {
	OrderHash                string              `json:"orderHash"`
	SignedOrder              *zeroex.SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount string              `json:"fillableTakerAssetAmount"`
	IsNew                    bool                `json:"isNew"`
}

// MarshalJSON is a custom Marshaler for AcceptedOrderInfo
func (a AcceptedOrderInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"orderHash":                a.OrderHash.Hex(),
		"signedOrder":              a.SignedOrder,
		"fillableTakerAssetAmount": a.FillableTakerAssetAmount.String(),
		"isNew":                    a.IsNew,
	})
}

// UnmarshalJSON implements a custom JSON unmarshaller for the OrderEvent type
func (a *AcceptedOrderInfo) UnmarshalJSON(data []byte) error {
	var acceptedOrderInfoJSON acceptedOrderInfoJSON
	err := json.Unmarshal(data, &acceptedOrderInfoJSON)
	if err != nil {
		return err
	}

	a.OrderHash = common.HexToHash(acceptedOrderInfoJSON.OrderHash)
	a.SignedOrder = acceptedOrderInfoJSON.SignedOrder
	a.IsNew = acceptedOrderInfoJSON.IsNew
	var ok bool
	a.FillableTakerAssetAmount, ok = math.ParseBig256(acceptedOrderInfoJSON.FillableTakerAssetAmount)
	if !ok {
		return errors.New("Invalid uint256 number encountered for FillableTakerAssetAmount")
	}
	return nil
}

// RejectedOrderStatus enumerates all the unique reasons for an orders rejection
type RejectedOrderStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RejectedOrderStatus values
var (
	ROEthRPCRequestFailed = RejectedOrderStatus{
		Code:    "EthRPCRequestFailed",
		Message: "network request to Ethereum RPC endpoint failed",
	}
	ROInvalidMakerAssetAmount = RejectedOrderStatus{
		Code:    "OrderHasInvalidMakerAssetAmount",
		Message: "order makerAssetAmount cannot be 0",
	}
	ROInvalidTakerAssetAmount = RejectedOrderStatus{
		Code:    "OrderHasInvalidTakerAssetAmount",
		Message: "order takerAssetAmount cannot be 0",
	}
	ROExpired = RejectedOrderStatus{
		Code:    "OrderExpired",
		Message: "order expired according to latest block timestamp",
	}
	ROFullyFilled = RejectedOrderStatus{
		Code:    "OrderFullyFilled",
		Message: "order already fully filled",
	}
	ROCancelled = RejectedOrderStatus{
		Code:    "OrderCancelled",
		Message: "order cancelled",
	}
	ROUnfunded = RejectedOrderStatus{
		Code:    "OrderUnfunded",
		Message: "maker has insufficient balance or allowance for this order to be filled",
	}
	ROInvalidMakerAssetData = RejectedOrderStatus{
		Code:    "OrderHasInvalidMakerAssetData",
		Message: "order makerAssetData must encode a supported assetData type",
	}
	ROInvalidMakerFeeAssetData = RejectedOrderStatus{
		Code:    "OrderHasInvalidMakerFeeAssetData",
		Message: "order makerFeeAssetData must encode a supported assetData type",
	}
	ROInvalidTakerAssetData = RejectedOrderStatus{
		Code:    "OrderHasInvalidTakerAssetData",
		Message: "order takerAssetData must encode a supported assetData type",
	}
	ROInvalidTakerFeeAssetData = RejectedOrderStatus{
		Code:    "OrderHasInvalidTakerFeeAssetData",
		Message: "order takerFeeAssetData must encode a supported assetData type",
	}
	ROInvalidSignature = RejectedOrderStatus{
		Code:    "OrderHasInvalidSignature",
		Message: "order signature must be valid",
	}
	ROMaxExpirationExceeded = RejectedOrderStatus{
		Code:    "OrderMaxExpirationExceeded",
		Message: "order expiration too far in the future",
	}
	ROInternalError = RejectedOrderStatus{
		Code:    "InternalError",
		Message: "an unexpected internal error has occurred",
	}
	ROMaxOrderSizeExceeded = RejectedOrderStatus{
		Code:    "MaxOrderSizeExceeded",
		Message: fmt.Sprintf("order exceeds the maximum encoded size of %d bytes", constants.MaxOrderSizeInBytes),
	}
	ROOrderAlreadyStoredAndUnfillable = RejectedOrderStatus{
		Code:    "OrderAlreadyStoredAndUnfillable",
		Message: "order is already stored and is unfillable. Mesh keeps unfillable orders in storage for a little while incase a block re-org makes them fillable again",
	}
	ROIncorrectChain = RejectedOrderStatus{
		Code:    "OrderForIncorrectChain",
		Message: "order was created for a different chain than the one this Mesh node is configured to support",
	}
	ROIncorrectExchangeAddress = RejectedOrderStatus{
		Code:    "IncorrectExchangeAddress",
		Message: "the exchange address for the order does not match the chain ID/network ID",
	}
	ROSenderAddressNotAllowed = RejectedOrderStatus{
		Code:    "SenderAddressNotAllowed",
		Message: "orders with a senderAddress are not currently supported",
	}
	RODatabaseFullOfOrders = RejectedOrderStatus{
		Code:    "DatabaseFullOfOrders",
		Message: "database is full of pinned orders and no orders can be deleted to make space (consider increasing MAX_ORDERS_IN_STORAGE)",
	}
	ROTakerAddressNotAllowed = RejectedOrderStatus{
		Code:    "TakerAddressNotAllowed",
		Message: "the taker address is not a whitelisted address",
	}
)

// ROInvalidSchemaCode is the RejectedOrderStatus emitted if an order doesn't conform to the order schema
const ROInvalidSchemaCode = "InvalidSchema"

// ConvertRejectOrderCodeToOrderEventEndState converts an RejectOrderCode to an OrderEventEndState type
func ConvertRejectOrderCodeToOrderEventEndState(rejectedOrderStatus RejectedOrderStatus) (zeroex.OrderEventEndState, bool) {
	switch rejectedOrderStatus {
	case ROExpired:
		return zeroex.ESOrderExpired, true
	case ROFullyFilled:
		return zeroex.ESOrderFullyFilled, true
	case ROCancelled:
		return zeroex.ESOrderCancelled, true
	case ROUnfunded:
		return zeroex.ESOrderBecameUnfunded, true
	default:
		// Catch-all returns Invalid OrderEventEndState
		return zeroex.ESInvalid, false
	}
}

// RejectedOrderKind enumerates all kinds of reasons an order could be rejected by Mesh
type RejectedOrderKind string

// RejectedOrderKind values
const (
	ZeroExValidation = RejectedOrderKind("ZEROEX_VALIDATION")
	MeshError        = RejectedOrderKind("MESH_ERROR")
	MeshValidation   = RejectedOrderKind("MESH_VALIDATION")
)

// ValidationResults defines the validation results returned from BatchValidate
// Within this context, an order is `Accepted` if it passes all the 0x schema tests
// and is fillable for a non-zero amount. An order is `Rejected` if it does not
// satisfy these conditions OR if we were unable to complete the validation process
// for whatever reason
type ValidationResults struct {
	Accepted []*AcceptedOrderInfo `json:"accepted"`
	Rejected []*RejectedOrderInfo `json:"rejected"`
}

// OrderValidator validates 0x orders
type OrderValidator struct {
	maxRequestContentLength      int
	devUtils                     *wrappers.DevUtilsCaller
	assetDataDecoder             *zeroex.AssetDataDecoder
	chainID                      int
	cachedFeeRecipientToEndpoint map[common.Address]string
	contractAddresses            ethereum.ContractAddresses
}

// New instantiates a new order validator
func New(contractCaller bind.ContractCaller, chainID int, maxRequestContentLength int, contractAddresses ethereum.ContractAddresses) (*OrderValidator, error) {
	devUtils, err := wrappers.NewDevUtilsCaller(contractAddresses.DevUtils, contractCaller)
	if err != nil {
		return nil, err
	}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

	return &OrderValidator{
		maxRequestContentLength:      maxRequestContentLength,
		devUtils:                     devUtils,
		assetDataDecoder:             assetDataDecoder,
		chainID:                      chainID,
		cachedFeeRecipientToEndpoint: map[common.Address]string{},
		contractAddresses:            contractAddresses,
	}, nil
}

// BatchValidate retrieves all the information needed to validate the supplied orders.
// It splits the orders into chunks of `chunkSize`, and makes no more then `concurrencyLimit`
// requests concurrently. If a request fails, re-attempt it up to four times before giving up.
// If some requests fail, this method still returns whatever order information it was able to
// retrieve up until the failure.
// The `validationBlock` parameter lets the caller specify a specific block at which to validate
// the orders. This can be set to the `latest` block or any other historical block.
func (o *OrderValidator) BatchValidate(ctx context.Context, signedOrders []*zeroex.SignedOrder, areNewOrders bool, validationBlock *types.MiniHeader) *ValidationResults {
	if len(signedOrders) == 0 {
		return &ValidationResults{}
	}
	offchainValidSignedOrders, rejectedOrderInfos := o.BatchOffchainValidation(signedOrders)
	validationResults := &ValidationResults{
		Accepted: []*AcceptedOrderInfo{},
		Rejected: rejectedOrderInfos,
	}

	signedOrderChunks := [][]*zeroex.SignedOrder{}
	chunkSizes := o.computeOptimalChunkSizes(offchainValidSignedOrders)
	for _, chunkSize := range chunkSizes {
		signedOrderChunks = append(signedOrderChunks, offchainValidSignedOrders[:chunkSize])
		offchainValidSignedOrders = offchainValidSignedOrders[chunkSize:]
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for _, signedOrders := range signedOrderChunks {
		wg.Add(1)
		go func(signedOrders []*zeroex.SignedOrder) {
			// FIXME - Is this needed?
			// trimmedOrders := []wrappers.LibOrderOrder{}
			// for _, signedOrder := range signedOrders {
			// 	trimmedOrders = append(trimmedOrders, signedOrder.Trim())
			// }
			// signatures := [][]byte{}
			// for _, signedOrder := range signedOrders {
			// 	signatures = append(signatures, signedOrder.Signature)
			// }

			defer wg.Done()

			select {
			case <-ctx.Done():
			// Blocks until a slot opens up in the semaphore. We read off of the
			// semaphore whenever onchain validation completes to allow another
			// goroutine to begin processing.
			case semaphoreChan <- struct{}{}:
				defer func() { <-semaphoreChan }()
				o.batchOnchainValidation(ctx, signedOrders, validationBlock, areNewOrders, validationResults)
			}
		}(signedOrders)
	}

	wg.Wait()
	return validationResults
}

// BatchOffchainValidation performs all off-chain validation checks on a batch of 0x orders.
// These checks include:
// - `MakerAssetAmount` and `TakerAssetAmount` cannot be 0
// - `AssetData` fields contain properly encoded, and currently supported assetData (ERC20 & ERC721 for now)
// - `Signature` contains a properly encoded 0x signature
// - Validate that order isn't expired
// Returns the signedOrders that are off-chain valid along with an array of orderInfo for the rejected orders
func (o *OrderValidator) BatchOffchainValidation(signedOrders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, []*RejectedOrderInfo) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	offchainValidSignedOrders := []*zeroex.SignedOrder{}
	for _, signedOrder := range signedOrders {
		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.WithError(err).WithField("signedOrder", signedOrder).Error("Computing the orderHash failed unexpectedly")
		}
		if !signedOrder.ExpirationTimeSeconds.IsInt64() {
			// Shouldn't happen because we separately enforce a max expiration time.
			// See core/validation.go.
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        MeshValidation,
				Status:      ROMaxExpirationExceeded,
			})
			continue
		}

		if signedOrder.MakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROInvalidMakerAssetAmount,
			})
			continue
		}
		if signedOrder.TakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROInvalidTakerAssetAmount,
			})
			continue
		}

		isMakerAssetDataSupported := o.isSupportedAssetData(signedOrder.MakerAssetData)
		if !isMakerAssetDataSupported {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROInvalidMakerAssetData,
			})
			continue
		}
		isTakerAssetDataSupported := o.isSupportedAssetData(signedOrder.TakerAssetData)
		if !isTakerAssetDataSupported {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROInvalidTakerAssetData,
			})
			continue
		}
		// If the MakerFee is zero, the fee asset data will not affect the
		// validity of the signed order.
		// https://github.com/0xProject/0x-monorepo/blob/development/contracts/exchange/contracts/src/MixinAssetProxyDispatcher.sol#L90
		if signedOrder.MakerFee.Cmp(big.NewInt(0)) == 1 && len(signedOrder.MakerFeeAssetData) != 0 {
			isMakerFeeAssetDataSupported := o.isSupportedAssetData(signedOrder.MakerFeeAssetData)
			if !isMakerFeeAssetDataSupported {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        ZeroExValidation,
					Status:      ROInvalidMakerFeeAssetData,
				})
				continue
			}
		}
		// If the TakerFee is zero, the fee asset data will not affect the
		// validity of the signed order.
		// https://github.com/0xProject/0x-monorepo/blob/development/contracts/exchange/contracts/src/MixinAssetProxyDispatcher.sol#L90
		if signedOrder.TakerFee.Cmp(big.NewInt(0)) == 1 && len(signedOrder.TakerFeeAssetData) != 0 {
			isTakerFeeAssetDataSupported := o.isSupportedAssetData(signedOrder.TakerFeeAssetData)
			if !isTakerFeeAssetDataSupported {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        ZeroExValidation,
					Status:      ROInvalidTakerFeeAssetData,
				})
				continue
			}
		}

		isSupportedSignature := isSupportedSignature(signedOrder.Signature)
		if !isSupportedSignature {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROInvalidSignature,
			})
			continue
		}

		offchainValidSignedOrders = append(offchainValidSignedOrders, signedOrder)
	}

	return offchainValidSignedOrders, rejectedOrderInfos
}

// batchOnchainValidation validates a list of signed orders using the deployed
// DevUtils contract. This validation performs signature validation, checks balances
// and allowances, and identifies other issues in asset data (for example, DevUtils
// will invalidate MultiAssetProxy orders that contain duplicate ERC721 asset data).
func (o *OrderValidator) batchOnchainValidation(
	ctx context.Context,
	signedOrders []*zeroex.SignedOrder,
	validationBlock *types.MiniHeader,
	areNewOrders bool,
	validationResults *ValidationResults,
) {
	trimmedOrders := []wrappers.LibOrderOrder{}
	for _, signedOrder := range signedOrders {
		trimmedOrders = append(trimmedOrders, signedOrder.Trim())
	}
	signatures := [][]byte{}
	for _, signedOrder := range signedOrders {
		signatures = append(signatures, signedOrder.Signature)
	}

	// Attempt to make the eth_call request 4 times with an exponential back-off.
	maxDuration := 4 * time.Second
	b := &backoff.Backoff{
		Min:    250 * time.Millisecond, // First back-off length
		Max:    maxDuration,            // Longest back-off length
		Factor: 2,                      // Factor to multiple each successive back-off
	}

	for {
		opts := &bind.CallOpts{
			// HACK(albrow): From field should not be required for eth_call but
			// including it here is a workaround for a bug in Ganache. Removing
			// this line causes Ganache to crash.
			From:    constants.GanacheDummyERC721TokenAddress,
			Pending: false,
			Context: ctx,
		}
		opts.BlockNumber = validationBlock.Number

		results, err := o.devUtils.GetOrderRelevantStates(opts, trimmedOrders, signatures)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"attempt":   b.Attempt(),
				"numOrders": len(trimmedOrders),
			}).Info("GetOrderRelevantStates request failed")
			d := b.Duration()
			if d == maxDuration {
				var fields log.Fields
				match, regexpErr := regexp.MatchString("abi: improperly formatted output", err.Error())
				if regexpErr != nil {
					log.WithField("error", regexpErr).Error("Unexpectedly failed to test regexp on error")
				}
				if err.Error() == "VM execution error." || match {
					fields = log.Fields{
						"error":     err.Error(),
						"numOrders": len(trimmedOrders),
						"orders":    trimmedOrders,
					}
				} else {
					fields = log.Fields{
						"error":     err.Error(),
						"numOrders": len(trimmedOrders),
					}
				}
				log.WithFields(fields).Warning("Gave up on GetOrderRelevantStates request after backoff limit reached")
				for _, signedOrder := range signedOrders {
					orderHash, err := signedOrder.ComputeOrderHash()
					if err != nil {
						log.WithField("error", err).Error("Unexpectedly failed to generate orderHash")
						continue
					}
					validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
						OrderHash:   orderHash,
						SignedOrder: signedOrder,
						Kind:        MeshError,
						Status:      ROEthRPCRequestFailed,
					})
				}
				return // Give up after 4 attempts
			}
			time.Sleep(d)
			continue
		}

		for j, orderInfo := range results.OrdersInfo {
			isValidSignature := results.IsValidSignature[j]
			fillableTakerAssetAmount := results.FillableTakerAssetAmounts[j]
			orderHash := common.Hash(orderInfo.OrderHash)
			signedOrder := signedOrders[j]
			orderStatus := zeroex.OrderStatus(orderInfo.OrderStatus)
			if !isValidSignature {
				orderStatus = zeroex.OSSignatureInvalid
			}
			switch orderStatus {
			case zeroex.OSExpired, zeroex.OSFullyFilled, zeroex.OSCancelled, zeroex.OSSignatureInvalid:
				var status RejectedOrderStatus
				switch orderStatus {
				case zeroex.OSExpired:
					status = ROExpired
				case zeroex.OSFullyFilled:
					status = ROFullyFilled
				case zeroex.OSCancelled:
					status = ROCancelled
				case zeroex.OSSignatureInvalid:
					status = ROInvalidSignature
				}
				validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        ZeroExValidation,
					Status:      status,
				})
				continue
			case zeroex.OSFillable:
				remainingTakerAssetAmount := big.NewInt(0).Sub(signedOrder.TakerAssetAmount, orderInfo.OrderTakerAssetFilledAmount)
				// If `fillableTakerAssetAmount` != `remainingTakerAssetAmount`, the order is partially fillable. We consider
				// partially fillable orders as invalid
				if fillableTakerAssetAmount.Cmp(remainingTakerAssetAmount) != 0 {
					validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
						OrderHash:   orderHash,
						SignedOrder: signedOrder,
						Kind:        ZeroExValidation,
						Status:      ROUnfunded,
					})
				} else {
					validationResults.Accepted = append(validationResults.Accepted, &AcceptedOrderInfo{
						OrderHash:                orderHash,
						SignedOrder:              signedOrder,
						FillableTakerAssetAmount: fillableTakerAssetAmount,
						IsNew:                    areNewOrders,
					})
				}
				continue
			}
		}
		return
	}
}

func (o *OrderValidator) isSupportedAssetData(assetData []byte) bool {
	assetDataName, err := o.assetDataDecoder.GetName(assetData)
	if err != nil {
		return false
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
		return o.isSupportedStaticCallData(decodedAssetData)
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "ERC20Bridge":
		var decodedAssetData zeroex.ERC20BridgeAssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
		// We currently restrict ERC20Bridge orders to those referencing the
		// Chai bridge. If the ChaiBridge is not deployed on the selected network
		// we also reject the ERC20Bridge asset.
		if o.contractAddresses.ChaiBridge == constants.NullAddress || decodedAssetData.BridgeAddress != o.contractAddresses.ChaiBridge {
			return false
		}
	default:
		return false
	}
	return true
}

func (o *OrderValidator) isSupportedStaticCallData(staticCallAssetData zeroex.StaticCallAssetData) bool {
	staticCallDataName, err := o.assetDataDecoder.GetName(staticCallAssetData.StaticCallData)
	if err != nil {
		return false
	}
	switch staticCallDataName {
	case "checkGasPrice":
		var decodedStaticCallData zeroex.CheckGasPriceStaticCallData
		err := o.assetDataDecoder.Decode(staticCallAssetData.StaticCallData, &decodedStaticCallData)
		if err != nil {
			return false
		}
		// We currently restrict the `checkGasPrice` staticcall to the known MaximumGasPrice contract.
		if o.contractAddresses.MaximumGasPrice == constants.NullAddress || staticCallAssetData.StaticCallTargetAddress != o.contractAddresses.MaximumGasPrice {
			return false
		}
	default:
		return false
	}
	return true
}

// jsonRPCPayloadByteLength is the number of bytes occupied by the default call to `getOrderRelevantStates` with 0 signedOrders
// passed in. The `data` includes the empty `getOrderRelevantStates` calldata.
/*
{
    "id": 2,
    "jsonrpc": "2.0",
    "method": "eth_call",
    "params": [
        {
            "data": "0x7f46448d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
            "from": "0x0000000000000000000000000000000000000000",
            "to": "0x0000000000000000000000000000000000000000"
        },
        "latest"
    ]
}
*/
const jsonRPCPayloadByteLength = 444

func numWords(bytes []byte) int {
	return (len(bytes) + 31) / 32
}

func computeABIEncodedSignedOrderStringLength(signedOrder *zeroex.SignedOrder) int {
	// The fixed size fields in a SignedOrder take up 1536 bytes. The variable length fields take up 64 characters per 256-bit word. This is because each byte in signedOrder is as two bytes in the JSON string.
	return 1536 + 64*(numWords(signedOrder.Order.MakerAssetData)+
		numWords(signedOrder.Order.TakerAssetData)+
		numWords(signedOrder.Order.MakerFeeAssetData)+
		numWords(signedOrder.Order.MakerFeeAssetData))
}

// computeOptimalChunkSizes splits the signedOrders into chunks where the payload size of each chunk
// is beneath the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizes(signedOrders []*zeroex.SignedOrder) []int {
	chunkSizes := []int{}

	payloadLength := jsonRPCPayloadByteLength
	nextChunkSize := 0
	for _, signedOrder := range signedOrders {
		encodedSignedOrderByteLength := computeABIEncodedSignedOrderStringLength(signedOrder)
		if payloadLength+encodedSignedOrderByteLength < o.maxRequestContentLength {
			payloadLength += encodedSignedOrderByteLength
			nextChunkSize++
		} else {
			if nextChunkSize == 0 {
				// This case should never be hit since we enforce that EthereumRPCMaxContentLength >= maxOrderSizeInBytes
				log.WithField("signedOrder", signedOrder).Panic("EthereumRPCMaxContentLength is set so low, a single 0x order cannot fit beneath the payload limit")
			}
			chunkSizes = append(chunkSizes, nextChunkSize)
			nextChunkSize = 1
			payloadLength = jsonRPCPayloadByteLength + encodedSignedOrderByteLength
		}
	}
	if nextChunkSize != 0 {
		chunkSizes = append(chunkSizes, nextChunkSize)
	}

	return chunkSizes
}

func isSupportedSignature(signature []byte) bool {
	if len(signature) == 0 {
		return false
	}
	signatureType := zeroex.SignatureType(signature[len(signature)-1])

	switch signatureType {
	case zeroex.InvalidSignature, zeroex.IllegalSignature:
		return false

	case zeroex.EIP712Signature:
		if len(signature) != 66 {
			return false
		}
		// TODO(fabio): Do further validation by splitting into r,s,v and do ECRecover

	case zeroex.EthSignSignature:
		if len(signature) != 66 {
			return false
		}
		// TODO(fabio): Do further validation by splitting into r,s,v, add prefix to hash
		// and do ECRecover

	case zeroex.ValidatorSignature:
		if len(signature) < 21 {
			return false
		}

	case zeroex.PreSignedSignature, zeroex.WalletSignature, zeroex.EIP1271WalletSignature:
		return true

	default:
		return false

	}

	return true
}
