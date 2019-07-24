package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEthereumNetworkDetection(t *testing.T) {
	meshDB, err := meshdb.NewMeshDB("/tmp/meshdb_testing/" + uuid.New().String())
	require.NoError(t, err)

	// simulate starting up on mainnet
	err = initNetworkId(1, meshDB)
	require.NoError(t, err)

	// simulate restart on same network
	err = initNetworkId(1, meshDB)
	require.NoError(t, err)

	// should error when attempting to start on different network
	err = initNetworkId(2, meshDB)
	assert.Error(t, err)
}
