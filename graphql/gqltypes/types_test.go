package gqltypes

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type marshalUnmarshalTestCase struct {
	// original is the value we start with.
	original interface{}
	// holderType is the type of the value.
	holderType reflect.Type
}

var marshalUnmarshalTestCases = []marshalUnmarshalTestCase{
	{
		original:   BigNumber(*math.MaxBig256),
		holderType: reflect.TypeOf(BigNumber{}),
	},
	{
		original:   Bytes([]byte("abcdefg")),
		holderType: reflect.TypeOf(Bytes{}),
	},
	{
		original:   Bytes([]byte{}),
		holderType: reflect.TypeOf(Bytes{}),
	},
	{
		original:   Hash(common.HexToHash("0x6a7e632d4cf18534b7dd85c43a5b819324c9cf9640a74a227570921ca99efb56")),
		holderType: reflect.TypeOf(Hash{}),
	},
	{
		original:   Address(common.HexToAddress("0x61935cbdd02287b511119ddb11aeb42f1593b7ef")),
		holderType: reflect.TypeOf(Address{}),
	},
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	// For each test case, we marshal the original value to JSON and
	// then unmarshal it back into a value of type holderType. Then
	// we check that the original value is equal to the holder value.
	for _, tc := range marshalUnmarshalTestCases {
		buf := &bytes.Buffer{}
		require.NoError(t, json.NewEncoder(buf).Encode(tc.original))
		holderRef := reflect.New(tc.holderType).Interface()
		require.NoError(t, json.NewDecoder(buf).Decode(holderRef))
		holder := reflect.ValueOf(holderRef).Elem().Interface()
		assert.Equal(t, tc.original, holder)
	}
}

func TestMarshalUnmarshalGQL(t *testing.T) {
	// For each test case, we marshal the original value to GraphQL and
	// then unmarshal it back into a value of type holderType. Then
	// we check that the original value is equal to the holder value.
	for _, tc := range marshalUnmarshalTestCases {
		buf := &bytes.Buffer{}
		marshaler, ok := tc.original.(graphql.Marshaler)
		if !ok {
			t.Errorf("Type %T does not implement graphql.Marshaler", tc.original)
		}
		marshaler.MarshalGQL(buf)
		holderRef := reflect.New(tc.holderType).Interface()
		unmarshaler, ok := holderRef.(graphql.Unmarshaler)
		if !ok {
			t.Errorf("Type %T does not implement graphql.Unmarshaler", holderRef)
		}
		require.NoError(t, unmarshaler.UnmarshalGQL(buf.String()))
		holder := reflect.ValueOf(holderRef).Elem().Interface()
		assert.Equal(t, tc.original, holder)
	}
}
