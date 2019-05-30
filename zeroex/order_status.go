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

// ConvertRejectOrderCodeToOrderEventKind converts an RejectOrderCode to an OrderEventKind type
func ConvertRejectOrderCodeToOrderEventKind(rejectedOrderCode RejectedOrderCode) OrderEventKind {
	switch rejectedOrderCode {
	case ROExpired:
		return EKOrderExpired
	case ROFullyFilled:
		return EKOrderFullyFilled
	case ROCancelled:
		return EKOrderCancelled
	case ROUnfunded:
		return EKOrderBecameUnfunded
	default:
		panicMessage := "No OrderEventKind corresponding to supplied RejectOrderCode"
		log.WithField("rejectedOrderCode", rejectedOrderCode).Panic(panicMessage)
		// HACK(fabio): Go annoyingly complains about this function missing a return at the end
		// unless I add this never-to-be-hit panic here because it cannot discern that Logrus
		// panics when calling `.Panic()`
		panic(panicMessage)
	}
}

// ConvertOrderStatusToRejectOrderCode converts an OrderStatus returned from the 0x smart contract
// validation into a Mesh-specific RejectOrderCode
func ConvertOrderStatusToRejectOrderCode(orderStatus OrderStatus) RejectedOrderCode {
	switch orderStatus {
	case OSInvalidMakerAssetAmount:
		return ROInvalidMakerAssetAmount
	case OSInvalidTakerAssetAmount:
		return ROInvalidTakerAssetAmount
	case OSExpired:
		return ROExpired
	case OSFullyFilled:
		return ROFullyFilled
	case OSCancelled:
		return ROCancelled
	case OSSignatureInvalid:
		return ROSignatureInvalid
	case OSInvalidMakerAssetData:
		return ROInvalidMakerAssetData
	case OSInvalidTakerAssetData:
		return ROInvalidTakerAssetData
	case OSFillable:
		return ROUnfunded
	default:
		panicMessage := "No RejectOrderCode corresponding to supplied OrderStatus"
		log.WithField("orderStatus", orderStatus).Panic(panicMessage)
		// HACK(fabio): Go annoyingly complains about this function missing a return at the end
		// unless I add this never-to-be-hit panic here because it cannot discern that Logrus
		// panics when calling `.Panic()`
		panic(panicMessage)
	}
}
