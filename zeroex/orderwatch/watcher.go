package orderwatch

import (
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	blockWatcher *blockwatch.Watcher
	decoder      *Decoder
}

// New instantiates a new order watcher
func New(pollingInterval time.Duration, startBlockDepth rpc.BlockNumber, rpcClient blockwatch.Client) (*Watcher, error) {
	blockRetentionLimit := 20
	withLogs := true
	topics := []common.Hash{}
	for _, signature := range EVENT_SIGNATURES {
		topic := common.BytesToHash(crypto.Keccak256([]byte(signature)))
		topics = append(topics, topic)
	}
	blockWatcher := blockwatch.New(pollingInterval, startBlockDepth, blockRetentionLimit, withLogs, topics, rpcClient)
	decoder, err := NewDecoder()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		blockWatcher: blockWatcher,
		decoder:      decoder,
	}, nil
}

// Start starts the order watcher
func (w *Watcher) Start() {
	w.setupEventWatcher()

	// TODO(fabio): Implement and instantiate expirationwatch

	// TODO(fabio): Implement and instantiate the cleanup worker

	// Everything has been set up. Let's start the block poller.
	w.blockWatcher.StartPolling()
}

// Watch adds a 0x order to the ones being tracked for order-relevant state changes
// func (w *Watcher) Watch(signedOrder SignedOrder) {
// 	// TODO(fabio): Add order's exchange addresses to decoder
// 	// TODO(fabio): Decode assetDatas and add ERC20/ERC721's addresses to decoder
// 	// TODO(fabio): Add expiration & hash to expiration watcher
// }

func (w *Watcher) setupEventWatcher() {
	// TODO(fabio): Add all ERC20 & ERC721 token addresses involved in 0x orders to the decoder
	w.decoder.AddKnownERC20(common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498"))
	w.decoder.AddKnownERC20(common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"))
	w.decoder.AddKnownERC721(common.HexToAddress("0x5d00d312e171be5342067c09bae883f9bcb2003b"))
	w.decoder.AddKnownERC721(common.HexToAddress("0xf5b0a3efb8e8e4c201e2a935f110eaaf3ffecb8d"))
	w.decoder.AddKnownExchange(common.HexToAddress("0x4f833a24e1f95d70f028921e27040ca56e09ab0b"))

	blockEvents := make(chan []*blockwatch.Event, 10)
	sub := w.blockWatcher.Subscribe(blockEvents)
	defer sub.Unsubscribe()

	for events := range blockEvents {
		for _, event := range events {
			decodedLogs, err := w.decodeLogs(event.BlockHeader.Logs)
			if err != nil {
				panic(err) // TODO(fabio): Should we panic here?
			}
			for _, decodedLog := range decodedLogs {
				switch decodedLog.(type) {
				case ERC20TransferEvent:
					transferEvent := decodedLog.(ERC20TransferEvent)
					fmt.Printf("%+v\n", transferEvent)
					// TODO(fabio): Handle this event
				case ERC20ApprovalEvent:
					approvalEvent := decodedLog.(ERC20ApprovalEvent)
					fmt.Printf("%+v\n", approvalEvent)
					// TODO(fabio): Handle this event
				case ERC721TransferEvent:
					transferEvent := decodedLog.(ERC721TransferEvent)
					fmt.Printf("%+v\n", transferEvent)
					// TODO(fabio): Handle this event
				case ERC721ApprovalEvent:
					approvalEvent := decodedLog.(ERC721ApprovalEvent)
					fmt.Printf("%+v\n", approvalEvent)
					// TODO(fabio): Handle this event
				case ERC721ApprovalForAllEvent:
					approvalForAllEvent := decodedLog.(ERC721ApprovalForAllEvent)
					fmt.Printf("%+v\n", approvalForAllEvent)
					// TODO(fabio): Handle this event
				case ExchangeFillEvent:
					fillEvent := decodedLog.(ExchangeFillEvent)
					fmt.Printf("%+v\n", fillEvent)
					// TODO(fabio): Handle this event
				case ExchangeCancelEvent:
					cancelEvent := decodedLog.(ExchangeCancelEvent)
					fmt.Printf("%+v\n", cancelEvent)
					// TODO(fabio): Handle this event
				case ExchangeCancelUpToEvent:
					cancelUpToEvent := decodedLog.(ExchangeCancelUpToEvent)
					fmt.Printf("%+v\n", cancelUpToEvent)
					// TODO(fabio): Handle this event
				case WethDepositEvent:
					depositEvent := decodedLog.(WethDepositEvent)
					fmt.Printf("%+v\n", depositEvent)
					// TODO(fabio): Handle this event
				case WethWithdrawalEvent:
					withdrawalEvent := decodedLog.(WethWithdrawalEvent)
					fmt.Printf("%+v\n", withdrawalEvent)
					// TODO(fabio): Handle this event

				case nil:
					// We were unable to decode the event, ignore.

				default:
					panic(fmt.Sprintf("Unrecognized event returned: %+v", decodedLog))
				}
			}
		}
	}
}

func (w *Watcher) decodeLogs(logs []types.Log) ([]interface{}, error) {
	decodedLogs := []interface{}{}
	for _, log := range logs {
		decodedLog, err := w.decoder.Decode(log)
		// Ignore unsupported events
		if err != nil && err.Error() != "Unsupported event" {
			return nil, err
		}
		decodedLogs = append(decodedLogs, decodedLog)
	}

	return decodedLogs, nil
}
