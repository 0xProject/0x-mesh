package orderwatch

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var EVENT_SIGNATURES = [...]string{
	"Transfer(address,address,uint256)",    // ERC20 & ERC721
	"Approval(address,address,uint256)",    // ERC20 & ERC721
	"ApprovalForAll(address,address,bool)", // ERC721
	"Deposit(address,uint256)",             // WETH9
	"Withdrawal(address,uint256)",          // WETH9
	"Fill(address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes,bytes)", // Exchange
	"Cancel(address,address,address,bytes32,bytes,bytes)",                                       // Exchange
	"CancelUpTo(address,address,uint256)",                                                       // Exchange
}

// Includes ERC20 `Transfer` & `Approval` events as well as WETH `Deposit` & `Withdraw` events
const erc20EventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"

// Includes ERC721 `Transfer`, `Approval` & `ApprovalForAll` events
const erc721EventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// Includes Exchange `Fill`, `Cancel`, `CancelUpTo` events
const exchangeEventsAbi = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"}]"

// ERC20TransferEvent represents an ERC20 Transfer event
type ERC20TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ERC20ApprovalEvent represents an ERC20 Approval event
type ERC20ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

// ERC721TransferEvent represents an ERC721 Transfer event
type ERC721TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

// ERC721ApprovalEvent represents an ERC721 Approval event
type ERC721ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

// ERC721ApprovalForAllEvent represents an ERC721 ApprovalForAll event
type ERC721ApprovalForAllEvent struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
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

// ExchangeCancelEvent represents a 0x Exchange Cancel event
type ExchangeCancelEvent struct {
	MakerAddress        common.Address
	FeeRecipientAddress common.Address
	SenderAddress       common.Address
	OrderHash           common.Hash
	MakerAssetData      []byte
	TakerAssetData      []byte
}

// ExchangeCancelUpToEvent represents a 0x Exchange CancelUpTo event
type ExchangeCancelUpToEvent struct {
	MakerAddress  common.Address
	SenderAddress common.Address
	OrderEpoch    *big.Int
}

// WethWithdrawalEvent represents a wrapped Ether Withdraw event
type WethWithdrawalEvent struct {
	Owner common.Address
	Value *big.Int
}

// WethDepositEvent represents a wrapped Ether Deposit event
type WethDepositEvent struct {
	Owner common.Address
	Value *big.Int
}

// UnsupportedEventError is thrown when an unsupported topic is encountered
type UnsupportedEventError struct {
	Topics []common.Hash
}

// Error returns the error string
func (e UnsupportedEventError) Error() string {
	return fmt.Sprintf("unsupported event: topics: %d", e.Topics)
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
	knownERC20Addresses      map[common.Address]bool
	knownERC721Addresses     map[common.Address]bool
	knownExchangeAddresses   map[common.Address]bool
	erc20ABI                 abi.ABI
	erc721ABI                abi.ABI
	exchangeABI              abi.ABI
	erc20TopicToEventName    map[common.Hash]string
	erc721TopicToEventName   map[common.Hash]string
	exchangeTopicToEventName map[common.Hash]string
}

// NewDecoder instantiates a new 0x order-relevant events decoder
func NewDecoder() (*Decoder, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(erc20EventsAbi))
	if err != nil {
		return nil, err
	}

	erc721ABI, err := abi.JSON(strings.NewReader(erc721EventsAbi))
	if err != nil {
		return nil, err
	}

	exchangeABI, err := abi.JSON(strings.NewReader(exchangeEventsAbi))
	if err != nil {
		return nil, err
	}

	erc20TopicToEventName := map[common.Hash]string{}
	for _, event := range erc20ABI.Events {
		erc20TopicToEventName[event.Id()] = event.Name
	}
	erc721TopicToEventName := map[common.Hash]string{}
	for _, event := range erc721ABI.Events {
		erc721TopicToEventName[event.Id()] = event.Name
	}
	exchangeTopicToEventName := map[common.Hash]string{}
	for _, event := range exchangeABI.Events {
		exchangeTopicToEventName[event.Id()] = event.Name
	}

	return &Decoder{
		knownERC20Addresses:      make(map[common.Address]bool),
		knownERC721Addresses:     make(map[common.Address]bool),
		knownExchangeAddresses:   make(map[common.Address]bool),
		erc20ABI:                 erc20ABI,
		erc721ABI:                erc721ABI,
		exchangeABI:              exchangeABI,
		erc20TopicToEventName:    erc20TopicToEventName,
		erc721TopicToEventName:   erc721TopicToEventName,
		exchangeTopicToEventName: exchangeTopicToEventName,
	}, nil
}

// AddKnownERC20 registers the supplied contract address as an ERC20 contract. If an event is found
// from this contract address, the decoder will properly decode the `Transfer` and `Approve` events
// including the correct event parameter names.
func (d *Decoder) AddKnownERC20(address common.Address) {
	d.knownERC20Addresses[address] = true
}

// AddKnownERC721 registers the supplied contract address as an ERC721 contract. If an event is found
// from this contract address, the decoder will properly decode the `Transfer` and `Approve` events
// including the correct event parameter names.
func (d *Decoder) AddKnownERC721(address common.Address) {
	d.knownERC721Addresses[address] = true
}

// AddKnownExchange registers the supplied contract address as a 0x Exchange contract. If an event is found
// from this contract address, the decoder will properly decode it's events including the correct event
// parameter names.
func (d *Decoder) AddKnownExchange(address common.Address) {
	d.knownExchangeAddresses[address] = true
}

// FindEventType returns to event type contained in the supplied log. It looks both at the registered
// contract addresses and the log topic.
func (d *Decoder) FindEventType(log types.Log) (string, error) {
	firstTopic := log.Topics[0]
	if _, exists := d.knownERC20Addresses[log.Address]; exists {
		eventName, ok := d.erc20TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics}
		}
		if eventName == "Deposit" || eventName == "Withdraw" {
			return fmt.Sprintf("Weth%sEvent", eventName), nil
		}
		return fmt.Sprintf("ERC20%sEvent", eventName), nil
	}
	if _, exists := d.knownERC721Addresses[log.Address]; exists {
		eventName, ok := d.erc721TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics}
		}
		return fmt.Sprintf("ERC721%sEvent", eventName), nil
	}
	if _, exists := d.knownExchangeAddresses[log.Address]; exists {
		eventName, ok := d.erc721TopicToEventName[firstTopic]
		if !ok {
			return "", UnsupportedEventError{Topics: log.Topics}
		}
		return fmt.Sprintf("Exchange%sEvent", eventName), nil
	}

	return "", UntrackedTokenError{Topic: firstTopic, TokenAddress: log.Address}
}

// Decode attempts to decode the supplied log given the event types relevant to 0x orders. The
// decoded result is stored in the value pointed to by supplied `decodedLog` struct.
func (d *Decoder) Decode(log types.Log, decodedLog interface{}) error {
	if _, exists := d.knownERC20Addresses[log.Address]; exists {
		return d.decodeERC20(log, decodedLog)
	}
	if _, exists := d.knownERC721Addresses[log.Address]; exists {
		return d.decodeERC721(log, decodedLog)
	}
	if _, exists := d.knownExchangeAddresses[log.Address]; exists {
		return d.decodeExchange(log, decodedLog)
	}

	return UntrackedTokenError{Topic: log.Topics[0], TokenAddress: log.Address}
}

func (d *Decoder) decodeERC20(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.erc20TopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics}
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
		return UnsupportedEventError{Topics: log.Topics}
	}

	err := unpackLog(decodedLog, eventName, log, d.erc721ABI)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) decodeExchange(log types.Log, decodedLog interface{}) error {
	eventName, ok := d.exchangeTopicToEventName[log.Topics[0]]
	if !ok {
		return UnsupportedEventError{Topics: log.Topics}
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
		return UnsupportedEventError{Topics: log.Topics}
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
