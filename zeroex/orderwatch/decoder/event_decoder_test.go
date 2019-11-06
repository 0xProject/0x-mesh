package decoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

var erc20TokenAddress common.Address = common.HexToAddress("0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4")

const erc20TransferLog string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000090cf64cbb199523c893a1d519243e214b8e0b472\",\"0x000000000000000000000000fe5854255eb1eb921525fa856a3947ed2412a1d7\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000edf3e3c60\",\"blockNumber\":\"0x72628d\",\"transactionHash\":\"0xca38a891272ae2ff4654f8fa7f98bc8b2ef66cb6320745670849e91f208a228b\",\"transactionIndex\":\"0x57\",\"blockHash\":\"0xbf02aa44901301f2c7ea862a539d1ee6a2a4ae261e491a65c89f355334b3645f\",\"logIndex\":\"0x92\",\"removed\":false}"
const erc20ApprovalLog string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x000000000000000000000000cf67fdd3c580f148d20a26844b2169d52e2326db\",\"0x000000000000000000000000448a5065aebb8e423f0896e6c5d525c040f59af3\"],\"data\":\"0x0000000000000000000000000000000000000000000000000de0b6b3a7640000\",\"blockNumber\":\"0x72637c\",\"transactionHash\":\"0x7a4bb56fb212a7ef9ea5fff2010fcd905b583562a2187e3e4206d09c293f374b\",\"transactionIndex\":\"0x59\",\"blockHash\":\"0x84b4628be9d77715151dae165003eaff0bdc5f09f3d09fb736ccee7598889cdf\",\"logIndex\":\"0x57\",\"removed\":false}"
const wethWithdrawalLog string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65\",\"0x000000000000000000000000b3fa5ba98fdb56e493c4c362920289a42948294e\"],\"data\":\"0x00000000000000000000000000000000000000000000000004e8b5d353f6e400\",\"blockNumber\":\"0x726c3c\",\"transactionHash\":\"0xce1bfaad43cfb1a24cc3c85aa86c4bf867ff545cb13b3d947a2290a6890e27ac\",\"transactionIndex\":\"0x29\",\"blockHash\":\"0xd087cf26990c7d216925f07a0e3745aa4a193842e65e2215275231b069e23dfc\",\"logIndex\":\"0x38\",\"removed\":false}"
const wethDepositLog string = "{\"address\":\"0x02b3c88b805f1c6982e38ea1d40a1d83f159c3d4\",\"topics\":[\"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c\",\"0x00000000000000000000000081228ea33d680b0f51271abab1105886ecd01c2c\"],\"data\":\"0x00000000000000000000000000000000000000000000000002c68af0bb140000\",\"blockNumber\":\"0x726c20\",\"transactionHash\":\"0xd321c2d2aabe50187740b31bb4078c76c01075281816b3039af0a43f91ea9467\",\"transactionIndex\":\"0x2e\",\"blockHash\":\"0x151d07e1b6099fc4ef1f2281eec9edba0ce8df9c4e2e5bab1c6b5fcd1c09dd97\",\"logIndex\":\"0x23\",\"removed\":false}"

var erc721TokenAddress common.Address = common.HexToAddress("0x5d00d312e171be5342067c09bae883f9bcb2003b")

const erc721TransferLog string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x000000000000000000000000d8c67d024db85b271b6f6eeac5234e29c4d6bbb5\",\"0x000000000000000000000000f13685a175b95faa79db765631483ac79fb3d8e8\",\"0x000000000000000000000000000000000000000000000000000000000000c5b1\"],\"data\":\"0x\",\"blockNumber\":\"0x6f503c\",\"transactionHash\":\"0x9f2b5ef09d2cebd36ee2accd8a95eb3def06c59d984f177c134b34fa5444b102\",\"transactionIndex\":\"0x20\",\"blockHash\":\"0x8c65e77bde1be54e4ca53c1eaf0936ae136a67afe58a4a0e482560f5f98a5cab\",\"logIndex\":\"0x2d\",\"removed\":false}"
const erc721ApprovalLog string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x000000000000000000000000f4985070ce32b6b1994329df787d1acc9a2dd9e2\",\"0x0000000000000000000000000000000000000000000000000000000000000000\", \"0x000000000000000000000000000000000000000000000000000000000000a986\"],\"data\":\"0x\",\"blockNumber\":\"0x726650\",\"transactionHash\":\"0x8bf55be2fddbe9a941fd376e571cc0d6270f7b7bb87cb3c7c4476d8ed6e51bb0\",\"transactionIndex\":\"0x43\",\"blockHash\":\"0x2c14bdc4f78019146ca5fa7aeac6211c055059a00468867c2ccde1b66120e1dc\",\"logIndex\":\"0x19\",\"removed\":false}"
const incorrectERC721ApprovalLog string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\", \"0x000000000000000000000000f4985070ce32b6b1994329df787d1acc9a2dd9e2\", \"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000017c86\",\"blockNumber\":\"0x726650\",\"transactionHash\":\"0x8bf55be2fddbe9a941fd376e571cc0d6270f7b7bb87cb3c7c4476d8ed6e51bb0\",\"transactionIndex\":\"0x43\",\"blockHash\":\"0x2c14bdc4f78019146ca5fa7aeac6211c055059a00468867c2ccde1b66120e1dc\",\"logIndex\":\"0x19\",\"removed\":false}"
const incorrectERC721TransferLog string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\", \"0x0000000000000000000000000000000000000000000000000000000000000000\", \"0x00000000000000000000000052f1f2957b7ee88ee66b5b67cb765d6762304f30\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000001a740\",\"blockNumber\":\"0x726650\",\"transactionHash\":\"0x8bf55be2fddbe9a941fd376e571cc0d6270f7b7bb87cb3c7c4476d8ed6e51bb0\",\"transactionIndex\":\"0x43\",\"blockHash\":\"0x2c14bdc4f78019146ca5fa7aeac6211c055059a00468867c2ccde1b66120e1dc\",\"logIndex\":\"0x19\",\"removed\":false}"
const erc721ApprovalForAllLog string = "{\"address\":\"0x5d00d312e171be5342067c09bae883f9bcb2003b\",\"topics\":[\"0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31\",\"0x0000000000000000000000006aa0fc9fc46acb60e98439f9f89782ca78fb0990\",\"0x000000000000000000000000185b257aa51fdc45176cf1ffac6a0bfb5cf28afd\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000000001\",\"blockNumber\":\"0x725f70\",\"transactionHash\":\"0x0145607687ed9156c62abe5f42bdb8bf35ba7e4c05e0fb6f4d1addff0ff78619\",\"transactionIndex\":\"0x76\",\"blockHash\":\"0x86acc4d742f16e9a427906c1a21d68de7e26274dee9645ad84e6b3fe1e37d161\",\"logIndex\":\"0x43\",\"removed\":false}"

var erc1155TokenAddress common.Address = common.HexToAddress("0x1dc4c1cefef38a777b15aa20260a54e584b16c48")

const erc1155TransferSingleLog string = "{\"logIndex\":\"0x1d\",\"transactionIndex\":\"0x3f\",\"transactionHash\":\"0xaf8c9ead387b4ccf18b906e16e98bcda7d090f8fbd8e82e6df61f3675000bc24\",\"blockHash\":\"0x022af48193c9ae95c9fb7b0f6b132c891ecf86061afa9a34d729e95f518c07d3\",\"blockNumber\":\"0x725f88\",\"address\":\"0x1dc4c1cefef38a777b15aa20260a54e584b16c48\",\"data\":\"0x000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fa\",\"topics\":[\"0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x0000000000000000000000001d7022f5b17d2f8b695918fb48fa1089c9f85401\"],\"type\":\"mined\",\"event\":\"TransferSingle\",\"args\":{\"operator\":\"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"from\":\"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"to\":\"0x1d7022f5b17d2f8b695918fb48fa1089c9f85401\",\"id\":340282366920938463463374607431768211456,\"value\":250}}"
const erc1155TransferBatchLog string = "{\"logIndex\":\"0x1d\",\"transactionIndex\":\"0x3f\",\"transactionHash\":\"0x1eafccb276f41f6fad122eade75ee9a871a3b2d9353c5cac9e7e09c106524e24\",\"blockHash\":\"0xaa0b09ebd9b425e7f3fc5ee9a8d00e1492e798a14bf56b043c391b181b9c38b6\",\"blockNumber\":\"0x725f88\",\"address\":\"0x1dc4c1cefef38a777b15aa20260a54e584b16c48\",\"data\":\"0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000020000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001\",\"topics\":[\"0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x0000000000000000000000001d7022f5b17d2f8b695918fb48fa1089c9f85401\"],\"type\":\"mined\",\"event\":\"TransferBatch\",\"args\":{\"operator\":\"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"from\":\"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"to\":\"0x1d7022f5b17d2f8b695918fb48fa1089c9f85401\",\"ids\":[{\"_hex\":\"0x8000000000000000000000000000000200000000000000000000000000000001\"}],\"values\":[{\"_hex\":\"0x01\"}]}}"
const erc1155ApprovalForAllLog string = "{\"logIndex\":\"0x1d\",\"transactionIndex\":\"0x3f\",\"transactionHash\":\"0xf16a225332fa1de3a42c697e01065fade994d9987fff745f2faaec0bb9cbbf4b\",\"blockHash\":\"0xbd26cffbf58d54e4281297d02c1da9737e34366b857e914cf70f132b70da7ca2\",\"blockNumber\":\"0x725f88\",\"address\":\"0x1dc4c1cefef38a777b15aa20260a54e584b16c48\",\"data\":\"0x0000000000000000000000000000000000000000000000000000000000000001\",\"topics\":[\"0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x000000000000000000000000e36ea790bc9d7ab70c55260c66d52b1eca985f84\"],\"type\":\"mined\",\"event\":\"ApprovalForAll\",\"args\":{\"owner\":\"0x6ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"operator\":\"0xe36ea790bc9d7ab70c55260c66d52b1eca985f84\",\"approved\":true}}"

var exchangeAddress common.Address = common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788")

const exchangeFillLog string = "{\"address\":\"0x48bacb9266a570d521063ef5dd96e61686dbe788\",\"topics\":[\"0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x000000000000000000000000a258b39954cef5cb142fd567a46cddb31a670124\",\"0xddb8be9f6fed5209693ecce4eb127252827c1c331d661ae7a2491c80355f3fdd\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000000000000000000000000e36ea790bc9d7ab70c55260c66d52b1eca985f84000000000000000000000000e36ea790bc9d7ab70c55260c66d52b1eca985f840000000000000000000000000000000000000000000000056bc75e2d63100000000000000000000000000000000000000000000000000002b5e3af16b18800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c00000000000000000000000000000000000000000000000000000000\",\"blockNumber\":\"0x2f\",\"transactionHash\":\"0xedc057022ac01b0014f7eb921fe990d6997be58174aa31cb5af4be773c8f28ab\",\"transactionIndex\":\"0x0\",\"blockHash\":\"0x4baee7852a94e0e029d5d4e9ce6f9c953e970db021c751002807d346b301eaaa\",\"logIndex\":\"0x0\",\"removed\":false}"
const exchangeCancelLog string = "{\"address\":\"0x48bacb9266a570d521063ef5dd96e61686dbe788\",\"topics\":[\"0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x000000000000000000000000a258b39954cef5cb142fd567a46cddb31a670124\",\"0x0bd69c50d82412baa611657851a5cd4cbec05205fb204c2548289d6bd11d4ffd\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb0000000000000000000000000000000000000000000000000000000000000024f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e808200000000000000000000000000000000000000000000000000000000\",\"blockNumber\":\"0x2f\",\"transactionHash\":\"0x53c2c32ad2ee450295b1c5464ead8270cf2af6f74ebde08ad9bf3dc7712972ec\",\"transactionIndex\":\"0x0\",\"blockHash\":\"0x53dacef15e6dd06a15379a6bf4647a731661863e7fd9e4ceb941896d8f51d478\",\"logIndex\":\"0x0\",\"removed\":false}"
const exchangeCancelUpToLog string = "{\"address\":\"0x48bacb9266a570d521063ef5dd96e61686dbe788\",\"topics\":[\"0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0\",\"0x0000000000000000000000006ecbe1db9ef729cbe972c83fb886247691fb6beb\",\"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000016890e4e0eb\",\"blockNumber\":\"0x2f\",\"transactionHash\":\"0x6c53a519cf31c3bf86162f3a46037979e2a2f6d1ab917275e5f64e5a7e2a0671\",\"transactionIndex\":\"0x0\",\"blockHash\":\"0xfb0cfe4f64f2c5294b0f458a1343590ccf6af465270140d64378336d90781ff5\",\"logIndex\":\"0x0\",\"removed\":false}"

func TestDecodeERC20Transfer(t *testing.T) {
	var transferLog types.Log
	err := unmarshalLogStr(erc20TransferLog, &transferLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC20(erc20TokenAddress)
	var actualEvent ERC20TransferEvent
	err = decoder.Decode(transferLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC20TransferEvent{
		From:  common.HexToAddress("0x90CF64CbB199523C893A1D519243E214b8e0b472"),
		To:    common.HexToAddress("0xFE5854255eb1Eb921525fa856a3947Ed2412A1D7"),
		Value: big.NewInt(63874940000),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Transfer event decode")
}

func TestDecodeERC20Approval(t *testing.T) {
	var approvalLog types.Log
	err := unmarshalLogStr(erc20ApprovalLog, &approvalLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC20(erc20TokenAddress)
	var actualEvent ERC20ApprovalEvent
	err = decoder.Decode(approvalLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC20ApprovalEvent{
		Owner:   common.HexToAddress("0xcf67fdd3c580f148d20a26844b2169d52e2326db"),
		Spender: common.HexToAddress("0x448a5065aebb8e423f0896e6c5d525c040f59af3"),
		Value:   big.NewInt(1000000000000000000),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Approval event decode")

}

func TestDecodeERC721Transfer(t *testing.T) {
	var transferLog types.Log
	err := unmarshalLogStr(erc721TransferLog, &transferLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC721(erc721TokenAddress)
	var actualEvent ERC721TransferEvent
	err = decoder.Decode(transferLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC721TransferEvent{
		From:    common.HexToAddress("0xD8c67d024Db85B271b6F6EeaC5234E29C4D6bbB5"),
		To:      common.HexToAddress("0xF13685a175B95FAa79DB765631483ac79fB3D8E8"),
		TokenId: big.NewInt(50609),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Transfer event decode")
}

func TestDecodeERC721Approval(t *testing.T) {
	var approvalLog types.Log
	err := unmarshalLogStr(erc721ApprovalLog, &approvalLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC721(erc721TokenAddress)
	var actualEvent ERC721ApprovalEvent
	err = decoder.Decode(approvalLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC721ApprovalEvent{
		Owner:    common.HexToAddress("0xF4985070Ce32b6B1994329DF787D1aCc9a2dd9e2"),
		Approved: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TokenId:  big.NewInt(43398),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Approval event decode")
}

func TestDecodeIncorrectERC721Approval(t *testing.T) {
	var approvalLog types.Log
	err := unmarshalLogStr(incorrectERC721ApprovalLog, &approvalLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC721(erc721TokenAddress)
	var actualEvent ERC721ApprovalEvent
	err = decoder.Decode(approvalLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC721ApprovalEvent{
		Owner:    common.HexToAddress("0xF4985070Ce32b6B1994329DF787D1aCc9a2dd9e2"),
		Approved: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TokenId:  big.NewInt(97414),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Approval event decode")
}

func TestDecodeIncorrectERC721Transfer(t *testing.T) {
	var transferLog types.Log
	err := unmarshalLogStr(incorrectERC721TransferLog, &transferLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC721(erc721TokenAddress)
	var actualEvent ERC721TransferEvent
	err = decoder.Decode(transferLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC721TransferEvent{
		From:    common.HexToAddress("0x0000000000000000000000000000000000000000"),
		To:      common.HexToAddress("0x52F1F2957b7Ee88eE66B5b67cB765D6762304F30"),
		TokenId: big.NewInt(108352),
	}

	assert.Equal(t, expectedEvent, actualEvent, "Transfer event decode")
}

func TestDecodeERC721ApprovalForAll(t *testing.T) {
	var approvalForAllLog types.Log
	err := unmarshalLogStr(erc721ApprovalForAllLog, &approvalForAllLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC721(erc721TokenAddress)
	var actualEvent ERC721ApprovalForAllEvent
	err = decoder.Decode(approvalForAllLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC721ApprovalForAllEvent{
		Owner:    common.HexToAddress("0x6aA0FC9fc46Acb60E98439f9F89782ca78fB0990"),
		Operator: common.HexToAddress("0x185b257AA51Fdc45176cF1FfaC6a0bFB5cF28afD"),
		Approved: true,
	}

	assert.Equal(t, expectedEvent, actualEvent, "ApprovalForAll event decode")
}

func TestDecodeERC1155TransferSingle(t *testing.T) {
	var transferLog types.Log
	err := unmarshalLogStr(erc1155TransferSingleLog, &transferLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC1155(erc1155TokenAddress)
	var actualEvent ERC1155TransferSingleEvent
	err = decoder.Decode(transferLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	idStr := "340282366920938463463374607431768211456"
	id, ok := math.ParseBig256(idStr)
	if !ok {
		t.Fatal(fmt.Sprintf("Failed to parse `id` into Big.Int: %s", idStr))
	}
	expectedEvent := ERC1155TransferSingleEvent{
		Operator: common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		From:     common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		To:       common.HexToAddress("0x1D7022f5B17d2F8B695918FB48fa1089C9f85401"),
		Id:       id,
		Value:    big.NewInt(250),
	}

	assert.Equal(t, expectedEvent, actualEvent, "TransferSingle event decode")
}

func TestDecodeERC1155TransferBatch(t *testing.T) {
	var transferLog types.Log
	err := unmarshalLogStr(erc1155TransferBatchLog, &transferLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC1155(erc1155TokenAddress)
	var actualEvent ERC1155TransferBatchEvent
	err = decoder.Decode(transferLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	idStr := "57896044618658097711785492504343953927315557066662158946655541218820101242881"
	id, ok := math.ParseBig256(idStr)
	if !ok {
		t.Fatal(fmt.Sprintf("Failed to parse `id` into Big.Int: %s", idStr))
	}
	expectedEvent := ERC1155TransferBatchEvent{
		Operator: common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		From:     common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		To:       common.HexToAddress("0x1D7022f5B17d2F8B695918FB48fa1089C9f85401"),
		Ids:      []*big.Int{id},
		Values:   []*big.Int{big.NewInt(1)},
	}

	assert.Equal(t, expectedEvent, actualEvent, "TransferBatch event decode")
}

func TestDecodeERC1155ApprovalForAll(t *testing.T) {
	var approvalForAllLog types.Log
	err := unmarshalLogStr(erc1155ApprovalForAllLog, &approvalForAllLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC1155(erc1155TokenAddress)
	var actualEvent ERC1155ApprovalForAllEvent
	err = decoder.Decode(approvalForAllLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ERC1155ApprovalForAllEvent{
		Owner:    common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		Operator: common.HexToAddress("0xE36Ea790bc9d7AB70C55260C66D52b1eca985f84"),
		Approved: true,
	}

	assert.Equal(t, expectedEvent, actualEvent, "ApprovalForAll event decode")
}

func TestDecodeExchangeFill(t *testing.T) {
	var fillLog types.Log
	err := unmarshalLogStr(exchangeFillLog, &fillLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownExchange(exchangeAddress)
	var actualEvent ExchangeFillEvent
	err = decoder.Decode(fillLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ExchangeFillEvent{
		MakerAddress:           common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"),
		TakerAddress:           common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"),
		SenderAddress:          common.HexToAddress("0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"),
		FeeRecipientAddress:    common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetFilledAmount: actualEvent.MakerAssetFilledAmount,
		TakerAssetFilledAmount: actualEvent.TakerAssetFilledAmount,
		MakerFeePaid:           actualEvent.MakerFeePaid,
		TakerFeePaid:           actualEvent.TakerFeePaid,
		OrderHash:              common.HexToHash("0xddb8be9f6fed5209693ecce4eb127252827c1c331d661ae7a2491c80355f3fdd"),
		MakerAssetData:         common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:         common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		MakerFeeAssetData:      common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerFeeAssetData:      common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		ProtocolFeePaid:        actualEvent.ProtocolFeePaid,
	}
	assert.Equal(t, expectedEvent, actualEvent, "Exchange Fill event decode")

	// We do big.Int assertions separately due to the fact that the underlying
	// struct fields sometimes aren't identical.
	expectedMakerAssetFilledAmount := new(big.Int).Mul(big.NewInt(100), eighteenDecimalsInBaseUnits)
	expectedTakerAssetFilledAmount := new(big.Int).Mul(big.NewInt(50), eighteenDecimalsInBaseUnits)
	expectedMakerFeePaid := big.NewInt(0)
	expectedTakerFeePaid := big.NewInt(0)
	expectedProtocolFeePaid := big.NewInt(0)
	assertBigIntEqual(t, expectedMakerAssetFilledAmount, actualEvent.MakerAssetFilledAmount, "MakerAssetFilledAmount was not equal")
	assertBigIntEqual(t, expectedTakerAssetFilledAmount, actualEvent.TakerAssetFilledAmount, "TakerAssetFilledAmount was not equal")
	assertBigIntEqual(t, expectedMakerFeePaid, actualEvent.MakerFeePaid, "MakerFeePaid was not equal")
	assertBigIntEqual(t, expectedTakerFeePaid, actualEvent.TakerFeePaid, "TakerFeePaid was not equal")
	assertBigIntEqual(t, expectedProtocolFeePaid, actualEvent.ProtocolFeePaid, "ProtocolFeePaid was not equal")
}

func TestDecodeExchangeCancel(t *testing.T) {
	var cancelLog types.Log
	err := unmarshalLogStr(exchangeCancelLog, &cancelLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownExchange(exchangeAddress)
	var actualEvent ExchangeCancelEvent
	err = decoder.Decode(cancelLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ExchangeCancelEvent{
		MakerAddress:        common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"),
		SenderAddress:       common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"),
		FeeRecipientAddress: common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		OrderHash:           common.HexToHash("0x0bd69c50d82412baa611657851a5cd4cbec05205fb204c2548289d6bd11d4ffd"),
		MakerAssetData:      common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:      common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
	}
	assert.Equal(t, expectedEvent, actualEvent, "Exchange Cancel event decode")
}
func TestDecodeExchangeCancelUpTo(t *testing.T) {
	var cancelUpToLog types.Log
	err := unmarshalLogStr(exchangeCancelUpToLog, &cancelUpToLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownExchange(exchangeAddress)
	var actualEvent ExchangeCancelUpToEvent
	err = decoder.Decode(cancelUpToLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := ExchangeCancelUpToEvent{
		MakerAddress:       common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"),
		OrderSenderAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		OrderEpoch:         big.NewInt(1548619145451),
	}
	assert.Equal(t, expectedEvent, actualEvent, "Exchange CancelUpTo event decode")
}
func TestDecodeWethDeposit(t *testing.T) {
	var depositLog types.Log
	err := unmarshalLogStr(wethDepositLog, &depositLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC20(erc20TokenAddress)
	var actualEvent WethDepositEvent
	err = decoder.Decode(depositLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := WethDepositEvent{
		Owner: common.HexToAddress("0x81228eA33D680B0F51271aBAb1105886eCd01C2c"),
		Value: big.NewInt(200000000000000000),
	}
	assert.Equal(t, expectedEvent, actualEvent, "WETH Deposit event decode")
}
func TestDecodeWethWithdrawal(t *testing.T) {
	var withdrawalLog types.Log
	err := unmarshalLogStr(wethWithdrawalLog, &withdrawalLog)
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}
	decoder.AddKnownERC20(erc20TokenAddress)
	var actualEvent WethWithdrawalEvent
	err = decoder.Decode(withdrawalLog, &actualEvent)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedEvent := WethWithdrawalEvent{
		Owner: common.HexToAddress("0xb3fa5bA98fdB56E493C4C362920289A42948294e"),
		Value: big.NewInt(353732490000000000),
	}
	assert.Equal(t, expectedEvent, actualEvent, "WETH Withdrawal event decode")
}

func unmarshalLogStr(logStr string, out interface{}) error {
	err := json.Unmarshal([]byte(logStr), &out)
	if err != nil {
		return err
	}
	return nil
}

func assertBigIntEqual(t *testing.T, expected *big.Int, actual *big.Int, msgAndArgs ...interface{}) {
	// We need to account the fact that the represented value can be the same even
	// if the struct fields are different (e.g. 5 * 10^3 and 50 * 10^2 are the
	// same number).
	assert.Equal(t, expected.String(), actual.String(), msgAndArgs...)
}

func TestJSONMarshalUnmarshalERC20Transfer(t *testing.T) {
	expectedEvent := ERC20TransferEvent{
		From:  common.HexToAddress("0x90CF64CbB199523C893A1D519243E214b8e0b472"),
		To:    common.HexToAddress("0xFE5854255eb1Eb921525fa856a3947Ed2412A1D7"),
		Value: big.NewInt(63874940000),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC20TransferEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC20Approval(t *testing.T) {
	expectedEvent := ERC20ApprovalEvent{
		Owner:   common.HexToAddress("0xcf67fdd3c580f148d20a26844b2169d52e2326db"),
		Spender: common.HexToAddress("0x448a5065aebb8e423f0896e6c5d525c040f59af3"),
		Value:   big.NewInt(1000000000000000000),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC20ApprovalEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC721Transfer(t *testing.T) {
	expectedEvent := ERC721TransferEvent{
		From:    common.HexToAddress("0xD8c67d024Db85B271b6F6EeaC5234E29C4D6bbB5"),
		To:      common.HexToAddress("0xF13685a175B95FAa79DB765631483ac79fB3D8E8"),
		TokenId: big.NewInt(50609),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC721TransferEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC721Approval(t *testing.T) {
	expectedEvent := ERC721ApprovalEvent{
		Owner:    common.HexToAddress("0xF4985070Ce32b6B1994329DF787D1aCc9a2dd9e2"),
		Approved: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TokenId:  big.NewInt(43398),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC721ApprovalEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC721ApprovalForAll(t *testing.T) {
	expectedEvent := ERC721ApprovalForAllEvent{
		Owner:    common.HexToAddress("0x6aA0FC9fc46Acb60E98439f9F89782ca78fB0990"),
		Operator: common.HexToAddress("0x185b257AA51Fdc45176cF1FfaC6a0bFB5cF28afD"),
		Approved: true,
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC721ApprovalForAllEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC1155TransferSingle(t *testing.T) {
	expectedEvent := ERC1155TransferSingleEvent{
		Operator: common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		From:     common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		To:       common.HexToAddress("0x1D7022f5B17d2F8B695918FB48fa1089C9f85401"),
		Id:       big.NewInt(3402823669209384),
		Value:    big.NewInt(250),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC1155TransferSingleEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC1155TransferBatch(t *testing.T) {
	expectedEvent := ERC1155TransferBatchEvent{
		Operator: common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		From:     common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		To:       common.HexToAddress("0x1D7022f5B17d2F8B695918FB48fa1089C9f85401"),
		Ids:      []*big.Int{big.NewInt(3402823669209384)},
		Values:   []*big.Int{big.NewInt(1)},
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC1155TransferBatchEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalERC1155ApprovalForAll(t *testing.T) {
	expectedEvent := ERC1155ApprovalForAllEvent{
		Owner:    common.HexToAddress("0x6Ecbe1DB9EF729CBe972C83Fb886247691Fb6beb"),
		Operator: common.HexToAddress("0xE36Ea790bc9d7AB70C55260C66D52b1eca985f84"),
		Approved: true,
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ERC1155ApprovalForAllEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalExchangeFill(t *testing.T) {
	expectedEvent := ExchangeFillEvent{
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

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ExchangeFillEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalExchangeCancel(t *testing.T) {
	expectedEvent := ExchangeCancelEvent{
		MakerAddress:        common.HexToAddress("0x504a2ee3558612dB56c90186A73e690eCd57FE9E"),
		SenderAddress:       common.HexToAddress("0x504a2ee3558612dB56c90186A73e690eCd57FE9E"),
		FeeRecipientAddress: common.HexToAddress("0xA258b39954ceF5cB142fd567A46cDdB31a670124"),
		OrderHash:           common.HexToHash("0xdd50b0eec7425c3a365037a1bdeae9e12b59e06075b2bf2bdbfff8976f7419aa"),
		MakerAssetData:      common.Hex2Bytes("f47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
		TakerAssetData:      common.Hex2Bytes("f47261b000000000000000000000000089d24a6b4ccb1b6faa2625fe562bdd9a23260359"),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ExchangeCancelEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalExchangeCancelUpTo(t *testing.T) {
	expectedEvent := ExchangeCancelUpToEvent{
		MakerAddress:       common.HexToAddress("0x638C1eF824ACD48E63E6ACC84948f8eAD46f08De"),
		OrderSenderAddress: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		OrderEpoch:         big.NewInt(1554341123041),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent ExchangeCancelUpToEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}

func TestJSONMarshalUnmarshalWethDeposit(t *testing.T) {
	expectedEvent := WethDepositEvent{
		Owner: common.HexToAddress("0x81228eA33D680B0F51271aBAb1105886eCd01C2c"),
		Value: big.NewInt(200000000000000000),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent WethDepositEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}
func TestJSONMarshalUnmarshalWethWithdrawal(t *testing.T) {
	expectedEvent := WethWithdrawalEvent{
		Owner: common.HexToAddress("0xb3fa5bA98fdB56E493C4C362920289A42948294e"),
		Value: big.NewInt(353732490000000000),
	}

	buf := bytes.Buffer{}
	require.NoError(t, json.NewEncoder(&buf).Encode(expectedEvent))
	var unmarshaledEvent WethWithdrawalEvent
	require.NoError(t, json.NewDecoder(&buf).Decode(&unmarshaledEvent))
	assert.Equal(t, expectedEvent, unmarshaledEvent)
}
