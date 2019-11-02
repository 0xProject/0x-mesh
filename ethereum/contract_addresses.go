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
	// TODO(albrow): Uncomment these after we have deployed V3 to each network.
	// // Mainnet
	// 1: ContractAddresses{
	// 	ERC20Proxy:          common.HexToAddress("0x95e6f48254609a6ee006f7d493c8e5fb97094cef"),
	// 	ERC721Proxy:         common.HexToAddress("0xefc70a1b18c432bdc64b596838b4d138f6bc6cad"),
	// 	Exchange:            common.HexToAddress("0x080bf510fcbf18b91105470639e9561022937712"),
	// 	Coordinator:         common.HexToAddress("0xa14857e8930acd9a882d33ec20559beb5479c8a6"),
	// 	CoordinatorRegistry: common.HexToAddress("0x45797531b873fd5e519477a070a955764c1a5b07"),
	// 	DevUtils:            common.HexToAddress("0x92d9a4d50190ae04e03914db2ee650124af844e6"),
	// 	WETH9:               common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	// 	ZRXToken:            common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
	// },
	// Ropsten
	3: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
		ERC721Proxy:         common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
		Exchange:            common.HexToAddress("0xc56388332ddfc98701fefed94535100c6166956c"),
		ERC1155Proxy:        common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
		Coordinator:         common.HexToAddress("0xad8464022213a618c96a1178a927a5ed15ad6949"),
		CoordinatorRegistry: common.HexToAddress("0x403cc23e88c17c4652fb904784d1af640a6722d9"),
		DevUtils:            common.HexToAddress("0x9a8590eebcfc53f0cc7ab5ebb8c079e9e7d4e0f5"),
		WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
		ZRXToken:            common.HexToAddress("0xff67881f8d12f372d91baae9752eb3631ff0ed00"),
	},
	// Rinkeby
	4: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0x2f5ae4f6106e89b4147651688a92256885c5f410"),
		ERC721Proxy:         common.HexToAddress("0x7656d773e11ff7383a14dcf09a9c50990481cd10"),
		ERC1155Proxy:        common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
		Exchange:            common.HexToAddress("0x3afe8aa355e086d898447732cfa5d931cfb2a792"),
		Coordinator:         common.HexToAddress("0x9ae7a6e4e4d58c36b7aa573fc06ce46dd3cb0d44"),
		CoordinatorRegistry: common.HexToAddress("0x1084b6a398e47907bae43fec3ff4b677db6e4fee"),
		DevUtils:            common.HexToAddress("0xfcbb258112485f18dd68f4b1016e48c23542fdc5"),
		WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
		ZRXToken:            common.HexToAddress("0x8080c7e4b81ecf23aa6f877cfbfd9b0c228c6ffa"),
	},
	// Kovan
	42: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
		ERC721Proxy:         common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
		Exchange:            common.HexToAddress("0xca8b1626b3b7a0da722ca9f264c4630c7d34d3b8"),
		ERC1155Proxy:        common.HexToAddress("0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"),
		Coordinator:         common.HexToAddress("0x10e0b1c2e6065ec7f290c7e3731264f9a2bf2b2d"),
		CoordinatorRegistry: common.HexToAddress("0x09fb99968c016a3ff537bf58fb3d9fe55a7975d5"),
		DevUtils:            common.HexToAddress("0x58c4fbdf9222f10ad2bef8f4d374f209135e71a5"),
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
