package zeroex

import (
	log "github.com/sirupsen/logrus"
)

// OrderStatus represents the status of an order as returned from the 0x smart contracts
// as part of OrderInfo
type OrderStatus uint8

// OrderStatus values
const (
	OSInvalid OrderStatus = iota
	OSInvalidMakerAssetAmount
	OSInvalidTakerAssetAmount
	OSFillable
	OSExpired
	OSFullyFilled
	OSCancelled
	OSSignatureInvalid
	OSInvalidMakerAssetData
	OSInvalidTakerAssetData
)

// ConvertOrderStatusToRejectOrderCode converts an OrderStatus returned from the 0x smart contract
// validation into a Mesh-specific RejectOrderCode
func ConvertOrderStatusToRejectOrderCode(orderStatus OrderStatus) (RejectedOrderCode, bool) {
	switch orderStatus {
	case OSInvalidMakerAssetAmount:
		return ROInvalidMakerAssetAmount, true
	case OSInvalidTakerAssetAmount:
		return ROInvalidTakerAssetAmount, true
	case OSExpired:
		return ROExpired, true
	case OSFullyFilled:
		return ROFullyFilled, true
	case OSCancelled:
		return ROCancelled, true
	case OSSignatureInvalid:
		return ROSignatureInvalid, true
	case OSInvalidMakerAssetData:
		return ROInvalidMakerAssetData, true
	case OSInvalidTakerAssetData:
		return ROInvalidTakerAssetData, true
	case OSFillable:
		return ROUnfunded, true
	default:
		// Catch-all returns Invalid RejectedOrderCode
		return ROInvalid, false
	}
}
