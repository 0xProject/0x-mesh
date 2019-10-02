// +build js,wasm

package zeroex

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
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

func (c ContractEvent) JSValue() js.Value {
	m := map[string]interface{}{
		"blockHash": c.BlockHash.Hex(),
		"txHash":    c.TxHash.Hex(),
		"txIndex":   c.TxIndex,
		"logIndex":  c.LogIndex,
		"isRemoved": c.IsRemoved,
		"kind":      c.Kind,
	}
	switch c.Kind {
	case "ERC20TransferEvent":
		m["parameters"] = c.Parameters.(decoder.ERC20TransferEvent).JSValue()

	case "ERC20ApprovalEvent":
		m["parameters"] = c.Parameters.(decoder.ERC20ApprovalEvent).JSValue()

	case "ERC721TransferEvent":
		m["parameters"] = c.Parameters.(decoder.ERC721TransferEvent).JSValue()

	case "ERC721ApprovalEvent":
		m["parameters"] = c.Parameters.(decoder.ERC721ApprovalEvent).JSValue()

	case "ERC721ApprovalForAllEvent":
		m["parameters"] = c.Parameters.(decoder.ERC721ApprovalForAllEvent).JSValue()

	case "WethWithdrawalEvent":
		m["parameters"] = c.Parameters.(decoder.WethWithdrawalEvent).JSValue()

	case "WethDepositEvent":
		m["parameters"] = c.Parameters.(decoder.WethDepositEvent).JSValue()

	case "ExchangeFillEvent":
		m["parameters"] = c.Parameters.(decoder.ExchangeFillEvent).JSValue()

	case "ExchangeCancelEvent":
		m["parameters"] = c.Parameters.(decoder.ExchangeCancelEvent).JSValue()

	case "ExchangeCancelUpToEvent":
		m["parameters"] = c.Parameters.(decoder.ExchangeCancelUpToEvent).JSValue()

	default:
		panic(fmt.Sprintf("Unrecognized event encountered: %s", c.Kind))
	}
	return js.ValueOf(m)
}
