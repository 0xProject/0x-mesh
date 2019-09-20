// +build !js
// HACK(fabio): We currently don't run these tests in WASM because the `roundtrip_js.go` file in
// Go's net/http special-cases requests originating from test files, routing the requests to a fake
// in-process network handler. This causes the tests to fail.
// Source: https://github.com/golang/go/issues/31559

package ethereum

import (
	"context"
	"flag"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Values taken from Ganache snapshot
var firstAccount = constants.GanacheAccount0
var firstAccountBalance, _ = math.ParseBig256("99723678048000000000")
var secondAccount = constants.GanacheAccount1
var secondAccountBalance, _ = math.ParseBig256("100000000000000000000")
var hundredEth, _ = math.ParseBig256("100000000000000000000")
var randomAccount = common.HexToAddress("0x49fea72f146d41bfc5b9329e4ebbc3c382589f37")

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	flag.Parse()
}

var ethAccountToBalance = map[common.Address]*big.Int{
	firstAccount:              firstAccountBalance,
	secondAccount:             secondAccountBalance,
	constants.GanacheAccount2: hundredEth,
	constants.GanacheAccount3: hundredEth,
	common.HexToAddress("0xa8dda8d7f5310e4a9e24f8eba77e091ac264f872"): hundredEth,
	common.HexToAddress("0x06cef8e666768cc40cc78cf93d9611019ddcb628"): hundredEth,
	common.HexToAddress("0x4404ac8bd8f9618d27ad2f1485aa1b2cfd82482d"): hundredEth,
	common.HexToAddress("0x7457d5e02197480db681d3fdf256c7aca21bdc12"): hundredEth,
	common.HexToAddress("0x91c987bf62d25945db517bdaa840a6c661374402"): hundredEth,
}

var pollingInterval = 100 * time.Millisecond

func TestAddingAddressesToETHWatcher(t *testing.T) {
	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	addresses := []common.Address{}
	for address := range ethAccountToBalance {
		addresses = append(addresses, address)
	}
	addressToBalance, failedAddresses := ethWatcher.Add(addresses)

	assert.Len(t, failedAddresses, 0)
	assert.Equal(t, ethAccountToBalance, addressToBalance)
}

func TestUpdateBalancesETHWatcher(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	blockchainLifecycle, err := NewBlockchainLifecycle(rpcClient)
	require.NoError(t, err)
	blockchainLifecycle.Start(t)
	defer blockchainLifecycle.Revert(t)

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	ethWatcher, err := NewETHWatcher(pollingInterval, ethClient, constants.TestNetworkID)
	require.NoError(t, err)

	addresses := []common.Address{}
	for address := range ethAccountToBalance {
		addresses = append(addresses, address)
	}
	addressToInitialBalance, failedAddresses := ethWatcher.Add(addresses)
	assert.Len(t, failedAddresses, 0)
	assert.Len(t, addressToInitialBalance, len(ethAccountToBalance))

	amount := big.NewInt(int64(1000000))
	transferFunds(t, ethClient, constants.GanacheAccount0, randomAccount, amount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		require.NoError(t, ethWatcher.Watch(ctx))
	}()

	select {
	case balance := <-ethWatcher.BalanceUpdates():
		assert.Equal(t, constants.GanacheAccount0, balance.Address)
		assert.NotEqual(t, addressToInitialBalance[balance.Address], balance.Amount, "wrong balance for account: %s", balance.Address.Hex())

	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for balance channel to deliver expected balances")
	}
}

func transferFunds(t *testing.T, client *ethclient.Client, from, to common.Address, value *big.Int) {
	pkBytes, ok := constants.GanacheAccountToPrivateKey[from]
	if !ok {
		t.Errorf("no corresponding private key found for: %s", from.Hex())
	}
	privateKey, err := crypto.ToECDSA(pkBytes)

	nonce, err := client.NonceAt(context.Background(), from, nil)
	require.NoError(t, err)
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(10000000)
	rawTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, []byte{})
	require.NoError(t, err)

	signer := types.HomesteadSigner{}
	signature, err := crypto.Sign(signer.Hash(rawTx).Bytes(), privateKey)
	require.NoError(t, err)
	signedTx, err := rawTx.WithSignature(signer, signature)
	require.NoError(t, err)

	err = client.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)

	// Make sure transaction did not revert
	txReceipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	assert.Equal(t, uint64(1), txReceipt.Status)
}
