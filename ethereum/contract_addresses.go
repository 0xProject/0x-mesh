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
	ERC20Proxy          common.Address
	ERC721Proxy         common.Address
	Exchange            common.Address
	Coordinator         common.Address
	CoordinatorRegistry common.Address
	DevUtils            common.Address
	WETH9               common.Address
	ZRXToken            common.Address
}

// NetworkIDToContractAddresses maps networkId to a mapping of contract name to Ethereum address
// on that given network
var NetworkIDToContractAddresses = map[int]ContractAddresses{
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
	// // Ropsten
	// 3: ContractAddresses{
	// 	ERC20Proxy:          common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
	// 	ERC721Proxy:         common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
	// 	Exchange:            common.HexToAddress("0xbff9493f92a3df4b0429b6d00743b3cfb4c85831"),
	// 	Coordinator:         common.HexToAddress("0x2ba02e03ee0029311e0f43715307870a3e701b53"),
	// 	CoordinatorRegistry: common.HexToAddress("0x403cc23e88c17c4652fb904784d1af640a6722d9"),
	// 	DevUtils:            common.HexToAddress("0x3e0b46bad8e374e4a110c12b832cb120dbe4a479"),
	// 	WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
	// 	ZRXToken:            common.HexToAddress("0xff67881f8d12f372d91baae9752eb3631ff0ed00"),
	// },
	// // Rinkeby
	// 4: ContractAddresses{
	// 	ERC20Proxy:          common.HexToAddress("0x3e809c563c15a295e832e37053798ddc8d6c8dab"),
	// 	ERC721Proxy:         common.HexToAddress("0x8e1ff02637cb5e39f2fa36c14706aa348b065b09"),
	// 	Exchange:            common.HexToAddress("0xbff9493f92a3df4b0429b6d00743b3cfb4c85831"),
	// 	Coordinator:         common.HexToAddress("0x2ba02e03ee0029311e0f43715307870a3e701b53"),
	// 	CoordinatorRegistry: common.HexToAddress("0x1084b6a398e47907bae43fec3ff4b677db6e4fee"),
	// 	DevUtils:            common.HexToAddress("0x2d4a9abda7b8b3605c8dbd34e3550a7467c78287"),
	// 	WETH9:               common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
	// 	ZRXToken:            common.HexToAddress("0x2727e688b8fd40b198cd5fe6e408e00494a06f07"),
	// },
	// Kovan
	42: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
		ERC721Proxy:         common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
		Exchange:            common.HexToAddress("0x30589010550762d2f0d06f650d8e8b6ade6dbf4b"),
		Coordinator:         common.HexToAddress("0x2ba02e03ee0029311e0f43715307870a3e701b53"),
		CoordinatorRegistry: common.HexToAddress("0x09fb99968c016a3ff537bf58fb3d9fe55a7975d5"),
		DevUtils:            common.HexToAddress("0x6387a62a340de79f2f0353bd05d9567fe0aca955"),
		WETH9:               common.HexToAddress("0xd0a1e359811322d97991e03f863a0c30c2cf029c"),
		ZRXToken:            common.HexToAddress("0x2002d3812f58e35f0ea1ffbf80a75a38c32175fa"),
	},
	// Ganache snapshot
	50: ContractAddresses{
		ERC20Proxy:          common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:         common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		Exchange:            common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		Coordinator:         common.HexToAddress("0x0d8b0dd11f5d34ed41d556def5f841900d5b1c6b"),
		CoordinatorRegistry: common.HexToAddress("0x1941ff73d1154774d87521d2d0aaad5d19c8df60"),
		DevUtils:            common.HexToAddress("0x38ef19fdf8e8415f18c307ed71967e19aac28ba1"),
		WETH9:               common.HexToAddress("0x0b1ba0af832d7c05fd64161e0db78e85978e8082"),
		ZRXToken:            common.HexToAddress("0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
	},
}
