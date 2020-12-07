// FIXME(jalextowle) Add any new contract addresses needed for v4
package ethereum

import (
	"fmt"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/common"
)

// ContractAddresses maps a contract's name to it's Ethereum address
type ContractAddresses struct {
	ERC20Proxy   common.Address `json:"erc20Proxy"`
	ERC721Proxy  common.Address `json:"erc721Proxy"`
	ERC1155Proxy common.Address `json:"erc1155Proxy"`
	ExchangeV3   common.Address `json:"exchangeV3"`
	ExchangeV4   common.Address `json:"exchangeV4"`
	// TODO(jalextowle): This should be removed when 0x v4 is released.
	ExchangeProxyFlashWallet common.Address `json:"exchangeProxyFlashWallet"`
	DevUtils                 common.Address `json:"devUtils"`
	WETH9                    common.Address `json:"weth9"`
	ZRXToken                 common.Address `json:"zrxToken"`
	ChaiBridge               common.Address `json:"chaiBridge"`
	ChaiToken                common.Address `json:"chaiToken"`
	MaximumGasPrice          common.Address `json:"maximumGasPrice"`
}

// GanacheAddresses The addresses that the 0x contracts were deployed to on the Ganache snapshot (chainID = 1337).
var GanacheAddresses = ganacheAddresses()

// NewContractAddressesForChainID The default contract addresses for the standard chainIDs.
func NewContractAddressesForChainID(chainID int) (ContractAddresses, error) {
	switch chainID {
	case 1:
		return ContractAddresses{
			ERC20Proxy:               common.HexToAddress("0x95e6f48254609a6ee006f7d493c8e5fb97094cef"),
			ERC721Proxy:              common.HexToAddress("0xefc70a1b18c432bdc64b596838b4d138f6bc6cad"),
			ERC1155Proxy:             common.HexToAddress("0x7eefbd48fd63d441ec7435d024ec7c5131019add"),
			ExchangeV3:               common.HexToAddress("0x61935cbdd02287b511119ddb11aeb42f1593b7ef"),
			ExchangeProxyFlashWallet: common.HexToAddress("0x22f9dcf4647084d6c31b2765f6910cd85c178c18"),
			DevUtils:                 common.HexToAddress("0xb1a3d901bad1df7d710fc8d008db7cdd6bbbffe6"),
			WETH9:                    common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
			ZRXToken:                 common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"),
			ChaiBridge:               common.HexToAddress("0x77c31eba23043b9a72d13470f3a3a311344d7438"),
			ChaiToken:                common.HexToAddress("0x06af07097c9eeb7fd685c692751d5c66db49c215"),
			MaximumGasPrice:          common.HexToAddress("0xe2bfd35306495d11e3c9db0d8de390cda24563cf"),
		}, nil
	case 3:
		return ContractAddresses{
			ERC20Proxy:               common.HexToAddress("0xb1408f4c245a23c31b98d2c626777d4c0d766caa"),
			ERC721Proxy:              common.HexToAddress("0xe654aac058bfbf9f83fcaee7793311dd82f6ddb4"),
			ERC1155Proxy:             common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
			ExchangeV3:               common.HexToAddress("0xfb2dd2a1366de37f7241c83d47da58fd503e2c64"),
			ExchangeProxyFlashWallet: common.HexToAddress("0x22f9dcf4647084d6c31b2765f6910cd85c178c18"),
			DevUtils:                 common.HexToAddress("0xb1a3d901bad1df7d710fc8d008db7cdd6bbbffe6"),
			WETH9:                    common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
			ZRXToken:                 common.HexToAddress("0xff67881f8d12f372d91baae9752eb3631ff0ed00"),
			ChaiBridge:               common.HexToAddress("0x0000000000000000000000000000000000000000"),
			ChaiToken:                common.HexToAddress("0x0000000000000000000000000000000000000000"),
			MaximumGasPrice:          common.HexToAddress("0x407b4128e9ecad8769b2332312a9f655cb9f5f3a"),
		}, nil
	case 4:
		return ContractAddresses{
			ERC20Proxy:               common.HexToAddress("0x2f5ae4f6106e89b4147651688a92256885c5f410"),
			ERC721Proxy:              common.HexToAddress("0x7656d773e11ff7383a14dcf09a9c50990481cd10"),
			ERC1155Proxy:             common.HexToAddress("0x19bb6caa3bc34d39e5a23cedfa3e6c7e7f3c931d"),
			ExchangeV3:               common.HexToAddress("0x198805e9682fceec29413059b68550f92868c129"),
			ExchangeProxyFlashWallet: common.HexToAddress("0x22f9dcf4647084d6c31b2765f6910cd85c178c18"),
			DevUtils:                 common.HexToAddress("0xb1a3d901bad1df7d710fc8d008db7cdd6bbbffe6"),
			WETH9:                    common.HexToAddress("0xc778417e063141139fce010982780140aa0cd5ab"),
			ZRXToken:                 common.HexToAddress("0x8080c7e4b81ecf23aa6f877cfbfd9b0c228c6ffa"),
			ChaiBridge:               common.HexToAddress("0x0000000000000000000000000000000000000000"),
			ChaiToken:                common.HexToAddress("0x0000000000000000000000000000000000000000"),
			MaximumGasPrice:          common.HexToAddress("0x47697b44bd89051e93b4d5857ba8e024800a74ac"),
		}, nil
	case 42:
		return ContractAddresses{
			ERC20Proxy:               common.HexToAddress("0xf1ec01d6236d3cd881a0bf0130ea25fe4234003e"),
			ERC721Proxy:              common.HexToAddress("0x2a9127c745688a165106c11cd4d647d2220af821"),
			ERC1155Proxy:             common.HexToAddress("0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"),
			ExchangeV3:               common.HexToAddress("0x4eacd0af335451709e1e7b570b8ea68edec8bc97"),
			ExchangeProxyFlashWallet: common.HexToAddress("0x22f9dcf4647084d6c31b2765f6910cd85c178c18"),
			DevUtils:                 common.HexToAddress("0xb1a3d901bad1df7d710fc8d008db7cdd6bbbffe6"),
			WETH9:                    common.HexToAddress("0xd0a1e359811322d97991e03f863a0c30c2cf029c"),
			ZRXToken:                 common.HexToAddress("0x2002d3812f58e35f0ea1ffbf80a75a38c32175fa"),
			ChaiBridge:               common.HexToAddress("0x0000000000000000000000000000000000000000"),
			ChaiToken:                common.HexToAddress("0x0000000000000000000000000000000000000000"),
			MaximumGasPrice:          common.HexToAddress("0x67a094cf028221ffdd93fc658f963151d05e2a74"),
		}, nil
	case 1337:
		return ganacheAddresses(), nil
	default:
		return ContractAddresses{}, fmt.Errorf("Cannot create contract addresses for non-standard chainID")
	}
}

func ValidateContractAddressesForChainID(chainID int, addresses ContractAddresses) error {
	if chainID == 1 {
		return fmt.Errorf("cannot add contract addresses for chainID 1: addresses for mainnet are hard-coded and cannot be changed")
	}
	if addresses.ExchangeV3 == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ExchangeV3 address is required", chainID)
	}
	if addresses.ExchangeV4 == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ExchangeV4 address is required", chainID)
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
	return nil
}

// ganacheAddresses Returns the addresses of the deployed contracts on the Ganache snapshot. This
// function allows these addresses to only be defined in one place.
func ganacheAddresses() ContractAddresses {
	return ContractAddresses{
		ERC20Proxy:               common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48"),
		ERC721Proxy:              common.HexToAddress("0x1d7022f5b17d2f8b695918fb48fa1089c9f85401"),
		ERC1155Proxy:             common.HexToAddress("0x6a4a62e5a7ed13c361b176a5f62c2ee620ac0df8"),
		ExchangeV3:               common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788"),
		ExchangeProxyFlashWallet: common.HexToAddress("0x22f9dcf4647084d6c31b2765f6910cd85c178c18"),
		DevUtils:                 common.HexToAddress("0xb23672f74749bf7916ba6827c64111a4d6de7f11"),
		WETH9:                    common.HexToAddress("0x0b1ba0af832d7c05fd64161e0db78e85978e8082"),
		ZRXToken:                 common.HexToAddress("0x871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		ChaiBridge:               common.HexToAddress("0x0000000000000000000000000000000000000000"),
		ChaiToken:                common.HexToAddress("0x0000000000000000000000000000000000000000"),
		MaximumGasPrice:          common.HexToAddress("0x2c530e4ecc573f11bd72cf5fdf580d134d25f15f"),
	}
}
