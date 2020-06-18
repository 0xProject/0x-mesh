// +build js, wasm

package ordervalidator

import (
	"context"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/packages/browser/go/jsutil"
	"github.com/0xProject/0x-mesh/zeroex"
	log "github.com/sirupsen/logrus"
)

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
		"signedOrder":              a.SignedOrder.JSValue(),
		"fillableTakerAssetAmount": a.FillableTakerAssetAmount.String(),
		"isNew":                    a.IsNew,
	})
}

func (r RejectedOrderInfo) JSValue() js.Value {
	return js.ValueOf(map[string]interface{}{
		"orderHash":   r.OrderHash.String(),
		"signedOrder": r.SignedOrder.JSValue(),
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

// computeOptimalChunkBatchSize is the number of computations that we do in
// computeOptimalChunkSizes before waiting for the next tick of the event loop.
const computeOptimalChunkBatchSize = 10

// computeOptimalChunkSizes splits the signedOrders into chunks where the payload size of each chunk
// is beneath the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizes(signedOrders []*zeroex.SignedOrder) []int {
	chunkSizes := []int{}

	payloadLength := jsonRPCPayloadByteLength
	nextChunkSize := 0
	batchIdx := 0
	for _, signedOrder := range signedOrders {
		// NOTE(jalextowle): We need to occasionally wait for the next tick
		// of the Javascript event loop in order to break up the execution
		// of this function. Doing this too often will significantly slow
		// down the execution of the Mesh node, so we sleep after executing
		// "batches" of computations.
		if batchIdx%computeOptimalChunkBatchSize == computeOptimalChunkBatchSize-1 {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			jsutil.NextTick(ctx)
			cancel()
		}
		encodedSignedOrderByteLength, _ := o.computeABIEncodedSignedOrderByteLength(signedOrder)
		if payloadLength+encodedSignedOrderByteLength < o.maxRequestContentLength {
			payloadLength += encodedSignedOrderByteLength
			nextChunkSize++
		} else {
			if nextChunkSize == 0 {
				// This case should never be hit since we enforce that EthereumRPCMaxContentLength >= maxOrderSizeInBytes
				log.WithField("signedOrder", signedOrder).Panic("EthereumRPCMaxContentLength is set so low, a single 0x order cannot fit beneath the payload limit")
			}
			chunkSizes = append(chunkSizes, nextChunkSize)
			nextChunkSize = 1
			payloadLength = jsonRPCPayloadByteLength + encodedSignedOrderByteLength
		}
		batchIdx = batchIdx + 1
	}
	if nextChunkSize != 0 {
		chunkSizes = append(chunkSizes, nextChunkSize)
	}

	return chunkSizes
}
