package orderwatch

import (
	"fmt"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	blockWatcher     *blockwatch.Watcher
	eventDecoder     *Decoder
	assetDataDecoder *zeroex.AssetDataDecoder
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
	assetDataDecoder, err := zeroex.NewAssetDataDecoder()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		blockWatcher:     blockWatcher,
		eventDecoder:     decoder,
		assetDataDecoder: assetDataDecoder,
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
func (w *Watcher) Watch(signedOrder *zeroex.SignedOrder) error {
	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	decodedMakerAssetData, err := w.assetDataDecoder.Decode(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}
	w.addAddressFromAssetDataToEventDecoder(decodedMakerAssetData)
	decodedTakerAssetData, err := w.assetDataDecoder.Decode(signedOrder.TakerAssetData)
	if err != nil {
		return err
	}
	w.addAddressFromAssetDataToEventDecoder(decodedTakerAssetData)

	// TODO(fabio): Add expiration & hash to expiration watcher

	return nil
}

func (w *Watcher) setupEventWatcher() {
	blockEvents := make(chan []*blockwatch.Event, 10)
	sub := w.blockWatcher.Subscribe(blockEvents)
	defer sub.Unsubscribe()

	for events := range blockEvents {
		for _, event := range events {
			for _, log := range event.BlockHeader.Logs {
				eventType, err := w.eventDecoder.findEventType(log)
				if err != nil {
					if err.Error() == unsupportedEvent {
						continue
					}
					// The decoder is very lenient, so if another error is returned,
					// it must be for an unrecoverable error and we should panic
					panic(err)
				}
				switch eventType {
				case "ERC20TransferEvent":
					var transferEvent ERC20TransferEvent
					err = w.eventDecoder.Decode(log, &transferEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ERC20ApprovalEvent":
					var approvalEvent ERC20ApprovalEvent
					err = w.eventDecoder.Decode(log, &approvalEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ERC721TransferEvent":
					var transferEvent ERC721TransferEvent
					err = w.eventDecoder.Decode(log, &transferEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ERC721ApprovalEvent":
					var approvalEvent ERC721ApprovalEvent
					err = w.eventDecoder.Decode(log, &approvalEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ERC721ApprovalForAllEvent":
					var approvalForAllEvent ERC721ApprovalForAllEvent
					err = w.eventDecoder.Decode(log, &approvalForAllEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "WethWithdrawalEvent":
					var withdrawalEvent WethWithdrawalEvent
					err = w.eventDecoder.Decode(log, &withdrawalEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "WethDepositEvent":
					var depositEvent WethDepositEvent
					err = w.eventDecoder.Decode(log, &depositEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ExchangeFillEvent":
					var exchangeFillEvent ExchangeFillEvent
					err = w.eventDecoder.Decode(log, &exchangeFillEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ExchangeCancelEvent":
					var exchangeCancelEvent ExchangeCancelEvent
					err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				case "ExchangeCancelUpToEvent":
					var exchangeCancelUpToEvent ExchangeCancelUpToEvent
					err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
					if err != nil {
						panic(err)
					}
					// TODO(fabio): Handle this event
				default:
					panic(fmt.Sprintf("Did not handle event %s\n", eventType))
				}
			}
		}
	}
}

func (w *Watcher) addAddressFromAssetDataToEventDecoder(decodedAssetData interface{}) error {
	switch decodedAssetData.(type) {
	case zeroex.ERC20AssetData:
		w.eventDecoder.AddKnownERC20(decodedAssetData.(zeroex.ERC20AssetData).Address)
	case zeroex.ERC721AssetData:
		w.eventDecoder.AddKnownERC721(decodedAssetData.(zeroex.ERC721AssetData).Address)
	case zeroex.MultiAssetData:
		multiAssetData := decodedAssetData.(zeroex.MultiAssetData)
		// Recursively add the nested assetData to the event decoder
		for _, assetData := range multiAssetData.NestedAssetData {
			d, err := w.assetDataDecoder.Decode(assetData)
			if err != nil {
				return err
			}
			w.addAddressFromAssetDataToEventDecoder(d)
		}
	}
	return nil
}
