package orderwatch

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/common"
	logger "github.com/sirupsen/logrus"
)

// ValidateAndStoreValidOrdersV4 applies general 0x validation and Mesh-specific validation to
// the given v4 orders and if they are valid, adds them to the OrderWatcher
func (w *Watcher) ValidateAndStoreValidOrdersV4(ctx context.Context, orders []*zeroex.SignedOrderV4, chainID int, pinned bool, opts *types.AddOrdersOpts) (*ordervalidator.ValidationResults, error) {
	if len(orders) == 0 {
		return &ordervalidator.ValidationResults{}, nil
	}
	results, validMeshOrders, err := w.meshSpecificOrderValidationV4(orders, chainID, pinned)
	if err != nil {
		return nil, err
	}

	validationBlock, zeroexResults, err := w.onchainOrderValidationV4(ctx, validMeshOrders)

	if err != nil {
		return nil, err
	}
	results.Accepted = append(results.Accepted, zeroexResults.Accepted...)
	results.Rejected = append(results.Rejected, zeroexResults.Rejected...)

	// Filter out only the new orders.
	newOrderInfos := []*ordervalidator.AcceptedOrderInfo{}
	for _, acceptedOrderInfo := range results.Accepted {
		// If the order isn't new, we don't add to OrderWatcher.
		if acceptedOrderInfo.IsNew {
			newOrderInfos = append(newOrderInfos, acceptedOrderInfo)
		}
	}

	if opts.KeepCancelled || opts.KeepExpired || opts.KeepFullyFilled || opts.KeepUnfunded {
		for _, rejectedOrderInfo := range zeroexResults.Rejected {
			// NOTE(jalextowle): We can use the rejectedOrderInfo.Status
			// field to see whether or not the order is new or not. If
			// the order has already been stored, the rejectedOrderInfo.Status
			// field will be ordervalidator.ROOrderAlreadyStoredAndUnfillable.
			// If the rejection reason involves on-chain validation, then the
			// order is new.
			if (opts.KeepCancelled && rejectedOrderInfo.Status.Code == ordervalidator.ROCancelled.Code) ||
				(opts.KeepExpired && rejectedOrderInfo.Status.Code == ordervalidator.ROExpired.Code) ||
				(opts.KeepFullyFilled && rejectedOrderInfo.Status.Code == ordervalidator.ROFullyFilled.Code) ||
				(opts.KeepUnfunded && rejectedOrderInfo.Status.Code == ordervalidator.ROUnfunded.Code) {
				newOrderInfos = append(newOrderInfos, &ordervalidator.AcceptedOrderInfo{
					OrderHash:   rejectedOrderInfo.OrderHash,
					SignedOrder: rejectedOrderInfo.SignedOrder,
					// TODO(jalextowle): Verify that this is consistent with the OrderWatcher
					FillableTakerAssetAmount: big.NewInt(0),
					IsNew:                    true,
				})
			}
		}
	}

	// Add the order to the OrderWatcher. This also saves the order in the
	// database.
	allOrderEvents := []*zeroex.OrderEvent{}
	orderEvents, err := w.add(newOrderInfos, validationBlock, pinned, opts)
	if err != nil {
		return nil, err
	}
	allOrderEvents = append(allOrderEvents, orderEvents...)

	if len(allOrderEvents) > 0 {
		// NOTE(albrow): Send can block if the subscriber(s) are slow. Blocking here can cause problems when Mesh is
		// shutting down, so to prevent that, we call Send in a goroutine and return immediately if the context
		// is done.
		done := make(chan interface{})
		go func() {
			w.orderFeed.Send(allOrderEvents)
			done <- struct{}{}
		}()
		select {
		case <-done:
			return results, nil
		case <-ctx.Done():
			return results, nil
		}
	}

	return results, nil
}

func (w *Watcher) meshSpecificOrderValidationV4(orders []*zeroex.SignedOrderV4, chainID int, pinned bool) (*ordervalidator.ValidationResults, []*zeroex.SignedOrderV4, error) {
	results := &ordervalidator.ValidationResults{}
	validMeshOrders := []*zeroex.SignedOrderV4{}

	// Calculate max expiration time based on number of orders stored.
	// This value is *exclusive*. Any incoming orders with an expiration time
	// greater or equal to this will be rejected.
	//
	// Note(albrow): Technically speaking this is sub-optimal. We are assuming
	// that we need to have space in the database for the entire slice of orders,
	// but some of them could be invalid and therefore not actually get stored.
	// However, the optimal implementation would be less efficient and could
	// result in sending more ETH RPC requests than necessary. The edge case
	// where potentially valid orders are rejected should be rare in practice, and
	// would affect at most len(orders)/2 orders.
	maxExpirationTime := constants.UnlimitedExpirationTime
	if !pinned {
		orderCount, err := w.db.CountOrdersV4(nil)
		if err != nil {
			return nil, nil, err
		}
		if orderCount+len(orders) > w.maxOrders {
			storedMaxExpirationTime, err := w.db.GetCurrentMaxExpirationTime()
			if err != nil {
				return nil, nil, err
			}
			maxExpirationTime = storedMaxExpirationTime
		}
	}

	validOrderHashes := []common.Hash{}
	for _, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			logger.WithField("error", err).Error("could not compute order hash")
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshError,
				Status:        ordervalidator.ROInternalError,
			})
			continue
		}
		if !pinned && order.Expiry.Cmp(maxExpirationTime) != -1 {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROMaxExpirationExceeded,
			})
			continue
		}
		// Note(albrow): Orders with a sender address can be canceled or invalidated
		// off-chain which is difficult to support since we need to prune
		// canceled/invalidated orders from the database. We can special-case some
		// sender addresses over time.
		if order.Sender != constants.NullAddress {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROSenderAddressNotAllowed,
			})
			continue
		}
		// NOTE(jalextowle): Orders with a taker address are only accessible to
		// one taker, which complicates more sophisticated pruning technology.
		// With this in mind, we only allow whitelisted taker addresses to be
		// propogated throughout the network. This whitelist should only include
		// addresses that correspond to contracts allow anyone to fill these
		// orders.
		// TODO(jalextowle): If any other addresses are whitelisted, create
		// a isTakerAddressWhitelisted function.
		if order.Taker != constants.NullAddress && order.Taker != w.contractAddresses.ExchangeProxyFlashWallet {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROTakerAddressNotAllowed,
			})
			continue
		}
		if order.ChainID.Cmp(big.NewInt(int64(chainID))) != 0 {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROIncorrectChain,
			})
			continue
		}
		// Only check the ExchangeAddress if we know the expected address for the
		// given chainID/networkID. If we don't know it, the order could still be
		// valid.
		expectedExchangeAddress := w.contractAddresses.ExchangeProxy
		if order.ExchangeAddress != expectedExchangeAddress {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROIncorrectExchangeAddress,
			})
			continue
		}

		if err := validateOrderSizeV4(order); err != nil {
			if err == constants.ErrMaxOrderSize {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:     orderHash,
					SignedOrderV4: order,
					Kind:          ordervalidator.MeshValidation,
					Status:        ordervalidator.ROMaxOrderSizeExceeded,
				})
				continue
			} else {
				logger.WithField("error", err).Error("could not validate order size")
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:     orderHash,
					SignedOrderV4: order,
					Kind:          ordervalidator.MeshError,
					Status:        ordervalidator.ROInternalError,
				})
				continue
			}
		}

		validOrderHashes = append(validOrderHashes, orderHash)
		validMeshOrders = append(validMeshOrders, order)
	}

	newValidOrders := []*zeroex.SignedOrderV4{}
	storedOrderStatuses, err := w.db.GetOrderStatuses(validOrderHashes)
	if err != nil {
		logger.WithField("error", err).Error("could not get stored order statuses")
		return nil, nil, err
	}
	if len(storedOrderStatuses) != len(validOrderHashes) {
		return nil, nil, errors.New("could not get stored order statuses for all orders")
	}
	for i, order := range validMeshOrders {
		orderStatus := storedOrderStatuses[i]
		orderHash := validOrderHashes[i]
		if !orderStatus.IsStored {
			// If not stored, add the order to a set of new orders.
			newValidOrders = append(newValidOrders, order)
		} else if orderStatus.IsMarkedRemoved || orderStatus.IsMarkedUnfillable {
			// If stored but marked as removed or unfillable, reject the order.
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: order,
				Kind:          ordervalidator.MeshValidation,
				Status:        ordervalidator.ROOrderAlreadyStoredAndUnfillable,
			})
		} else {
			// If stored but not marked as removed or unfillable, accept the order without re-validation
			results.Accepted = append(results.Accepted, &ordervalidator.AcceptedOrderInfo{
				OrderHash:                orderHash,
				SignedOrderV4:            order,
				FillableTakerAssetAmount: orderStatus.FillableTakerAssetAmount,
				IsNew:                    false,
			})
		}
	}

	return results, newValidOrders, nil
}

func (w *Watcher) onchainOrderValidationV4(ctx context.Context, orders []*zeroex.SignedOrderV4) (*types.MiniHeader, *ordervalidator.ValidationResults, error) {
	// HACK(fabio): While we wait for EIP-1898 support in Parity, we have no choice but to do the `eth_call`
	// at the latest known block _number_. As outlined in the `Rationale` section of EIP-1898, this approach cannot account
	// for the block being re-org'd out before the `eth_call` and then back in before the `eth_getBlockByNumber`
	// call (an unlikely but possible situation leading to an incorrect view of the world for these orders).
	// Unfortunately, this is the best we can do until EIP-1898 support in Parity.
	// Source: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1898.md#rationale
	latestBlock, err := w.getLatestBlock()
	if err != nil {
		return nil, nil, err
	}
	// This timeout of 1min is for limiting how long this call should block at the ETH RPC rate limiter
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	areNewOrders := true
	zeroexResults := w.orderValidator.BatchValidateV4(ctx, orders, areNewOrders, latestBlock)
	return latestBlock, zeroexResults, nil
}
