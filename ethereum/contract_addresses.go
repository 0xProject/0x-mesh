package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// GetContractAddressesForNetworkID returns the contract name mapping for the
// given network. It returns an error if the network doesn't exist.
func GetContractAddressesForNetworkID(networkID int) (ContractAddresses, error) {
	if contractAddresses, ok := NetworkIDToContractAddresses[networkID]; ok {
		return contractAddresses, nil
	}
	return ContractAddresses{}, fmt.Errorf("invalid network: %d", networkID)
}

// ContractAddresses maps a contract's name to it's Ethereum address
type ContractAddresses struct {
	ERC20Proxy           common.Address
	ERC721Proxy          common.Address
	OrderValidationUtils common.Address
	Exchange             common.Address
	EthBalanceChecker    common.Address
	Coordinator          common.Address
	CoordinatorRegistry  common.Address
}

// NetworkIDToContractAddresses maps networkId to a mapping of contract name to Ethereum address
// on that given network
var NetworkIDToContractAddresses = map[int]ContractAddresses{
	// Mainnet
	1: ContractAddresses{
		ERC20Proxy:           common.HexToAddress("0x2240dab907db71e64d3e0dba4800c83b5c502d4e"),
		ERC721Proxy:          common.HexToAddress("0x208e41fb445f1bb1b6780d58356e81405f3e6127"),
		Exchange:             common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b"),
		OrderValidationUtils: common.HexToAddress("0xaccc650ae2669b9a48c525c40f28c497ff88d839"),
		EthBalanceChecker:    common.HexToAddress("0x9bc2c6ae8b1a8e3c375b6ccb55eb4273b2c3fbde"),
		Coordinator:          common.HexToAddress("0x25aae5b981ce6683cc5aeea1855d927e0b59066f"),
		CoordinatorRegistry:  common.HexToAddress("0x45797531b873fd5e519477a070a955764c1a5b07"),
	},
	// Ropsten
	3: ContractAddresses{
		ERC20Proxy:           common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
		ERC721Proxy:          common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
		Exchange:             common.HexToAddress("0x4530c0483a1633c7a1c97d2c53721caff2caaaaf"),
		OrderValidationUtils: common.HexToAddress("0x078b18930c5c16f73b3d9e8e4911d1720c98e292"),
		EthBalanceChecker:    common.HexToAddress("0xd5d960219af544b6f2f3e14a8bfd03dec12292fa"),
		Coordinator:          common.HexToAddress("0x25aae5b981ce6683cc5aeea1855d927e0b59066f"),
		CoordinatorRegistry:  common.HexToAddress("0x403cc23e88c17c4652fb904784d1af640a6722d9"),
	},
	// Rinkeby
	4: ContractAddresses{
		ERC20Proxy:           common.HexToAddress("0x3e809c563c15a295e832e37053798ddc8d6c8dab"),
		ERC721Proxy:          common.HexToAddress("0x8e1ff02637cb5e39f2fa36c14706aa348b065b09"),
		Exchange:             common.HexToAddress("0x22ebc052f43a88efa06379426120718170f2204e"),
		OrderValidationUtils: common.HexToAddress("0x2d645bce387b25bdd807be48e57f374226eadff9"),
		EthBalanceChecker:    common.HexToAddress("0x08b71282431009022eda2dda8af0fbee535e1507"),
		Coordinator:          common.HexToAddress("0x25aae5b981ce6683cc5aeea1855d927e0b59066f"),
		CoordinatorRegistry:  common.HexToAddress("0x1084b6a398e47907bae43fec3ff4b677db6e4fee"),
	},
	// Kovan
	42: ContractAddresses{
		ERC20Proxy:           common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
		ERC721Proxy:          common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
		Exchange:             common.HexToAddress("0x35dd2932454449b14cee11a94d3674a936d5d7b2"),
		OrderValidationUtils: common.HexToAddress("0x98cafab5518485d84fe3c2c495f0fcc807f1a7c0"),
		EthBalanceChecker:    common.HexToAddress("0x505aa534485bf80ee919339717cff90eb2e3364c"),
		Coordinator:          common.HexToAddress("0x25aae5b981ce6683cc5aeea1855d927e0b59066f"),
		CoordinatorRegistry:  common.HexToAddress("0x09fb99968c016a3ff537bf58fb3d9fe55a7975d5"),
	},
	// Ganache snapshot
	50: ContractAddresses{
		ERC20Proxy:           common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:          common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		Exchange:             common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		OrderValidationUtils: common.HexToAddress("0x8d42e38980ce74736c21c059b2240df09958d3c8"),
		EthBalanceChecker:    common.HexToAddress("0xa31e64ea55b9b6bbb9d6a676738e9a5b23149f84"),
		Coordinator:          common.HexToAddress("0x4d3d5c850dd5bd9d6f4adda3dd039a3c8054ca29"),
		CoordinatorRegistry:  common.HexToAddress("0xaa86dda78e9434aca114b6676fc742a18d15a1cc"),
	},
}
