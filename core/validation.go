// +build !js

package core

import (
	"container/heap"
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
	results.Rejected = append(results.Rejected, zeroexResults.Rejected...)
	results.Accepted = zeroexResults.Accepted
	return results, nil
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

type makerInfo struct {
	ethBacking *meshdb.ETHBacking
	orders     []*zeroex.SignedOrder
}

// TODO(albrow): If we need to optimize further, we can look into reducing the
// number of new map allocations. After DB lookups and network requests,
// memory/GC pressure is likely the bottleneck.
func (app *App) validateETHBacking(orders []*zeroex.SignedOrder) (ordersWithSufficientBacking []*zeroex.SignedOrder, rejectedOrders []*zeroex.RejectedOrderInfo) {
	if len(orders) == 0 {
		return ordersWithSufficientBacking, rejectedOrders
	}

	totalExistingOrders, err := app.db.Orders.Count()
	if err != nil {
		return nil, nil
	}
	if totalExistingOrders+len(orders) <= app.config.MaxOrdersInStorage {
		// If we haven't yet reached our storage limit, all orders are considered
		// valid in terms of their ETH backing. We still need to insert ETH backings
		// in the database for any new maker addresses among the incoming orders. In
		// this case, the only way an order can be considered "invalid" is if we
		// could not check the ETH balance or insert the corresponding ETH backing.
		valid, rejected := app.safeInsertETHBackings(orders)
		return valid, rejected
	} else if totalExistingOrders > app.config.MaxOrdersInStorage {
		// Can only ever happen if MaxOrdersInStorage is changed between runs.
		// Should be rare. Easy workaround is to consider all incoming orders
		// invalid.
		log.WithFields(map[string]interface{}{
			"totalExistingOrders": totalExistingOrders,
			"maxOrdersInStorage":  app.config.MaxOrdersInStorage,
		}).Error("total existing orders exceeds config.MaxOrdersInStorage")
		allRejected := []*zeroex.RejectedOrderInfo{}
		allRejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, allRejected, orders)
		return nil, allRejected
	}

	// Set up a map of maker address to maker info.
	makerInfos := make(map[common.Address]*makerInfo, len(orders))
	for _, order := range orders {
		info, found := makerInfos[order.MakerAddress]
		if found {
			info.orders = append(info.orders, order)
		} else {
			makerInfos[order.MakerAddress] = &makerInfo{
				orders: []*zeroex.SignedOrder{order},
			}
		}
	}

	// Open a transaction for checking and inserting ETHBackings.
	txn := app.db.ETHBackings.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// Get the lowest N ETH backings where N is the number of incoming orders. Add
	// them to makerInfos.
	lowestETHBackings, err := app.db.GetETHBackingsWithLowestETHPerOrder(len(orders))
	if err != nil {
		log.WithField("error", err).Error("Could not get ETH backings from database")
		rejected := []*zeroex.RejectedOrderInfo{}
		appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, orders)
		return nil, rejected
	}
	for _, ethBacking := range lowestETHBackings {
		info, found := makerInfos[ethBacking.MakerAddress]
		if found {
			info.ethBacking = ethBacking
		} else {
			makerInfos[ethBacking.MakerAddress] = &makerInfo{
				ethBacking: ethBacking,
			}
		}
	}

	// Lookup the ETH backing in the database (if it exists) for any remaining
	// makers in the set of incoming orders.
	rejected := []*zeroex.RejectedOrderInfo{}
	makerAddressesWithoutKnownBalance := []common.Address{}
	for makerAddress, info := range makerInfos {
		if info.ethBacking != nil {
			continue
		}
		var ethBackingFromDB meshdb.ETHBacking
		if err := app.db.ETHBackings.FindByID(makerAddress.Bytes(), &ethBackingFromDB); err != nil {
			if _, ok := err.(db.NotFoundError); ok {
				// If there was a NotFoundError it just means we haven't stored the ETH
				// backing for this maker yet. Add it to the list of makers for which we
				// don't know the balance.
				makerAddressesWithoutKnownBalance = append(makerAddressesWithoutKnownBalance, makerAddress)
				continue
			}
			// If there was a different kind of db error, we have no choice but to
			// consider all its orders invalid. We might get a chance to try again in
			// the future. This shouldn't happen often.
			log.WithFields(map[string]interface{}{
				"error":        err,
				"makerAddress": makerAddress,
			}).Error("Unexpected error while looking for ETHBacking in database")
			rejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, makerInfos[makerAddress].orders)
			delete(makerInfos, makerAddress)
		}
	}

	// For any makers for which there is not an existing ETH backing in the
	// database, use ethWatcher to get the current ETH balance and create a new
	// ETHBacking with a starting order count of 0. We also insert this backing
	// into the database even if there are ultimately no valid orders for this
	// maker. (Doing so makes it faster to validate orders from this maker in the
	// future.)
	makerAddressToBalance, failedBalanceMakerAddresses := app.ethWatcher.Add(makerAddressesWithoutKnownBalance)
	for makerAddress, makerBalance := range makerAddressToBalance {
		amountFloat, accuracy := new(big.Float).SetInt(makerBalance).Float64()
		if accuracy != big.Exact {
			// It should be fine if we can't represent balances with exact precision,
			// as long as we are close enough. We still want to log a warning to
			// understand how often this happens (if ever).
			log.WithFields(map[string]interface{}{
				"bigInt":  makerBalance,
				"float64": amountFloat,
			}).Warn("Could not accurately represent maker balance as a float64")
		}
		ethBacking := &meshdb.ETHBacking{
			MakerAddress: makerAddress,
			OrderCount:   0,
			AmountInWei:  amountFloat,
		}
		makerInfos[makerAddress].ethBacking = ethBacking
		if err := txn.Insert(ethBacking); err != nil {
			// If we can't save the ETH backing for this maker, we have no choice but
			// to consider all its orders invalid. We might get a chance to try again
			// in the future. This shouldn't happen often.
			log.WithFields(map[string]interface{}{
				"error":        err,
				"makerAddress": ethBacking.MakerAddress,
			}).Error("Unexpected error while inserting new ETHBacking in database")
			rejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, makerInfos[makerAddress].orders)
			delete(makerInfos, makerAddress)
		} else {
			log.WithField("makerAddress", ethBacking.MakerAddress).Info("Inserted ETHBacking for makerAddress")
		}
	}

	// At this point we can go ahead and commit the transaction. We won't be
	// inserting any additional ETH backings and the ones we need to retreive are
	// already in memory.
	if err := txn.Commit(); err != nil {
		// If we can't save the ETH backings at all, just bail and consider all
		// incoming orders invalid. We might get a chance to try again in the
		// future. This shouldn't happen often.
		log.WithField("error", err).Error("Could not commit transaction for inserting ETHBackings")
		allRejected := []*zeroex.RejectedOrderInfo{}
		allRejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, orders)
		return nil, allRejected
	}

	// Add any orders for maker addresses for which we failed to get the balance
	// to the set of rejected orders.
	for _, failedAddress := range failedBalanceMakerAddresses {
		log.WithFields(map[string]interface{}{
			"makerAddress": failedAddress,
		}).Error("Could not get balance for maker address")
		orders := makerInfos[failedAddress].orders
		rejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, zeroex.ROEthRPCRequestFailed, rejected, orders)
		delete(makerInfos, failedAddress)
	}

	// Get all the remaining orders and ETH backings. (Some entries in makerInfos
	// may have been deleted above.)
	remainingOrders := []*zeroex.SignedOrder{}
	remainingETHBackings := []*meshdb.ETHBacking{}
	for _, info := range makerInfos {
		remainingOrders = append(remainingOrders, info.orders...)
		remainingETHBackings = append(remainingETHBackings, info.ethBacking)
	}

	// Call the core algorithm for validating ETH backings.
	spareCapacity := app.config.MaxOrdersInStorage - totalExistingOrders
	valid, rejected := validateETHBackingsWithHeap(spareCapacity, remainingETHBackings, remainingOrders)

	return valid, rejected
}

// safeInsertETHBackings *only* inserts an ETHBacking in the database for any
// new makerAddresses among the given orders and does not do any further
// validation. It assumes all the given orders are valid. It looks up the
// current balance via app.ethWatcher as needed.
// TODO(albrow): de-dupe this code if possible.
func (app *App) safeInsertETHBackings(orders []*zeroex.SignedOrder) (valid []*zeroex.SignedOrder, rejected []*zeroex.RejectedOrderInfo) {
	// Open a transaction for checking and inserting ETHBackings.
	txn := app.db.ETHBackings.OpenTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	makerAddressToOrders := map[common.Address][]*zeroex.SignedOrder{}
	for _, order := range orders {
		if _, found := makerAddressToOrders[order.MakerAddress]; found {
			makerAddressToOrders[order.MakerAddress] = append(makerAddressToOrders[order.MakerAddress], order)
		} else {
			makerAddressToOrders[order.MakerAddress] = []*zeroex.SignedOrder{order}
		}
	}

	makerAddressesWithUnkownBalance := []common.Address{}
	for makerAddress, ordersForMaker := range makerAddressToOrders {
		var existingBacking meshdb.ETHBacking
		if err := app.db.ETHBackings.FindByID(makerAddress.Bytes(), &existingBacking); err != nil {
			if _, ok := err.(db.NotFoundError); ok {
				makerAddressesWithUnkownBalance = append(makerAddressesWithUnkownBalance, makerAddress)
				continue
			}
			log.WithFields(map[string]interface{}{
				"error":        err,
				"makerAddress": makerAddress,
			}).Error("Could not look up existing ETHBacking")
			appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, ordersForMaker)
			continue
		}
		// The ETH backing is already stored. Consider all orders for this maker
		// "valid".
		valid = append(valid, ordersForMaker...)
	}

	makerAddressToBalance, failedBalanceMakerAddresses := app.ethWatcher.Add(makerAddressesWithUnkownBalance)
	for makerAddress, makerBalance := range makerAddressToBalance {
		amountFloat, accuracy := new(big.Float).SetInt(makerBalance).Float64()
		if accuracy != big.Exact {
			// It should be fine if we can't represent balances with exact precision,
			// as long as we are close enough. We still want to log a warning to
			// understand how often this happens (if ever).
			log.WithFields(map[string]interface{}{
				"bigInt":  makerBalance,
				"float64": amountFloat,
			}).Warn("Could not accurately represent maker balance as a float64")
		}
		ethBacking := &meshdb.ETHBacking{
			MakerAddress: makerAddress,
			OrderCount:   0,
			AmountInWei:  amountFloat,
		}
		if err := txn.Insert(ethBacking); err != nil {
			// If we can't save the ETH backing for this maker, we have no choice but
			// to consider all its orders invalid. We might get a chance to try again
			// in the future. This shouldn't happen often.
			log.WithFields(map[string]interface{}{
				"error":        err,
				"makerAddress": ethBacking.MakerAddress,
			}).Error("Unexpected error while inserting new ETHBacking in database")
			appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, rejected, makerAddressToOrders[makerAddress])
			continue
		} else {
			log.WithField("makerAddress", ethBacking.MakerAddress.Hex()).Info("Inserted ETHBacking for makerAddress")
			valid = append(valid, makerAddressToOrders[makerAddress]...)
		}
	}

	// Add any orders for maker addresses for which we failed to get the balance
	// to the set of rejected orders.
	for _, failedAddress := range failedBalanceMakerAddresses {
		log.WithFields(map[string]interface{}{
			"makerAddress": failedAddress,
		}).Error("Could not get balance for maker address")
		rejected = appendRejectedOrderInfoForOrders(zeroex.MeshError, zeroex.ROEthRPCRequestFailed, rejected, makerAddressToOrders[failedAddress])
	}

	if err := txn.Commit(); err != nil {
		log.WithField("error", err).Error("Unexpected error while inserting new ETHBackings")
		allRejected := []*zeroex.RejectedOrderInfo{}
		appendRejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, allRejected, orders)
		return nil, allRejected
	}

	return valid, rejected
}

// validateETHBackingsWithHeap is the core algorithm for validating ETH
// backings. It is a pure function whose output depends only on its input. It
// doesn't make any network requests or read from/write to the database.
func validateETHBackingsWithHeap(spareCapacity int, ethBackings []*meshdb.ETHBacking, incomingOrders []*zeroex.SignedOrder) (ordersWithSufficientBacking []*zeroex.SignedOrder, rejectedOrders []*zeroex.RejectedOrderInfo) {
	if len(incomingOrders) == 0 {
		return ordersWithSufficientBacking, rejectedOrders
	}

	// Initialize a heap which will keep track of the maker address with the
	// lowest ETH per order.
	ethBackingHeap := ETHBackingHeap(ethBackings)
	heap.Init(&ethBackingHeap)

	// rejected will store the RejectedOrderInfo for any incoming orders that are rejected.
	rejected := []*zeroex.RejectedOrderInfo{}

	// Create a mapping of maker address to valid orders (so far) for that maker
	// address.
	makerAddressToValidOrders := map[common.Address][]*zeroex.SignedOrder{}

	// If we have spare capacity left, consider all orders valid for now. Some of
	// them may become invalid.
	for i := 0; i < spareCapacity; i++ {
		order := incomingOrders[i]
		if _, found := makerAddressToValidOrders[order.MakerAddress]; found {
			makerAddressToValidOrders[order.MakerAddress] = append(makerAddressToValidOrders[order.MakerAddress], order)
		} else {
			makerAddressToValidOrders[order.MakerAddress] = []*zeroex.SignedOrder{order}
		}
		ethBackingHeap.UpdateByMakerAddress(order.MakerAddress, 1)
	}

	// Group remaining orders by makerAddress. Some of them may become valid.
	remainingOrders := incomingOrders[spareCapacity:]
	makerAddressToOrders := map[common.Address][]*zeroex.SignedOrder{}
	for _, order := range remainingOrders {
		if _, found := makerAddressToOrders[order.MakerAddress]; found {
			makerAddressToOrders[order.MakerAddress] = append(makerAddressToOrders[order.MakerAddress], order)
		} else {
			makerAddressToOrders[order.MakerAddress] = []*zeroex.SignedOrder{order}
		}
	}

	for makerAddress, orders := range makerAddressToOrders {
		backingForMaker, _ := ethBackingHeap.FindByMakerAddress(makerAddress)
		validOrdersForMaker := []*zeroex.SignedOrder{}
		for i, incomingOrder := range orders {
			// If we don't have any spare capacity, check if the ETH backing per order
			// corresponding to this order's maker address is greater than the current
			// lowest ETH backing per order.
			lowestBacking := ethBackingHeap.Peek()
			potentialETHPerOrder := backingForMaker.AmountInWei / float64(backingForMaker.OrderCount+1)
			if potentialETHPerOrder <= lowestBacking.WeiPerOrder() {
				// If we don't have the required ETH backing, this order and all other
				// orders for this maker are considered invalid. We don't need to update
				// the heap.
				log.WithFields(map[string]interface{}{
					"makerAddress":       incomingOrder.MakerAddress,
					"invalidOrdersCount": len(orders[i:]),
				}).Trace("Maker has insufficient ETH backing")
				rejected = appendRejectedOrderInfoForOrders(MeshValidation, ROInsufficientETHBacking, rejected, orders[i:])
				break
			}

			// If we do have the required ETH backing, we need to remove one order
			// from the maker with the lowest backing.
			if ordersForLowestMaker, found := makerAddressToValidOrders[lowestBacking.MakerAddress]; found && len(ordersForLowestMaker) != 0 {
				// If the maker with the lowest backing was previously included in the
				// set of valid orders, we need to remove one order from it.
				makerAddressToValidOrders[lowestBacking.MakerAddress] = ordersForLowestMaker[1:]
				rejected = append(rejected, rejectedOrderInfoForOrder(MeshValidation, ROInsufficientETHBacking, ordersForLowestMaker[0]))
			}

			// Add this order to the set of valid orders for this maker and update the
			// heap.
			validOrdersForMaker = append(validOrdersForMaker, incomingOrder)
			ethBackingHeap.UpdateByMakerAddress(incomingOrder.MakerAddress, 1)
			ethBackingHeap.UpdateLowest(-1)
		}

		// Add this makers orders to the set of valid orders.
		if _, found := makerAddressToValidOrders[makerAddress]; found {
			makerAddressToValidOrders[makerAddress] = append(makerAddressToValidOrders[makerAddress], validOrdersForMaker...)
		} else {
			makerAddressToValidOrders[makerAddress] = validOrdersForMaker
		}
	}

	// Add the valid orders for each maker to the final result.
	allValid := []*zeroex.SignedOrder{}
	for _, validOrders := range makerAddressToValidOrders {
		allValid = append(allValid, validOrders...)
	}
	return allValid, rejected
}

func rejectedOrderInfoForOrder(kind zeroex.RejectedOrderKind, status zeroex.RejectedOrderStatus, order *zeroex.SignedOrder) *zeroex.RejectedOrderInfo {
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		log.WithField("error", err).Panic("Could not compute order hash")
	}
	return &zeroex.RejectedOrderInfo{
		OrderHash:   orderHash,
		SignedOrder: order,
		Kind:        kind,
		Status:      status,
	}
}

func appendRejectedOrderInfoForOrders(kind zeroex.RejectedOrderKind, status zeroex.RejectedOrderStatus, rejected []*zeroex.RejectedOrderInfo, orders []*zeroex.SignedOrder) []*zeroex.RejectedOrderInfo {
	for _, order := range orders {
		rejected = append(rejected, rejectedOrderInfoForOrder(kind, status, order))
	}
	return rejected
}
