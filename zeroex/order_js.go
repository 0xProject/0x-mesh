// +build js,wasm

package zeroex

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/ethereum/go-ethereum/common"
)

func (o OrderEvent) JSValue() js.Value {
	contractEventsJS := make([]interface{}, len(o.ContractEvents))
	for i, contractEvent := range o.ContractEvents {
		contractEventsJS[i] = contractEvent.JSValue()
	}
	return js.ValueOf(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder.JSValue(),
		"endState":                 string(o.EndState),
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           contractEventsJS,
	})
}

func (s SignedOrder) JSValue() js.Value {
	makerAssetData := "0x"
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	// Note(albrow): Because of how our smart contracts work, most fields of an
	// order cannot be null. However, makerAssetFeeData and takerAssetFeeData are
	// the exception. For these fields, "0x" is used to indicate a null value.
	makerFeeAssetData := "0x"
	if len(s.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerFeeAssetData))
	}
	takerAssetData := "0x"
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	takerFeeAssetData := "0x"
	if len(s.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerFeeAssetData))
	}
	signature := "0x"
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	return js.ValueOf(map[string]interface{}{
		"chainId":               s.ChainID.Int64(),
		"exchangeAddress":       strings.ToLower(s.ExchangeAddress.Hex()),
		"makerAddress":          strings.ToLower(s.MakerAddress.Hex()),
		"makerAssetData":        makerAssetData,
		"makerFeeAssetData":     makerFeeAssetData,
		"makerAssetAmount":      s.MakerAssetAmount.String(),
		"makerFee":              s.MakerFee.String(),
		"takerAddress":          strings.ToLower(s.TakerAddress.Hex()),
		"takerAssetData":        takerAssetData,
		"takerFeeAssetData":     takerFeeAssetData,
		"takerAssetAmount":      s.TakerAssetAmount.String(),
		"takerFee":              s.TakerFee.String(),
		"senderAddress":         strings.ToLower(s.SenderAddress.Hex()),
		"feeRecipientAddress":   strings.ToLower(s.FeeRecipientAddress.Hex()),
		"expirationTimeSeconds": s.ExpirationTimeSeconds.String(),
		"salt":                  s.Salt.String(),
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
