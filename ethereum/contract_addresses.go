package ethereum

import (
	"fmt"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
)

// GetContractAddressesForChainID returns the contract name mapping for the
// given chain. It returns an error if the chain doesn't exist.
func GetContractAddressesForChainID(chainID int) (ContractAddresses, error) {
	if contractAddresses, ok := ChainIDToContractAddresses[chainID]; ok {
		return contractAddresses, nil
	}
	return ContractAddresses{}, fmt.Errorf("invalid chain: %d", chainID)
}

func AddContractAddressesForChainID(chainID int, addresses ContractAddresses) error {
	if _, alreadExists := ChainIDToContractAddresses[chainID]; alreadExists {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: addresses for this chain id are already defined", chainID)
	}
	if addresses.Exchange == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: Exchange address is required", chainID)
	}
	if addresses.DevUtils == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: DevUtils address is required", chainID)
	}
	if addresses.ERC20Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC20Proxy address is required", chainID)
	}
	if addresses.ERC721Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC721Proxy address is required", chainID)
	}
	if addresses.ERC1155Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC1155Proxy address is required", chainID)
	}
	// TODO(albrow): Uncomment this if we re-add coordinator support.
	// if addresses.CoordinatorRegistry == constants.NullAddress {
	// 	return fmt.Errorf("cannot add contract addresses for chain ID %d: CoordinatorRegistry address is required", chainID)
	// }
	ChainIDToContractAddresses[chainID] = addresses
	return nil
}

// ContractAddresses maps a contract's name to it's Ethereum address
type ContractAddresses struct {
	ERC20Proxy          common.Address `json:"erc20Proxy"`
	ERC721Proxy         common.Address `json:"erc721Proxy"`
	ERC1155Proxy        common.Address `json:"erc1155Proxy"`
	Exchange            common.Address `json:"exchange"`
	Coordinator         common.Address `json:"coordinator"`
	CoordinatorRegistry common.Address `json:"coordinatorRegistry"`
	DevUtils            common.Address `json:"devUtils"`
	WETH9               common.Address `json:"weth9"`
	ZRXToken            common.Address `json:"zrxToken"`
}

// ChainIDToContractAddresses maps chainId to a mapping of contract name to Ethereum address
// on that given chain
var ChainIDToContractAddresses = map[int]ContractAddresses{
	// // Mainnet
	1: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0x95e6f48254609a6ee006f7d493c8e5fb97094cef"),
		ERC721Proxy:         common.HexToAddress("0xefc70a1b18c432bdc64b596838b4d138f6bc6cad"),
		Exchange:            common.HexToAddress("0x61935cbdd02287b511119ddb11aeb42f1593b7ef"),
		ERC1155Proxy:        common.HexToAddress("0x7eefbd48fd63d441ec7435d024ec7c5131019add"),
		Coordinator:         common.HexToAddress("0x38a795580d0f687e399913a00ddef6a17612c722"),
		CoordinatorRegistry: common.HexToAddress("0x45797531b873fd5e519477a070a955764c1a5b07"),
		DevUtils:            common.HexToAddress("0xccc2431a7335f21d9268ba62f0b32b0f2efc463f"),
		WETH9:               common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		ZRXToken:            common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
	},
	// Ropsten
	3: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
		ERC721Proxy:         common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
		Exchange:            common.HexToAddress("0xfb2dd2a1366de37f7241c83d47da58fd503e2c64"),
		ERC1155Proxy:        common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
		Coordinator:         common.HexToAddress("0x6ff734d96104965c9c1b0108f83abc46e6e501df"),
		CoordinatorRegistry: common.HexToAddress("0x403cc23e88c17c4652fb904784d1af640a6722d9"),
		DevUtils:            common.HexToAddress("0xa3858baf73430c2fa2339c731e3ffb04ea5e359c"),
		WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
		ZRXToken:            common.HexToAddress("0xff67881f8d12f372d91baae9752eb3631ff0ed00"),
	},
	// Rinkeby
	4: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0x2f5ae4f6106e89b4147651688a92256885c5f410"),
		ERC721Proxy:         common.HexToAddress("0x7656d773e11ff7383a14dcf09a9c50990481cd10"),
		ERC1155Proxy:        common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
		Exchange:            common.HexToAddress("0x198805e9682fceec29413059b68550f92868c129"),
		Coordinator:         common.HexToAddress("0x70c5385ee5ee4629ef72abd169e888c8b4a12238"),
		CoordinatorRegistry: common.HexToAddress("0x1084b6a398e47907bae43fec3ff4b677db6e4fee"),
		DevUtils:            common.HexToAddress("0xfe3b67d48c723988c81fe918da6dd4fb91459d22"),
		WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
		ZRXToken:            common.HexToAddress("0x8080c7e4b81ecf23aa6f877cfbfd9b0c228c6ffa"),
	},
	// Kovan
	42: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
		ERC721Proxy:         common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
		Exchange:            common.HexToAddress("0x4eacd0af335451709e1e7b570b8ea68edec8bc97"),
		ERC1155Proxy:        common.HexToAddress("0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"),
		Coordinator:         common.HexToAddress("0xd29e59e51e8ab5f94121efaeebd935ca4214e257"),
		CoordinatorRegistry: common.HexToAddress("0x09fb99968c016a3ff537bf58fb3d9fe55a7975d5"),
		DevUtils:            common.HexToAddress("0x80101f90446d9e19a80109ac2767c8c957debb29"),
		WETH9:               common.HexToAddress("0xd0a1e359811322d97991e03f863a0c30c2cf029c"),
		ZRXToken:            common.HexToAddress("0x2002d3812f58e35f0ea1ffbf80a75a38c32175fa"),
	},
	// Ganache snapshot
	1337: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:         common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		ERC1155Proxy:        common.HexToAddress("0x6a4a62e5a7ed13c361b176a5f62c2ee620ac0df8"),
		Exchange:            common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		Coordinator:         common.HexToAddress("0x0d8b0dd11f5d34ed41d556def5f841900d5b1c6b"),
		CoordinatorRegistry: common.HexToAddress("0x1941ff73d1154774d87521d2d0aaad5d19c8df60"),
		DevUtils:            common.HexToAddress("0x38ef19fdf8e8415f18c307ed71967e19aac28ba1"),
		WETH9:               common.HexToAddress("0x0b1ba0af832d7c05fd64161e0db78e85978e8082"),
		ZRXToken:            common.HexToAddress("0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
	},
}
