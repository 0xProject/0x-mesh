package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// MainnetEthBalanceCheckerAddress is the mainnet EthBalanceChecker contract address
var MainnetEthBalanceCheckerAddress = common.HexToAddress("0x9bc2c6ae8b1a8e3c375b6ccb55eb4273b2c3fbde")

// GanacheEthBalanceCheckerAddress is the ganache snapshot EthBalanceChecker contract address
var GanacheEthBalanceCheckerAddress = common.HexToAddress("0xaa86dda78e9434aca114b6676fc742a18d15a1cc")

// The most addresses we can fetch balances for in a single CALL without going over the block gas
// limit. One of Geth/Parity caps the gas limit for `eth_call`s at the block gas limit.
// Block gas limit on 19th April 2019: 7,600,889
const chunkSize = 3500 // 7,475,648 gas

// Balance represents a single Ethereum addresses Ether balance
type Balance struct {
	Address common.Address
	Amount  *big.Int
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
	addressToBalanceMu sync.Mutex
}

// NewETHWatcher creates a new instance of ETHWatcher
func NewETHWatcher(minPollingInterval time.Duration, ethClient *ethclient.Client, networkID int) (*ETHWatcher, error) {
	contractNameToAddress, err := getContractAddressesForNetworkID(networkID)
	if err != nil {
		return nil, err
	}
	ethBalanceChecker, err := wrappers.NewEthBalanceChecker(contractNameToAddress.EthBalanceChecker, ethClient)
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
	e.isWatching = false
}

// Add adds new Ethereum addresses we'd like to track for balance changes and returns a map of added
// address to balance, and an array of addresses that it failed to add due to failed network requests
func (e *ETHWatcher) Add(addresses []common.Address) (addressToBalance map[common.Address]*big.Int, failedAddresses []common.Address) {
	e.addressToBalanceMu.Lock()
	defer e.addressToBalanceMu.Unlock()
	newAddresses := []common.Address{}
	for _, address := range addresses {
		if _, ok := e.addressToBalance[address]; !ok {
			newAddresses = append(newAddresses, address)
		}
	}
	addressToBalance, failedAddresses = e.getBalances(newAddresses)
	for address, balance := range addressToBalance {
		e.addressToBalance[address] = balance
	}
	return addressToBalance, failedAddresses
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

// getContractAddressesForNetworkID returns the contract name mapping for the given network.
// It returns an error if the network doesn't exist.
func getContractAddressesForNetworkID(networkID int) (constants.ContractNameToAddress, error) {
	if contractNameToAddress, ok := constants.NetworkIDToContractAddresses[networkID]; ok {
		return contractNameToAddress, nil
	}
	return constants.ContractNameToAddress{}, fmt.Errorf("invalid network: %d", networkID)
}

func (e *ETHWatcher) updateBalances() {
	e.addressToBalanceMu.Lock()
	defer e.addressToBalanceMu.Unlock()
	addresses := []common.Address{}
	for address := range e.addressToBalance {
		addresses = append(addresses, address)
	}
	// Intentionally ignore addresses we failed to fetch balances for
	// and simply attempt them again at the next polling interval
	addressToAmount, _ := e.getBalances(addresses)
	for address, newAmount := range addressToAmount {
		if cachedBalance, ok := e.addressToBalance[address]; ok {
			if cachedBalance.Cmp(newAmount) != 0 {
				e.addressToBalance[address] = newAmount
				updatedBalance := Balance{
					Address: address,
					Amount:  newAmount,
				}
				go func() {
					e.balanceChan <- updatedBalance
				}()
			}
		} else {
			// Due to the asynchronous nature of the ethWatcher, there are race-conditions
			// where we try to update the balance of an address after it has been removed from the
			// ethWatcher.
			log.WithFields(log.Fields{
				"address": address,
				"balance": newAmount,
			}).Trace("Attempted to update an ETH balance from ethWatcher that is no longer tracked")
		}
	}
}

func (e *ETHWatcher) getBalances(addresses []common.Address) (map[common.Address]*big.Int, []common.Address) {
	chunks := [][]common.Address{}
	// Chunk into groups of chunkSize addresses for the call
	for len(addresses) > chunkSize {
		chunks = append(chunks, addresses[:chunkSize])
		addresses = addresses[chunkSize:]
	}
	if len(addresses) > 0 {
		chunks = append(chunks, addresses)
	}

	addressToBalance := map[common.Address]*big.Int{}
	failedAddresses := []common.Address{}

	wg := &sync.WaitGroup{}
	for _, chunk := range chunks {
		// Call contract for each chunk of addresses in parallel
		wg.Add(1)
		go func(chunk []common.Address) {
			defer wg.Done()

			// Pass a context with a 20 second timeout to `GetEthBalances` in order to avoid
			// any one request from taking longer then 20 seconds and as a consequence, hold
			// up the polling loop for more then 20 seconds.
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			opts := &bind.CallOpts{
				Pending: false,
				Context: ctx,
			}
			balances, err := e.ethBalanceChecker.GetEthBalances(opts, chunk)
			if err != nil {
				for _, address := range chunk {
					failedAddresses = append(failedAddresses, address)
				}
				log.WithFields(log.Fields{
					"error":     err.Error(),
					"addresses": chunk,
				}).Info("ether batch balance check failed")
				return // Noop on failure
			}
			for i, address := range chunk {
				newBalance := balances[i]
				addressToBalance[address] = newBalance
			}
		}(chunk)
	}
	wg.Wait()
	return addressToBalance, failedAddresses
}
