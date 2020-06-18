// +build js, wasm

package ordervalidator

import (
	"sync"
	"syscall/js"
	"time"

	"github.com/0xProject/0x-mesh/zeroex"
	log "github.com/sirupsen/logrus"
)

// sleepTime is used to force processes that block the event loop to give the
// event loop time to continue. This should be as low as possible.
const sleepTime = 175 * time.Microsecond

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

const computeOptimalChunkBatchSize = 3

// computeOptimalChunkSizes splits the signedOrders into chunks where the payload size of each chunk
// is beneath the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizes(signedOrders []*zeroex.SignedOrder) []int {
	chunkSizes := []int{}

	payloadLength := jsonRPCPayloadByteLength
	nextChunkSize := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		batchIdx := 0
		for _, signedOrder := range signedOrders {
			if batchIdx%computeOptimalChunkBatchSize == 2 {
				time.Sleep(sleepTime)
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
		wg.Done()
	}()
	wg.Wait()
	if nextChunkSize != 0 {
		chunkSizes = append(chunkSizes, nextChunkSize)
	}

	return chunkSizes
}
