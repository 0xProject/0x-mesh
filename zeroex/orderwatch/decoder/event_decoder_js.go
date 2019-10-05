// +build js,wasm

package decoder

import (
	"fmt"
	"syscall/js"

	"github.com/ethereum/go-ethereum/common"
)

func (e ERC20TransferEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"from":                e.From.Hex(),
		"to":              	   e.To.Hex(),
		"value":               e.Value.String(),
	})
}

func (e ERC20ApprovalEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":                e.Owner.Hex(),
		"spender":              e.Spender.Hex(),
		"value":                e.Value.String(),
	})
}

func (e ERC721TransferEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"from":                e.From.Hex(),
		"to":              	   e.To.Hex(),
		"tokenId":			   e.TokenId.String(),
	})
}

func (e ERC721ApprovalEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":                e.Owner.Hex(),
		"approved":             e.Approved.Hex(),
		"tokenId":              e.TokenId.String(),
	})
}

func (e ERC721ApprovalForAllEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":                e.Owner.Hex(),
		"operator":             e.Operator.Hex(),
		"approved":             e.Approved,
	})
}

func (e ExchangeFillEvent) JSValue() js.Value {
	makerAssetData := ""
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := ""
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	return js.ValueOf(map[string]interface{}{
		"makerAddress":           e.MakerAddress.Hex(),
		"takerAddress":           e.TakerAddress.Hex(),
		"senderAddress":          e.SenderAddress.Hex(),
		"feeRecipientAddress":    e.FeeRecipientAddress.Hex(),
		"makerAssetFilledAmount": e.MakerAssetFilledAmount.String(),
		"takerAssetFilledAmount": e.TakerAssetFilledAmount.String(),
		"makerFeePaid":           e.MakerFeePaid.String(),
		"takerFeePaid":           e.TakerFeePaid.String(),
		"orderHash":              e.OrderHash.Hex(),
		"makerAssetData":         makerAssetData,
		"takerAssetData":         takerAssetData,
	})
}

func (e ExchangeCancelEvent) JSValue() js.Value {
	makerAssetData := ""
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := ""
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	return js.ValueOf(map[string]interface{}{
		"makerAddress":        e.MakerAddress.Hex(),
		"senderAddress":       e.SenderAddress.Hex(),
		"feeRecipientAddress": e.FeeRecipientAddress.Hex(),
		"orderHash":           e.OrderHash.Hex(),
		"makerAssetData":      makerAssetData,
		"takerAssetData":      takerAssetData,
	})
}

func (e ExchangeCancelUpToEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"makerAddress":  e.MakerAddress.Hex(),
		"senderAddress": e.SenderAddress.Hex(),
		"orderEpoch":    e.OrderEpoch.String(),
	})
}

func (w WethWithdrawalEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner": w.Owner.Hex(),
		"value": w.Value.String(),
	})
}

func (w WethDepositEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner": w.Owner.Hex(),
		"value": w.Value.String(),
	})
}

