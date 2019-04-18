package ethereum

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/0xproject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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
	addressToBalance         map[common.Address]*big.Int
	pollingInterval          time.Duration
	isWatching               bool
	balanceChan              chan Balance
	ethBalanceCheckerAddress common.Address
	ethClient                *ethclient.Client
	ticker                   *time.Ticker
	addressToBalanceMu       sync.Mutex
}

// NewETHWatcher creates a new instance of ETHWatcher
func NewETHWatcher(pollingInterval time.Duration, ethClient *ethclient.Client, ethBalanceCheckerAddress common.Address) *ETHWatcher {
	return &ETHWatcher{
		addressToBalance:         make(map[common.Address]*big.Int),
		balanceChan:              make(chan Balance, 100),
		pollingInterval:          pollingInterval,
		ethBalanceCheckerAddress: ethBalanceCheckerAddress,
		isWatching:               false,
		ethClient:                ethClient,
	}
}

// Start begins the ETH balance poller
func (e *ETHWatcher) Start() error {
	if e.isWatching {
		return errors.New("Watcher already started")
	}
	e.isWatching = true
	e.ticker = time.NewTicker(e.pollingInterval)
	go func() {
		for {
			<-e.ticker.C

			// TODO(fabio): Currently if `updateBalance` takes longer then the ticker interval,
			// we would kick off an additional call before the previous one completes. This might
			// not be desirable, and we might want to noop if a previous request is in progress.
			e.updateBalances()
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
	if _, ok := e.addressToBalance[address]; ok {
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

	ethBalanceChecker, err := wrappers.NewEthBalanceChecker(e.ethBalanceCheckerAddress, e.ethClient)
	if err != nil {
		return err
	}
	for _, chunk := range chunks {
		// Call contract for each chunk of addresses in parallel
		go func() {
			balances, err := ethBalanceChecker.GetEthBalances(nil, chunk)
			if err != nil {
				// TODO(fabio): Log error
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
				}
				e.addressToBalanceMu.Unlock()
			}
		}()
	}
	return nil
}
