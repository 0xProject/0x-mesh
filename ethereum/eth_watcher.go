package ethereum

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/0xproject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// MainnetEthBalanceCheckerAddress is the mainnet EthBalanceChecker contract address
var MainnetEthBalanceCheckerAddress = common.HexToAddress("0x9bc2c6ae8b1a8e3c375b6ccb55eb4273b2c3fbde")

// GanacheEthBalanceCheckerAddress is the ganache snapshot EthBalanceChecker contract address
var GanacheEthBalanceCheckerAddress = common.HexToAddress("0xaa86dda78e9434aca114b6676fc742a18d15a1cc")

// The most addresses we can fetch balances for in a single CALL without having Parity nor Geth
// timeout
const chunkSize = 4000

// Balance represents a single Ethereum addresses Ether balance
type Balance struct {
	Address common.Address
	Balance *big.Int
}

// ETHWatcher allows for watching a set of Ethereum addresses for ETH balance
// changes
type ETHWatcher struct {
	addressToBalance   map[common.Address]*big.Int
	minPollingInterval time.Duration
	isWatching         bool
	balanceChan        chan Balance
	ethBalanceChecker  *wrappers.EthBalanceChecker
	ethClient          *ethclient.Client
	ticker             *time.Ticker
	addressToBalanceMu sync.Mutex
}

// NewETHWatcher creates a new instance of ETHWatcher
func NewETHWatcher(minPollingInterval time.Duration, ethClient *ethclient.Client, ethBalanceCheckerAddress common.Address) (*ETHWatcher, error) {
	ethBalanceChecker, err := wrappers.NewEthBalanceChecker(ethBalanceCheckerAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &ETHWatcher{
		addressToBalance:   make(map[common.Address]*big.Int),
		balanceChan:        make(chan Balance, 100),
		minPollingInterval: minPollingInterval,
		isWatching:         false,
		ethClient:          ethClient,
		ethBalanceChecker:  ethBalanceChecker,
	}, nil
}

// Start begins the ETH balance poller
func (e *ETHWatcher) Start() error {
	if e.isWatching {
		return errors.New("Watcher already started")
	}
	e.isWatching = true
	e.ticker = time.NewTicker(e.minPollingInterval)
	go func() {
		for {
			start := time.Now()

			e.updateBalances()

			// Wait minPollingInterval before calling updateBalances again. Since
			// we only start sleeping _after_ updateBalances completes, we will never
			// have multiple calls to updateBalances running in parallel
			time.Sleep(e.minPollingInterval - time.Since(start))
		}
	}()
	return nil
}

// Stop stops the ETH balance poller
func (e *ETHWatcher) Stop() {
	if !e.isWatching {
		return // noop
	}
	e.ticker.Stop()
	e.isWatching = false
}

// Add adds a new Ethereum address we'd like to track for balance changes
func (e *ETHWatcher) Add(address common.Address, initialBalance *big.Int) {
	e.addressToBalanceMu.Lock()
	defer e.addressToBalanceMu.Unlock()
	if existingBalance, ok := e.addressToBalance[address]; ok {
		log.WithFields(log.Fields{
			"address":         address.Hex(),
			"initialBalance":  initialBalance,
			"existingBalance": existingBalance,
		}).Warn("tried to add address to ETHWatcher that already exists")
		return // Noop. Already exists and we bias towards our existing balance
	}
	e.addressToBalance[address] = initialBalance
}

// Remove removes a new Ethereum address we no longer want to track for balance changes
func (e *ETHWatcher) Remove(address common.Address) {
	e.addressToBalanceMu.Lock()
	defer e.addressToBalanceMu.Unlock()
	if _, ok := e.addressToBalance[address]; ok {
		delete(e.addressToBalance, address)
	}
}

// GetBalance returns the Ether balance for a particular Ethereum address
func (e *ETHWatcher) GetBalance(address common.Address) (*big.Int, error) {
	e.addressToBalanceMu.Lock()
	defer e.addressToBalanceMu.Unlock()
	if balance, ok := e.addressToBalance[address]; ok {
		return balance, nil
	}
	return nil, errors.New("Supplied address not tracked")
}

// Receive returns a read-only channel that can be used to listen for modified ETH balances
func (e *ETHWatcher) Receive() <-chan Balance {
	return e.balanceChan
}

func (e *ETHWatcher) updateBalances() error {
	e.addressToBalanceMu.Lock()
	addresses := []common.Address{}
	for address := range e.addressToBalance {
		addresses = append(addresses, address)
	}
	e.addressToBalanceMu.Unlock()

	chunks := [][]common.Address{}
	// Chunk into groups of chunkSize addresses for the call
	for len(addresses) > chunkSize {
		chunks = append(chunks, addresses[:chunkSize])
		addresses = addresses[chunkSize:]
	}
	if len(addresses) > 0 {
		chunks = append(chunks, addresses)
	}

	wg := &sync.WaitGroup{}
	for _, chunk := range chunks {
		// Call contract for each chunk of addresses in parallel
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Pass a context with a 10 second timeout to `GetEthBalances` in order to avoid
			// any one request from taking longer then 10 seconds and as a consequence, hold
			// up the polling loop for more then 10 seconds.
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			opts := &bind.CallOpts{
				Pending: false,
				Context: ctx,
			}
			balances, err := e.ethBalanceChecker.GetEthBalances(opts, chunk)
			if err != nil {
				log.WithFields(log.Fields{
					"error":     err.Error(),
					"addresses": chunk,
				}).Info("ether batch balance check failed")
				return // Noop on failure, simply wait for next polling interval
			}
			for i, address := range chunk {
				e.addressToBalanceMu.Lock()
				if balance, ok := e.addressToBalance[address]; ok {
					if balance != balances[i] {
						e.addressToBalance[address] = balances[i]
						updatedBalance := Balance{
							Address: address,
							Balance: balances[i],
						}
						go func() {
							e.balanceChan <- updatedBalance
						}()
					}
				} else {
					log.WithFields(log.Fields{
						"address": address,
					}).Error("address unexpectedly missing from addressToBalance map")
				}
				e.addressToBalanceMu.Unlock()
			}
		}()
	}
	wg.Wait()
	return nil
}
