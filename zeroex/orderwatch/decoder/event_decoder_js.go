// +build js,wasm

package decoder

import (
	"fmt"
	"syscall/js"

	"github.com/ethereum/go-ethereum/common"
)

func (e ERC20TransferEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"from":  e.From.Hex(),
		"to":    e.To.Hex(),
		"value": e.Value.String(),
	})
}

func (e ERC20ApprovalEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":   e.Owner.Hex(),
		"spender": e.Spender.Hex(),
		"value":   e.Value.String(),
	})
}

func (e ERC721TransferEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"from":    e.From.Hex(),
		"to":      e.To.Hex(),
		"tokenId": e.TokenId.String(),
	})
}

func (e ERC721ApprovalEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"approved": e.Approved.Hex(),
		"tokenId":  e.TokenId.String(),
	})
}

func (e ERC721ApprovalForAllEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"operator": e.Operator.Hex(),
		"approved": e.Approved,
	})
}

func (e ERC1155ApprovalForAllEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"operator": e.Operator.Hex(),
		"approved": e.Approved,
	})
}

func (e ERC1155TransferSingleEvent) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"operator": e.Operator.Hex(),
		"from":     e.From.Hex(),
		"to":       e.To.Hex(),
		"id":       e.Id.String(),
		"value":    e.Value.String(),
	})
}

func (e ERC1155TransferBatchEvent) JSValue() js.Value {
	// NOTE(jalextowle): Both ids and values must be interface slices because
	// `ValueOf` is only defined for slices of interfaces.
	ids := []interface{}{}
	for _, id := range e.Ids {
		ids = append(ids, id.String())
	}
	values := []interface{}{}
	for _, value := range e.Values {
		values = append(values, value.String())
	}
	return js.ValueOf(map[string]interface{}{
		"operator": e.Operator.Hex(),
		"from":     e.From.Hex(),
		"to":       e.To.Hex(),
		"ids":      ids,
		"values":   values,
	})
}

func (e ExchangeFillEvent) JSValue() js.Value {
	makerAssetData := "0x"
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := "0x"
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	makerFeeAssetData := "0x"
	if len(e.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerFeeAssetData))
	}
	takerFeeAssetData := "0x"
	if len(e.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerFeeAssetData))
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
		"protocolFeePaid":        e.ProtocolFeePaid.String(),
		"orderHash":              e.OrderHash.Hex(),
		"makerAssetData":         makerAssetData,
		"takerAssetData":         takerAssetData,
		"makerFeeAssetData":      makerFeeAssetData,
		"takerFeeAssetData":      takerFeeAssetData,
	})
}

func (e ExchangeCancelEvent) JSValue() js.Value {
	makerAssetData := "0x"
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := "0x"
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
		"makerAddress":       e.MakerAddress.Hex(),
		"orderSenderAddress": e.OrderSenderAddress.Hex(),
		"orderEpoch":         e.OrderEpoch.String(),
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
