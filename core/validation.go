// +build !js

package core

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

// maxOrderSizeInBytes is the maximum number of bytes allowed for encoded orders. It
// is more than 10x the size of a typical ERC20 order to account for multiAsset orders.
const maxOrderSizeInBytes = 8192

// maxOrderExpirationDuration is the maximum duration between the current time and the expiration
// set on an order that will be accepted by Mesh.
const maxOrderExpirationDuration = 9 * 30 * 24 * time.Hour // 9 months

var errMaxSize = fmt.Errorf("message exceeds maximum size of %d bytes", maxOrderSizeInBytes)

// RejectedOrderStatus values
var (
	ROInternalError = zeroex.RejectedOrderStatus{
		Code:    "InternalError",
		Message: "an unexpected internal error has occurred",
	}
	ROMaxOrderSizeExceeded = zeroex.RejectedOrderStatus{
		Code:    "MaxOrderSizeExceeded",
		Message: fmt.Sprintf("order exceeds the maximum encoded size of %d bytes", maxOrderSizeInBytes),
	}
	ROOrderAlreadyStored = zeroex.RejectedOrderStatus{
		Code:    "OrderAlreadyStored",
		Message: "order is already stored",
	}
	ROMaxExpirationExceeded = zeroex.RejectedOrderStatus{
		Code:    "OrderMaxExpirationExceeded",
		Message: "order expiration too far in the future",
	}
	ROIncorrectNetwork = zeroex.RejectedOrderStatus{
		Code:    "OrderForIncorrectNetwork",
		Message: "order was created for a different network than the one this Mesh node is configured to support",
	}
	ROInsufficientETHBacking = zeroex.RejectedOrderStatus{
		Code:    "InsufficientETHBacking",
		Message: "the maker address does not meet the minimum required ETH backing for storing the order",
	}
)

// RejectedOrderKind values
const (
	MeshValidation = zeroex.RejectedOrderKind("MESH_VALIDATION")
)

// validateOrders applies general 0x validation and Mesh-specific validation to
// the given orders.
func (app *App) validateOrders(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error) {
	results := &zeroex.ValidationResults{}
	validMeshOrders := []*zeroex.SignedOrder{}
	contractAddresses, err := ethereum.GetContractAddressesForNetworkID(app.networkID)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			log.WithField("error", err).Error("could not compute order hash")
			results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        zeroex.MeshError,
				Status:      ROInternalError,
			})
			continue
		}
		if order.ExchangeAddress != contractAddresses.Exchange {
			results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        MeshValidation,
				Status:      ROIncorrectNetwork,
			})
			continue
		}
		maxExpiration := big.NewInt(time.Now().Add(maxOrderExpirationDuration).Unix())
		if order.ExpirationTimeSeconds.Cmp(maxExpiration) > 0 {
			results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        MeshValidation,
				Status:      ROMaxExpirationExceeded,
			})
			continue
		}
		if err := validateOrderSize(order); err != nil {
			if err == errMaxSize {
				results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        MeshValidation,
					Status:      ROMaxOrderSizeExceeded,
				})
				continue
			} else {
				log.WithField("error", err).Error("could not validate order size")
				results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        zeroex.MeshError,
					Status:      ROInternalError,
				})
				continue
			}
		}
		alreadyStored, err := app.orderAlreadyStored(orderHash)
		if err != nil {
			log.WithField("error", err).Error("could not check if order was already stored")
			return nil, err
		}
		if alreadyStored {
			results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        MeshValidation,
				Status:      ROOrderAlreadyStored,
			})
			continue
		}
		validMeshOrders = append(validMeshOrders, order)
	}

	// Perform basic 0x off-chain validation
	validOffchainOrders, rejectedOffchainOrders := app.orderValidator.BatchOffchainValidation(validMeshOrders)
	results.Rejected = append(results.Rejected, rejectedOffchainOrders...)

	// Check ETH Balance for each maker in the set of remaining orders.
	ordersWithSufficientBacking, ordersWithInsufficientBacking := app.validateETHBacking(validOffchainOrders)
	results.Rejected = append(results.Rejected, ordersWithInsufficientBacking...)

	// Perform 0x on-chain validation.
	zeroexResults := app.orderValidator.BatchValidate(ordersWithSufficientBacking)
	zeroexResults.Rejected = append(zeroexResults.Rejected, results.Rejected...)
	return zeroexResults, nil
}

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > maxOrderSizeInBytes {
		return errMaxSize
	}
	return nil
}

func validateOrderSize(order *zeroex.SignedOrder) error {
	encoded, err := encodeOrder(order)
	if err != nil {
		return err
	}
	if len(encoded) > maxOrderSizeInBytes {
		return errMaxSize
	}
	return nil
}

// TODO(albrow): Use the more efficient Exists method instead of FindByID.
func (app *App) orderAlreadyStored(orderHash common.Hash) (bool, error) {
	var order meshdb.Order
	err := app.db.Orders.FindByID(orderHash.Bytes(), &order)
	if err == nil {
		return true, nil
	}
	if _, ok := err.(db.NotFoundError); ok {
		return false, nil
	}
	return false, err
}

// TODO(albrow): This function needs to be rigorously tested.
func (app *App) validateETHBacking(orders []*zeroex.SignedOrder) (ordersWithSufficientBacking []*zeroex.SignedOrder, rejectedOrders []*zeroex.RejectedOrderInfo) {
	totalExistingOrders, err := app.db.Orders.Count()
	if err != nil {
		return nil, nil
	}
	if totalExistingOrders+len(orders) <= app.config.MaxOrdersInStorage {
		// If we haven't yet reached our storage limit, all orders are considered
		// valid in terms of their ETH backing.
		return orders, nil
	}

	// Group orders by their maker address.
	makerAddressToOrders := map[common.Address][]*zeroex.SignedOrder{}
	for _, order := range orders {
		if _, found := makerAddressToOrders[order.MakerAddress]; found {
			makerAddressToOrders[order.MakerAddress] = append(makerAddressToOrders[order.MakerAddress], order)
		} else {
			makerAddressToOrders[order.MakerAddress] = []*zeroex.SignedOrder{order}
		}
	}
	var makerAddresses []common.Address
	for makerAddress := range makerAddressToOrders {
		makerAddresses = append(makerAddresses, makerAddress)
	}

	// TODO(albrow): Uncomment this.
	// Get the current ETH balance for each maker.
	// makerAddressToBalance, failedBalanceMakerAddresses := app.ethWathcher.Add(makerAddresses)

	// Call the helper function.
	// TODO(albrow): Actually build the heap and pass in the correct arguments. This is just a placeholder.
	valid, rejected := validateETHBackingsWithHeap(0, &ETHBackingHeap{}, orders)

	// TODO(albrow): Uncomment this.
	// Add any orders for maker addresses for which we failed to get the balance
	// to the set of rejected orders.
	// for _, failedAddress := range failedBalanceMakerAddresses {
	// 	orders := makerAddressToOrders[failedAddress]
	// 	rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, zeroex.ROEthRPCRequestFailed, orders)...)
	// }

	return valid, rejected
}

func validateETHBackingsWithHeap(spareCapacity int, ethBackingHeap *ETHBackingHeap, orders []*zeroex.SignedOrder) (ordersWithSufficientBacking []*zeroex.SignedOrder, rejectedOrders []*zeroex.RejectedOrderInfo) {
	valid := []*zeroex.SignedOrder{}
	rejected := []*zeroex.RejectedOrderInfo{}
	for i := len(valid); i < len(orders); i++ {
		incomingOrder := orders[i]

		// If we have spare capacity left, consider all orders valid for now.
		if len(valid) < spareCapacity {
			valid = append(valid, incomingOrder)
			ethBackingHeap.UpdateByMakerAddress(incomingOrder.MakerAddress, 1)
			continue
		}

		// If we don't have any spare capacity, check if the ETH backing per order
		// corresponding to this order's maker address is greater than the current
		// lowest ETH backing per order.
		lowestBacking := ethBackingHeap.Peek()
		thisBacking, _ := ethBackingHeap.FindByMakerAddress(incomingOrder.MakerAddress)
		potentialBacking := copyETHBacking(thisBacking)
		potentialBacking.OrderCount += 1
		potentialETHPerOrder := potentialBacking.ETHPerOrder()
		if potentialETHPerOrder.Cmp(lowestBacking.ETHPerOrder()) != 1 {
			// If we don't have the required ETH backing, this order is considered
			// invalid. We don't need to update the heap.
			orderHash, err := incomingOrder.ComputeOrderHash()
			if err != nil {
				log.WithField("error", err).Panic("Could not compute order hash")
			}
			rejectedOrderInfo := &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: incomingOrder,
				Kind:        MeshValidation,
				Status:      ROInsufficientETHBacking,
			}
			rejected = append(rejected, rejectedOrderInfo)
			continue
		}

		// If we do have the required ETH backing, we need to replace an order
		// which was previously considered valid with this one and update the heap.
		orderToRemove, indexToReplace := findOrderByMakerAddress(lowestBacking.MakerAddress, valid)
		if indexToReplace != -1 {
			// If the order which would be removed is in the set of valid orders so
			// far, we need to replace it.
			valid[indexToReplace] = incomingOrder
			// In this case we also need to add the order to the set of
			// rejectedOrderInfos.
			orderHash, err := orderToRemove.ComputeOrderHash()
			if err != nil {
				log.WithField("error", err).Panic("Could not compute order hash")
			}
			rejectedOrderInfo := &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: orderToRemove,
				Kind:        MeshValidation,
				Status:      ROInsufficientETHBacking,
			}
			rejected = append(rejected, rejectedOrderInfo)
		} else {
			// If the order which would be removed is *not* in the set of valid orders
			// so far, don't remove any. Just add the incoming order to the valid
			// orders so far.
			valid = append(valid, incomingOrder)
		}
		ethBackingHeap.UpdateByMakerAddress(incomingOrder.MakerAddress, 1)
		ethBackingHeap.UpdateLowest(-1)
	}

	return valid, rejected
}

func copyETHBacking(backing *meshdb.ETHBacking) *meshdb.ETHBacking {
	return &meshdb.ETHBacking{
		MakerAddress: backing.MakerAddress,
		OrderCount:   backing.OrderCount,
		ETHAmount:    backing.ETHAmount,
	}
}

func findOrderByMakerAddress(makerAddress common.Address, orders []*zeroex.SignedOrder) (*zeroex.SignedOrder, int) {
	// TODO(albrow): We could potentially optimize this by maintaining a map of
	// maker address to orders for that maker.
	for i, order := range orders {
		if order.MakerAddress.Hex() == makerAddress.Hex() {
			return order, i
		}
	}
	return nil, -1
}

func rejectedOrderInfoForOrders(kind zeroex.RejectedOrderKind, status zeroex.RejectedOrderStatus, orders []*zeroex.SignedOrder) []*zeroex.RejectedOrderInfo {
	rejected := make([]*zeroex.RejectedOrderInfo, len(orders))
	for i, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			log.WithField("error", err).Panic("Could not compute order hash")
		}
		rejected[i] = &zeroex.RejectedOrderInfo{
			OrderHash:   orderHash,
			SignedOrder: order,
			Kind:        kind,
			Status:      status,
		}
	}
	return rejected
}
