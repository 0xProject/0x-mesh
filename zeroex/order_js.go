// +build js,wasm

package zeroex

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (o OrderEvent) JSValue() js.Value {
	contractEventsJSValues := []js.Value{}
	for _, contractEvent := range o.ContractEvents {
		contractEventsJSValues = append(contractEventsJSValues, contractEvent.JSValue())
	}
	return js.ValueOf(map[string]interface{}{
		"orderHash":                o.OrderHash.Hex(),
		"signedOrder":              o.SignedOrder.JSValue(),
		"kind":                     string(o.Kind),
		"fillableTakerAssetAmount": o.FillableTakerAssetAmount.String(),
		"contractEvents":           contractEventsJSValues,
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
		m["parameters"] = c.Parameters.(ethereum.ERC20TransferEvent)

	case "ERC20ApprovalEvent":
		m["parameters"] = c.Parameters.(ethereum.ERC20ApprovalEvent)

	case "ERC721TransferEvent":
		m["parameters"] = c.Parameters.(ethereum.ERC721TransferEvent)

	case "ERC721ApprovalEvent":
		m["parameters"] = c.Parameters.(ethereum.ERC721ApprovalEvent)

	case "ERC721ApprovalForAllEvent":
		m["parameters"] = c.Parameters.(ethereum.ERC721ApprovalForAllEvent)

	case "WethWithdrawalEvent":
		m["parameters"] = c.Parameters.(ethereum.WethWithdrawalEvent)

	case "WethDepositEvent":
		m["parameters"] = c.Parameters.(ethereum.WethDepositEvent)

	case "ExchangeFillEvent":
		m["parameters"] = c.Parameters.(ethereum.ExchangeFillEvent)

	case "ExchangeCancelEvent":
		m["parameters"] = c.Parameters.(ethereum.ExchangeCancelEvent)

	case "ExchangeCancelUpToEvent":
		m["parameters"] = c.Parameters.(ethereum.ExchangeCancelUpToEvent)

	default:
		panic(fmt.Sprintf("Unrecognized event encountered: %s", c.Kind))
	}
	return js.ValueOf(m)
}
