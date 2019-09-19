package core

import (
	"fmt"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

var errMaxSize = fmt.Errorf("message exceeds maximum size of %d bytes", ordervalidator.MaxOrderSizeInBytes)

// JSON-schema schemas
var (
	addressSchemaLoader     = gojsonschema.NewStringLoader(`{"id":"/addressSchema","type":"string","pattern":"^0x[0-9a-fA-F]{40}$"}`)
	wholeNumberSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/wholeNumberSchema","anyOf":[{"type":"string","pattern":"^\\d+$"},{"type":"integer"}]}`)
	hexSchemaLoader         = gojsonschema.NewStringLoader(`{"id":"/hexSchema","type":"string","pattern":"^0x(([0-9a-fA-F][0-9a-fA-F])+)?$"}`)
	orderSchemaLoader       = gojsonschema.NewStringLoader(`{"id":"/orderSchema","properties":{"makerAddress":{"$ref":"/addressSchema"},"takerAddress":{"$ref":"/addressSchema"},"makerFee":{"$ref":"/wholeNumberSchema"},"takerFee":{"$ref":"/wholeNumberSchema"},"senderAddress":{"$ref":"/addressSchema"},"makerAssetAmount":{"$ref":"/wholeNumberSchema"},"takerAssetAmount":{"$ref":"/wholeNumberSchema"},"makerAssetData":{"$ref":"/hexSchema"},"takerAssetData":{"$ref":"/hexSchema"},"salt":{"$ref":"/wholeNumberSchema"},"exchangeAddress":{"$ref":"/addressSchema"},"feeRecipientAddress":{"$ref":"/addressSchema"},"expirationTimeSeconds":{"$ref":"/wholeNumberSchema"}},"required":["makerAddress","takerAddress","makerFee","takerFee","senderAddress","makerAssetAmount","takerAssetAmount","makerAssetData","takerAssetData","salt","exchangeAddress","feeRecipientAddress","expirationTimeSeconds"],"type":"object"}`)
	signedOrderSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/signedOrderSchema","allOf":[{"$ref":"/orderSchema"},{"properties":{"signature":{"$ref":"/hexSchema"}},"required":["signature"]}]}`)
	meshMessageSchemaLoader = gojsonschema.NewStringLoader(`{"id":"/meshMessageSchema","properties":{"MessageType":{"type":"string"},"Order":{"$ref":"/signedOrderSchema"}},"required":["MessageType","Order"]}`)
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
func (app *App) validateOrders(orders []*zeroex.SignedOrder) (*ordervalidator.ValidationResults, error) {
	results := &ordervalidator.ValidationResults{}
	validMeshOrders := []*zeroex.SignedOrder{}
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			log.WithField("error", err).Error("could not compute order hash")
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshError,
				Status:      ordervalidator.ROInternalError,
			})
			continue
		}
		// Note(albrow): Orders with a sender address can be canceled or invalidated
		// off-chain which is difficult to support since we need to prune
		// canceled/invalidated orders from the database. We can special-case some
		// sender addresses over time. (For example we already have support for
		// validating Coordinator orders. What we're missing is a way to effeciently
		// remove orders that are soft-canceled via the Coordinator API).
		if order.SenderAddress != constants.NullAddress {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROSenderAddressNotAllowed,
			})
			continue
		}
		// TODO(fabio): This check now also checks that the order is for the right protocol version.
		// Should we update the status code to reflect this?
		expectedDomainHash := constants.NetworkIdToDomainHash[app.networkID]
		if order.DomainHash != expectedDomainHash {
			results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
				OrderHash:   orderHash,
				SignedOrder: order,
				Kind:        ordervalidator.MeshValidation,
				Status:      ordervalidator.ROIncorrectNetwork,
			})
			continue
		}
		if err := validateOrderSize(order); err != nil {
			if err == errMaxSize {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshValidation,
					Status:      ordervalidator.ROMaxOrderSizeExceeded,
				})
				continue
			} else {
				log.WithField("error", err).Error("could not validate order size")
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshError,
					Status:      ordervalidator.ROInternalError,
				})
				continue
			}
		}

		// Check if order is already stored in DB
		var dbOrder meshdb.Order
		err = app.db.Orders.FindByID(orderHash.Bytes(), &dbOrder)
		if err != nil {
			if _, ok := err.(db.NotFoundError); !ok {
				log.WithField("error", err).Error("could not check if order was already stored")
				return nil, err
			}
		} else {
			// If stored but flagged for removal, reject it
			if dbOrder.IsRemoved {
				results.Rejected = append(results.Rejected, &ordervalidator.RejectedOrderInfo{
					OrderHash:   orderHash,
					SignedOrder: order,
					Kind:        ordervalidator.MeshValidation,
					Status:      ordervalidator.ROOrderAlreadyStoredAndUnfillable,
				})
				continue
			} else {
				// If stored but not flagged for removal, accept it without re-validation
				results.Accepted = append(results.Accepted, &ordervalidator.AcceptedOrderInfo{
					OrderHash:                orderHash,
					SignedOrder:              order,
					FillableTakerAssetAmount: dbOrder.FillableTakerAssetAmount,
					IsNew:                    false,
				})
				continue
			}
		}

		validMeshOrders = append(validMeshOrders, order)
	}
	areNewOrders := true
	zeroexResults := app.orderValidator.BatchValidate(validMeshOrders, areNewOrders, rpc.LatestBlockNumber)
	zeroexResults.Accepted = append(zeroexResults.Accepted, results.Accepted...)
	zeroexResults.Rejected = append(zeroexResults.Rejected, results.Rejected...)
	return zeroexResults, nil
}

func validateMessageSize(message *p2p.Message) error {
	if len(message.Data) > ordervalidator.MaxOrderSizeInBytes {
		return errMaxSize
	}
	return nil
}

func validateOrderSize(order *zeroex.SignedOrder) error {
	encoded, err := encodeOrder(order)
	if err != nil {
		return err
	}
	if len(encoded) > ordervalidator.MaxOrderSizeInBytes {
		return errMaxSize
	}
	return nil
}
