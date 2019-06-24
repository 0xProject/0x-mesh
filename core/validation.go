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
		Message: "the order's maker address does not meet the minimum required ETH backing for storing the order",
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
	ordersByMakerAddress := map[common.Address][]*zeroex.SignedOrder{}
	for _, order := range orders {
		if _, found := ordersByMakerAddress[order.MakerAddress]; found {
			ordersByMakerAddress[order.MakerAddress] = append(ordersByMakerAddress[order.MakerAddress], order)
		} else {
			ordersByMakerAddress[order.MakerAddress] = []*zeroex.SignedOrder{order}
		}
	}
	var makerAddresses []common.Address
	for makerAddress := range ordersByMakerAddress {
		makerAddresses = append(makerAddresses, makerAddress)
	}

	// Get the current ETH balance for each maker.
	addressToBalance, failedAddresses := app.ethWathcher.Add(makerAddresses)

	// Add any failedAddresses to RejectedOrderInfo.
	rejected := []*zeroex.RejectedOrderInfo{}
	for _, failedAddress := range failedAddresses {
		orders := ordersByMakerAddress[failedAddress]
		rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, zeroex.ROEthRPCRequestFailed, orders)...)
	}

	// We need to open a global transaction in order to ensure that the balances
	// and number of orders per maker don't change while we're making our
	// calculations. We're also going to potentially insert some ETHBackings for
	// new maker addresses and don't want to overwrite any ETHBackings that get
	// inserted in the middle of our calculations.
	txn := app.db.Database.OpenGlobalTransaction()
	defer func() {
		_ = txn.Discard()
	}()

	// Find the current lowest ETHBacking by amount and amount per order.
	lowestETHBacking, err := app.db.GetETHBackingWithLowestETHPerOrder()
	if err != nil {
		log.WithField("error", err).Error("Could not get lowest ETH backing")
	}
	lowestETHPerOrder := lowestETHBacking.ETHPerOrder()

	// We iterate through each maker and classify all orders for that maker as
	// either valid or invalid.
	valid := []*zeroex.SignedOrder{}
	for makerAddress, ethAmount := range addressToBalance {
		// First check if there is an ETH backing for this maker already stored in
		// the database.
		var ethBacking meshdb.ETHBacking
		if err := app.db.ETHBackings.FindByID(makerAddress.Bytes(), &ethBacking); err != nil {
			if _, ok := err.(db.NotFoundError); !ok {
				// Some unexpected error occurred (not a NotFoundError). Log an error
				// and add the orders for this maker to the set of rejected orders.
				log.WithField("error", err).Error("Could not find existing ETH backing for maker address")
				orders := ordersByMakerAddress[makerAddress]
				rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, orders)...)
				continue
			}

			// If the ETHBacking was not found, insert a new one with the current
			// balance for this maker. Note that we do this regardless of whether the
			// maker has the minimum required ETH backing.
			ethBacking = meshdb.ETHBacking{
				MakerAddress: makerAddress,
				// Start with an OrderCount of 0. This will only change if and when we
				// actually insert an order.
				OrderCount: 0,
				ETHAmount:  ethAmount,
			}
			if err := txn.Insert(app.db.ETHBackings.Collection, &ethBacking); err != nil {
				// Log an error and add the orders for this maker to the set of rejected
				// orders.
				log.WithField("error", err).Error("Could not store ETH backing")
				orders := ordersByMakerAddress[makerAddress]
				rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, orders)...)
				continue
			}
		}

		// Check whether the maker would have sufficient backing if we were to
		// insert one additional order.
		ethBacking.OrderCount += 1
		if ethBacking.ETHPerOrder().Cmp(lowestETHPerOrder) != 1 {
			// The ETH per order for this maker if we inserted one more order would
			// not be greater than the lowest currently stored ETH backing. I.e. the
			// maker doesn't have sufficient balance and all their orders are invalid.
			orders := ordersByMakerAddress[makerAddress]
			rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, ROInsufficientETHBacking, orders)...)
		}

		// If we reached here, all orders for the maker should be considered valid.
		// Note that this doesn't mean that the maker has enough backing to store
		// all orders. It only means that it has enough backing to store *at least
		// one* order. We need to keep going through the validation process in order
		// to determine which orders to store (if any).
		orders := ordersByMakerAddress[makerAddress]
		valid = append(valid, orders...)
	}

	if err := txn.Commit(); err != nil {
		// If we couldn't save any new ETHBackings, we have no choice but to
		// consider all incoming orders invalid.
		log.WithField("error", err).Error("Could not commit transaction to insert new ETH backings")
		rejected = append(rejected, rejectedOrderInfoForOrders(zeroex.MeshError, ROInternalError, orders)...)
	}
	return valid, rejected
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
