package decoder

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
)

var EVENT_SIGNATURES = [...]string{
	"Transfer(address,address,uint256)",                          // ERC20 & ERC721
	"Approval(address,address,uint256)",                          // ERC20 & ERC721
	"TransferSingle(address,address,address,uint256,uint256)",    // ERC1155
	"TransferBatch(address,address,address,uint256[],uint256[])", // ERC1155
	"ApprovalForAll(address,address,bool)",                       // ERC721 & ERC1155
	"Deposit(address,uint256)",                                   // WETH9
	"Withdrawal(address,uint256)",                                // WETH9
	"Fill(address,address,bytes,bytes,bytes,bytes,bytes32,address,address,uint256,uint256,uint256,uint256,uint256)", // Exchange
	"Cancel(address,address,bytes,bytes,address,bytes32)",                                                           // Exchange
	"CancelUpTo(address,address,uint256)",
}

// Includes ERC20 `Transfer` & `Approval` events as well as WETH `Deposit` & `Withdraw` events
const erc20EventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"

// Includes ERC721 `Transfer`, `Approval` & `ApprovalForAll` events
const erc721EventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// Includes ERC721 `Transfer` and `Approval` as specified in Axie Infinity contract (without index on TokenID)
const erc721EventsAbiWithoutTokenIDIndexStr = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// Includes ERC1155 `TransferSingle`, `TransferBatch` & `ApprovalForAll` events
const erc1155EventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"}]"

// Includes Exchange `Fill`, `Cancel`, `CancelUpTo` events
const exchangeEventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"name\":\"TransactionExecution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isApproved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldProtocolFeeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"ProtocolFeeMultiplier\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldProtocolFeeCollector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"ProtocolFeeCollectorAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"orderSenderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"}]"

// ERC20TransferEvent represents an ERC20 Transfer event
type ERC20TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type erc20TransferEventJSON struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC20TransferEvent type
func (e ERC20TransferEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(erc20TransferEventJSON{
		From:  e.From.Hex(),
		To:    e.To.Hex(),
		Value: e.Value.String(),
	})
}

func (e *ERC20TransferEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc20TransferEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.From = common.HexToAddress(eventJSON.From)
	e.To = common.HexToAddress(eventJSON.To)
	var ok bool
	e.Value, ok = math.ParseBig256(eventJSON.Value)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC20TransferEvent.Value: %q", eventJSON.Value)
	}
	return nil
}

// ERC20ApprovalEvent represents an ERC20 Approval event
type ERC20ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

type erc20ApprovalEventJSON struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
	Value   string `json:"value"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC20ApprovalEvent type
func (e ERC20ApprovalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(erc20ApprovalEventJSON{
		Owner:   e.Owner.Hex(),
		Spender: e.Spender.Hex(),
		Value:   e.Value.String(),
	})
}

func (e *ERC20ApprovalEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc20ApprovalEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Owner = common.HexToAddress(eventJSON.Owner)
	e.Spender = common.HexToAddress(eventJSON.Spender)
	var ok bool
	e.Value, ok = math.ParseBig256(eventJSON.Value)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC20ApprovalEvent.Value: %q", eventJSON.Value)
	}
	return nil
}

// ERC721TransferEvent represents an ERC721 Transfer event
type ERC721TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type erc721TransferEventJSON struct {
	From    string `json:"from"`
	To      string `json:"to"`
	TokenId string `json:"tokenId"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC721TransferEvent type
func (e ERC721TransferEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(erc721TransferEventJSON{
		From:    e.From.Hex(),
		To:      e.To.Hex(),
		TokenId: e.TokenId.String(),
	})
}

func (e *ERC721TransferEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc721TransferEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.From = common.HexToAddress(eventJSON.From)
	e.To = common.HexToAddress(eventJSON.To)
	var ok bool
	e.TokenId, ok = math.ParseBig256(eventJSON.TokenId)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC20ApprovalEvent.TokenId: %q", eventJSON.TokenId)
	}
	return nil
}

// ERC721ApprovalEvent represents an ERC721 Approval event
type ERC721ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

type erc721ApprovalEventJSON struct {
	Owner    string `json:"owner"`
	Approved string `json:"approved"`
	TokenId  string `json:"tokenId"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC721ApprovalEvent type
func (e ERC721ApprovalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(erc721ApprovalEventJSON{
		Owner:    e.Owner.Hex(),
		Approved: e.Approved.Hex(),
		TokenId:  e.TokenId.String(),
	})
}

func (e *ERC721ApprovalEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc721ApprovalEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Owner = common.HexToAddress(eventJSON.Owner)
	e.Approved = common.HexToAddress(eventJSON.Approved)
	var ok bool
	e.TokenId, ok = math.ParseBig256(eventJSON.TokenId)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC721ApprovalEvent.TokenId: %q", eventJSON.TokenId)
	}
	return nil
}

// ERC721ApprovalForAllEvent represents an ERC721 ApprovalForAll event
type ERC721ApprovalForAllEvent struct {
	Owner    common.Address `json:"owner"`
	Operator common.Address `json:"operator"`
	Approved bool           `json:"approved"`
}

// ERC1155ApprovalForAllEvent represents an ERC1155 ApprovalForAll event
type ERC1155ApprovalForAllEvent struct {
	Owner    common.Address `json:"owner"`
	Operator common.Address `json:"operator"`
	Approved bool           `json:"approved"`
}

// ERC1155TransferSingleEvent represents an ERC1155 TransfeSingler event
type ERC1155TransferSingleEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
}

type erc1155TransferSingleEventJSON struct {
	Operator string `json:"operator"`
	From     string `json:"from"`
	To       string `json:"to"`
	Id       string `json:"id"`
	Value    string `json:"value"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC1155TransferSingleEvent type
func (e ERC1155TransferSingleEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(erc1155TransferSingleEventJSON{
		Operator: e.Operator.Hex(),
		From:     e.From.Hex(),
		To:       e.To.Hex(),
		Id:       e.Id.String(),
		Value:    e.Value.String(),
	})
}

func (e *ERC1155TransferSingleEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc1155TransferSingleEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Operator = common.HexToAddress(eventJSON.Operator)
	e.From = common.HexToAddress(eventJSON.From)
	e.To = common.HexToAddress(eventJSON.To)
	var ok bool
	e.Id, ok = math.ParseBig256(eventJSON.Id)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC1155TransferSingleEvent.Id: %q", eventJSON.Id)
	}
	e.Value, ok = math.ParseBig256(eventJSON.Value)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ERC1155TransferSingleEvent.Value: %q", eventJSON.Value)
	}
	return nil
}

// ERC1155TransferBatchEvent represents an ERC1155 TransfeSingler event
type ERC1155TransferBatchEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
}

type erc1155TransferBatchEventJSON struct {
	Operator string   `json:"operator"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Ids      []string `json:"ids"`
	Values   []string `json:"values"`
}

// MarshalJSON implements a custom JSON marshaller for the ERC1155TransferBatchEvent type
func (e ERC1155TransferBatchEvent) MarshalJSON() ([]byte, error) {
	ids := make([]string, len(e.Ids))
	for i, id := range e.Ids {
		ids[i] = id.String()
	}
	values := make([]string, len(e.Values))
	for i, value := range e.Values {
		values[i] = value.String()
	}
	return json.Marshal(erc1155TransferBatchEventJSON{
		Operator: e.Operator.Hex(),
		From:     e.From.Hex(),
		To:       e.To.Hex(),
		Ids:      ids,
		Values:   values,
	})
}

func (e *ERC1155TransferBatchEvent) UnmarshalJSON(data []byte) error {
	var eventJSON erc1155TransferBatchEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Operator = common.HexToAddress(eventJSON.Operator)
	e.From = common.HexToAddress(eventJSON.From)
	e.To = common.HexToAddress(eventJSON.To)
	e.Ids = make([]*big.Int, len(eventJSON.Ids))
	var ok bool
	for i, idString := range eventJSON.Ids {
		e.Ids[i], ok = math.ParseBig256(idString)
		if !ok {
			return fmt.Errorf("Invalid uint256 number for ERC1155TransferBatchEvent.Ids: %v", eventJSON.Ids)
		}
	}
	e.Values = make([]*big.Int, len(eventJSON.Values))
	for i, valString := range eventJSON.Values {
		e.Values[i], ok = math.ParseBig256(valString)
		if !ok {
			return fmt.Errorf("Invalid uint256 number for ERC1155TransferBatchEvent.Values: %v", eventJSON.Values)
		}
	}

	return nil
}

// ExchangeFillEvent represents a 0x Exchange Fill event
type ExchangeFillEvent struct {
	MakerAddress           common.Address
	TakerAddress           common.Address
	SenderAddress          common.Address
	FeeRecipientAddress    common.Address
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	ProtocolFeePaid        *big.Int
	OrderHash              common.Hash
	MakerAssetData         []byte
	TakerAssetData         []byte
	MakerFeeAssetData      []byte
	TakerFeeAssetData      []byte
}

type exchangeFillEventJSON struct {
	MakerAddress           string `json:"makerAddress"`
	TakerAddress           string `json:"takerAddress"`
	SenderAddress          string `json:"senderAddress"`
	FeeRecipientAddress    string `json:"feeRecipientAddress"`
	MakerAssetFilledAmount string `json:"makerAssetFilledAmount"`
	TakerAssetFilledAmount string `json:"takerAssetFilledAmount"`
	MakerFeePaid           string `json:"makerFeePaid"`
	TakerFeePaid           string `json:"takerFeePaid"`
	ProtocolFeePaid        string `json:"protocolFeePaid"`
	OrderHash              string `json:"orderHash"`
	MakerAssetData         string `json:"makerAssetData"`
	TakerAssetData         string `json:"takerAssetData"`
	MakerFeeAssetData      string `json:"makerFeeAssetData"`
	TakerFeeAssetData      string `json:"takerFeeAssetData"`
}

// MarshalJSON implements a custom JSON marshaller for the ExchangeFillEvent type
func (e ExchangeFillEvent) MarshalJSON() ([]byte, error) {
	makerAssetData := "0x"
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := "0x"
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	makerFeeAssetData := "0x"
	if len(e.MakerFeeAssetData) != 0 {
		makerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerFeeAssetData))
	}
	takerFeeAssetData := "0x"
	if len(e.TakerFeeAssetData) != 0 {
		takerFeeAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerFeeAssetData))
	}
	return json.Marshal(exchangeFillEventJSON{
		MakerAddress:           e.MakerAddress.Hex(),
		TakerAddress:           e.TakerAddress.Hex(),
		SenderAddress:          e.SenderAddress.Hex(),
		FeeRecipientAddress:    e.FeeRecipientAddress.Hex(),
		MakerAssetFilledAmount: e.MakerAssetFilledAmount.String(),
		TakerAssetFilledAmount: e.TakerAssetFilledAmount.String(),
		MakerFeePaid:           e.MakerFeePaid.String(),
		TakerFeePaid:           e.TakerFeePaid.String(),
		ProtocolFeePaid:        e.ProtocolFeePaid.String(),
		OrderHash:              e.OrderHash.Hex(),
		MakerAssetData:         makerAssetData,
		TakerAssetData:         takerAssetData,
		MakerFeeAssetData:      makerFeeAssetData,
		TakerFeeAssetData:      takerFeeAssetData,
	})
}

func (e *ExchangeFillEvent) UnmarshalJSON(data []byte) error {
	var eventJSON exchangeFillEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.MakerAddress = common.HexToAddress(eventJSON.MakerAddress)
	e.TakerAddress = common.HexToAddress(eventJSON.TakerAddress)
	e.SenderAddress = common.HexToAddress(eventJSON.SenderAddress)
	e.FeeRecipientAddress = common.HexToAddress(eventJSON.FeeRecipientAddress)
	var ok bool
	e.MakerAssetFilledAmount, ok = math.ParseBig256(eventJSON.MakerAssetFilledAmount)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeFillEvent.MakerAssetFilledAmount: %q", eventJSON.MakerAssetFilledAmount)
	}
	e.TakerAssetFilledAmount, ok = math.ParseBig256(eventJSON.TakerAssetFilledAmount)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeFillEvent.TakerAssetFilledAmount: %q", eventJSON.TakerAssetFilledAmount)
	}
	e.MakerFeePaid, ok = math.ParseBig256(eventJSON.MakerFeePaid)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeFillEvent.MakerFeePaid: %q", eventJSON.MakerFeePaid)
	}
	e.TakerFeePaid, ok = math.ParseBig256(eventJSON.TakerFeePaid)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeFillEvent.TakerFeePaid: %q", eventJSON.TakerFeePaid)
	}
	e.ProtocolFeePaid, ok = math.ParseBig256(eventJSON.ProtocolFeePaid)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeFillEvent.ProtocolFeePaid: %q", eventJSON.ProtocolFeePaid)
	}
	e.OrderHash = common.HexToHash(eventJSON.OrderHash)
	e.MakerAssetData = common.FromHex(eventJSON.MakerAssetData)
	e.TakerAssetData = common.FromHex(eventJSON.TakerAssetData)
	e.MakerFeeAssetData = common.FromHex(eventJSON.MakerFeeAssetData)
	e.TakerFeeAssetData = common.FromHex(eventJSON.TakerFeeAssetData)

	return nil
}

// ExchangeCancelEvent represents a 0x Exchange Cancel event
type ExchangeCancelEvent struct {
	MakerAddress        common.Address
	FeeRecipientAddress common.Address
	SenderAddress       common.Address
	OrderHash           common.Hash
	MakerAssetData      []byte
	TakerAssetData      []byte
}

type exchangeCancelEventJSON struct {
	MakerAddress        string `json:"makerAddress"`
	FeeRecipientAddress string `json:"feeRecipientAddress"`
	SenderAddress       string `json:"senderAddress"`
	OrderHash           string `json:"orderHash"`
	MakerAssetData      string `json:"makerAssetData"`
	TakerAssetData      string `json:"takerAssetData"`
}

// MarshalJSON implements a custom JSON marshaller for the ExchangeCancelEvent type
func (e ExchangeCancelEvent) MarshalJSON() ([]byte, error) {
	makerAssetData := "0x"
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := "0x"
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	return json.Marshal(exchangeCancelEventJSON{
		MakerAddress:        e.MakerAddress.Hex(),
		SenderAddress:       e.SenderAddress.Hex(),
		FeeRecipientAddress: e.FeeRecipientAddress.Hex(),
		OrderHash:           e.OrderHash.Hex(),
		MakerAssetData:      makerAssetData,
		TakerAssetData:      takerAssetData,
	})
}

func (e *ExchangeCancelEvent) UnmarshalJSON(data []byte) error {
	var eventJSON exchangeCancelEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.MakerAddress = common.HexToAddress(eventJSON.MakerAddress)
	e.FeeRecipientAddress = common.HexToAddress(eventJSON.FeeRecipientAddress)
	e.SenderAddress = common.HexToAddress(eventJSON.SenderAddress)
	e.OrderHash = common.HexToHash(eventJSON.OrderHash)
	e.MakerAssetData = common.FromHex(eventJSON.MakerAssetData)
	e.TakerAssetData = common.FromHex(eventJSON.TakerAssetData)

	return nil
}

// ExchangeCancelUpToEvent represents a 0x Exchange CancelUpTo event
type ExchangeCancelUpToEvent struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	OrderEpoch         *big.Int
}

type exchangeCancelUpToEventJSON struct {
	MakerAddress       string `json:"makerAddress"`
	OrderSenderAddress string `json:"orderSenderAddress"`
	OrderEpoch         string `json:"orderEpoch"`
}

// MarshalJSON implements a custom JSON marshaller for the ExchangeCancelUpToEvent type
func (e ExchangeCancelUpToEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(exchangeCancelUpToEventJSON{
		MakerAddress:       e.MakerAddress.Hex(),
		OrderSenderAddress: e.OrderSenderAddress.Hex(),
		OrderEpoch:         e.OrderEpoch.String(),
	})
}

func (e *ExchangeCancelUpToEvent) UnmarshalJSON(data []byte) error {
	var eventJSON exchangeCancelUpToEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.MakerAddress = common.HexToAddress(eventJSON.MakerAddress)
	e.OrderSenderAddress = common.HexToAddress(eventJSON.OrderSenderAddress)
	var ok bool
	e.OrderEpoch, ok = math.ParseBig256(eventJSON.OrderEpoch)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for ExchangeCancelUpToEvent.OrderEpoch: %q", eventJSON.OrderEpoch)
	}

	return nil
}

// WethWithdrawalEvent represents a wrapped Ether Withdraw event
type WethWithdrawalEvent struct {
	Owner common.Address
	Value *big.Int
}

type wethWithdrawalEventJSON struct {
	Owner string `json:"owner"`
	Value string `json:"value"`
}

// MarshalJSON implements a custom JSON marshaller for the WethWithdrawalEvent type
func (w WethWithdrawalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(wethWithdrawalEventJSON{
		Owner: w.Owner.Hex(),
		Value: w.Value.String(),
	})
}

func (e *WethWithdrawalEvent) UnmarshalJSON(data []byte) error {
	var eventJSON wethWithdrawalEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Owner = common.HexToAddress(eventJSON.Owner)
	var ok bool
	e.Value, ok = math.ParseBig256(eventJSON.Value)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for WethWithdrawalEvent.Value: %q", eventJSON.Value)
	}

	return nil
}

// WethDepositEvent represents a wrapped Ether Deposit event
type WethDepositEvent struct {
	Owner common.Address
	Value *big.Int
}

type wethDepositEventJSON struct {
	Owner string `json:"owner"`
	Value string `json:"value"`
}

// MarshalJSON implements a custom JSON marshaller for the WethDepositEvent type
func (w WethDepositEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(wethDepositEventJSON{
		Owner: w.Owner.Hex(),
		Value: w.Value.String(),
	})
}

func (e *WethDepositEvent) UnmarshalJSON(data []byte) error {
	var eventJSON wethDepositEventJSON
	if err := json.Unmarshal(data, &eventJSON); err != nil {
		return err
	}
	e.Owner = common.HexToAddress(eventJSON.Owner)
	var ok bool
	e.Value, ok = math.ParseBig256(eventJSON.Value)
	if !ok {
		return fmt.Errorf("Invalid uint256 number for WethDepositEvent.Value: %q", eventJSON.Value)
	}

	return nil
}

// UnsupportedEventError is thrown when an unsupported topic is encountered
type UnsupportedEventError struct {
	Topics          []common.Hash
	ContractAddress common.Address
}

// Error returns the error string
func (e UnsupportedEventError) Error() string {
	hexTopics := []string{}
	for _, topic := range e.Topics {
		hexTopics = append(hexTopics, topic.Hex())
	}
	return fmt.Sprintf("unsupported event: contract address: %s, topics: %v", e.ContractAddress.Hex(), hexTopics)
}

type UntrackedTokenError struct {
	Topic        common.Hash
	TokenAddress common.Address
}

// Error returns the error string
func (e UntrackedTokenError) Error() string {
	return fmt.Sprintf("event for an untracked token: contract address: %s, topic: %s", e.TokenAddress.Hex(), e.Topic.Hex())
}

// AbiParserError is thrown when the decoder fails to parse topics from a
// retrieved log.
// NOTE(oskar): this error occurs during abi.ParseTopics and can be a result of
// certain tokens not conforming to the ERC20 event standard, for example not
// using indexed log parameters.
type AbiTopicParserError struct {
	Topics          []common.Hash
	ContractAddress common.Address
	parserError     error
}

func (e AbiTopicParserError) Error() string {
	return fmt.Sprintf("abi parser error: %s for contract address: %s, topics: %v", e.parserError.Error(), e.ContractAddress, e.Topics)
}

// Decoder decodes events relevant to the fillability of 0x orders. Since ERC20 & ERC721 events
// have the same signatures, but different meanings, all ERC20 & ERC721 contract addresses must
// be added to the decoder ahead of time.
type Decoder struct {
	knownERC20AddressesMu              sync.RWMutex
	knownERC721AddressesMu             sync.RWMutex
	knownERC1155AddressesMu            sync.RWMutex
	knownExchangeAddressesMu           sync.RWMutex
	knownERC20Addresses                map[common.Address]bool
	knownERC721Addresses               map[common.Address]bool
	knownERC1155Addresses              map[common.Address]bool
	knownExchangeAddresses             map[common.Address]bool
	erc20ABI                           abi.ABI
	erc721ABI                          abi.ABI
	erc721EventsAbiWithoutTokenIDIndex abi.ABI
	erc1155ABI                         abi.ABI
	exchangeABI                        abi.ABI
	erc20TopicToEventName              map[common.Hash]string
	erc721TopicToEventName             map[common.Hash]string
	erc1155TopicToEventName            map[common.Hash]string
	exchangeTopicToEventName           map[common.Hash]string
}

// New instantiates a new 0x order-relevant events decoder
func New() (*Decoder, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20EventsAbi))
	if err != nil {
		return nil, err
	}

	erc721ABI, err := abi.JSON(strings.NewReader(erc721EventsAbi))
	if err != nil {
		return nil, err
	}

	erc721EventsAbiWithoutTokenIDIndex, err := abi.JSON(strings.NewReader(erc721EventsAbiWithoutTokenIDIndexStr))
	if err != nil {
		return nil, err
	}

	erc1155ABI, err := abi.JSON(strings.NewReader(erc1155EventsAbi))
	if err != nil {
		return nil, err
	}

	exchangeABI, err := abi.JSON(strings.NewReader(exchangeEventsAbi))
	if err != nil {
		return nil, err
	}

	erc20TopicToEventName := map[common.Hash]string{}
	for _, event := range erc20ABI.Events {
		erc20TopicToEventName[event.ID] = event.Name
	}
	erc721TopicToEventName := map[common.Hash]string{}
	for _, event := range erc721ABI.Events {
		erc721TopicToEventName[event.ID] = event.Name
	}
	erc1155TopicToEventName := map[common.Hash]string{}
	for _, event := range erc1155ABI.Events {
		erc1155TopicToEventName[event.ID] = event.Name
	}
	exchangeTopicToEventName := map[common.Hash]string{}
	for _, event := range exchangeABI.Events {
		exchangeTopicToEventName[event.ID] = event.Name
	}

	return &Decoder{
		knownERC20Addresses:                make(map[common.Address]bool),
		knownERC721Addresses:               make(map[common.Address]bool),
		knownERC1155Addresses:              make(map[common.Address]bool),
		knownExchangeAddresses:             make(map[common.Address]bool),
		erc20ABI:                           erc20ABI,
		erc721ABI:                          erc721ABI,
		erc721EventsAbiWithoutTokenIDIndex: erc721EventsAbiWithoutTokenIDIndex,
		erc1155ABI:                         erc1155ABI,
		exchangeABI:                        exchangeABI,
		erc20TopicToEventName:              erc20TopicToEventName,
		erc721TopicToEventName:             erc721TopicToEventName,
		erc1155TopicToEventName:            erc1155TopicToEventName,
		exchangeTopicToEventName:           exchangeTopicToEventName,
	}, nil
}

// AddKnownERC20 registers the supplied contract address as an ERC20 contract. If an event is found
// from this contract address, the decoder will properly decode the `Transfer` and `Approve` events
// including the correct event parameter names.
func (d *Decoder) AddKnownERC20(address common.Address) {
	d.knownERC20AddressesMu.Lock()
	defer d.knownERC20AddressesMu.Unlock()
	d.knownERC20Addresses[address] = true
}

// RemoveKnownERC20 removes an ERC20 address from the list of known addresses. We will no longer decode
// events for this token.
func (d *Decoder) RemoveKnownERC20(address common.Address) {
	d.knownERC20AddressesMu.Lock()
	defer d.knownERC20AddressesMu.Unlock()
	delete(d.knownERC20Addresses, address)
}

// isKnownERC20 checks if the supplied address is a known ERC20 contract
func (d *Decoder) isKnownERC20(address common.Address) bool {
	d.knownERC20AddressesMu.RLock()
	defer d.knownERC20AddressesMu.RUnlock()
	_, exists := d.knownERC20Addresses[address]
	return exists
}

// AddKnownERC721 registers the supplied contract address as an ERC721 contract. If an event is found
// from this contract address, the decoder will properly decode the `Transfer` and `Approve` events
// including the correct event parameter names.
func (d *Decoder) AddKnownERC721(address common.Address) {
	d.knownERC721AddressesMu.Lock()
	defer d.knownERC721AddressesMu.Unlock()
	d.knownERC721Addresses[address] = true
}

// RemoveKnownERC721 removes an ERC721 address from the list of known addresses. We will no longer decode
// events for this token.
func (d *Decoder) RemoveKnownERC721(address common.Address) {
	d.knownERC721AddressesMu.Lock()
	defer d.knownERC721AddressesMu.Unlock()
	delete(d.knownERC721Addresses, address)
}

// isKnownERC721 checks if the supplied address is a known ERC721 contract
func (d *Decoder) isKnownERC721(address common.Address) bool {
	d.knownERC721AddressesMu.RLock()
	defer d.knownERC721AddressesMu.RUnlock()
	_, exists := d.knownERC721Addresses[address]
	return exists
}

// AddKnownERC1155 registers the supplied contract address as an ERC1155 contract. If an event is found
// from this contract address, the decoder will properly decode the `Transfer` and `Approve` events
// including the correct event parameter names.
func (d *Decoder) AddKnownERC1155(address common.Address) {
	d.knownERC1155AddressesMu.Lock()
	defer d.knownERC1155AddressesMu.Unlock()
	d.knownERC1155Addresses[address] = true
}

// RemoveKnownERC1155 removes an ERC1155 address from the list of known addresses. We will no longer decode
// events for this token.
func (d *Decoder) RemoveKnownERC1155(address common.Address) {
	d.knownERC1155AddressesMu.Lock()
	defer d.knownERC1155AddressesMu.Unlock()
	delete(d.knownERC1155Addresses, address)
}

// isKnownERC1155 checks if the supplied address is a known ERC1155 contract
func (d *Decoder) isKnownERC1155(address common.Address) bool {
	d.knownERC1155AddressesMu.RLock()
	defer d.knownERC1155AddressesMu.RUnlock()
	_, exists := d.knownERC1155Addresses[address]
	return exists
}

// AddKnownExchange registers the supplied contract address as a 0x Exchange contract. If an event is found
// from this contract address, the decoder will properly decode it's events including the correct event
// parameter names.
func (d *Decoder) AddKnownExchange(address common.Address) {
	d.knownExchangeAddressesMu.Lock()
	defer d.knownExchangeAddressesMu.Unlock()
	d.knownExchangeAddresses[address] = true
}

// RemoveKnownExchange removes an Exchange address from the list of known addresses. We will no longer decode
// events for this contract.
func (d *Decoder) RemoveKnownExchange(address common.Address) {
	d.knownExchangeAddressesMu.Lock()
	defer d.knownExchangeAddressesMu.Unlock()
	delete(d.knownExchangeAddresses, address)
}

// isKnownExchange checks if the supplied address is a known Exchange contract address
func (d *Decoder) isKnownExchange(address common.Address) bool {
	d.knownExchangeAddressesMu.RLock()
	defer d.knownExchangeAddressesMu.RUnlock()
	_, exists := d.knownExchangeAddresses[address]
	return exists
}

// FindEventType returns to event type contained in the supplied log. It looks both at the registered
// contract addresses and the log topic.
func (d *Decoder) FindEventType(log types.Log) (string, error) {
	firstTopic := log.Topics[0]
	if isKnown := d.isKnownERC20(log.Address); isKnown {
		eventName, ok := d.erc20TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
		}
		if eventName == "Deposit" || eventName == "Withdrawal" {
			return fmt.Sprintf("Weth%sEvent", eventName), nil
		}
		return fmt.Sprintf("ERC20%sEvent", eventName), nil
	}
	if isKnown := d.isKnownERC721(log.Address); isKnown {
		eventName, ok := d.erc721TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
		}
		return fmt.Sprintf("ERC721%sEvent", eventName), nil
	}
	if isKnown := d.isKnownERC1155(log.Address); isKnown {
		eventName, ok := d.erc1155TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
		}
		return fmt.Sprintf("ERC1155%sEvent", eventName), nil
	}
	if isKnown := d.isKnownExchange(log.Address); isKnown {
		eventName, ok := d.exchangeTopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
		}
		return fmt.Sprintf("Exchange%sEvent", eventName), nil
	}

	return "", UntrackedTokenError{Topic: firstTopic, TokenAddress: log.Address}
}

// Decode attempts to decode the supplied log given the event types relevant to 0x orders. The
// decoded result is stored in the value pointed to by supplied `decodedLog` struct.
func (d *Decoder) Decode(log types.Log, decodedLog interface{}) error {
	if isKnown := d.isKnownERC20(log.Address); isKnown {
		return d.decodeERC20(log, decodedLog)
	}
	if isKnown := d.isKnownERC721(log.Address); isKnown {
		return d.decodeERC721(log, decodedLog)
	}
	if isKnown := d.isKnownERC1155(log.Address); isKnown {
		return d.decodeERC1155(log, decodedLog)
	}
	if isKnown := d.isKnownExchange(log.Address); isKnown {
		return d.decodeExchange(log, decodedLog)
	}

	return UntrackedTokenError{Topic: log.Topics[0], TokenAddress: log.Address}
}

func (d *Decoder) decodeERC20(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.erc20TopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
	}

	err := unpackLog(decodedLog, eventName, log, d.erc20ABI)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decodeERC721(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.erc721TopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
	}

	erc721Err := unpackLog(decodedLog, eventName, log, d.erc721ABI)
	if _, ok := erc721Err.(UnsupportedEventError); ok {
		// Try unpacking using the incorrect ERC721 event ABIs
		fallbackErr := unpackLog(decodedLog, eventName, log, d.erc721EventsAbiWithoutTokenIDIndex)
		if fallbackErr != nil {
			// We return the original attempt's error if the fallback fails
			return erc721Err
		}
	}
	return nil
}

func (d *Decoder) decodeERC1155(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.erc1155TopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
	}

	err := unpackLog(decodedLog, eventName, log, d.erc1155ABI)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decodeExchange(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.exchangeTopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
	}

	err := unpackLog(decodedLog, eventName, log, d.exchangeABI)
	if err != nil {
		return err
	}
	return nil
}

// unpackLog unpacks a retrieved log into the provided output structure.
func unpackLog(decodedEvent interface{}, event string, log types.Log, _abi abi.ABI) error {
	if len(log.Data) > 0 {
		if err := _abi.Unpack(decodedEvent, event, log.Data); err != nil {
			if strings.Contains(err.Error(), "Unpack(no-values unmarshalled") {
				return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
			}
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range _abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	err := abi.ParseTopics(decodedEvent, indexed, log.Topics[1:])
	if err != nil {
		return AbiTopicParserError{Topics: log.Topics, ContractAddress: log.Address, parserError: err}
	}

	return nil
}
