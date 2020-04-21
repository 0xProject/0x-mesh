package meshdb

import "errors"

const (
	// The default miniHeaderRetentionLimit used by Mesh. This default only gets overwritten in tests.
	defaultMiniHeaderRetentionLimit = 20
	// The maximum MiniHeaders to query per page when deleting MiniHeaders
	miniHeadersMaxPerPage = 1000
)

var ErrDBFilledWithPinnedOrders = errors.New("the database is full of pinned orders; no orders can be removed in order to make space")
