package graphql

import (
	"github.com/0xProject/0x-mesh/core"
	graphql "github.com/graph-gophers/graphql-go"
)

// TODO(albrow): Consider moving this schema to a file and loading with gobindata.
// That way we would benefit from syntax highlighting and static analysis.
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
  fillableTakerAssetAmount: String!
}

"""
An enum containing all the order fields for which filters and/or sorting is supported.
"""
enum OrderField {
  hash
  chainId
  exchangeAddress
  makerAddress
  makerAssetData
  makerAssetAmount
  makerFeeAssetData
  makerFee
  takerAddress
  takerAssetData
  takerAssetAmount
  takerFeeAssetData
  takerFee
  senderAddress
  feeRecipientAddress
  expirationTimeSeconds
  salt
  fillableTakerAssetAmount
}

"""
The kind of comparison to be used in a filter.
"""
enum FilterKind {
  EQUAL
  NOT_EQUAL
  GREATER
  GREATER_OR_EQUAL
  LESS
  LESS_OR_EQUAL
}

"""
The direction to sort in. Ascending means lowest to highest. Descending means highest to lowest.
"""
enum SortDirection {
  ASC
  DESC
}

"""
The value to filter with. Must be the same type as the field you are filtering by.
"""
scalar FilterValue

"""
A filter on orders. Can be used in queries to only return orders that meet certain criteria.
"""
input OrderFilter {
  field: OrderField!
  kind: FilterKind!
  value: FilterValue!
}

"""
A sort ordering for orders. Can be used in queries to control the order in which results are returned.
"""
input OrderSort {
  field: OrderField!
  direction: SortDirection!
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
  Returns an array of orders that satisfy certain criteria.
  """
  orders(
    """
    Determines the order of the results. If more than one sort option is provided, results we be sorted by the
    first option first, then by any subsequent options. By default, orders are sorted by hash in ascending order.
    """
    sort: [OrderSort!] = [{ field: hash, direction: ASC }]
    """
    A set of filters. Only the orders that match all filters will be included in the results. By default no
    filters are used.
    """
    filters: [OrderFilter!] = []
    """
    The maximum number of orders to be included in the results. Defaults to 20.
    """
    limit: Int = 20
  ): [OrderWithMetadata!]!
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