package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/graphql/gqltypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

type resolver struct {
	app *core.App
}

func (r *resolver) Stats() (*gqltypes.Stats, error) {
	stats, err := r.app.GetStats()
	if err != nil {
		return nil, err
	}
	return gqltypes.StatsFromCommonType(stats), nil
}

type orderArgs struct {
	Hash string
}

func (r *resolver) Order(args orderArgs) (*gqltypes.OrderWithMetadata, error) {
	hash := common.HexToHash(args.Hash)
	if len(hash) == 0 {
		return nil, errors.New("invalid order hash")
	}
	order, err := r.app.GetOrder(hash)
	if err != nil {
		return nil, err
	}
	return gqltypes.OrderWithMetadataFromCommonType(order), nil
}

// TODO(albrow): Consider moving some conversion code to gqltypes package.

type ordersArgs struct {
	Sort    []orderSort
	Filters []orderFilter
	Limit   int32
}

type orderSort struct {
	Field     string
	Direction string
}

type orderFilter struct {
	Field string
	Kind  string
	Value filterValue
}

// filterValue corresponds to the FilterValue scalar type in the GraphQL Schema.
// We need this custom type because GraphQL doesn't ship with an "any" type.
type filterValue struct {
	value interface{}
}

func (fv *filterValue) ImplementsGraphQLType(name string) bool {
	return name == "FilterValue"
}

func (fv *filterValue) UnmarshalGraphQL(input interface{}) error {
	fv.value = input
	return nil
}

func (fv *filterValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(fv.value)
}

func (f orderFilter) getConvertedFilterValue() (interface{}, error) {
	switch f.Field {
	case "chainID", "makerAssetAmount", "makerFee", "takerAssetAmount", "takerFee", "expirationTimeSeconds", "salt", "fillableTakerAssetAmount":
		return stringToBigInt(f.Value.value)
	case "hash":
		return stringToHash(f.Value.value)
	case "exchangeAddress", "makerAddress", "takerAddress", "senderAddress", "feeRecipientAddress":
		return stringToAddress(f.Value.value)
	case "makerAssetData", "makerFeeAssetData", "takerAssetData", "takerFeeAssetData":
		return stringToBytes(f.Value.value)
	default:
		return "", fmt.Errorf("invalid filter field: %q", f.Field)
	}
}

func filterValueToString(value interface{}) (string, error) {
	valueString, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("invalid type for filter value (expected string but got %T)", value)
	}
	return valueString, nil
}

func stringToBigInt(value interface{}) (*big.Int, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return nil, err
	}
	result, valid := math.ParseBig256(valueString)
	if !valid {
		return nil, fmt.Errorf("could not convert %q to *big.Int", value)
	}
	return result, nil
}

func stringToHash(value interface{}) (common.Hash, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(valueString), nil
}

func stringToAddress(value interface{}) (common.Address, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(valueString), nil
}

func stringToBytes(value interface{}) ([]byte, error) {
	valueString, err := filterValueToString(value)
	if err != nil {
		return nil, err
	}
	return common.FromHex(valueString), nil
}

func (r *resolver) Orders(args ordersArgs) ([]*gqltypes.OrderWithMetadata, error) {
	// TODO(albrow): More validation of query args. We can assume
	//               basic structure is correct but may need to validate
	//               some of the semantics.
	query := &db.OrderQuery{
		// We never include orders that are marked as removed.
		Filters: []db.OrderFilter{
			{
				Field: db.OFIsRemoved,
				Kind:  db.Equal,
				Value: false,
			},
		},
		Limit: uint(args.Limit),
	}
	for _, filter := range args.Filters {
		kind, err := gqltypes.FilterKindToDBType(filter.Kind)
		if err != nil {
			return nil, err
		}
		value, err := filter.getConvertedFilterValue()
		if err != nil {
			return nil, err
		}
		query.Filters = append(query.Filters, db.OrderFilter{
			Field: db.OrderField(filter.Field),
			Kind:  kind,
			Value: value,
		})
	}
	for _, sort := range args.Sort {
		direction, err := gqltypes.SortDirectionToDBType(sort.Direction)
		if err != nil {
			return nil, err
		}
		query.Sort = append(query.Sort, db.OrderSort{
			Field:     db.OrderField(sort.Field),
			Direction: direction,
		})
	}

	orders, err := r.app.FindOrders(query)
	if err != nil {
		return nil, err
	}
	return gqltypes.OrdersWithMetadataFromCommonType(orders), nil
}
