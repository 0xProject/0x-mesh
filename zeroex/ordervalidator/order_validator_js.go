// +build js, wasm

package ordervalidator

import "syscall/js"

func (v ValidationResults) JSValue() js.Value {
	accepted := make([]interface{}, len(v.Accepted))
	for i, info := range v.Accepted {
		accepted[i] = info
	}
	rejected := make([]interface{}, len(v.Rejected))
	for i, info := range v.Rejected {
		rejected[i] = info
	}
	return js.ValueOf(map[string]interface{}{
		"accepted": accepted,
		"rejected": rejected,
	})
}

func (a AcceptedOrderInfo) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"orderHash":                a.OrderHash.Hex(),
		"signedOrder":              a.SignedV3Order.JSValue(),
		"fillableTakerAssetAmount": a.FillableTakerAssetAmount.String(),
		"isNew":                    a.IsNew,
	})
}

func (r RejectedOrderInfo) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"orderHash":   r.OrderHash.String(),
		"signedOrder": r.SignedV3Order.JSValue(),
		"kind":        string(r.Kind),
		"status":      r.Status.JSValue(),
	})
}

func (s RejectedOrderStatus) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"code":    s.Code,
		"message": s.Message,
	})
}
