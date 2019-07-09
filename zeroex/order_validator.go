package zeroex

import (
	"bytes"
	"context"
	"encoding/json"
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

var getOrderRelevantStatesMethodName = "getOrderRelevantStates"

// MainnetOrderValidatorAddress is the mainnet OrderValidator contract address
var MainnetOrderValidatorAddress = common.HexToAddress("0x9463e518dea6810309563c81d5266c1b1d149138")

// GanacheOrderValidatorAddress is the ganache snapshot OrderValidator contract address
var GanacheOrderValidatorAddress = common.HexToAddress("0x32eecaf51dfea9618e9bc94e9fbfddb1bbdcba15")

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
	maxRequestContentLength      int
	orderValidationUtilsABI      abi.ABI
	orderValidationUtils         *wrappers.OrderValidationUtils
	coordinatorRegistry          *wrappers.CoordinatorRegistry
	assetDataDecoder             *AssetDataDecoder
	networkID                    int
	cachedFeeRecipientToEndpoint map[common.Address]string
	contractAddresses            ethereum.ContractAddresses
}

// NewOrderValidator instantiates a new order validator
func NewOrderValidator(ethClient *ethclient.Client, networkID int, maxRequestContentLength int) (*OrderValidator, error) {
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

	signedOrderChunks := [][]*SignedOrder{}
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

func (o *OrderValidator) computeDataSize(method string, params ...interface{}) (int, error) {
	data, err := o.orderValidationUtilsABI.Pack(method, params...)
	if err != nil {
		return 0, err
	}
	dataBytes := hexutil.Bytes(data)
	encodedData, err := json.Marshal(dataBytes)
	if err != nil {
		return 0, err
	}
	return len(encodedData), nil
}

// abiEncodingOverheadForMethodCall includes all the bytes in the encoded `calldata` field that isn't specific to
// the orders being encoded, but are overhead bytes required to call the getOrderRelevantStates method. As such
// this is what needs to get subtracted from the encoded `calldata` in order to remain with the bytes taken up by
// an order.
// abiEncodingOverheadForMethodCall = quotes + hexPrefix + methodID + (numParams * (abiHead + abiArrayLen))
// Source: https://solidity.readthedocs.io/en/v0.5.10/abi-spec.html#use-of-dynamic-types
const abiEncodingOverheadForMethodCall = 270


// jsonRPCPayloadSize is the number of bytes occupied by the default call to `getOrderRelevantStates` with 0 signedOrders
// passed in. The `data` still includes the methodID, and encodes the number of params, and both param array lengths of 0
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
const jsonRPCPayloadSize  = 444

// computeOptimalChunkSizes splits the signedOrders into chunks where the payload size of each chunk
// is before the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizes(signedOrders []*SignedOrder) []int {

	chunkSizes := []int{}

	payloadSize := jsonRPCPayloadSize
	nextChunkSize := 0
	for _, signedOrder := range signedOrders {
		orderWithExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
		encodedDataSize, _ := o.computeDataSize(
			getOrderRelevantStatesMethodName,
			[]wrappers.OrderWithoutExchangeAddress{orderWithExchangeAddress},
			[][]byte{signedOrder.Signature},
		)
		encodedOrderSize := encodedDataSize - abiEncodingOverheadForMethodCall
		if payloadSize+encodedOrderSize < o.maxRequestContentLength {
			payloadSize += encodedOrderSize
			nextChunkSize++
		} else {
			if nextChunkSize == 0 {
				// This case should never be hit since we enforce that EthereumRPCMaxContentLength >= maxOrderSizeInBytes
				log.WithField("signedOrder", signedOrder).Panic("EthereumRPCMaxContentLength is set so low, a single 0x order cannot fit beneath the payload limit")
			}
			chunkSizes = append(chunkSizes, nextChunkSize)
			nextChunkSize = 1
			payloadSize = jsonRPCPayloadSize + encodedOrderSize
		}
	}
	if nextChunkSize != 0 {
		chunkSizes = append(chunkSizes, nextChunkSize)
	}

	return chunkSizes
}

/**
- Call `o.orderValidationUtilsABI.Pack` on every order one at a time.
- Subtract common elements of result (methodID, head, len)
- Add 'em to total that starts with "restOfPayload" + (methodId, head, len)
- Once total bytes as close to max as possible, take that many orders and ABI encode them into final payload
*/

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
