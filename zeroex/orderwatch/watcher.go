package orderwatch

import (
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/zeroex"
)

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	blockWatcher     *blockwatch.Watcher
	eventDecoder     *Decoder
	assetDataDecoder *zeroex.AssetDataDecoder
}

// New instantiates a new order watcher
func New(blockWatcher *blockwatch.Watcher, rpcClient blockwatch.Client) (*Watcher, error) {
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

// Setup sets up the event & expiration watchers as well as the cleanup worker
func (w *Watcher) Setup() {
	w.setupEventWatcher()

	// TODO(fabio): Implement and instantiate expirationwatch

	// TODO(fabio): Implement and instantiate the cleanup worker
}

// Watch adds a 0x order to the ones being tracked for order-relevant state changes
func (w *Watcher) Watch(signedOrder *zeroex.SignedOrder) error {
	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	err := w.addAssetDataAddressToEventDecoder(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}
	err = w.addAssetDataAddressToEventDecoder(signedOrder.TakerAssetData)
	if err != nil {
		return err
	}

	// TODO(fabio): Add expiration & hash to expiration watcher

	return nil
}

func (w *Watcher) setupEventWatcher() {
	blockEvents := make(chan []*blockwatch.Event, 10)
	sub := w.blockWatcher.Subscribe(blockEvents)

	go func() {
		defer sub.Unsubscribe()
		for events := range blockEvents {
			for _, event := range events {
				for _, log := range event.BlockHeader.Logs {
					eventType, err := w.eventDecoder.FindEventType(log)
					w.handleEventDecodingError(err)
					switch eventType {
					case "ERC20TransferEvent":
						var transferEvent ERC20TransferEvent
						err = w.eventDecoder.Decode(log, &transferEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ERC20ApprovalEvent":
						var approvalEvent ERC20ApprovalEvent
						err = w.eventDecoder.Decode(log, &approvalEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ERC721TransferEvent":
						var transferEvent ERC721TransferEvent
						err = w.eventDecoder.Decode(log, &transferEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ERC721ApprovalEvent":
						var approvalEvent ERC721ApprovalEvent
						err = w.eventDecoder.Decode(log, &approvalEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ERC721ApprovalForAllEvent":
						var approvalForAllEvent ERC721ApprovalForAllEvent
						err = w.eventDecoder.Decode(log, &approvalForAllEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "WethWithdrawalEvent":
						var withdrawalEvent WethWithdrawalEvent
						err = w.eventDecoder.Decode(log, &withdrawalEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "WethDepositEvent":
						var depositEvent WethDepositEvent
						err = w.eventDecoder.Decode(log, &depositEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ExchangeFillEvent":
						var exchangeFillEvent ExchangeFillEvent
						err = w.eventDecoder.Decode(log, &exchangeFillEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ExchangeCancelEvent":
						var exchangeCancelEvent ExchangeCancelEvent
						err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					case "ExchangeCancelUpToEvent":
						var exchangeCancelUpToEvent ExchangeCancelUpToEvent
						err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
						w.handleEventDecodingError(err)
						// TODO(fabio): Handle this event
					default:
						panic(fmt.Sprintf("Did not handle event %s\n", eventType))
					}
				}
			}
		}
	}()
}

func (w *Watcher) handleEventDecodingError(err error) {
	if err == nil {
		return
	}

	switch err := err.(type) {
	case UntrackedTokenError:
		return
	case UnsupportedEventError:
		// TODO(fabio): Write topics to a log for future inspection
		return
	default:
		panic(err) // Should never hit this
	}
}

func (w *Watcher) addAssetDataAddressToEventDecoder(assetData []byte) error {
	assetDataName, err := w.assetDataDecoder.GetName(assetData)
	if err != nil {
		return err
	}
	switch assetDataName {
	case "ERC20Token":
		var decodedAssetData zeroex.ERC20AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC20(decodedAssetData.Address)
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		w.eventDecoder.AddKnownERC721(decodedAssetData.Address)
	case "MultiAsset":
		var decodedAssetData zeroex.MultiAssetData
		err := w.assetDataDecoder.Decode(assetData, &decodedAssetData)
		if err != nil {
			return err
		}
		for _, assetData := range decodedAssetData.NestedAssetData {
			w.addAssetDataAddressToEventDecoder(assetData)
		}
	default:
		return errors.New(fmt.Sprintf("Unrecognized assetData type name found: %s\n", assetDataName))
	}
	return nil
}
