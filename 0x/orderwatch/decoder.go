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

// Also includes WETH `Deposit` & `Withdraw` events
const ERC20_EVENTS_ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]"

// TODO(fabio): Some ERC721 tokens don't comply 100% for the specification. E.g., Axie Infinity
// doesn't have an `index` on `tokenId` in their `Transfer` event. This kind of stuff will cause an
// error in the decoding. We should make this event decoder more robust such that is can handle
// differences in the number of `index`'s without crashing
const ERC721_EVENTS_ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

// Rest contains: Exchange `Fill`, `Cancel`, `CancelUpTo` and WETH `Deposit` and `Withdraw`
const EXCHANGE_EVENTS_ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"}]"

const UNSUPPORTED_EVENT = "Unsupported event"

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

// Decoder decodes events relevant to the fillability of 0x orders. Since ERC20 & ERC721 events
// have the same signatures, but different meanings, all ERC20 & ERC721 contract addresses must
// be added to the decoder ahead of time.
type Decoder struct {
	knownERC20Addresses    map[common.Address]bool
	knownERC721Addresses   map[common.Address]bool
	knownExchangeAddresses map[common.Address]bool
	erc20ABI               abi.ABI
	erc721ABI              abi.ABI
	exchangeABI            abi.ABI
	topicToNumTopics       map[common.Hash]int
}

// NewDecoder instantiates a new decoder
func NewDecoder() (*Decoder, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(ERC20_EVENTS_ABI))
	if err != nil {
		return nil, err
	}

	erc721ABI, err := abi.JSON(strings.NewReader(ERC721_EVENTS_ABI))
	if err != nil {
		return nil, err
	}

	exchangeABI, err := abi.JSON(strings.NewReader(EXCHANGE_EVENTS_ABI))
	if err != nil {
		return nil, err
	}

	return &Decoder{
		knownERC20Addresses:    make(map[common.Address]bool),
		knownERC721Addresses:   make(map[common.Address]bool),
		knownExchangeAddresses: make(map[common.Address]bool),
		erc20ABI:               erc20ABI,
		erc721ABI:              erc721ABI,
		exchangeABI:            exchangeABI,
		topicToNumTopics:       make(map[common.Hash]int),
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

// Decode attempts to decode the supplied log given the event types relevant to 0x orders
func (d *Decoder) Decode(log types.Log) (interface{}, error) {
	if _, ok := d.knownERC20Addresses[log.Address]; ok {
		return d.decodeERC20(log)
	}
	if _, ok := d.knownERC721Addresses[log.Address]; ok {
		return d.decodeERC721(log)
	}
	if _, ok := d.knownExchangeAddresses[log.Address]; ok {
		return d.decodeExchange(log)
	}

	return nil, errors.New(UNSUPPORTED_EVENT)
}

func (d *Decoder) decodeERC20(log types.Log) (interface{}, error) {
	topicToEventName := map[common.Hash]string{}
	for _, event := range d.erc20ABI.Events {
		topicToEventName[event.Id()] = event.Name
	}
	eventName, ok := topicToEventName[log.Topics[0]]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Could not find ERC20 event name with topic: %s", log.Topics[0].Hex()))
	}

	switch eventName {
	case "Transfer":
		var transferEvent ERC20TransferEvent
		err := unpackLog(&transferEvent, eventName, log, d.erc20ABI)
		if err != nil {
			return nil, err
		}
		return transferEvent, nil

	case "Approval":
		var approvalEvent ERC20ApprovalEvent
		err := unpackLog(&approvalEvent, eventName, log, d.erc20ABI)
		if err != nil {
			return nil, err
		}
		return approvalEvent, nil

	// WETH is an ERC20 token with `Withdraw` & `Deposit` events
	case "Withdrawal":
		var withdrawalEvent WethWithdrawalEvent
		err := unpackLog(&withdrawalEvent, eventName, log, d.erc20ABI)
		if err != nil {
			return nil, err
		}
		return withdrawalEvent, nil

	case "Deposit":
		var depositEvent WethDepositEvent
		err := unpackLog(&depositEvent, eventName, log, d.erc20ABI)
		if err != nil {
			return nil, err
		}
		return depositEvent, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized ERC20 Event: %s", eventName))
	}
}

func (d *Decoder) decodeERC721(log types.Log) (interface{}, error) {
	topicToEventName := map[common.Hash]string{}
	for _, event := range d.erc721ABI.Events {
		topicToEventName[event.Id()] = event.Name
	}
	eventName, ok := topicToEventName[log.Topics[0]]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Could not find ERC721 event name with topic: %s", log.Topics[0].Hex()))
	}

	switch eventName {
	case "Transfer":
		var transferEvent ERC721TransferEvent
		err := unpackLog(&transferEvent, eventName, log, d.erc721ABI)
		if err != nil {
			return nil, err
		}
		return transferEvent, nil

	case "Approval":
		var approvalEvent ERC721ApprovalEvent
		err := unpackLog(&approvalEvent, eventName, log, d.erc721ABI)
		if err != nil {
			return nil, err
		}
		return approvalEvent, nil

	case "ApprovalForAll":
		var approvalForAllEvent ERC721ApprovalForAllEvent
		err := unpackLog(&approvalForAllEvent, eventName, log, d.erc721ABI)
		if err != nil {
			return nil, err
		}
		return approvalForAllEvent, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized ERC721 Event: %s", eventName))
	}
}

func (d *Decoder) decodeExchange(log types.Log) (interface{}, error) {
	topicToEventName := map[common.Hash]string{}
	for _, event := range d.exchangeABI.Events {
		topicToEventName[event.Id()] = event.Name
	}
	eventName, ok := topicToEventName[log.Topics[0]]
	if !ok {
		return nil, errors.New(UNSUPPORTED_EVENT)
	}

	switch eventName {
	case "Fill":
		var fillEvent ExchangeFillEvent
		err := unpackLog(&fillEvent, eventName, log, d.exchangeABI)
		if err != nil {
			return nil, err
		}
		return fillEvent, nil

	case "Cancel":
		var cancelEvent ExchangeCancelEvent
		err := unpackLog(&cancelEvent, eventName, log, d.exchangeABI)
		if err != nil {
			return nil, err
		}
		return cancelEvent, nil

	case "CancelUpTo":
		var cancelUpToEvent ExchangeCancelUpToEvent
		err := unpackLog(&cancelUpToEvent, eventName, log, d.exchangeABI)
		if err != nil {
			return nil, err
		}
		return cancelUpToEvent, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unsupported Log Event: %s", eventName))
	}
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
		return errors.New(UNSUPPORTED_EVENT)
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

// capitalise makes a camel-case string which starts with an upper case character.
func capitalise(input string) string {
	return abi.ToCamelCase(input)
}

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
		field := reflect.ValueOf(out).Elem().FieldByName(capitalise(arg.Name))

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
