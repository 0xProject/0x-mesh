// +build !js

package core

import (
	"testing"

	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEthereumNetworkDetection(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/meshdb_testing/"+uuid.New().String(), 0)
	require.NoError(t, err)
	defer meshDB.Close()

	// simulate starting up on mainnet
	err = initNetworkID(1, meshDB)
	require.NoError(t, err)

	// simulate restart on same network
	err = initNetworkID(1, meshDB)
	require.NoError(t, err)

	// should error when attempting to start on different network
	err = initNetworkID(2, meshDB)
	assert.Error(t, err)
}
