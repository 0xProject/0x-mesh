package zeroex

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

// MainnetOrderValidatorAddress is the mainnet OrderValidator contract address
var MainnetOrderValidatorAddress = common.HexToAddress("0x9463e518dea6810309563c81d5266c1b1d149138")

// GanacheOrderValidatorAddress is the ganache snapshot OrderValidator contract address
var GanacheOrderValidatorAddress = common.HexToAddress("0x32eecaf51dfea9618e9bc94e9fbfddb1bbdcba15")

// The most orders we can validate in a single eth_call without having the request timeout
const chunkSize = 500

// The context timeout length to use for requests to getOrdersAndTradersInfoTimeout
const getOrdersAndTradersInfoTimeout = 15 * time.Second

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
	OrderHash   common.Hash
	SignedOrder *SignedOrder
	Message     string
	Kind        RejectedOrderKind
	Code        RejectedOrderCode
}

// AcceptedOrderInfo represents an fillable order and how much it could be filled for
type AcceptedOrderInfo struct {
	OrderHash                common.Hash
	SignedOrder              *SignedOrder
	FillableTakerAssetAmount *big.Int
}

// RejectedOrderCode enumerates all the unique reasons for an orders rejection
type RejectedOrderCode string

// RejectedOrderCode values
const (
	// TODO(fabio): Remove ROInvalid once we remove switch
	ROInvalid                 = RejectedOrderCode("INVALID_REJECTED_ORDER_CODE")
	RORequestFailed           = RejectedOrderCode("ETH_RPC_REQUEST_FAILED")
	ROInvalidMakerAssetAmount = RejectedOrderCode("ORDER_WITH_INVALID_MAKER_ASSET_AMOUNT")
	ROInvalidTakerAssetAmount = RejectedOrderCode("ORDER_WITH_INVALID_TAKER_ASSET_AMOUNT")
	ROExpired                 = RejectedOrderCode("ORDER_EXPIRED")
	ROFullyFilled             = RejectedOrderCode("ORDER_FULLY_FILLED")
	ROCancelled               = RejectedOrderCode("ORDER_CANCELLED")
	ROSignatureInvalid        = RejectedOrderCode("ORDER_SIGNATURE_INVALID")
	ROInvalidMakerAssetData   = RejectedOrderCode("ORDER_WITH_INVALID_MAKER_ASSET_DATA")
	ROInvalidTakerAssetData   = RejectedOrderCode("ORDER_WITH_INVALID_TAKER_ASSET_DATA")
	// ROUnfunded is a catch-all for when either the maker or taker have insufficient
	// balance or allowance set to fullfil the order.
	ROUnfunded = RejectedOrderCode("ORDER_UNFUNDED")
)

// ConvertRejectOrderCodeToOrderEventKind converts an RejectOrderCode to an OrderEventKind type
func ConvertRejectOrderCodeToOrderEventKind(rejectedOrderCode RejectedOrderCode) (OrderEventKind, bool) {
	switch rejectedOrderCode {
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

// GetMessageForRejectOrderCode returns the corresponding message for the RejectedOrderCode
func GetMessageForRejectOrderCode(rejectedOrderCode RejectedOrderCode) string {
	switch rejectedOrderCode {
	case RORequestFailed:
		return "network request to Ethereum RPC endpoint failed"
	case ROExpired:
		return "order expired"
	case ROFullyFilled:
		return "order fully filled"
	case ROCancelled:
		return "order cancelled"
	case ROUnfunded:
		return "maker has insufficient balance or allowance for this order to be filled"
	case ROInvalidMakerAssetAmount:
		return "order makerAssetAmount cannot be 0"
	case ROInvalidTakerAssetAmount:
		return "order takerAssetAmount cannot be 0"
	case ROInvalidMakerAssetData:
		return "order makerAssetData must encode a supported assetData type"
	case ROInvalidTakerAssetData:
		return "order makerAssetData must encode a supported assetData type"
	case ROSignatureInvalid:
		return "order signature must be valid"
	default:
		return "invalid RejectedOrderCode found"
	}
}

// RejectedOrderKind enumerates all kinds of reasons an order could be rejected by Mesh
type RejectedOrderKind string

// RejectedOrderKind values
const (
	ZeroExValidation = RejectedOrderKind("ZEROEX_VALIDATION")
	MeshError        = RejectedOrderKind("MESH_ERROR")
)

// ValidationResults defines the validation results returned from BatchValidate
type ValidationResults struct {
	Accepted []*AcceptedOrderInfo
	Rejected []*RejectedOrderInfo
}

// OrderValidator validates 0x orders
type OrderValidator struct {
	orderValidator   *wrappers.OrderValidator
	assetDataDecoder *AssetDataDecoder
}

// NewOrderValidator instantiates a new order validator
func NewOrderValidator(ethClient *ethclient.Client, networkID int) (*OrderValidator, error) {
	contractNameToAddress := constants.NetworkIDToContractAddresses[networkID]
	orderValidator, err := wrappers.NewOrderValidator(contractNameToAddress.OrderValidator, ethClient)
	if err != nil {
		return nil, err
	}
	assetDataDecoder := NewAssetDataDecoder()

	return &OrderValidator{
		orderValidator:   orderValidator,
		assetDataDecoder: assetDataDecoder,
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
	rejectedOrderInfos, offchainValidSignedOrders := o.BatchOffchainValidation(rawSignedOrders)
	validationResults := &ValidationResults{
		Accepted: []*AcceptedOrderInfo{},
		Rejected: rejectedOrderInfos,
	}

	// Chunk into groups of chunkSize signedOrders for each call to the smart contract
	signedOrderChunks := [][]*SignedOrder{}
	for len(offchainValidSignedOrders) > chunkSize {
		signedOrderChunks = append(signedOrderChunks, offchainValidSignedOrders[:chunkSize])
		offchainValidSignedOrders = offchainValidSignedOrders[chunkSize:]
	}
	if len(offchainValidSignedOrders) > 0 {
		signedOrderChunks = append(signedOrderChunks, offchainValidSignedOrders)
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for i, signedOrders := range signedOrderChunks {
		wg.Add(1)
		go func(signedOrders []*SignedOrder, i int) {
			takerAddresses := []common.Address{}
			for _, signedOrder := range offchainValidSignedOrders {
				takerAddresses = append(takerAddresses, signedOrder.TakerAddress)
			}
			orders := []wrappers.OrderWithoutExchangeAddress{}
			for _, signedOrder := range offchainValidSignedOrders {
				orders = append(orders, signedOrder.ConvertToOrderWithoutExchangeAddress())
			}
			signatures := [][]byte{}
			for _, signedOrder := range offchainValidSignedOrders {
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
				// Pass a context with a 15 second timeout to `GetOrdersAndTradersInfo` in order to avoid
				// any one request from taking longer then 15 seconds
				ctx, cancel := context.WithTimeout(context.Background(), getOrdersAndTradersInfoTimeout)
				defer cancel()
				opts := &bind.CallOpts{
					Pending: false,
					Context: ctx,
				}
				results, err := o.orderValidator.GetOrdersAndTradersInfo(opts, orders, signatures, takerAddresses)
				if err != nil {
					log.WithFields(log.Fields{
						"error":     err.Error(),
						"attempt":   b.Attempt(),
						"numOrders": len(orders),
					}).Info("GetOrdersAndTradersInfo request failed")
					d := b.Duration()
					if d == maxDuration {
						<-semaphoreChan
						log.WithFields(log.Fields{
							"error":     err.Error(),
							"numOrders": len(orders),
						}).Warning("Gave up on GetOrdersAndTradersInfo request after backoff limit reached")
						for _, signedOrder := range signedOrders {
							orderHash, err := signedOrder.ComputeOrderHash()
							if err != nil { // Should never happen
								log.WithField("error", err).Panic("Unexpectedly failed to generate orderHash")
								continue
							}
							code := RORequestFailed
							validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
								OrderHash:   orderHash,
								SignedOrder: signedOrder,
								Kind:        MeshError,
								Code:        code,
								Message:     GetMessageForRejectOrderCode(code),
							})
						}
						return // Give up after 4 attempts
					}
					time.Sleep(d)
					continue
				}

				for j, orderInfo := range results.OrdersInfo {
					traderInfo := results.TradersInfo[j]
					isValidSignature := results.IsValidSignature[j]
					orderHash := common.Hash(orderInfo.OrderHash)
					signedOrder := offchainValidSignedOrders[chunkSize*i+j]
					orderStatus := OrderStatus(orderInfo.OrderStatus)
					if !isValidSignature {
						orderStatus = OSSignatureInvalid
					}
					switch orderStatus {
					case OSExpired, OSFullyFilled, OSCancelled, OSSignatureInvalid:
						var code RejectedOrderCode
						switch orderStatus {
						case OSExpired:
							code = ROExpired
						case OSFullyFilled:
							code = ROFullyFilled
						case OSCancelled:
							code = ROCancelled
						case OSSignatureInvalid:
							code = ROSignatureInvalid
						}
						validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
							OrderHash:   orderHash,
							SignedOrder: signedOrder,
							Kind:        ZeroExValidation,
							Code:        code,
							Message:     GetMessageForRejectOrderCode(code),
						})
						continue
					case OSFillable:
						fillableTakerAssetAmount := calculateRemainingFillableTakerAmount(signedOrder, orderInfo, traderInfo)
						if fillableTakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
							code := ROUnfunded // Order with remaining fillable amount of 0 but that is otherwise fillable
							validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
								OrderHash:   orderHash,
								SignedOrder: signedOrder,
								Kind:        ZeroExValidation,
								Code:        code,
								Message:     GetMessageForRejectOrderCode(code),
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

// BatchOffchainValidation performs all off-chain validation checks on a batch of 0x orders.
// These checks include:
// - `MakerAssetAmount` and `TakerAssetAmount` cannot be 0
// - `AssetData` fields contain properly encoded, and currently supported assetData (ERC20 & ERC721 for now)
// - `Signature` contains a properly encoded 0x signature
// - Validate that order isn't expired
// Returns an orderHashToInfo mapping with all invalid orders added to it, and an array of the valid signedOrders
func (o *OrderValidator) BatchOffchainValidation(signedOrders []*SignedOrder) ([]*RejectedOrderInfo, []*SignedOrder) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	offchainValidSignedOrders := []*SignedOrder{}
	for _, signedOrder := range signedOrders {
		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.Panic("Computing the orderHash failed unexpectedly")
		}
		now := big.NewInt(time.Now().Unix())
		if signedOrder.ExpirationTimeSeconds.Cmp(now) == -1 {
			code := ROExpired
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}

		if signedOrder.MakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
			code := ROInvalidMakerAssetAmount
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}
		if signedOrder.TakerAssetAmount.Cmp(big.NewInt(0)) == 0 {
			code := ROInvalidTakerAssetAmount
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}

		isMakerAssetDataSupported := o.isSupportedAssetData(signedOrder.MakerAssetData)
		if !isMakerAssetDataSupported {
			code := ROInvalidMakerAssetData
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}
		isTakerAssetDataSupported := o.isSupportedAssetData(signedOrder.TakerAssetData)
		if !isTakerAssetDataSupported {
			code := ROInvalidTakerAssetData
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}

		isSupportedSignature := isSupportedSignature(signedOrder.Signature, orderHash)
		if !isSupportedSignature {
			code := ROSignatureInvalid
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: signedOrder,
				Kind:        ZeroExValidation,
				Code:        code,
				Message:     GetMessageForRejectOrderCode(code),
			})
			continue
		}

		offchainValidSignedOrders = append(offchainValidSignedOrders, signedOrder)
	}

	return rejectedOrderInfos, offchainValidSignedOrders
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
		// TODO(fabio): Once OrderValidator.sol supports validating orders involving multiAssetData,
		// refactor this to add support.
		return false
	default:
		return false
	}
	return true
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

func calculateRemainingFillableTakerAmount(signedOrder *SignedOrder, orderInfo wrappers.OrderInfo, traderInfo wrappers.TraderInfo) *big.Int {
	minSet := []*big.Int{}

	// Calculate min of balance & allowance of makers makerAsset -> translate into takerAsset amount
	var maxMakerAssetFillAmount *big.Int
	if traderInfo.MakerBalance.Cmp(traderInfo.MakerAllowance) == -1 {
		maxMakerAssetFillAmount = traderInfo.MakerBalance
	} else {
		maxMakerAssetFillAmount = traderInfo.MakerAllowance
	}
	maxTakerAssetFillAmountGivenMakerConstraints := new(big.Int).Div(new(big.Int).Mul(maxMakerAssetFillAmount, signedOrder.TakerAssetAmount), signedOrder.MakerAssetAmount)

	minSet = append(minSet, maxTakerAssetFillAmountGivenMakerConstraints)

	// Calculate min of balance & allowance of maker's ZRX -> translate into takerAsset amount
	if signedOrder.MakerFee.Cmp(big.NewInt(0)) != 0 {
		var makerZRXAvailable *big.Int
		if traderInfo.MakerZrxBalance.Cmp(traderInfo.MakerZrxAllowance) == -1 {
			makerZRXAvailable = traderInfo.MakerZrxBalance
		} else {
			makerZRXAvailable = traderInfo.MakerZrxAllowance
		}
		maxTakerAssetFillAmountGivenMakerZRXConstraints := new(big.Int).Div(new(big.Int).Mul(makerZRXAvailable, signedOrder.TakerAssetAmount), signedOrder.MakerFee)
		minSet = append(minSet, maxTakerAssetFillAmountGivenMakerZRXConstraints)
	}

	// Add the remaining takerAsset fill amount to the minSet
	remainingTakerAssetFillAmount := new(big.Int).Sub(signedOrder.TakerAssetAmount, orderInfo.OrderTakerAssetFilledAmount)
	minSet = append(minSet, remainingTakerAssetFillAmount)

	var maxTakerAssetFillAmount *big.Int
	for _, minVal := range minSet {
		if maxTakerAssetFillAmount == nil || maxTakerAssetFillAmount.Cmp(minVal) != -1 {
			maxTakerAssetFillAmount = minVal
		}
	}

	return maxTakerAssetFillAmount
}
