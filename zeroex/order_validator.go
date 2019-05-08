package zeroex

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

var nullAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// MainnetOrderValidatorAddress is the mainnet OrderValidator contract address
var MainnetOrderValidatorAddress = common.HexToAddress("0x9463e518dea6810309563c81d5266c1b1d149138")

// GanacheOrderValidatorAddress is the ganache snapshot OrderValidator contract address
var GanacheOrderValidatorAddress = common.HexToAddress("0x32eecaf51dfea9618e9bc94e9fbfddb1bbdcba15")

// The most orders we can validate in a single eth_call without having the request timeout
const chunkSize = 500

// Specifies the max number of eth_call requests we want to make concurrently.
// Additional requests will block until an ongoing request has completed.
const concurrencyLimit = 5

type getOrdersAndTradersInfoParams struct {
	TakerAddresses []common.Address
	Orders         []wrappers.OrderWithoutExchangeAddress
}

// OrderInfo represents the order information returned from OrderValidator methods.
type OrderInfo struct {
	OrderHash                common.Hash
	SignedOrder              *SignedOrder
	FillableTakerAssetAmount *big.Int
	OrderStatus              OrderStatus
}

// OrderValidator validates 0x orders
type OrderValidator struct {
	orderValidator *wrappers.OrderValidator
}

// NewOrderValidator instantiates a new order validator
func NewOrderValidator(orderValidatorAddress common.Address, ethClient *ethclient.Client) (*OrderValidator, error) {
	orderValidator, err := wrappers.NewOrderValidator(orderValidatorAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &OrderValidator{
		orderValidator: orderValidator,
	}, nil
}

// BatchValidate retrieves all the information needed to validate the supplied orders.
// It splits the orders into chunks of `chunkSize`, and makes no more then `concurrencyLimit`
// requests concurrently. If a request fails, re-attempt it up to four times before giving up.
// If it some requests fail, this method still returns whatever order information it was able to
// retrieve.
func (o *OrderValidator) BatchValidate(signedOrders []*SignedOrder) map[common.Hash]OrderInfo {
	takerAddresses := []common.Address{}
	for i := 0; i < len(signedOrders); i++ {
		takerAddresses = append(takerAddresses, signedOrders[i].TakerAddress)
	}
	orders := []wrappers.OrderWithoutExchangeAddress{}
	for i := 0; i < len(signedOrders); i++ {
		orders = append(orders, signedOrders[i].ConvertToOrderWithoutExchangeAddress())
	}

	// Chunk into groups of chunkSize orders/takerAddresses for each call
	chunks := []getOrdersAndTradersInfoParams{}
	for len(orders) > chunkSize {
		chunks = append(chunks, getOrdersAndTradersInfoParams{
			TakerAddresses: takerAddresses[:chunkSize],
			Orders:         orders[:chunkSize],
		})
		takerAddresses = takerAddresses[chunkSize:]
		orders = orders[chunkSize:]
	}
	if len(orders) > 0 {
		chunks = append(chunks, getOrdersAndTradersInfoParams{
			TakerAddresses: takerAddresses,
			Orders:         orders,
		})
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	orderHashToInfo := map[common.Hash]OrderInfo{}
	wg := &sync.WaitGroup{}
	for i, params := range chunks {
		wg.Add(1)
		go func(params getOrdersAndTradersInfoParams, i int) {
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
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				opts := &bind.CallOpts{
					Pending: false,
					Context: ctx,
				}
				results, err := o.orderValidator.GetOrdersAndTradersInfo(opts, params.Orders, params.TakerAddresses)
				if err != nil {
					log.WithFields(log.Fields{
						"err":            err.Error(),
						"attempt":        b.Attempt(),
						"orders":         params.Orders,
						"takerAddresses": params.TakerAddresses,
					}).Info("GetOrdersAndTradersInfo request failed")
					d := b.Duration()
					if d == maxDuration {
						<-semaphoreChan
						return // Give up after 4 attempts
					}
					time.Sleep(d)
					continue
				}

				for j, orderInfo := range results.OrdersInfo {
					traderInfo := results.TradersInfo[j]
					orderHash := common.Hash(orderInfo.OrderHash)
					signedOrder := signedOrders[chunkSize*i+j]
					switch OrderStatus(orderInfo.OrderStatus) {
					// TODO(fabio): A future optimization would be to check that both the maker & taker
					// amounts are non-zero locally rather then wait for the RPC call to catch these two
					// failure cases.
					case InvalidMakerAssetAmount, InvalidTakerAssetAmount, Expired, FullyFilled, Cancelled:
						orderHashToInfo[orderHash] = OrderInfo{
							OrderHash:                orderHash,
							SignedOrder:              signedOrder,
							FillableTakerAssetAmount: big.NewInt(0),
							OrderStatus:              OrderStatus(orderInfo.OrderStatus),
						}
						continue
					case Fillable:
						orderHashToInfo[orderHash] = OrderInfo{
							OrderHash:                orderHash,
							SignedOrder:              signedOrder,
							FillableTakerAssetAmount: calculateRemainingFillableTakerAmount(signedOrder, orderInfo, traderInfo),
							OrderStatus:              OrderStatus(orderInfo.OrderStatus),
						}
						continue
					}
				}

				<-semaphoreChan
				return
			}
		}(params, i)
	}

	wg.Wait()
	return orderHashToInfo
}

func calculateRemainingFillableTakerAmount(signedOrder *SignedOrder, orderInfo wrappers.OrderInfo, traderInfo wrappers.TraderInfo) *big.Int {
	minSet := []*big.Int{}

	// Calculate min of balance & allowance of taker's takerAsset
	if signedOrder.TakerAddress != nullAddress {
		var maxTakerAssetFillAmountGivenTakerConstraints *big.Int
		if traderInfo.TakerBalance.Cmp(traderInfo.TakerAllowance) == -1 {
			maxTakerAssetFillAmountGivenTakerConstraints = traderInfo.TakerBalance
		} else {
			maxTakerAssetFillAmountGivenTakerConstraints = traderInfo.TakerAllowance
		}
		minSet = append(minSet, maxTakerAssetFillAmountGivenTakerConstraints)
	}

	// Calculate min of balance & allowance of makers makerAsset -> translate into takerAsset amount
	var maxMakerAssetFillAmount *big.Int
	if traderInfo.MakerBalance.Cmp(traderInfo.MakerAllowance) == -1 {
		maxMakerAssetFillAmount = traderInfo.MakerBalance
	} else {
		maxMakerAssetFillAmount = traderInfo.MakerAllowance
	}
	maxTakerAssetFillAmountGivenMakerConstraints := new(big.Int).Div(new(big.Int).Mul(maxMakerAssetFillAmount, signedOrder.TakerAssetAmount), signedOrder.MakerAssetAmount)

	minSet = append(minSet, maxTakerAssetFillAmountGivenMakerConstraints)

	// Calculate min of balance & allowance of taker's ZRX -> translate into takerAsset amount
	if signedOrder.TakerFee.Cmp(big.NewInt(0)) != 0 {
		var takerZRXAvailable *big.Int
		if traderInfo.TakerZrxBalance.Cmp(traderInfo.TakerZrxAllowance) == -1 {
			takerZRXAvailable = traderInfo.TakerZrxBalance
		} else {
			takerZRXAvailable = traderInfo.TakerZrxAllowance
		}
		maxTakerAssetFillAmountGivenTakerZRXConstraints := new(big.Int).Div(new(big.Int).Mul(takerZRXAvailable, signedOrder.TakerAssetAmount), signedOrder.TakerFee)
		minSet = append(minSet, maxTakerAssetFillAmountGivenTakerZRXConstraints)
	}

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
