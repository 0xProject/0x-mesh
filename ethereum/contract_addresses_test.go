package ethereum

import (
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContractAddressesForNetworkID(t *testing.T) {
	// a valid network returns no error
	_, err := GetContractAddressesForNetworkID(constants.TestNetworkID)
	require.NoError(t, err)

	// an invalid network returns an error stating the desired network id
	_, err = GetContractAddressesForNetworkID(-1)
	assert.EqualError(t, err, "invalid network: -1")
}
