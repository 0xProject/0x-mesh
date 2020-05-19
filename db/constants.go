package db

import (
	"errors"
	"time"
)

const (
	// The default miniHeaderRetentionLimit used by Mesh. This default only gets overwritten in tests.
	defaultMiniHeaderRetentionLimit = 20
	// The maximum MiniHeaders to query per page when deleting MiniHeaders
	miniHeadersMaxPerPage = 1000
	// The amount of time to wait before timing out when connecting to the database for the first time.
	connectTimeout = 10 * time.Second
)

var (
	ErrDBFilledWithPinnedOrders = errors.New("the database is full of pinned orders; no orders can be removed in order to make space")
	ErrMetadataAlreadyExists    = errors.New("metadata already exists in the database (use UpdateMetadata instead?)")
	ErrMetadataNotFound         = errors.New("could not find existing metadata in database")
)
