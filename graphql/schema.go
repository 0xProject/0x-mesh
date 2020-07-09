package graphql

import (
	"github.com/0xProject/0x-mesh/core"
	graphql "github.com/graph-gophers/graphql-go"
)

const schemaString = `
"""
A signed 0x order along with some additional metadata about the order which is not part of the 0x protocol specification.
"""
type OrderWithMetadata {
  chainId: String!
  exchangeAddress: String!
  makerAddress: String!
  makerAssetData: String!
  makerAssetAmount: String!
  makerFeeAssetData: String!
  makerFee: String!
  takerAddress: String!
  takerAssetData: String!
  takerAssetAmount: String!
  takerFeeAssetData: String!
  takerFee: String!
  senderAddress: String!
  feeRecipientAddress: String!
  expirationTimeSeconds: String!
  salt: String!
  signature: String!
  """
  The hash, which can be used to uniquely identify an order.
  """
  hash: String!
  """
  The remaining amount of the maker asset which has not yet been filled.
  """
  remainingFillableTakerAssetAmount: String!
}

"""
The block number and block hash for the latest block that has been processed by Mesh.
"""
type LatestBlock {
  number: String!
  hash: String!
}

"""
Contains configuration options and various stats for Mesh.
"""
type Stats {
  version: String!
  pubSubTopic: String!
  rendezvous: String!
  peerID: String!
  ethereumChainID: Int!
  latestBlock: LatestBlock!
  numPeers: Int!
  numOrders: Int!
  numOrdersIncludingRemoved: Int!
  startOfCurrentUTCDay: String!
  ethRPCRequestsSentInCurrentUTCDay: Int!
  ethRPCRateLimitExpiredRequests: Int!
  maxExpirationTime: String!
}

type Query {
  """
  Returns the order with the specified hash, or null if no order is found with that hash.
  """
  order(hash: String!): OrderWithMetadata
	"""
	Returns the current stats.
	"""
	stats: Stats
}
`

func NewSchema(app *core.App) (*graphql.Schema, error) {
	// TODO(albrow): Look into more schema options.
	var opts = []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.UseStringDescriptions()}
	return graphql.ParseSchema(schemaString, &resolver{app: app}, opts...)
}
