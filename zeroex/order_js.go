// +build js,wasm

package zeroex

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/ethereum/go-ethereum/common"
)

func (o OrderEvent) JSValue() js.Value {
	stringifiedTxHashes := []interface{}{}
	for _, txHash := range o.TxHashes {
		stringifiedTxHashes = append(stringifiedTxHashes, txHash.Hex())
	}
	return js.ValueOf(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder.JSValue(),
		"kind":                     string(o.Kind),
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"txHashes":                 stringifiedTxHashes,
	})
}

func (s SignedOrder) JSValue() js.Value {
	makerAssetData := ""
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	takerAssetData := ""
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	signature := ""
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	return js.ValueOf(map[string]interface{}{
		"makerAddress":          strings.ToLower(s.MakerAddress.Hex()),
		"makerAssetData":        makerAssetData,
		"makerAssetAmount":      s.MakerAssetAmount.String(),
		"makerFee":              s.MakerFee.String(),
		"takerAddress":          strings.ToLower(s.TakerAddress.Hex()),
		"takerAssetData":        takerAssetData,
		"takerAssetAmount":      s.TakerAssetAmount.String(),
		"takerFee":              s.TakerFee.String(),
		"senderAddress":         strings.ToLower(s.SenderAddress.Hex()),
		"exchangeAddress":       strings.ToLower(s.ExchangeAddress.Hex()),
		"feeRecipientAddress":   strings.ToLower(s.FeeRecipientAddress.Hex()),
		"expirationTimeSeconds": s.ExpirationTimeSeconds.String(),
		"salt":                  s.Salt.String(),
		"signature":             signature,
	})
}
