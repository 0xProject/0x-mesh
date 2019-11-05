package decoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var EVENT_SIGNATURES = [...]string{
	"Transfer(address,address,uint256)",                                                         // ERC20 & ERC721
	"Approval(address,address,uint256)",                                                         // ERC20 & ERC721
	"TransferSingle(address,address,address,uint256,uint256)",                                   // ERC1155
	"TransferBatch(address,address,address,uint256[],uint256[])",                                // ERC1155
	"ApprovalForAll(address,address,bool)",                                                      // ERC721 & ERC1155
	"Deposit(address,uint256)",                                                                  // WETH9
	"Withdrawal(address,uint256)",                                                               // WETH9
	"Fill(address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes,bytes)", // Exchange
	"Cancel(address,address,address,bytes32,bytes,bytes)",                                       // Exchange
	"CancelUpTo(address,address,uint256)",                                                       // Exchange
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
const exchangeEventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"}]"

// ERC20TransferEvent represents an ERC20 Transfer event
type ERC20TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC20TransferEvent type
func (e ERC20TransferEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"from":  e.From.Hex(),
		"to":    e.To.Hex(),
		"value": e.Value.String(),
	})
}

// ERC20ApprovalEvent represents an ERC20 Approval event
type ERC20ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC20ApprovalEvent type
func (e ERC20ApprovalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner":   e.Owner.Hex(),
		"spender": e.Spender.Hex(),
		"value":   e.Value.String(),
	})
}

// ERC721TransferEvent represents an ERC721 Transfer event
type ERC721TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC721TransferEvent type
func (e ERC721TransferEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"from":    e.From.Hex(),
		"to":      e.To.Hex(),
		"tokenId": e.TokenId.String(),
	})
}

// ERC721ApprovalEvent represents an ERC721 Approval event
type ERC721ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC721ApprovalEvent type
func (e ERC721ApprovalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"approved": e.Approved.Hex(),
		"tokenId":  e.TokenId.String(),
	})
}

// ERC721ApprovalForAllEvent represents an ERC721 ApprovalForAll event
type ERC721ApprovalForAllEvent struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

// MarshalJSON implements a custom JSON marshaller for the ERC721ApprovalForAllEvent type
func (e ERC721ApprovalForAllEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"operator": e.Operator.Hex(),
		"approved": e.Approved,
	})
}

// ERC1155ApprovalForAllEvent represents an ERC1155 ApprovalForAll event
type ERC1155ApprovalForAllEvent struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

// MarshalJSON implements a custom JSON marshaller for the ERC1155ApprovalForAllEvent type
func (e ERC1155ApprovalForAllEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner":    e.Owner.Hex(),
		"operator": e.Operator.Hex(),
		"approved": e.Approved,
	})
}

// ERC1155TransferSingleEvent represents an ERC1155 TransfeSingler event
type ERC1155TransferSingleEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC1155TransferSingleEvent type
func (e ERC1155TransferSingleEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"operator": e.Operator.Hex(),
		"from":     e.From.Hex(),
		"to":       e.To.Hex(),
		"id":       e.Id.String(),
		"value":    e.Value.String(),
	})
}

// ERC1155TransferBatchEvent represents an ERC1155 TransfeSingler event
type ERC1155TransferBatchEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ERC1155TransferBatchEvent type
func (e ERC1155TransferBatchEvent) MarshalJSON() ([]byte, error) {
	ids := []string{}
	for _, id := range e.Ids {
		ids = append(ids, id.String())
	}
	values := []string{}
	for _, value := range e.Values {
		values = append(values, value.String())
	}
	return json.Marshal(map[string]interface{}{
		"operator": e.Operator.Hex(),
		"from":     e.From.Hex(),
		"to":       e.To.Hex(),
		"ids":      ids,
		"values":   values,
	})
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
	OrderHash              common.Hash
	MakerAssetData         []byte
	TakerAssetData         []byte
}

// MarshalJSON implements a custom JSON marshaller for the ExchangeFillEvent type
func (e ExchangeFillEvent) MarshalJSON() ([]byte, error) {
	makerAssetData := ""
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := ""
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	return json.Marshal(map[string]interface{}{
		"makerAddress":           e.MakerAddress.Hex(),
		"takerAddress":           e.TakerAddress.Hex(),
		"senderAddress":          e.SenderAddress.Hex(),
		"feeRecipientAddress":    e.FeeRecipientAddress.Hex(),
		"makerAssetFilledAmount": e.MakerAssetFilledAmount.String(),
		"takerAssetFilledAmount": e.TakerAssetFilledAmount.String(),
		"makerFeePaid":           e.MakerFeePaid.String(),
		"takerFeePaid":           e.TakerFeePaid.String(),
		"orderHash":              e.OrderHash.Hex(),
		"makerAssetData":         makerAssetData,
		"takerAssetData":         takerAssetData,
	})
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

// MarshalJSON implements a custom JSON marshaller for the ExchangeCancelEvent type
func (e ExchangeCancelEvent) MarshalJSON() ([]byte, error) {
	makerAssetData := ""
	if len(e.MakerAssetData) != 0 {
		makerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.MakerAssetData))
	}
	takerAssetData := ""
	if len(e.TakerAssetData) != 0 {
		takerAssetData = fmt.Sprintf("0x%s", common.Bytes2Hex(e.TakerAssetData))
	}
	return json.Marshal(map[string]interface{}{
		"makerAddress":        e.MakerAddress.Hex(),
		"senderAddress":       e.SenderAddress.Hex(),
		"feeRecipientAddress": e.FeeRecipientAddress.Hex(),
		"orderHash":           e.OrderHash.Hex(),
		"makerAssetData":      makerAssetData,
		"takerAssetData":      takerAssetData,
	})
}

// ExchangeCancelUpToEvent represents a 0x Exchange CancelUpTo event
type ExchangeCancelUpToEvent struct {
	MakerAddress  common.Address
	SenderAddress common.Address
	OrderEpoch    *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the ExchangeCancelUpToEvent type
func (e ExchangeCancelUpToEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"makerAddress":  e.MakerAddress.Hex(),
		"senderAddress": e.SenderAddress.Hex(),
		"orderEpoch":    e.OrderEpoch.String(),
	})
}

// WethWithdrawalEvent represents a wrapped Ether Withdraw event
type WethWithdrawalEvent struct {
	Owner common.Address
	Value *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the WethWithdrawalEvent type
func (w WethWithdrawalEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner": w.Owner.Hex(),
		"value": w.Value.String(),
	})
}

// WethDepositEvent represents a wrapped Ether Deposit event
type WethDepositEvent struct {
	Owner common.Address
	Value *big.Int
}

// MarshalJSON implements a custom JSON marshaller for the WethDepositEvent type
func (w WethDepositEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"owner": w.Owner.Hex(),
		"value": w.Value.String(),
	})
}

// UnsupportedEventError is thrown when an unsupported topic is encountered
type UnsupportedEventError struct {
	Topics          []common.Hash
	ContractAddress common.Address
}

// Error returns the error string
func (e UnsupportedEventError) Error() string {
	return fmt.Sprintf("unsupported event: contract address: %s, topics: %v", e.ContractAddress, e.Topics)
}

type UntrackedTokenError struct {
	Topic        common.Hash
	TokenAddress common.Address
}

// Error returns the error string
func (e UntrackedTokenError) Error() string {
	return fmt.Sprintf("event for an untracked token: contract address: %s, topic: %s", e.TokenAddress, e.Topic)
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
		erc20TopicToEventName[event.ID()] = event.Name
	}
	erc721TopicToEventName := map[common.Hash]string{}
	for _, event := range erc721ABI.Events {
		erc721TopicToEventName[event.ID()] = event.Name
	}
	erc1155TopicToEventName := map[common.Hash]string{}
	for _, event := range erc1155ABI.Events {
		erc1155TopicToEventName[event.ID()] = event.Name
	}
	exchangeTopicToEventName := map[common.Hash]string{}
	for _, event := range exchangeABI.Events {
		exchangeTopicToEventName[event.ID()] = event.Name
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
	if erc721Err != nil {
		if _, ok := erc721Err.(UnsupportedEventError); ok {
			// Try unpacking using the incorrect ERC721 event ABIs
			fallbackErr := unpackLog(decodedLog, eventName, log, d.erc721EventsAbiWithoutTokenIDIndex)
			if fallbackErr != nil {
				// We return the original attempt's error if the fallback fails
				return erc721Err
			}
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
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range _abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if len(indexed) != len(log.Topics[1:]) {
		return UnsupportedEventError{Topics: log.Topics, ContractAddress: log.Address}
	}
	return parseTopics(decodedEvent, indexed, log.Topics[1:])
}

/**
* HACK(fabio): The code below was pulled in from `go-ethereum/accounts/abi/bind` since it was
* unfortunately not exported.
**/

// Big batch of reflect types for topic reconstruction.
var (
	reflectHash    = reflect.TypeOf(common.Hash{})
	reflectAddress = reflect.TypeOf(common.Address{})
	reflectBigInt  = reflect.TypeOf(new(big.Int))
)

// parseTopics converts the indexed topic fields into actual log field values.
//
// Note, dynamic types cannot be reconstructed since they get mapped to Keccak256
// hashes as the topic value!
func parseTopics(out interface{}, fields abi.Arguments, topics []common.Hash) error {
	// Sanity check that the fields and topics match up
	if len(fields) != len(topics) {
		return errors.New("topic/field count mismatch")
	}
	// Iterate over all the fields and reconstruct them from topics
	for _, arg := range fields {
		if !arg.Indexed {
			return errors.New("non-indexed field in topic reconstruction")
		}
		field := reflect.ValueOf(out).Elem().FieldByName(abi.ToCamelCase(arg.Name))

		// Try to parse the topic back into the fields based on primitive types
		switch field.Kind() {
		case reflect.Bool:
			if topics[0][common.HashLength-1] == 1 {
				field.Set(reflect.ValueOf(true))
			}
		case reflect.Int8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int8(num.Int64())))

		case reflect.Int16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int16(num.Int64())))

		case reflect.Int32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int32(num.Int64())))

		case reflect.Int64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Int64()))

		case reflect.Uint8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint8(num.Uint64())))

		case reflect.Uint16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint16(num.Uint64())))

		case reflect.Uint32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint32(num.Uint64())))

		case reflect.Uint64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Uint64()))

		default:
			// Ran out of plain primitive types, try custom types
			switch field.Type() {
			case reflectHash: // Also covers all dynamic types
				field.Set(reflect.ValueOf(topics[0]))

			case reflectAddress:
				var addr common.Address
				copy(addr[:], topics[0][common.HashLength-common.AddressLength:])
				field.Set(reflect.ValueOf(addr))

			case reflectBigInt:
				num := new(big.Int).SetBytes(topics[0][:])
				field.Set(reflect.ValueOf(num))

			default:
				// Ran out of custom types, try the crazies
				switch {

				// static byte array
				case arg.Type.T == abi.FixedBytesTy:
					reflect.Copy(field, reflect.ValueOf(topics[0][:arg.Type.Size]))

				default:
					return fmt.Errorf("unsupported indexed type: %v", arg.Type)
				}
			}
		}
		topics = topics[1:]
	}
	return nil
}
