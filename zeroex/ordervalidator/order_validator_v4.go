package ordervalidator

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jpillora/backoff"
	log "github.com/sirupsen/logrus"
)

// V4 Orders have a 32 fields. Signatures add 4 more, making 15.
// each field is encoded as 256 bytes, in hex, so 64 characters each.
const signedOrderV4AbiHexLength = 16 * 64

// BatchValidateV4 is like BatchValidate but for V4 orders
func (o *OrderValidator) BatchValidateV4(ctx context.Context, signedOrders []*zeroex.SignedOrderV4, areNewOrders bool, validationBlock *types.MiniHeader) *ValidationResults {
	if len(signedOrders) == 0 {
		return &ValidationResults{}
	}
	offchainValidSignedOrders, rejectedOrderInfos := o.BatchOffchainValidationV4(signedOrders)
	validationResults := &ValidationResults{
		Accepted: []*AcceptedOrderInfo{},
		Rejected: rejectedOrderInfos,
	}

	signedOrderChunks := [][]*zeroex.SignedOrderV4{}
	chunkSizes := o.computeOptimalChunkSizesV4(offchainValidSignedOrders)
	for _, chunkSize := range chunkSizes {
		signedOrderChunks = append(signedOrderChunks, offchainValidSignedOrders[:chunkSize])
		offchainValidSignedOrders = offchainValidSignedOrders[chunkSize:]
	}

	semaphoreChan := make(chan struct{}, concurrencyLimit)
	defer close(semaphoreChan)

	wg := &sync.WaitGroup{}
	for _, signedOrders := range signedOrderChunks {
		wg.Add(1)
		go func(signedOrders []*zeroex.SignedOrderV4) {
			defer wg.Done()

			select {
			case <-ctx.Done():
			// Blocks until a slot opens up in the semaphore. We read off of the
			// semaphore whenever onchain validation completes to allow another
			// goroutine to begin processing.
			case semaphoreChan <- struct{}{}:
				defer func() { <-semaphoreChan }()
				o.batchOnchainValidationV4(ctx, signedOrders, validationBlock, areNewOrders, validationResults)
			}
		}(signedOrders)
	}

	wg.Wait()
	return validationResults
}

// BatchOffchainValidationV4 is like BatchOffchainValidation but for V4 orders
func (o *OrderValidator) BatchOffchainValidationV4(signedOrders []*zeroex.SignedOrderV4) ([]*zeroex.SignedOrderV4, []*RejectedOrderInfo) {
	rejectedOrderInfos := []*RejectedOrderInfo{}
	offchainValidSignedOrders := []*zeroex.SignedOrderV4{}
	for _, signedOrder := range signedOrders {
		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			log.WithError(err).WithField("signedOrder", signedOrder).Error("Computing the orderHash failed unexpectedly")
		}
		if !signedOrder.Expiry.IsInt64() {
			// Shouldn't happen because we separately enforce a max expiration time.
			// See core/validation.go.
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: signedOrder,
				Kind:          MeshValidation,
				Status:        ROMaxExpirationExceeded,
			})
			continue
		}

		if signedOrder.MakerAmount.Cmp(big.NewInt(0)) == 0 {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: signedOrder,
				Kind:          ZeroExValidation,
				Status:        ROInvalidMakerAssetAmount,
			})
			continue
		}
		if signedOrder.TakerAmount.Cmp(big.NewInt(0)) == 0 {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: signedOrder,
				Kind:          ZeroExValidation,
				Status:        ROInvalidTakerAssetAmount,
			})
			continue
		}

		isSupportedSignature := signedOrder.SignatureTypeV4 == zeroex.EIP712SignatureV4 || signedOrder.SignatureTypeV4 == zeroex.EthSignSignatureV4
		if !isSupportedSignature {
			rejectedOrderInfos = append(rejectedOrderInfos, &RejectedOrderInfo{
				OrderHash:     orderHash,
				SignedOrderV4: signedOrder,
				Kind:          ZeroExValidation,
				Status:        ROInvalidSignature,
			})
			continue
		}

		offchainValidSignedOrders = append(offchainValidSignedOrders, signedOrder)
	}

	return offchainValidSignedOrders, rejectedOrderInfos
}

// batchOnchainValidationV4 is like batchOnchainValidation but for V4 orders
func (o *OrderValidator) batchOnchainValidationV4(
	ctx context.Context,
	signedOrders []*zeroex.SignedOrderV4,
	validationBlock *types.MiniHeader,
	areNewOrders bool,
	validationResults *ValidationResults,
) {
	ethOrders := []wrappers.LibNativeOrderLimitOrder{}
	for _, signedOrder := range signedOrders {
		ethOrders = append(ethOrders, signedOrder.EthereumAbiLimitOrder())
	}
	signatures := []wrappers.LibSignatureSignature{}
	for _, signedOrder := range signedOrders {
		signatures = append(signatures, signedOrder.EthereumAbiSignature())
	}

	// Attempt to make the eth_call request 4 times with an exponential back-off.
	maxDuration := 4 * time.Second
	b := &backoff.Backoff{
		Min:    250 * time.Millisecond, // First back-off length
		Max:    maxDuration,            // Longest back-off length
		Factor: 2,                      // Factor to multiple each successive back-off
	}

	for {
		opts := &bind.CallOpts{
			// HACK(albrow): From field should not be required for eth_call but
			// including it here is a workaround for a bug in Ganache. Removing
			// this line causes Ganache to crash.
			From:    constants.GanacheDummyERC721TokenAddress,
			Pending: false,
			Context: ctx,
		}
		opts.BlockNumber = validationBlock.Number

		results, err := o.exchangeV4.BatchGetLimitOrderRelevantStates(opts, ethOrders, signatures)
		fmt.Printf("### result = %+v\n", results)
		if err != nil {
			log.WithFields(log.Fields{
				"error":     err.Error(),
				"attempt":   b.Attempt(),
				"numOrders": len(ethOrders),
			}).Info("BatchGetLimitOrderRelevantStates request failed")
			d := b.Duration()
			if d == maxDuration {
				var fields log.Fields
				match, regexpErr := regexp.MatchString("abi: improperly formatted output", err.Error())
				if regexpErr != nil {
					log.WithField("error", regexpErr).Error("Unexpectedly failed to test regexp on error")
				}
				if err.Error() == "VM execution error." || match {
					fields = log.Fields{
						"error":     err.Error(),
						"numOrders": len(ethOrders),
						"orders":    ethOrders,
					}
				} else {
					fields = log.Fields{
						"error":     err.Error(),
						"numOrders": len(ethOrders),
					}
				}
				log.WithFields(fields).Warning("Gave up on GetOrderRelevantStates request after backoff limit reached")
				for _, signedOrder := range signedOrders {
					orderHash, err := signedOrder.ComputeOrderHash()
					if err != nil {
						log.WithField("error", err).Error("Unexpectedly failed to generate orderHash")
						continue
					}
					validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
						OrderHash:     orderHash,
						SignedOrderV4: signedOrder,
						Kind:          MeshError,
						Status:        ROEthRPCRequestFailed,
					})
				}
				return // Give up after 4 attempts
			}
			time.Sleep(d)
			continue
		}

		for j, orderInfo := range results.OrderInfos {
			fmt.Printf("### orderInfo = %+v\n", orderInfo)
			isValidSignature := results.IsSignatureValids[j]
			fillableTakerAssetAmount := results.ActualFillableTakerTokenAmounts[j]
			orderHash := common.Hash(orderInfo.OrderHash)
			signedOrder := signedOrders[j]
			orderStatus := zeroex.OrderStatusV4(orderInfo.Status)
			if !isValidSignature {
				orderStatus = zeroex.OS4InvalidSignature
			}
			fmt.Printf("-=- orderStatus = %+v\n", orderStatus)
			fmt.Printf("-=- orderStatus = %+v\n", zeroex.OS4Fillable)
			fmt.Printf("-=- orderStatus = %+v\n", orderStatus == zeroex.OS4Fillable)
			switch orderStatus {
			case zeroex.OS4Fillable:
				fmt.Printf("signedOrder.TakerAmount = %v\n", signedOrder.TakerAmount)
				fmt.Printf("orderInfo.TakerTokenFilledAmount = %v\n", orderInfo.TakerTokenFilledAmount)
				remainingTakerAssetAmount := big.NewInt(0).Sub(signedOrder.TakerAmount, orderInfo.TakerTokenFilledAmount)
				fmt.Printf("fillableTakerAssetAmount = %v\n", fillableTakerAssetAmount)
				fmt.Printf("remainingTakerAssetAmount = %v\n", remainingTakerAssetAmount)
				// If `fillableTakerAssetAmount` != `remainingTakerAssetAmount`, the order is partially fillable. We consider
				// partially fillable orders as invalid
				if fillableTakerAssetAmount.Cmp(remainingTakerAssetAmount) != 0 {
					fmt.Printf("Rejecting!\n")
					validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
						OrderHash:     orderHash,
						SignedOrderV4: signedOrder,
						Kind:          ZeroExValidation,
						Status:        ROUnfunded,
					})
				} else {
					validationResults.Accepted = append(validationResults.Accepted, &AcceptedOrderInfo{
						OrderHash:                orderHash,
						SignedOrderV4:            signedOrder,
						FillableTakerAssetAmount: fillableTakerAssetAmount,
						IsNew:                    areNewOrders,
					})
				}
				continue
			default:
				var status RejectedOrderStatus
				switch orderStatus {
				case zeroex.OS4Invalid:
					// TODO: Add an ROInvalid constant
					status = ROInternalError
				case zeroex.OS4Expired:
					status = ROExpired
				case zeroex.OS4Filled:
					status = ROFullyFilled
				case zeroex.OS4Cancelled:
					status = ROCancelled
				case zeroex.OS4InvalidSignature:
					status = ROInvalidSignature
				default:
					log.Errorf("Unknown order status %v", orderStatus)
					status = ROInternalError
				}
				validationResults.Rejected = append(validationResults.Rejected, &RejectedOrderInfo{
					OrderHash:     orderHash,
					SignedOrderV4: signedOrder,
					Kind:          ZeroExValidation,
					Status:        status,
				})
				continue
			}
		}
		fmt.Printf("-=# validationResults = %+v\n", validationResults)
		return
	}
}

// computeOptimalChunkSizesV4 splits the signedOrders into chunks where the payload size of each chunk
// is beneath the maxRequestContentLength. It does this by implementing a greedy algorithm which ABI
// encodes signedOrders one at a time until the computed payload size is as close to the
// maxRequestContentLength as possible.
func (o *OrderValidator) computeOptimalChunkSizesV4(signedOrders []*zeroex.SignedOrderV4) []int {
	chunkSizes := []int{}

	payloadLength := jsonRPCPayloadByteLength
	nextChunkSize := 0
	for range signedOrders {
		// TODO: With this being constant, the whole chunking mechanism probably simplifies substantially.
		encodedSignedOrderByteLength := signedOrderV4AbiHexLength
		if payloadLength+encodedSignedOrderByteLength < o.maxRequestContentLength {
			payloadLength += encodedSignedOrderByteLength
			nextChunkSize++
		} else {
			if nextChunkSize == 0 {
				// This case should never be hit since we enforce that EthereumRPCMaxContentLength >= maxOrderSizeInBytes
				log.Panic("EthereumRPCMaxContentLength is set so low, a single 0x order v4 cannot fit beneath the payload limit")
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
