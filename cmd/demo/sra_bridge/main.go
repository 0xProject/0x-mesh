// +build !js

// demo/sra_bridge is a short program that fetches orders from the RadarRelay SRA endpoint
// and dumps them into a Mesh node which watching the order event stream.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

type clientEnvVars struct {
	// RPCAddress is the address of the 0x Mesh node to communicate with.
	RPCAddress string `envvar:"RPC_ADDRESS"`
}

func main() {
	env := clientEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	client, err := rpc.NewClient(env.RPCAddress)
	if err != nil {
		log.WithError(err).Fatal("could not create client")
	}

	ctx := context.Background()
	orderEventsChan := make(chan []*zeroex.OrderEvent, 8000)
	clientSubscription, err := client.SubscribeToOrders(ctx, orderEventsChan)
	_ = clientSubscription
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}

	orderHashToSignedOrder := map[common.Hash]*zeroex.SignedOrder{}

	go func() {
		for {
			select {
			case orderEvents := <-orderEventsChan:
				for _, orderEvent := range orderEvents {
					if orderEvent.Kind != "ADDED" && orderEvent.Kind != "EXPIRED" {
						log.Printf("%s - %s - remaining: %d -- txHash %s\n", orderEvent.Kind, orderEvent.OrderHash.Hex(), orderEvent.FillableTakerAssetAmount, orderEvent.TxHash.Hex())
					}
				}
			}
		}

		// clientSubscription.Unsubscribe()
	}()

	ticker := time.NewTicker(20 * time.Second)
	chunkSize := 200
	for {
		signedOrders := getRadarOrders()
		// signedOrderJSON := `{"exchangeAddress":"0x4530c0483a1633c7a1c97d2c53721caff2caaaaf","makerAddress":"0x8cff49b26d4d13e0601769f8a60fd697b713b9c6","takerAddress":"0x0000000000000000000000000000000000000000","senderAddress":"0x0000000000000000000000000000000000000000","feeRecipientAddress":"0x0000000000000000000000000000000000000000","expirationTimeSeconds":"1559826927","salt":"48128453606684653105952683301312821720867493716494911784363103883716429240740","makerAssetAmount":"100000000000000000","takerAssetAmount":"1000000000000000000","takerAssetData":"0xf47261b0000000000000000000000000ff67881f8d12f372d91baae9752eb3631ff0ed00","makerAssetData":"0xf47261b0000000000000000000000000c778417e063141139fce010982780140aa0cd5ab","makerFee":"0","takerFee":"0","signature":"0x1cf5839d9a0025e684c3663151b1db14533cc8c9e495fb92543a37a7fffc0677a23f3b6d66a1f56d3fda46eb5277b4a91c7b7faad4fdaaa5aac9a1185dd545a8a002"}`
		// var signedOrder zeroex.SignedOrder
		// err := json.Unmarshal([]byte(signedOrderJSON), &signedOrder)
		// if err != nil {
		// 	panic(err)
		// }
		// signedOrders := []*zeroex.SignedOrder{&signedOrder}
		unseenSignedOrder := []*zeroex.SignedOrder{}
		for _, signedorder := range signedOrders {
			orderHash, _ := signedorder.ComputeOrderHash()
			if _, ok := orderHashToSignedOrder[orderHash]; !ok {
				orderHashToSignedOrder[orderHash] = signedorder
				unseenSignedOrder = append(unseenSignedOrder, signedorder)
			}
		}
		log.Println("Found", len(unseenSignedOrder), " new orders!")
		chunks := [][]*zeroex.SignedOrder{}
		for len(unseenSignedOrder) > chunkSize {
			chunks = append(chunks, unseenSignedOrder[:chunkSize])
			unseenSignedOrder = unseenSignedOrder[chunkSize:]
		}
		if len(unseenSignedOrder) > 0 {
			chunks = append(chunks, unseenSignedOrder)
		}

		for _, chunk := range chunks {
			validationResults, err := client.AddOrders(chunk)
			if err != nil {
				log.WithError(err).Fatal("error from AddOrder")
			} else {
				log.Printf("submitted %d orders. Accepted: %d, Rejected: %d", len(chunk), len(validationResults.Accepted), len(validationResults.Rejected))
				for _, rejected := range validationResults.Rejected {
					fmt.Printf("%s - %s - %s\n", rejected.Status, rejected.Kind, rejected.OrderHash.Hex())
				}
			}
		}

		<-ticker.C
	}
}

var radarRelaySRAEndpoint = "https://api.radarrelay.com/0x/v2"

const perPage = 1000
const networkID = 1

type ordersResponse struct {
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"perPage"`
	Records []orderData `json:"records"`
}

type orderData struct {
	Order wireOrder `json:"order"`
}

type wireOrder struct {
	MakerAddress          string `json:"makerAddress"`
	MakerAssetData        string `json:"makerAssetData"`
	MakerAssetAmount      string `json:"makerAssetAmount"`
	MakerFee              string `json:"makerFee"`
	TakerAddress          string `json:"takerAddress"`
	TakerAssetData        string `json:"takerAssetData"`
	TakerAssetAmount      string `json:"takerAssetAmount"`
	TakerFee              string `json:"takerFee"`
	SenderAddress         string `json:"senderAddress"`
	ExchangeAddress       string `json:"exchangeAddress"`
	FeeRecipientAddress   string `json:"feeRecipientAddress"`
	ExpirationTimeSeconds string `json:"expirationTimeSeconds"`
	Salt                  string `json:"salt"`
	Signature             string `json:"signature"`
}

func (w *wireOrder) convertToSignedOrder() *zeroex.SignedOrder {
	makerAssetAmount, ok := math.ParseBig256(w.MakerAssetAmount)
	if !ok {
		panic("Failed to parse makerAssetAmount")
	}
	takerAssetAmount, ok := math.ParseBig256(w.TakerAssetAmount)
	if !ok {
		panic("Failed to parse takerAssetAmount")
	}
	makerFee, ok := math.ParseBig256(w.MakerFee)
	if !ok {
		panic("Failed to parse makerFee")
	}
	takerFee, ok := math.ParseBig256(w.TakerFee)
	if !ok {
		panic("Failed to parse takerFee")
	}
	salt, ok := math.ParseBig256(w.Salt)
	if !ok {
		panic("Failed to parse salt")
	}
	expirationTimeSeconds, ok := math.ParseBig256(w.ExpirationTimeSeconds)
	if !ok {
		panic("Failed to parse expirationTimeSeconds")
	}
	return &zeroex.SignedOrder{
		Order: zeroex.Order{
			MakerAddress:          common.HexToAddress(w.MakerAddress),
			MakerAssetData:        common.Hex2Bytes(w.MakerAssetData[2:]),
			MakerAssetAmount:      makerAssetAmount,
			MakerFee:              makerFee,
			TakerAddress:          common.HexToAddress(w.TakerAddress),
			TakerAssetData:        common.Hex2Bytes(w.TakerAssetData[2:]),
			TakerAssetAmount:      takerAssetAmount,
			TakerFee:              takerFee,
			SenderAddress:         common.HexToAddress(w.SenderAddress),
			ExchangeAddress:       common.HexToAddress(w.ExchangeAddress),
			FeeRecipientAddress:   common.HexToAddress(w.FeeRecipientAddress),
			ExpirationTimeSeconds: expirationTimeSeconds,
			Salt:                  salt,
		},
		Signature: common.Hex2Bytes(w.Signature[2:]),
	}
}

func getRadarOrders() []*zeroex.SignedOrder {
	requestURL := fmt.Sprintf("%s/%s?networkId=%d&perPage=%d", radarRelaySRAEndpoint, "orders", networkID, perPage)
	resp, err := http.Get(requestURL)
	if err != nil {
		log.WithError(err).Fatal("Couldn't fetch orders from RadarRelay")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.WithFields(map[string]interface{}{
			"statusCode": resp.StatusCode,
			"body":       string(body),
		}).Warn("Got non-200 status code back from RadarRelay")
		return []*zeroex.SignedOrder{}
	}
	var ordersResponse ordersResponse
	err = json.Unmarshal(body, &ordersResponse)
	if err != nil {
		log.WithError(err).Warn("Failed to parse body into OrdersResponse struct")
		return []*zeroex.SignedOrder{}
	}

	signedOrders := []*zeroex.SignedOrder{}
	for _, record := range ordersResponse.Records {
		signedOrder := record.Order.convertToSignedOrder()
		signedOrders = append(signedOrders, signedOrder)
	}
	return signedOrders
}
