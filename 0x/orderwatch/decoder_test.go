package orderwatch

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/core/types"
)

var ERC20_TOKEN_ADDRESS common.Address = common.HexToAddress("0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4")

const ERC20_TRANSFER_LOG_STR string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000090cf64cbb199523c893a1d519243e214b8e0b472\",\"0x000000000000000000000000fe5854255eb1eb921525fa856a3947ed2412a1d7\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000edf3e3c60\",\"blockNumber\":\"0x72628d\",\"transactionHash\":\"0xca38a891272ae2ff4654f8fa7f98bc8b2ef66cb6320745670849e91f208a228b\",\"transactionIndex\":\"0x57\",\"blockHash\":\"0xbf02aa44901301f2c7ea862a539d1ee6a2a4ae261e491a65c89f355334b3645f\",\"logIndex\":\"0x92\",\"removed\":false}"
const ERC20_APPROVAL_LOG_STR string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x000000000000000000000000cf67fdd3c580f148d20a26844b2169d52e2326db\",\"0x000000000000000000000000448a5065aebb8e423f0896e6c5d525c040f59af3\"],\"data\":\"0x0000000000000000000000000000000000000000000000000de0b6b3a7640000\",\"blockNumber\":\"0x72637c\",\"transactionHash\":\"0x7a4bb56fb212a7ef9ea5fff2010fcd905b583562a2187e3e4206d09c293f374b\",\"transactionIndex\":\"0x59\",\"blockHash\":\"0x84b4628be9d77715151dae165003eaff0bdc5f09f3d09fb736ccee7598889cdf\",\"logIndex\":\"0x57\",\"removed\":false}"
const WETH_WITHDRAWAL_LOG_STR string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65\",\"0x000000000000000000000000b3fa5ba98fdb56e493c4c362920289a42948294e\"],\"data\":\"0x00000000000000000000000000000000000000000000000004e8b5d353f6e400\",\"blockNumber\":\"0x726c3c\",\"transactionHash\":\"0xce1bfaad43cfb1a24cc3c85aa86c4bf867ff545cb13b3d947a2290a6890e27ac\",\"transactionIndex\":\"0x29\",\"blockHash\":\"0xd087cf26990c7d216925f07a0e3745aa4a193842e65e2215275231b069e23dfc\",\"logIndex\":\"0x38\",\"removed\":false}"
const WETH_DEPOSIT_LOG_STR string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c\",\"0x00000000000000000000000081228ea33d680b0f51271abab1105886ecd01c2c\"],\"data\":\"0x00000000000000000000000000000000000000000000000002c68af0bb140000\",\"blockNumber\":\"0x726c20\",\"transactionHash\":\"0xd321c2d2aabe50187740b31bb4078c76c01075281816b3039af0a43f91ea9467\",\"transactionIndex\":\"0x2e\",\"blockHash\":\"0x151d07e1b6099fc4ef1f2281eec9edba0ce8df9c4e2e5bab1c6b5fcd1c09dd97\",\"logIndex\":\"0x23\",\"removed\":false}"

var ERC721_TOKEN_ADDRESS common.Address = common.HexToAddress("0x5d00d312e171be5342067c09bae883f9bcb2003b")

const ERC721_TRANSFER_LOG_STR string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x000000000000000000000000d8c67d024db85b271b6f6eeac5234e29c4d6bbb5\",\"0x000000000000000000000000f13685a175b95faa79db765631483ac79fb3d8e8\",\"0x000000000000000000000000000000000000000000000000000000000000c5b1\"],\"data\":\"0x\",\"blockNumber\":\"0x6f503c\",\"transactionHash\":\"0x9f2b5ef09d2cebd36ee2accd8a95eb3def06c59d984f177c134b34fa5444b102\",\"transactionIndex\":\"0x20\",\"blockHash\":\"0x8c65e77bde1be54e4ca53c1eaf0936ae136a67afe58a4a0e482560f5f98a5cab\",\"logIndex\":\"0x2d\",\"removed\":false}"
const ERC721_APPROVAL_LOG_STR string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x000000000000000000000000f4985070ce32b6b1994329df787d1acc9a2dd9e2\",\"0x0000000000000000000000000000000000000000000000000000000000000000\", \"0x000000000000000000000000000000000000000000000000000000000000a986\"],\"data\":\"0x\",\"blockNumber\":\"0x726650\",\"transactionHash\":\"0x8bf55be2fddbe9a941fd376e571cc0d6270f7b7bb87cb3c7c4476d8ed6e51bb0\",\"transactionIndex\":\"0x43\",\"blockHash\":\"0x2c14bdc4f78019146ca5fa7aeac6211c055059a00468867c2ccde1b66120e1dc\",\"logIndex\":\"0x19\",\"removed\":false}"
const ERC721_APPROVAL_FOR_ALL_LOG_STR string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31\",\"0x0000000000000000000000006aa0fc9fc46acb60e98439f9f89782ca78fb0990\",\"0x000000000000000000000000185b257aa51fdc45176cf1ffac6a0bfb5cf28afd\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000000001\",\"blockNumber\":\"0x725f70\",\"transactionHash\":\"0x0145607687ed9156c62abe5f42bdb8bf35ba7e4c05e0fb6f4d1addff0ff78619\",\"transactionIndex\":\"0x76\",\"blockHash\":\"0x86acc4d742f16e9a427906c1a21d68de7e26274dee9645ad84e6b3fe1e37d161\",\"logIndex\":\"0x43\",\"removed\":false}"

var EXCHANGE_ADDRESS common.Address = common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b")

const EXCHANGE_FILL_LOG_STR string = "{\"address\":\"0x4f833a24e1f95d70f028921e27040ca56e09ab0b\",\"topics\":[\"0x0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129\",\"0x00000000000000000000000090079aabc47b5bea2dfc358d7114ade57ee39209\",\"0x00000000000000000000000061b9898c9b60a159fc91ae8026563cd226b7a0c1\",\"0xe5cd991e034cd4517cbf180307031074f3d560949fe9ddae9a06a829052dc759\"],\"data\":\"0x00000000000000000000000061b9898c9b60a159fc91ae8026563cd226b7a0c100000000000000000000000000000000000000000000000000563cd226b7a0c10000000000000000000000000000000000000000000000000082dace9d900000000000000000000000000000000000000000000000000000000081c19f850000000000000000000000000000000000000000000000000000000081c19f850000000000000000000000000000000000000000000000000000000081c19f850000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001600000000000000000000000000000000000000000000000000000000000000024f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000aa7427d8f17d87a28f5e1ba3adbb270badbe101100000000000000000000000000000000000000000000000000000000\",\"blockNumber\":\"0x725f88\",\"transactionHash\":\"0x9270762fe20a8a127d7acc386c04689ae2dda9a0d4c9ada59f9fe9c92c9fde76\",\"transactionIndex\":\"0x3f\",\"blockHash\":\"0x75f51d845afe56789c04e02681b5a1562896821a739def301583c49a9ee0dc6d\",\"logIndex\":\"0x26\",\"removed\":false}"
const EXCHANGE_CANCEL_LOG_STR string = "{\"address\":\"0x4f833a24e1f95d70f028921e27040ca56e09ab0b\",\"topics\":[\"0xdc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7\",\"0x000000000000000000000000504a2ee3558612db56c90186a73e690ecd57fe9e\",\"0x000000000000000000000000a258b39954cef5cb142fd567a46cddb31a670124\",\"0xdd50b0eec7425c3a365037a1bdeae9e12b59e06075b2bf2bdbfff8976f7419aa\"],\"data\":\"0x000000000000000000000000504a2ee3558612db56c90186a73e690ecd57fe9e000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b000000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a2326035900000000000000000000000000000000000000000000000000000000\",\"blockNumber\":\"0x725f9e\",\"transactionHash\":\"0x870afc7f4b550f621b4908c859d4c61e6740bcdd63b8969cf6d57104769a2852\",\"transactionIndex\":\"0x35\",\"blockHash\":\"0xb7e23f840464a73d2fa4b29a27864a1745cfbcc97ba735a747ec32cdd52a38da\",\"logIndex\":\"0x1d\",\"removed\":false}"
const EXCHANGE_CANCEL_UP_TO_LOG_STR string = "{\"address\":\"0x4f833a24e1f95d70f028921e27040ca56e09ab0b\",\"topics\":[\"0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0\",\"0x000000000000000000000000638c1ef824acd48e63e6acc84948f8ead46f08de\",\"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x00000000000000000000000000000000000000000000000000000169e5f353e1\",\"blockNumber\":\"0x726c1c\",\"transactionHash\":\"0x3c9f27e89e48dfa3854558ae8615979350b544330101e60d75d72f92050db0f8\",\"transactionIndex\":\"0x2c\",\"blockHash\":\"0xc631ddaaa39299998b62c2284717a56598ec86183eb64dad5434ea3aeb259a0b\",\"logIndex\":\"0x21\",\"removed\":false}"

func TestDecodeERC20Transfer(t *testing.T) {
	var transferLog types.Log
	err := json.Unmarshal([]byte(ERC20_TRANSFER_LOG_STR), &transferLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC20(ERC20_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(transferLog)
	if err != nil {
		panic(err)
	}

	expectedTransferEvent := ERC20TransferEvent{
		From:  common.HexToAddress("0x90CF64CbB199523C893A1D519243E214b8e0b472"),
		To:    common.HexToAddress("0xFE5854255eb1Eb921525fa856a3947Ed2412A1D7"),
		Value: big.NewInt(63874940000),
	}
	actualTransferEvent := decodedLog.(ERC20TransferEvent)

	assert.Equal(t, expectedTransferEvent, actualTransferEvent, "Transfer event decode")
}

func TestDecodeERC20Approval(t *testing.T) {
	var approvalLog types.Log
	err := json.Unmarshal([]byte(ERC20_APPROVAL_LOG_STR), &approvalLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC20(ERC20_TOKEN_ADDRESS)
	decodedLog, err := decoder.decodeERC20(approvalLog)
	if err != nil {
		panic(err)
	}

	expectedApprovalEvent := ERC20ApprovalEvent{
		Owner:   common.HexToAddress("0xcf67fdd3c580f148d20a26844b2169d52e2326db"),
		Spender: common.HexToAddress("0x448a5065aebb8e423f0896e6c5d525c040f59af3"),
		Value:   big.NewInt(1000000000000000000),
	}
	actualApprovalEvent := decodedLog.(ERC20ApprovalEvent)

	assert.Equal(t, expectedApprovalEvent, actualApprovalEvent, "Transfer event decode")

}

func TestDecodeERC721Transfer(t *testing.T) {
	var transferLog types.Log
	err := json.Unmarshal([]byte(ERC721_TRANSFER_LOG_STR), &transferLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC721(ERC721_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(transferLog)
	if err != nil {
		panic(err)
	}

	expectedTransferEvent := ERC721TransferEvent{
		From:    common.HexToAddress("0xD8c67d024Db85B271b6F6EeaC5234E29C4D6bbB5"),
		To:      common.HexToAddress("0xF13685a175B95FAa79DB765631483ac79fB3D8E8"),
		TokenId: big.NewInt(50609),
	}
	actualTransferEvent := decodedLog.(ERC721TransferEvent)

	assert.Equal(t, expectedTransferEvent, actualTransferEvent, "Transfer event decode")
}

func TestDecodeERC721Approval(t *testing.T) {
	var approvalLog types.Log
	err := json.Unmarshal([]byte(ERC721_APPROVAL_LOG_STR), &approvalLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC721(ERC721_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(approvalLog)
	if err != nil {
		panic(err)
	}

	expectedApprovalEvent := ERC721ApprovalEvent{
		Owner:    common.HexToAddress("0xF4985070Ce32b6B1994329DF787D1aCc9a2dd9e2"),
		Approved: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TokenId:  big.NewInt(43398),
	}
	actualApprovalEvent := decodedLog.(ERC721ApprovalEvent)

	assert.Equal(t, expectedApprovalEvent, actualApprovalEvent, "Approval event decode")
}

func TestDecodeERC721ApprovalForAll(t *testing.T) {
	var approvalForAllLog types.Log
	err := json.Unmarshal([]byte(ERC721_APPROVAL_FOR_ALL_LOG_STR), &approvalForAllLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC721(ERC721_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(approvalForAllLog)
	if err != nil {
		panic(err)
	}

	expectedApprovalForAllEvent := ERC721ApprovalForAllEvent{
		Owner:    common.HexToAddress("0x6aA0FC9fc46Acb60E98439f9F89782ca78fB0990"),
		Operator: common.HexToAddress("0x185b257AA51Fdc45176cF1FfaC6a0bFB5cF28afD"),
		Approved: true,
	}
	actualApprovalForAllEvent := decodedLog.(ERC721ApprovalForAllEvent)

	assert.Equal(t, expectedApprovalForAllEvent, actualApprovalForAllEvent, "ApprovalForAll event decode")
}

func TestDecodeExchangeFill(t *testing.T) {
	var fillLog types.Log
	err := json.Unmarshal([]byte(EXCHANGE_FILL_LOG_STR), &fillLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownExchange(EXCHANGE_ADDRESS)
	decodedLog, err := decoder.Decode(fillLog)
	if err != nil {
		panic(err)
	}

	actualFillEvent := decodedLog.(ExchangeFillEvent)
	expectedFillEvent := ExchangeFillEvent{
		MakerAddress:           common.HexToAddress("0x90079aABC47b5BeA2dFC358d7114Ade57Ee39209"),
		TakerAddress:           common.HexToAddress("0x61b9898C9b60A159fC91ae8026563cd226B7a0C1"),
		SenderAddress:          common.HexToAddress("0x00000000000000000000000000563cd226b7a0c1"),
		FeeRecipientAddress:    common.HexToAddress("0x61b9898C9b60A159fC91ae8026563cd226B7a0C1"),
		MakerAssetFilledAmount: big.NewInt(36832327913963520),
		TakerAssetFilledAmount: big.NewInt(142668604964864),
		MakerFeePaid:           big.NewInt(142668604964864),
		TakerFeePaid:           big.NewInt(142668604964864),
		OrderHash:              common.HexToHash("0xe5cd991e034cd4517cbf180307031074f3d560949fe9ddae9a06a829052dc759"),
		MakerAssetData:         common.Hex2Bytes("f47261b000000000000000000000000038ae374ecf4db50b0ff37125b591a04997106a32"),
		TakerAssetData:         common.Hex2Bytes("f47261b0000000000000000000000000aa7427d8f17d87a28f5e1ba3adbb270badbe1011"),
	}
	assert.Equal(t, expectedFillEvent, actualFillEvent, "Exchange Fill event decode")
}

func TestDecodeExchangeCancel(t *testing.T) {
	var cancelLog types.Log
	err := json.Unmarshal([]byte(EXCHANGE_CANCEL_LOG_STR), &cancelLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownExchange(EXCHANGE_ADDRESS)
	decodedLog, err := decoder.Decode(cancelLog)
	if err != nil {
		panic(err)
	}

	actualEvent := decodedLog.(ExchangeCancelEvent)
	expectedEvent := ExchangeCancelEvent{
		MakerAddress:        common.HexToAddress("0x504a2ee3558612dB56c90186A73e690eCd57FE9E"),
		SenderAddress:       common.HexToAddress("0x504a2ee3558612dB56c90186A73e690eCd57FE9E"),
		FeeRecipientAddress: common.HexToAddress("0xA258b39954ceF5cB142fd567A46cDdB31a670124"),
		OrderHash:           common.HexToHash("0xdd50b0eec7425c3a365037a1bdeae9e12b59e06075b2bf2bdbfff8976f7419aa"),
		MakerAssetData:      common.Hex2Bytes("f47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		TakerAssetData:      common.Hex2Bytes("f47261b000000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	}
	assert.Equal(t, expectedEvent, actualEvent, "Exchange Fill event decode")
}
func TestDecodeExchangeCancelUpTo(t *testing.T) {
	var cancelUpToLog types.Log
	err := json.Unmarshal([]byte(EXCHANGE_CANCEL_UP_TO_LOG_STR), &cancelUpToLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownExchange(EXCHANGE_ADDRESS)
	decodedLog, err := decoder.Decode(cancelUpToLog)
	if err != nil {
		panic(err)
	}

	actualEvent := decodedLog.(ExchangeCancelUpToEvent)
	expectedEvent := ExchangeCancelUpToEvent{
		MakerAddress:  common.HexToAddress("0x638C1eF824ACD48E63E6ACC84948f8eAD46f08De"),
		SenderAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		OrderEpoch:    big.NewInt(1554341123041),
	}
	assert.Equal(t, expectedEvent, actualEvent, "Exchange FillUpTo event decode")
}
func TestDecodeWethDeposit(t *testing.T) {
	var depositLog types.Log
	err := json.Unmarshal([]byte(WETH_DEPOSIT_LOG_STR), &depositLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC20(ERC20_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(depositLog)
	if err != nil {
		panic(err)
	}

	actualEvent := decodedLog.(WethDepositEvent)
	expectedEvent := WethDepositEvent{
		Owner: common.HexToAddress("0x81228eA33D680B0F51271aBAb1105886eCd01C2c"),
		Value: big.NewInt(200000000000000000),
	}
	assert.Equal(t, expectedEvent, actualEvent, "WETH Deposit event decode")
}
func TestDecodeWethWithdrawal(t *testing.T) {
	var withdrawalLog types.Log
	err := json.Unmarshal([]byte(WETH_WITHDRAWAL_LOG_STR), &withdrawalLog)
	if err != nil {
		panic(err)
	}

	decoder, err := NewDecoder()
	if err != nil {
		panic(err)
	}
	decoder.AddKnownERC20(ERC20_TOKEN_ADDRESS)
	decodedLog, err := decoder.Decode(withdrawalLog)
	if err != nil {
		panic(err)
	}

	actualEvent := decodedLog.(WethWithdrawalEvent)
	expectedEvent := WethWithdrawalEvent{
		Owner: common.HexToAddress("0xb3fa5bA98fdB56E493C4C362920289A42948294e"),
		Value: big.NewInt(353732490000000000),
	}
	assert.Equal(t, expectedEvent, actualEvent, "WETH Withdrawal	 event decode")
}
