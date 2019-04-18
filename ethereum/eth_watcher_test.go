package ethereum

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// var ethAccountToBalance = map[common.Address]*big.Int{
// 	common.HexToAddress("0x5409ed021d9299bf6814279a6a1411a7e866a631"): firstAccountBalance,
// 	common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"): hundredEth,
// 	common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"): hundredEth,
// 	common.HexToAddress("0xe834ec434daba538cd1b9fe1582052b880bd7e63"): hundredEth,
// 	common.HexToAddress("0x78dc5d2d739606d31509c31d654056a45185ecb6"): hundredEth,
// 	common.HexToAddress("0xa8dda8d7f5310e4a9e24f8eba77e091ac264f872"): hundredEth,
// 	common.HexToAddress("0x06cef8e666768cc40cc78cf93d9611019ddcb628"): hundredEth,
// 	common.HexToAddress("0x4404ac8bd8f9618d27ad2f1485aa1b2cfd82482d"): hundredEth,
// 	common.HexToAddress("0x7457d5e02197480db681d3fdf256c7aca21bdc12"): hundredEth,
// 	common.HexToAddress("0x91c987bf62d25945db517bdaa840a6c661374402"): hundredEth,
// }

// var firstAccountBalance, _ = math.ParseBig256("99943972190000000000")
// var hundredEth, _ = math.ParseBig256("100000000000000000000")

// var pollingInterval = 100 * time.Millisecond

// func TestAddingAddressToETHWatcher(t *testing.T) {
// 	ethClient, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	ethWatcher := NewETHWatcher(pollingInterval, ethClient, config.GanacheEthBalanceCheckerAddress)

// 	for address := range ethAccountToBalance {
// 		ethWatcher.Add(address, big.NewInt(0))
// 	}

// 	addresses := []common.Address{}
// 	for address := range ethWatcher.addressToBalance {
// 		addresses = append(addresses, address)
// 	}

// 	expectedCount := 10
// 	assert.Equal(t, expectedCount, len(addresses))
// }

// func TestUpdateBalancesETHWatcher(t *testing.T) {
// 	ethClient, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	ethWatcher := NewETHWatcher(pollingInterval, ethClient, config.GanacheEthBalanceCheckerAddress)

// 	for address := range ethAccountToBalance {
// 		ethWatcher.Add(address, big.NewInt(0))
// 	}

// 	go func() {
// 		ethWatcher.updateBalances()
// 	}()

// 	for i := 0; i < len(ethAccountToBalance); i++ {
// 		select {
// 		case balance := <-ethWatcher.Receive():
// 			assert.Equal(t, ethAccountToBalance[balance.Address], balance.Balance)

// 		case <-time.After(3 * time.Second):
// 			t.Fatal("Timed out waiting for balance channel to deliver expected balances")
// 		}
// 	}
// }
// func TestStartStopETHWatcher(t *testing.T) {
// 	ethClient, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	ethWatcher := NewETHWatcher(pollingInterval, ethClient, config.GanacheEthBalanceCheckerAddress)

// 	for address := range ethAccountToBalance {
// 		ethWatcher.Add(address, big.NewInt(0))
// 	}

// 	ethWatcher.Start()

// 	for i := 0; i < len(ethAccountToBalance); i++ {
// 		select {
// 		case balance := <-ethWatcher.Receive():
// 			assert.Equal(t, ethAccountToBalance[balance.Address], balance.Balance)

// 		case <-time.After(3 * time.Second):
// 			t.Fatal("Timed out waiting for balance channel to deliver expected balances")
// 		}
// 	}

// 	ethWatcher.Stop()
// }

func TestEthClient(t *testing.T) {
	fmt.Println("START!")
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://localhost:8545", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	msg := "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"eth_getBlockByNumber\",\"params\":[\"latest\",false]}"
	req.Body = ioutil.NopCloser(bytes.NewReader([]byte(msg)))
	req.ContentLength = int64(len(msg))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("data", data)

	// ethClient, err := ethclient.Dial("http://localhost:8545")
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// done := make(chan interface{})

	// go func() {
	// 	head, err := ethClient.HeaderByNumber(context.Background(), nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("head", head)
	// 	done <- 1
	// }()

	// <-done
}
