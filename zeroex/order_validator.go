package zeroex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

var getOrderRelevantStatesNumParams = 2

var getOrderRelevantStatesMethodName = "getOrderRelevantStates"

// MainnetOrderValidatorAddress is the mainnet OrderValidator contract address
var MainnetOrderValidatorAddress = common.HexToAddress("0x9463e518dea6810309563c81d5266c1b1d149138")

// GanacheOrderValidatorAddress is the ganache snapshot OrderValidator contract address
var GanacheOrderValidatorAddress = common.HexToAddress("0x32eecaf51dfea9618e9bc94e9fbfddb1bbdcba15")

// 1800 (data) - 8 (methodID) - ((64 (ABI head) + 64 (ABI array len val)) * getOrderRelevantStatesNumParams)
var encodedNonMultiAssetOrderSizeBytes = int64(1536)

// The context timeout length to use for requests to getOrdersAndTradersInfoTimeout
const getOrdersAndTradersInfoTimeout = 15 * time.Second

// The context timeout length to use for requests to getCoordinatorEndpoint
const getCoordinatorEndpointTimeout = 10 * time.Second

// Specifies the max number of eth_call requests we want to make concurrently.
// Additional requests will block until an ongoing request has completed.
const concurrencyLimit = 5

// OrderInfo represents the order information emitted from Mesh
type OrderInfo struct {
	OrderHash                common.Hash
	SignedOrder              *SignedOrder
	FillableTakerAssetAmount *big.Int
	OrderStatus              OrderStatus
	// The hash of the Ethereum transaction that caused the order status to change
	TxHash common.Hash
}

// RejectedOrderInfo encapsulates all the needed information to understand _why_ a 0x order
// was rejected (i.e. did not pass) order validation. Since there are many potential reasons, some
// Mesh-specific, others 0x-specific and others due to external factors (i.e., network
// disruptions, etc...), we categorize them into `Kind`s and uniquely identify the reasons for
// machines with a `Code`
type RejectedOrderInfo struct {
	OrderHash   common.Hash         `json:"orderHash"`
	SignedOrder *SignedOrder        `json:"signedOrder"`
	Kind        RejectedOrderKind   `json:"kind"`
	Status      RejectedOrderStatus `json:"status"`
}

// AcceptedOrderInfo represents an fillable order and how much it could be filled for
type AcceptedOrderInfo struct {
	OrderHash                common.Hash  `json:"orderHash"`
	SignedOrder              *SignedOrder `json:"signedOrder"`
	FillableTakerAssetAmount *big.Int     `json:"fillableTakerAssetAmount"`
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
		Message: "order already expired",
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
	ROInvalidTakerAssetData = RejectedOrderStatus{
		Code:    "OrderHasInvalidTakerAssetData",
		Message: "order takerAssetData must encode a supported assetData type",
	}
	ROInvalidSignature = RejectedOrderStatus{
		Code:    "OrderHasInvalidSignature",
		Message: "order signature must be valid",
	}
)

// ConvertRejectOrderCodeToOrderEventKind converts an RejectOrderCode to an OrderEventKind type
func ConvertRejectOrderCodeToOrderEventKind(rejectedOrderStatus RejectedOrderStatus) (OrderEventKind, bool) {
	switch rejectedOrderStatus {
	case ROExpired:
		return EKOrderExpired, true
	case ROFullyFilled:
		return EKOrderFullyFilled, true
	case ROCancelled:
		return EKOrderCancelled, true
	case ROUnfunded:
		return EKOrderBecameUnfunded, true
	default:
		// Catch-all returns Invalid OrderEventKind
		return EKInvalid, false
	}
}

// RejectedOrderKind enumerates all kinds of reasons an order could be rejected by Mesh
type RejectedOrderKind string

// RejectedOrderKind values
const (
	ZeroExValidation = RejectedOrderKind("ZEROEX_VALIDATION")
	MeshError        = RejectedOrderKind("MESH_ERROR")
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
	maxRequestContentLength      int64
	orderValidationUtilsABI      abi.ABI
	orderValidationUtils         *wrappers.OrderValidationUtils
	coordinatorRegistry          *wrappers.CoordinatorRegistry
	assetDataDecoder             *AssetDataDecoder
	networkID                    int
	cachedFeeRecipientToEndpoint map[common.Address]string
	contractAddresses            ethereum.ContractAddresses
}

// NewOrderValidator instantiates a new order validator
func NewOrderValidator(ethClient *ethclient.Client, networkID int, maxRequestContentLength int64) (*OrderValidator, error) {
	contractAddresses, err := ethereum.GetContractAddressesForNetworkID(networkID)
	if err != nil {
		return nil, err
	}
	orderValidationUtilsABI, err := abi.JSON(strings.NewReader(wrappers.OrderValidationUtilsABI))
	if err != nil {
		return nil, err
	}
	orderValidationUtils, err := wrappers.NewOrderValidationUtils(contractAddresses.OrderValidationUtils, ethClient)
	if err != nil {
		return nil, err
	}
	coordinatorRegistry, err := wrappers.NewCoordinatorRegistry(contractAddresses.CoordinatorRegistry, ethClient)
	if err != nil {
		return nil, err
	}
	assetDataDecoder := NewAssetDataDecoder()

	return &OrderValidator{
		maxRequestContentLength:      maxRequestContentLength,
		orderValidationUtilsABI:      orderValidationUtilsABI,
		orderValidationUtils:         orderValidationUtils,
		coordinatorRegistry:          coordinatorRegistry,
		assetDataDecoder:             assetDataDecoder,
		networkID:                    networkID,
		cachedFeeRecipientToEndpoint: map[common.Address]string{},
		contractAddresses:            contractAddresses,
	}, nil
}

// BatchValidate retrieves all the information needed to validate the supplied orders.
// It splits the orders into chunks of `chunkSize`, and makes no more then `concurrencyLimit`
// requests concurrently. If a request fails, re-attempt it up to four times before giving up.
// If it some requests fail, this method still returns whatever order information it was able to
// retrieve.
func (o *OrderValidator) BatchValidate(rawSignedOrders []*SignedOrder) *ValidationResults {
	if len(rawSignedOrders) == 0 {
		return &ValidationResults{}
	}
	offchainValidSignedOrders, rejectedOrderInfos := o.BatchOffchainValidation(rawSignedOrders)
	validationResults := &ValidationResults{
		Accepted: []*AcceptedOrderInfo{},
		Rejected: rejectedOrderInfos,
	}

	// Validate Coordinator orders for soft-cancels
	signedOrders, coordinatorRejectedOrderInfos := o.batchValidateSoftCancelled(offchainValidSignedOrders)
	for _, rejectedOrderInfo := range coordinatorRejectedOrderInfos {
		validationResults.Rejected = append(validationResults.Rejected, rejectedOrderInfo)
	}

	// Group into chunks of signedOrders for each call to the smart contract
	numBasicOrdersEncodable := o.computeNumBasicOrdersEncodable()
	signedOrderChunks := [][]*SignedOrder{}
	for len(signedOrders) > 0 {
		tentativeAmount := numBasicOrdersEncodable
		numOrders := int64(len(signedOrders))
		if numOrders < numBasicOrdersEncodable {
			tentativeAmount = numOrders
		}
		// If the tentativeAmount which is <= numBasicOrdersEncodable, is below the max content length,
		// we simply create a chunk of that size.
		if o.isBelowMaxContentLength(signedOrders[:tentativeAmount]) {
			signedOrderChunks = append(signedOrderChunks, signedOrders[:tentativeAmount])
			signedOrders = signedOrders[tentativeAmount:]
			continue
		}
		// If it's above the max content length, we fall-back to a binary-search between 1-tentativeAmount
		// to find a chunk size that is below the max content amount
		chunkSize := o.findChunkSize(signedOrders[:tentativeAmount])
		signedOrderChunks = append(signedOrderChunks, signedOrders[:chunkSize])
		signedOrders = signedOrders[chunkSize:]
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for i, signedOrders := range signedOrderChunks {
		wg.Add(1)
		go func(signedOrders []*SignedOrder, i int) {
			orders := []wrappers.OrderWithoutExchangeAddress{}
			for _, signedOrder := range signedOrders {
				orders = append(orders, signedOrder.ConvertToOrderWithoutExchangeAddress())
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
				// Pass a context with a 15 second timeout to `GetOrderRelevantStates` in order to avoid
				// any one request from taking longer then 15 seconds
				ctx, cancel := context.WithTimeout(context.Background(), getOrdersAndTradersInfoTimeout)
				defer cancel()
				opts := &bind.CallOpts{
					Pending: false,
					Context: ctx,
				}

				results, err := o.orderValidationUtils.GetOrderRelevantStates(opts, orders, signatures)
				if err != nil {
					log.WithFields(log.Fields{
						"error":     err.Error(),
						"attempt":   b.Attempt(),
						"numOrders": len(orders),
					}).Info("GetOrderRelevantStates request failed")
					d := b.Duration()
					if d == maxDuration {
						<-semaphoreChan
						log.WithFields(log.Fields{
							"error":     err.Error(),
							"numOrders": len(orders),
						}).Warning("Gave up on GetOrderRelevantStates request after backoff limit reached")
						for _, signedOrder := range signedOrders {
							orderHash, err := signedOrder.ComputeOrderHash()
							if err != nil { // Should never happen
								log.WithField("error", err).Panic("Unexpectedly failed to generate orderHash")
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
					orderStatus := OrderStatus(orderInfo.OrderStatus)
					if !isValidSignature {
						orderStatus = OSSignatureInvalid
					}
					switch orderStatus {
					case OSExpired, OSFullyFilled, OSCancelled, OSSignatureInvalid:
						var status RejectedOrderStatus
						switch orderStatus {
						case OSExpired:
							status = ROExpired
						case OSFullyFilled:
							status = ROFullyFilled
						case OSCancelled:
							status = ROCancelled
						case OSSignatureInvalid:
							status = ROInvalidSignature
						}
						validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
							OrderHash:   orderHash,
							SignedOrder: signedOrder,
							Kind:        ZeroExValidation,
							Status:      status,
						})
						continue
					case OSFillable:
						if fillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
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
func (o *OrderValidator) batchValidateSoftCancelled(signedOrders []*SignedOrder) ([]*SignedOrder, []*RejectedOrderInfo) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	validSignedOrders := []*SignedOrder{}

	endpointToSignedOrders := map[string][]*SignedOrder{}
	for _, signedOrder := range signedOrders {
		if signedOrder.SenderAddress != o.contractAddresses.Coordinator {
			validSignedOrders = append(validSignedOrders, signedOrder)
			continue
		}

		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.WithField("signedOrder", signedOrder).Panic("Computing the orderHash failed unexpectedly")
		}
		endpoint, ok := o.cachedFeeRecipientToEndpoint[signedOrder.FeeRecipientAddress]
		if !ok {
			ctx, cancel := context.WithTimeout(context.Background(), getCoordinatorEndpointTimeout)
			defer cancel()
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
			endpointToSignedOrders[endpoint] = []*SignedOrder{signedOrder}
		} else {
			endpointToSignedOrders[endpoint] = append(existingOrders, signedOrder)
		}
	}

	for endpoint, signedOrders := range endpointToSignedOrders {
		orderHashToSignedOrder := map[common.Hash]*SignedOrder{}
		orderHashes := []common.Hash{}
		for _, signedOrder := range signedOrders {
			orderHash, err := signedOrder.ComputeOrderHash()
			if err != nil {
				log.WithField("signedOrder", signedOrder).Panic("Computing the orderHash failed unexpectedly")
			}
			orderHashToSignedOrder[orderHash] = signedOrder
			orderHashes = append(orderHashes, orderHash)
		}
		payload := &bytes.Buffer{}
		err := json.NewEncoder(payload).Encode(orderHashes)
		if err != nil {
			log.WithField("orderHashes", orderHashes).Panic("Unable to marshal `orderHashes` into JSON")
		}
		// Check if the orders have been soft-cancelled by querying the Coordinator server
		requestURL := fmt.Sprintf("%s/v1/soft_cancels?networkId=%d", endpoint, o.networkID)
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
// Returns an orderHashToInfo mapping with all invalid orders added to it, and an array of the valid signedOrders
func (o *OrderValidator) BatchOffchainValidation(signedOrders []*SignedOrder) ([]*SignedOrder, []*RejectedOrderInfo) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	offchainValidSignedOrders := []*SignedOrder{}
	for _, signedOrder := range signedOrders {
		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.WithField("signedOrder", signedOrder).Panic("Computing the orderHash failed unexpectedly")
		}
		now := big.NewInt(time.Now().Unix())
		if signedOrder.ExpirationTimeSeconds.Cmp(now) == -1 {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Status:      ROExpired,
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
		var decodedAssetData ERC20AssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "ERC721Token":
		var decodedAssetData ERC721AssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	case "MultiAsset":
		var decodedAssetData MultiAssetData
		err := o.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return false
		}
	default:
		return false
	}
	return true
}

// findChunkSize uses binary-search to find a chunk size that is less then the max content length
func (o *OrderValidator) findChunkSize(signedOrders []*SignedOrder) int64 {
	// Determine the lowest and highest possible signedOrders limits to binary search in between
	lo := int64(0)
	hi := int64(len(signedOrders))
	min := int64(1)

	// Execute the binary search and hone in on an executable payload size of signedOrders
	maxBinarySearchDepth := 7
	depth := 0
	for lo+1 < hi {
		depth++
		mid := (hi + lo) / 2
		if o.isBelowMaxContentLength(signedOrders[:mid]) {
			lo = mid
		} else {
			hi = mid
		}
		if depth == maxBinarySearchDepth {
			break
		}
	}

	if !o.isBelowMaxContentLength(signedOrders[:hi]) {
		// Reject the signedOrders if it still fails at the lowest allowance
		if hi == min {
			// We don't ever expect this to happen since the maxOrderSizeInBytes is much smaller then the
			// maxRequestContentLength. This error requires a single order to be larger then the
			// maxRequestContentLength, which would only happen if someone set the maxRequestContentLength
			// to a tiny number.
			err := errors.New("A single signedOrder was found that was larger then the maxRequestContentLength. Increase the maxRequestContentLength")
			log.WithField("SignedOrder", signedOrders[0]).Panic(err)
		}
		hi = lo
	}

	return hi
}

// isBelowMaxContentLength checks if the signedOrders once encoded into a payload is
// below the max content length.
func (o *OrderValidator) isBelowMaxContentLength(signedOrders []*SignedOrder) bool {
	orders := []wrappers.OrderWithoutExchangeAddress{}
	for _, signedOrder := range signedOrders {
		orders = append(orders, signedOrder.ConvertToOrderWithoutExchangeAddress())
	}
	signatures := [][]byte{}
	for _, signedOrder := range signedOrders {
		signatures = append(signatures, signedOrder.Signature)
	}

	payloadSize, err := o.computePayloadSize(getOrderRelevantStatesMethodName, orders, signatures)
	if err != nil {
		return false
	}
	return payloadSize < o.maxRequestContentLength
}

func (o *OrderValidator) computePayloadSize(method string, params ...interface{}) (int64, error) {
	data, err := o.orderValidationUtilsABI.Pack(method, params...)
	if err != nil {
		return int64(0), err
	}
	payload := map[string]interface{}{
		"id":      2,
		"jsonrpc": "2.0",
		"method":  "eth_call",
	}
	payload["params"] = []interface{}{
		map[string]interface{}{
			"data": hexutil.Bytes(data),
			"from": constants.NullAddress,
			"to":   constants.NullAddress,
		},
		"latest",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return int64(0), err
	}
	return int64(len(payloadBytes)), nil
}

// computeNumBasicOrdersEncodable calculates the number of "basic" (e.g., non-multiAsset) 0x orders
// that can be sent in a payload of max content size. Unlike "basic" orders, those involving the
// MultiAssetProxy can be of arbitrary size.
func (o *OrderValidator) computeNumBasicOrdersEncodable() int64 {
	// We can safely ignore the following error because we are passing in known empty arrays as params
	restOfPayload, _ := o.computePayloadSize(getOrderRelevantStatesMethodName, []wrappers.OrderWithoutExchangeAddress{}, [][]byte{})
	numBasicOrders := (o.maxRequestContentLength - (int64(getOrderRelevantStatesNumParams) * 128) - restOfPayload) / encodedNonMultiAssetOrderSizeBytes
	return numBasicOrders
}

func isSupportedSignature(signature []byte, orderHash common.Hash) bool {
	signatureType := SignatureType(signature[len(signature)-1])

	switch signatureType {
	case IllegalSignature:
	case InvalidSignature:
		return false

	case EIP712Signature:
		if len(signature) != 66 {
			return false
		}
		// TODO(fabio): Do further validation by splitting into r,s,v and do ECRecover

	case EthSignSignature:
		if len(signature) != 66 {
			return false
		}
		// TODO(fabio): Do further validation by splitting into r,s,v, add prefix to hash
		// and do ECRecover

	case ValidatorSignature:
		if len(signature) < 21 {
			return false
		}

	case WalletSignature:
	case PreSignedSignature:
		return true

	default:
		return false

	}

	return true
}
