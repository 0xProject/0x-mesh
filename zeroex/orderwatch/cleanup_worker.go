package orderwatch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/configs"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
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

// The most orders we can validate in a single eth_call without having Parity nor Geth
// timeout
const chunkSize = 500

// Specifies the number of eth_call requests we want to make at any given point in time.
// Additional requests to block until an ongoing request has completed.
const concurrencyLimit = 5

type getOrdersAndTradersInfoParams struct {
	TakerAddresses []common.Address
	Orders         []wrappers.OrderWithoutExchangeAddress
}

// CleanupWorker is a cleanup worker that runs at a specified interval and re-validates
// all orders stored by a Mesh node in-case an order has become unfillable but this change
// in state was missed by the contract event watcher.
type CleanupWorker struct {
	isCleanupWorkerRunning bool
	orderValidator         *wrappers.OrderValidator
}

// NewCleanupWorker instantiates a new cleanup worker
func NewCleanupWorker(orderValidatorAddress common.Address, ethClient *ethclient.Client) (*CleanupWorker, error) {
	orderValidator, err := wrappers.NewOrderValidator(orderValidatorAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &CleanupWorker{
		isCleanupWorkerRunning: false,
		orderValidator:         orderValidator,
	}, nil
}

// Start starts the cleanup workers polling loop
func (c *CleanupWorker) Start() {
	c.isCleanupWorkerRunning = true
	go func() {
		for {
			if !c.isCleanupWorkerRunning {
				return
			}

			start := time.Now()

			// TODO: Get all orders from DB where lastUpdated field is > X
			orders := []zeroex.SignedOrder{}
			c.RevalidateOrders(orders)

			// Wait MinCleanupInterval before calling RevalidateOrders again. Since
			// we only start sleeping _after_ RevalidateOrders completes, we will never
			// have multiple calls to RevalidateOrders running in parallel
			time.Sleep(configs.MinCleanupInterval - time.Since(start))
		}
	}()
}

// Stop stops the cleanup worker
func (c *CleanupWorker) Stop() {
	c.isCleanupWorkerRunning = false
}

// RevalidateOrders revalidates all the supplied orders in chunks of chunkSize, with never more then
// concurrencyLimit number of requests in parallel. If a request fails, re-attempt it
// up to four times and then give up.
func (c *CleanupWorker) RevalidateOrders(signedOrders []zeroex.SignedOrder) {
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

	// Make concurrencyLimit eth_call requests in parallel
	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for _, params := range chunks {
		wg.Add(1)
		go func(params getOrdersAndTradersInfoParams) {
			defer wg.Done()

			// Add one to the semaphore chan. If it already has concurrencyLimit values,
			// the request blocks here until one frees up.
			semaphoreChan <- struct{}{}

			// Attempt to make the eth_call request 4 times with an exponential back-off.
			maxDuration := 4 * time.Second
			b := &backoff.Backoff{
				//These are the defaults
				Min:    250 * time.Millisecond,
				Max:    maxDuration,
				Factor: 2,
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
				// Make eth_call request
				results, err := c.orderValidator.GetOrdersAndTradersInfo(opts, params.Orders, params.TakerAddresses)
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
						// TODO(fabio): Do we want to re-schedule the cleanup job immediately if
						// this happens?
						return // Give up after 4 attempts
					}
					time.Sleep(d)
					continue
				}

				// TODO: Process results
				fmt.Printf("%+v\n", results)

				<-semaphoreChan
				return
			}
		}(params)
	}

	wg.Wait()
}
