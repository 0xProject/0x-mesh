// +build js,wasm

package zeroex

import (
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (o OrderEvent) JSValue() js.Value {
	contractEventsJS := make([]interface{}, len(o.ContractEvents))
	for i, contractEvent := range o.ContractEvents {
		contractEventsJS[i] = contractEvent.JSValue()
	}
	return js.ValueOf(map[string]interface{}{
		"timestamp":                o.Timestamp.Format(time.RFC3339),
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder.JSValue(),
		"endState":                 string(o.EndState),
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           contractEventsJS,
	})
}

func (s SignedOrder) JSValue() js.Value {
	// FIXME
	o, ok := s.Order.(*OrderV3)
	if !ok {
		panic("can't use non-v3 orders")
	}
	makerAssetData := "0x"
	if len(o.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.MakerAssetData))
	}
	// Note(albrow): Because of how our smart contracts work, most fields of an
	// order cannot be null. However, makerAssetFeeData and takerAssetFeeData are
	// the exception. For these fields, "0x" is used to indicate a null value.
	makerFeeAssetData := "0x"
	if len(o.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.MakerFeeAssetData))
	}
	takerAssetData := "0x"
	if len(o.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.TakerAssetData))
	}
	takerFeeAssetData := "0x"
	if len(o.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(o.TakerFeeAssetData))
	}
	signature := "0x"
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	return js.ValueOf(map[string]interface{}{
		"chainId":               o.ChainID.Int64(),
		"exchangeAddress":       strings.ToLower(o.ExchangeAddress.Hex()),
		"makerAddress":          strings.ToLower(o.MakerAddress.Hex()),
		"makerAssetData":        makerAssetData,
		"makerFeeAssetData":     makerFeeAssetData,
		"makerAssetAmount":      o.MakerAssetAmount.String(),
		"makerFee":              o.MakerFee.String(),
		"takerAddress":          strings.ToLower(o.TakerAddress.Hex()),
		"takerAssetData":        takerAssetData,
		"takerFeeAssetData":     takerFeeAssetData,
		"takerAssetAmount":      o.TakerAssetAmount.String(),
		"takerFee":              o.TakerFee.String(),
		"senderAddress":         strings.ToLower(o.SenderAddress.Hex()),
		"feeRecipientAddress":   strings.ToLower(o.FeeRecipientAddress.Hex()),
		"expirationTimeSeconds": o.ExpirationTimeSeconds.String(),
		"salt":                  o.Salt.String(),
		"signature":             signature,
	})
}

func (c ContractEvent) JSValue() js.Value {
	m := map[string]interface{}{
		"address":    c.Address.Hex(),
		"blockHash":  c.BlockHash.Hex(),
		"txHash":     c.TxHash.Hex(),
		"txIndex":    c.TxIndex,
		"logIndex":   c.LogIndex,
		"isRemoved":  c.IsRemoved,
		"kind":       c.Kind,
		"parameters": c.Parameters.(js.Wrapper).JSValue(),
	}
	return js.ValueOf(m)
}
