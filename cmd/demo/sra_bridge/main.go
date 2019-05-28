// +build !js

// demo/sra_bridge is a short program that fetches orders from the RadarRelay SRA endpoint
// and dumps them into a Mesh node which watching the order event stream.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
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
	orderInfosChan := make(chan []*zeroex.OrderInfo, 8000)
	clientSubscription, err := client.SubscribeToOrderStream(ctx, orderInfosChan)
	_ = clientSubscription
	if err != nil {
		log.WithError(err).Fatal("Couldn't set up OrderStream subscription")
	}

	go func() {
		for {
			select {
			case orderInfos := <-orderInfosChan:
				for _, orderInfo := range orderInfos {
					var orderStatus string
					switch orderInfo.OrderStatus {
					case zeroex.Invalid:
						orderStatus = "Invalid"
					case zeroex.Fillable:
						orderStatus = "Fillable"
					case zeroex.Expired:
						orderStatus = "Expired"
					case zeroex.FullyFilled:
						orderStatus = "FullyFilled"
					case zeroex.Cancelled:
						orderStatus = "Cancelled"
					case zeroex.SignatureInvalid:
						orderStatus = "SignatureInvalid"
					case zeroex.InvalidMakerAssetAmount:
						orderStatus = "InvalidMakerAssetAmount"
					case zeroex.InvalidTakerAssetAmount:
						orderStatus = "InvalidTakerAssetAmount"
					}
					log.Printf("Order event: Hash: %s status: %s remaining: %d\n", orderInfo.OrderHash.Hex(), orderStatus, orderInfo.FillableTakerAssetAmount)
				}
			}
		}

		// clientSubscription.Unsubscribe()
	}()

	orderHashToWasSeen := map[common.Hash]bool{}
	ticker := time.NewTicker(20 * time.Second)
	chunkSize := 200
	for {
		signedOrders := getRadarOrders()
		unseenSignedOrder := []*zeroex.SignedOrder{}
		for _, signedorder := range signedOrders {
			orderHash, _ := signedorder.ComputeOrderHash()
			if _, ok := orderHashToWasSeen[orderHash]; !ok {
				orderHashToWasSeen[orderHash] = true
				unseenSignedOrder = append(unseenSignedOrder, signedorder)
			}
		}
		fmt.Println("Found", len(unseenSignedOrder), " new orders!")
		chunks := [][]*zeroex.SignedOrder{}
		for len(unseenSignedOrder) > chunkSize {
			chunks = append(chunks, unseenSignedOrder[:chunkSize])
			unseenSignedOrder = unseenSignedOrder[chunkSize:]
		}
		if len(unseenSignedOrder) > 0 {
			chunks = append(chunks, unseenSignedOrder)
		}

		for _, chunk := range chunks {
			orderHashToSuccinctOrderInfo, err := client.AddOrders(chunk)
			if err != nil {
				log.WithError(err).Fatal("error from AddOrder")
			} else {
				log.Printf("submitted %d orders", len(orderHashToSuccinctOrderInfo))
				invalidOrderHashes := []common.Hash{}
				for orderHash, succinctOrderInfo := range orderHashToSuccinctOrderInfo {
					if succinctOrderInfo.FillableTakerAssetAmount == big.NewInt(0) {
						invalidOrderHashes = append(invalidOrderHashes, orderHash)
					}
				}
				log.Println(len(invalidOrderHashes), "invalid orders found:", invalidOrderHashes)
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
		Order: &zeroex.Order{
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
