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
		"endState":                     string(o.EndState),
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           contractEventsJS,
	})
}

func (s SignedOrder) JSValue() js.Value {
	makerAssetData := ""
	if len(s.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerAssetData))
	}
	makerFeeAssetData := ""
	if len(s.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.MakerFeeAssetData))
	}
	takerAssetData := ""
	if len(s.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerAssetData))
	}
	takerFeeAssetData := ""
	if len(s.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(s.TakerFeeAssetData))
	}
	signature := ""
	if len(s.Signature) != 0 {
		signature = fmt.Sprintf("0x%s", common.Bytes2Hex(s.Signature))
	}

	return js.ValueOf(map[string]interface{}{
		"makerAddress":          strings.ToLower(s.MakerAddress.Hex()),
		"makerAssetData":        makerAssetData,
		"makerFeeAssetData":        makerFeeAssetData,
		"makerAssetAmount":      s.MakerAssetAmount.String(),
		"makerFee":              s.MakerFee.String(),
		"takerAddress":          strings.ToLower(s.TakerAddress.Hex()),
		"takerAssetData":        takerAssetData,
		"takerFeeAssetData":        takerFeeAssetData,
		"takerAssetAmount":      s.TakerAssetAmount.String(),
		"takerFee":              s.TakerFee.String(),
		"senderAddress":         strings.ToLower(s.SenderAddress.Hex()),
		"verifyingContractAddress":       strings.ToLower(s.verifyingContractAddress.Hex()),
		"feeRecipientAddress":   strings.ToLower(s.FeeRecipientAddress.Hex()),
		"expirationTimeSeconds": s.ExpirationTimeSeconds.String(),
		"salt":                  s.Salt.String(),
		"chainId":				s.ChainId.String(),
		"signature":             signature,
	})
}

func (c ContractEvent) JSValue() js.Value {
	m := map[string]interface{}{
		"blockHash": c.BlockHash.Hex(),
		"txHash":    c.TxHash.Hex(),
		"txIndex":   c.TxIndex,
		"logIndex":  c.LogIndex,
		"isRemoved": c.IsRemoved,
		"kind":      c.Kind,
		"parameters": c.Parameters.(js.Wrapper).JSValue(),
	}
	return js.ValueOf(m)
}
