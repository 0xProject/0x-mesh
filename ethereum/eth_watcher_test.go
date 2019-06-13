// +build !js
// HACK(fabio): We currently don't run these tests in WASM because the `roundtrip_js.go` file in
// Go's net/http special-cases requests originating from test files, routing the requests to a fake
// in-process network handler. This causes the tests to fail.
// Source: https://github.com/golang/go/issues/31559

package ethereum

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Values taken from Ganache snapshot
var firstAccount = constants.GanacheAccount0
var firstAccountBalance, _ = math.ParseBig256("99931715680000000000")
var secondAccount = common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb")
var secondAccountBalance, _ = math.ParseBig256("49999822428000000000")
var hundredEth, _ = math.ParseBig256("100000000000000000000")

var ethAccountToBalance = map[common.Address]*big.Int{
	firstAccount:  firstAccountBalance,
	secondAccount: secondAccountBalance,
	common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"): hundredEth,
	common.HexToAddress("0xe834ec434daba538cd1b9fe1582052b880bd7e63"): hundredEth,
	common.HexToAddress("0x78dc5d2d739606d31509c31d654056a45185ecb6"): hundredEth,
	common.HexToAddress("0xa8dda8d7f5310e4a9e24f8eba77e091ac264f872"): hundredEth,
	common.HexToAddress("0x06cef8e666768cc40cc78cf93d9611019ddcb628"): hundredEth,
	common.HexToAddress("0x4404ac8bd8f9618d27ad2f1485aa1b2cfd82482d"): hundredEth,
	common.HexToAddress("0x7457d5e02197480db681d3fdf256c7aca21bdc12"): hundredEth,
	common.HexToAddress("0x91c987bf62d25945db517bdaa840a6c661374402"): hundredEth,
}

var pollingInterval = 100 * time.Millisecond

func TestAddingAddressToETHWatcher(t *testing.T) {
	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	for address := range ethAccountToBalance {
		ethWatcher.Add(address, big.NewInt(0))
	}

	addresses := []common.Address{}
	for address := range ethWatcher.addressToBalance {
		addresses = append(addresses, address)
	}

	expectedCount := 10
	assert.Equal(t, expectedCount, len(addresses))
}

func TestUpdateBalancesETHWatcher(t *testing.T) {
	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	for address := range ethAccountToBalance {
		// Set initial balances to 0 so that when the watcher is started an event will be emitted
		// for each tracked address with the correct balance.
		ethWatcher.Add(address, big.NewInt(0))
	}

	go func() {
		require.NoError(t, ethWatcher.updateBalances())
	}()

	for i := 0; i < len(ethAccountToBalance); i++ {
		select {
		case balance := <-ethWatcher.Receive():
			assert.Equal(t, ethAccountToBalance[balance.Address], balance.Balance, "wrong balance for account: %s", balance.Address.Hex())

		case <-time.After(3 * time.Second):
			t.Fatal("Timed out waiting for balance channel to deliver expected balances")
		}
	}
}
func TestUpdateChangedBalancesOnlyETHWatcher(t *testing.T) {
	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	// Add the first account with the correct initial balance. We expect no event to be emitted for
	// this account since the watcher already has the most up-to-date
	ethWatcher.Add(firstAccount, firstAccountBalance)

	// Add the second account with an incorrect balance, so that an event will be emitted for it
	// when running updateBalances()
	ethWatcher.Add(secondAccount, big.NewInt(0))

	go func() {
		require.NoError(t, ethWatcher.updateBalances())
	}()

	for i := 0; i < 1; i++ {
		select {
		case balance := <-ethWatcher.Receive():
			assert.Equal(t, secondAccountBalance, balance.Balance, "wrong balance for account: %s", balance.Address.Hex())
			assert.Equal(t, secondAccount, balance.Address)

		case <-time.After(3 * time.Second):
			t.Fatal("Timed out waiting for balance channel to deliver expected balances")
		}
	}
}
func TestStartStopETHWatcher(t *testing.T) {
	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	for address := range ethAccountToBalance {
		// Set initial balances to 0 so that when the watcher is started an event will be emitted
		// for each tracked address with the correct balance.
		ethWatcher.Add(address, big.NewInt(0))
	}

	require.NoError(t, ethWatcher.Start())

	for i := 0; i < len(ethAccountToBalance); i++ {
		select {
		case balance := <-ethWatcher.Receive():
			assert.Equal(t, ethAccountToBalance[balance.Address], balance.Balance, "wrong balance for account: %s", balance.Address.Hex())

		case <-time.After(3 * time.Second):
			t.Fatal("Timed out waiting for balance channel to deliver expected balances")
		}
	}

	ethWatcher.Stop()
	assert.False(t, ethWatcher.isWatching, "Calling Stop() should stop the ethWatcher poller")
}
