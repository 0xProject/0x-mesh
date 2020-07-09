package graphql

import (
	"github.com/0xProject/0x-mesh/core"
	graphql "github.com/graph-gophers/graphql-go"
)

const schemaString = `
"""
A 32-byte Keccak256 hash encoded as a hexadecimal string.
"""
scalar Hash
"""
An Ethereum address encoded as a hexadecimal string.
"""
scalar Address
"""
A BigNumber or uint256 value encoded as a numerical string.
"""
scalar BigNumber
"""
An array of arbitrary bytes encoded as a hexadecimal string.
"""
scalar Bytes
"""
A time encoded as a string using the RFC3339 standard.
"""
scalar Timestamp
"""
An arbitrary set of key-value pairs. Encoded as a JSON object.
"""
scalar Object

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
