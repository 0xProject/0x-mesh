// +build !js

package ordervalidator

import (
	"github.com/0xProject/0x-mesh/zeroex"
	log "github.com/sirupsen/logrus"
)

// computeOptimalChunkSizes splits the signedOrders into chunks where the payload size of each chunk
// is beneath the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizes(signedOrders []*zeroex.SignedOrder) []int {
	chunkSizes := []int{}

	payloadLength := jsonRPCPayloadByteLength
	nextChunkSize := 0
	for _, signedOrder := range signedOrders {
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
	}
	if nextChunkSize != 0 {
		chunkSizes = append(chunkSizes, nextChunkSize)
	}

	return chunkSizes
}
