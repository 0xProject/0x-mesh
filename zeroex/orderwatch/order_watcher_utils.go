package orderwatch

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type contractAddressesSeenCounter struct {
	contractAddressToSeenCount map[common.Address]uint

	mu sync.RWMutex
}

func NewContractAddressesSeenCounter() *contractAddressesSeenCounter {
	return &contractAddressesSeenCounter{
		contractAddressToSeenCount: map[common.Address]uint{},
	}
}

func (ca *contractAddressesSeenCounter) Inc(address common.Address) uint {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	newValue := ca.contractAddressToSeenCount[address] + 1
	ca.contractAddressToSeenCount[address] = newValue
	return newValue
}

func (ca *contractAddressesSeenCounter) Dec(address common.Address) uint {
	ca.mu.Lock()
	defer ca.mu.Unlock()
	newValue := ca.contractAddressToSeenCount[address] - 1
	ca.contractAddressToSeenCount[address] = newValue
	return newValue
}

func (ca *contractAddressesSeenCounter) Get(address common.Address) uint {
	ca.mu.RLock()
	defer ca.mu.RUnlock()
	return ca.contractAddressToSeenCount[address]
}
