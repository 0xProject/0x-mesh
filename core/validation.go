// +build !js

package core

import (
	"fmt"
	"math/big"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
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
	ROSenderAddressNotAllowed = zeroex.RejectedOrderStatus{
		Code:    "SenderAddressNotAllowed",
		Message: "orders with a senderAddress are not currently supported",
	}
)

const ROInvalidSchemaCode = "InvalidSchema"

// JSON-schema schemas
var (
	addressSchemaLoader     = gojsonschema.NewStringLoader(`{"id":"/addressSchema","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`)
	wholeNumberSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/wholeNumberSchema","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`)
	hexSchemaLoader         = gojsonschema.NewStringLoader(`{"id":"/hexSchema","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`)
	orderSchemaLoader       = gojsonschema.NewStringLoader(`{"id":"/orderSchema","properties":{"makerAddress":{"$ref":"/addressSchema"},"takerAddress":{"$ref":"/addressSchema"},"makerFee":{"$ref":"/wholeNumberSchema"},"takerFee":{"$ref":"/wholeNumberSchema"},"senderAddress":{"$ref":"/addressSchema"},"makerAssetAmount":{"$ref":"/wholeNumberSchema"},"takerAssetAmount":{"$ref":"/wholeNumberSchema"},"makerAssetData":{"$ref":"/hexSchema"},"takerAssetData":{"$ref":"/hexSchema"},"salt":{"$ref":"/wholeNumberSchema"},"exchangeAddress":{"$ref":"/addressSchema"},"feeRecipientAddress":{"$ref":"/addressSchema"},"expirationTimeSeconds":{"$ref":"/wholeNumberSchema"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","salt","exchangeAddress","feeRecipientAddress","expirationTimeSeconds"],"type":"object"}`)
	signedOrderSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/signedOrderSchema","allOf":[{"$ref":"/orderSchema"},{"properties":{"signature":{"$ref":"/hexSchema"}},"required":["signature"]}]}`)
	meshMessageSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/meshMessageSchema","properties":{"MessageType":{"type":"string"},"Order":{"$ref":"/signedOrderSchema"}},"required":["MessageType","Order"]}`)
)

// RejectedOrderKind values
const (
	MeshValidation = zeroex.RejectedOrderKind("MESH_VALIDATION")
)

func setupMeshMessageSchemaValidator() (*gojsonschema.Schema, error) {
	sl := gojsonschema.NewSchemaLoader()
	if err := sl.AddSchemas(addressSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(wholeNumberSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(hexSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(orderSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(signedOrderSchemaLoader); err != nil {
		return nil, err
	}
	schema, err := sl.Compile(meshMessageSchemaLoader)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func setupOrderSchemaValidator() (*gojsonschema.Schema, error) {
	sl := gojsonschema.NewSchemaLoader()
	if err := sl.AddSchemas(addressSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(wholeNumberSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(hexSchemaLoader); err != nil {
		return nil, err
	}
	if err := sl.AddSchemas(orderSchemaLoader); err != nil {
		return nil, err
	}
	schema, err := sl.Compile(signedOrderSchemaLoader)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

func (app *App) schemaValidateOrder(o []byte) (*gojsonschema.Result, error) {
	orderLoader := gojsonschema.NewBytesLoader(o)

	result, err := app.orderJSONSchema.Validate(orderLoader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (app *App) schemaValidateMeshMessage(o []byte) (*gojsonschema.Result, error) {
	messageLoader := gojsonschema.NewBytesLoader(o)

	result, err := app.meshMessageJSONSchema.Validate(messageLoader)
	if err != nil {
		return nil, err
	}

	return result, nil
}

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
		// Note(albrow): Orders with a sender address can be canceled or invalidated
		// off-chain which is difficult to suppport since we need to purne
		// canceled/invalidated orders from the database. We can special case some
		// sender addresses over time. (For example we already have support for
		// validating Coordinator orders. What we're missing is a way to effeciently
		// remove orders that are soft-canceled via the Coordinator API).
		if order.SenderAddress != constants.NullAddress {
			results.Rejected = append(results.Rejected, &zeroex.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        MeshValidation,
				Status:      ROSenderAddressNotAllowed,
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
	zeroexResults := app.orderValidator.BatchValidate(validMeshOrders)
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
