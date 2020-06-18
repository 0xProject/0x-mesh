package ordervalidator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	ROCoordinatorRequestFailed = RejectedOrderStatus{
		Code:    "CoordinatorRequestFailed",
		Message: "network request to coordinator server endpoint failed",
	}
	ROCoordinatorSoftCancelled = RejectedOrderStatus{
		Code:    "CoordinatorSoftCancelled",
		Message: "order was soft-cancelled via the coordinator server",
	}
	ROCoordinatorEndpointNotFound = RejectedOrderStatus{
		Code:    "CoordinatorEndpointNotFound",
		Message: "corresponding coordinator endpoint not found in CoordinatorRegistry contract",
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
	CoordinatorError = RejectedOrderKind("COORDINATOR_ERROR")
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
	devUtilsABI                  abi.ABI
	devUtils                     *wrappers.DevUtilsCaller
	coordinatorRegistry          *wrappers.CoordinatorRegistryCaller
	assetDataDecoder             *zeroex.AssetDataDecoder
	chainID                      int
	cachedFeeRecipientToEndpoint map[common.Address]string
	contractAddresses            ethereum.ContractAddresses
}

// New instantiates a new order validator
func New(contractCaller bind.ContractCaller, chainID int, maxRequestContentLength int, contractAddresses ethereum.ContractAddresses) (*OrderValidator, error) {
	devUtilsABI, err := abi.JSON(strings.NewReader(wrappers.DevUtilsABI))
	if err != nil {
		return nil, err
	}
	devUtils, err := wrappers.NewDevUtilsCaller(contractAddresses.DevUtils, contractCaller)
	if err != nil {
		return nil, err
	}
	coordinatorRegistry, err := wrappers.NewCoordinatorRegistryCaller(contractAddresses.CoordinatorRegistry, contractCaller)
	if err != nil {
		return nil, err
	}
	assetDataDecoder := zeroex.NewAssetDataDecoder()

	return &OrderValidator{
		maxRequestContentLength:      maxRequestContentLength,
		devUtilsABI:                  devUtilsABI,
		devUtils:                     devUtils,
		coordinatorRegistry:          coordinatorRegistry,
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
// The `blockNumber` parameter lets the caller specify a specific block height at which to validate
// the orders. This can be set to the `latest` block or any other historical block number.
func (o *OrderValidator) BatchValidate(ctx context.Context, rawSignedOrders []*zeroex.SignedOrder, areNewOrders bool, blockNumber *big.Int) *ValidationResults {
	if len(rawSignedOrders) == 0 {
		return &ValidationResults{}
	}
	offchainValidSignedOrders, rejectedOrderInfos := o.BatchOffchainValidation(rawSignedOrders)
	validationResults := &ValidationResults{
		Accepted: []*AcceptedOrderInfo{},
		Rejected: rejectedOrderInfos,
	}

	// Validate Coordinator orders for soft-cancels
	signedOrders, coordinatorRejectedOrderInfos := o.batchValidateSoftCancelled(ctx, offchainValidSignedOrders)
	for _, rejectedOrderInfo := range coordinatorRejectedOrderInfos {
		validationResults.Rejected = append(validationResults.Rejected, rejectedOrderInfo)
	}

	signedOrderChunks := [][]*zeroex.SignedOrder{}
	chunkSizes := o.computeOptimalChunkSizes(signedOrders)
	for _, chunkSize := range chunkSizes {
		signedOrderChunks = append(signedOrderChunks, signedOrders[:chunkSize])
		signedOrders = signedOrders[chunkSize:]
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for i, signedOrders := range signedOrderChunks {
		wg.Add(1)
		go func(signedOrders []*zeroex.SignedOrder, i int) {
			trimmedOrders := []wrappers.LibOrderOrder{}
			for _, signedOrder := range signedOrders {
				trimmedOrders = append(trimmedOrders, signedOrder.Trim())
			}
			signatures := [][]byte{}
			for _, signedOrder := range signedOrders {
				signatures = append(signatures, signedOrder.Signature)
			}

			defer wg.Done()

			// Add one to the semaphore chan. If it already has concurrencyLimit values,
			// the request blocks here until one frees up.
			semaphoreChan <- struct{}{}

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
				opts.BlockNumber = blockNumber

				results, err := o.devUtils.GetOrderRelevantStates(opts, trimmedOrders, signatures)
				if err != nil {
					log.WithFields(log.Fields{
						"error":     err.Error(),
						"attempt":   b.Attempt(),
						"numOrders": len(trimmedOrders),
					}).Info("GetOrderRelevantStates request failed")
					d := b.Duration()
					if d == maxDuration {
						<-semaphoreChan
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

				<-semaphoreChan
				return
			}
		}(signedOrders, i)
	}

	wg.Wait()
	return validationResults
}

type softCancelResponse struct {
	OrderHashes []common.Hash `json:"orderHashes"`
}

// batchValidateSoftCancelled validates any order specifying the Coordinator contract as the `senderAddress` to ensure
// that it hasn't been cancelled off-chain (soft cancellation). It does this by looking up the Coordinator server endpoint
// given the `feeRecipientAddress` specified in the order, and then hitting that endpoint to query whether the orders have
// been soft cancelled.
func (o *OrderValidator) batchValidateSoftCancelled(ctx context.Context, signedOrders []*zeroex.SignedOrder) ([]*zeroex.SignedOrder, []*RejectedOrderInfo) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	validSignedOrders := []*zeroex.SignedOrder{}

	endpointToSignedOrders := map[string][]*zeroex.SignedOrder{}
	for _, signedOrder := range signedOrders {
		if signedOrder.SenderAddress != o.contractAddresses.Coordinator {
			validSignedOrders = append(validSignedOrders, signedOrder)
			continue
		}

		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.WithError(err).WithField("signedOrder", signedOrder).Error("Computing the orderHash failed unexpectedly")
		}
		endpoint, ok := o.cachedFeeRecipientToEndpoint[signedOrder.FeeRecipientAddress]
		if !ok {
			opts := &bind.CallOpts{
				Pending: false,
				Context: ctx,
			}
			var err error
			// Look-up the coordinator endpoint in the CoordinatorRegistry by the order's `feeRecipientAddress`
			endpoint, err = o.coordinatorRegistry.GetCoordinatorEndpoint(opts, signedOrder.FeeRecipientAddress)
			if err != nil {
				log.WithFields(log.Fields{
					"error":               err.Error(),
					"feeRecipientAddress": signedOrder.FeeRecipientAddress,
				}).Info("GetCoordinatorEndpoint request failed")
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        MeshError,
					Status:      ROEthRPCRequestFailed,
				})
				continue
			}
			// CoordinatorRegistry lookup returns empty string if endpoint not found for the feeRecipientAddress
			if endpoint == "" {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        CoordinatorError,
					Status:      ROCoordinatorEndpointNotFound,
				})
				continue
			}
		}
		existingOrders, ok := endpointToSignedOrders[endpoint]
		if !ok {
			endpointToSignedOrders[endpoint] = []*zeroex.SignedOrder{signedOrder}
		} else {
			endpointToSignedOrders[endpoint] = append(existingOrders, signedOrder)
		}
	}

	for endpoint, signedOrders := range endpointToSignedOrders {
		orderHashToSignedOrder := map[common.Hash]*zeroex.SignedOrder{}
		orderHashes := []common.Hash{}
		for _, signedOrder := range signedOrders {
			orderHash, err := signedOrder.ComputeOrderHash()
			if err != nil {
				log.WithError(err).WithField("signedOrder", signedOrder).Error("Computing the orderHash failed unexpectedly")
			}
			orderHashToSignedOrder[orderHash] = signedOrder
			orderHashes = append(orderHashes, orderHash)
		}
		payload := &bytes.Buffer{}
		err := json.NewEncoder(payload).Encode(orderHashes)
		if err != nil {
			log.WithError(err).WithField("orderHashes", orderHashes).Error("Unable to marshal `orderHashes` into JSON")
		}
		// Check if the orders have been soft-cancelled by querying the Coordinator server
		requestURL := fmt.Sprintf("%s/v1/soft_cancels?networkId=%d", endpoint, o.chainID)
		resp, err := http.Post(requestURL, "application/json", payload)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"endpoint":  endpoint,
				"requstURL": requestURL,
				"payload":   orderHashes,
			}).Warn("failed to send request to Coordinator server")
			for orderHash, signedOrder := range orderHashToSignedOrder {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        MeshError,
					Status:      ROCoordinatorRequestFailed,
				})
			}
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"endpoint":   endpoint,
				"statusCode": resp.StatusCode,
				"requstURL":  requestURL,
			}).Warn("Failed to read body received from Coordinator server")
			for orderHash, signedOrder := range orderHashToSignedOrder {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        MeshError,
					Status:      ROCoordinatorRequestFailed,
				})
			}
			continue
		}
		if resp.StatusCode != 200 {
			log.WithFields(map[string]interface{}{
				"endpoint":   endpoint,
				"statusCode": resp.StatusCode,
				"requstURL":  requestURL,
				"body":       string(body),
			}).Warn("Got non-200 status code from Coordinator server")
			for orderHash, signedOrder := range orderHashToSignedOrder {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        MeshError,
					Status:      ROCoordinatorRequestFailed,
				})
			}
			continue
		}
		var response softCancelResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.WithFields(map[string]interface{}{
				"endpoint":   endpoint,
				"statusCode": resp.StatusCode,
				"requstURL":  requestURL,
				"body":       string(body),
			}).Warn("Unable to unmarshal body returned from Coordinator server")
			for orderHash, signedOrder := range orderHashToSignedOrder {
				rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: signedOrder,
					Kind:        MeshError,
					Status:      ROCoordinatorRequestFailed,
				})
			}
			continue
		}
		softCancelledOrderHashes := response.OrderHashes
		softCancelledOrderHashMap := map[common.Hash]interface{}{}
		for _, orderHash := range softCancelledOrderHashes {
			softCancelledOrderHashMap[orderHash] = struct{}{}
			signedOrder := orderHashToSignedOrder[orderHash]
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        MeshError,
				Status:      ROCoordinatorSoftCancelled,
			})
		}
		for orderHash, signedOrder := range orderHashToSignedOrder {
			// If order hasn't been soft-cancelled
			if _, ok := softCancelledOrderHashMap[orderHash]; !ok {
				validSignedOrders = append(validSignedOrders, signedOrder)
			}
		}
	}
	return validSignedOrders, rejectedOrderInfos
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

		if len(signedOrder.MakerFeeAssetData) != 0 {
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
		if len(signedOrder.TakerFeeAssetData) != 0 {
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

		isSupportedSignature := isSupportedSignature(signedOrder.Signature, orderHash)
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

// emptyGetOrderRelevantStatesCallDataByteLength is all the boilerplate ABI encoding required when calling
// `getOrderRelevantStates` that does not include the encoded SignedOrder. By subtracting this amount from the
// calldata length returned from encoding a call to `getOrderRelevantStates` involving a single SignedOrder, we
// get the number of bytes taken up by the SignedOrder alone.
// i.e.: len(`"0x7f46448d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`)
const emptyGetOrderRelevantStatesCallDataByteLength = 268

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

func (o *OrderValidator) computeABIEncodedSignedOrderByteLength(signedOrder *zeroex.SignedOrder) (int, error) {
	trimmedOrder := signedOrder.Trim()
	data, err := o.devUtilsABI.Pack(
		"getOrderRelevantStates",
		[]wrappers.LibOrderOrder{trimmedOrder},
		[][]byte{signedOrder.Signature},
	)
	if err != nil {
		return 0, err
	}
	dataBytes := hexutil.Bytes(data)
	encodedData, err := json.Marshal(dataBytes)
	if err != nil {
		return 0, err
	}
	encodedSignedOrderByteLength := len(encodedData) - emptyGetOrderRelevantStatesCallDataByteLength
	return encodedSignedOrderByteLength, nil
}

func isSupportedSignature(signature []byte, orderHash common.Hash) bool {
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
