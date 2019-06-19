package ethereum

import (
	"fmt"

	"github.com/0xProject/0x-mesh/constants"
)

// GetContractAddressesForNetworkID returns the contract name mapping for the
// given network. It returns an error if the network doesn't exist.
func GetContractAddressesForNetworkID(networkID int) (constants.ContractNameToAddress, error) {
	if contractNameToAddress, ok := constants.NetworkIDToContractAddresses[networkID]; ok {
		return contractNameToAddress, nil
	}
	return constants.ContractNameToAddress{}, fmt.Errorf("invalid network: %d", networkID)
}
