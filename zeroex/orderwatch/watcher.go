package orderwatch

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// Watcher watches all order-relevant state and handles the state transitions
type Watcher struct {
	blockWatcher      *blockwatch.Watcher
	eventDecoder      *Decoder
	assetDataDecoder  *zeroex.AssetDataDecoder
	blockSubscription event.Subscription
	expirationWatcher *ExpirationWatcher
	isSetup           bool
	setupMux          sync.RWMutex
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
	var expirationBuffer int64 = 0
	return &Watcher{
		blockWatcher:      blockWatcher,
		expirationWatcher: NewExpirationWatcher(expirationBuffer),
		eventDecoder:      decoder,
		assetDataDecoder:  assetDataDecoder,
	}, nil
}

// Setup sets up the event & expiration watchers as well as the cleanup worker.
func (w *Watcher) Setup(expirationPollingInterval time.Duration) error {
	w.setupMux.Lock()
	defer w.setupMux.Unlock()
	if w.isSetup {
		return errors.New("Setup can only be called once")
	}

	w.setupEventWatcher()

	w.setupExpirationWatcher(expirationPollingInterval)

	// TODO(fabio): Implement and instantiate the cleanup worker

	w.isSetup = true
	return nil
}

// Watch adds a 0x order to the ones being tracked for order-relevant state changes
func (w *Watcher) Watch(signedOrder *zeroex.SignedOrder, orderHash common.Hash) error {
	w.setupMux.Lock()
	defer w.setupMux.Unlock()
	if !w.isSetup {
		return errors.New("Cannot watch orders before calling Setup()")
	}
	if !w.blockWatcher.IsWatching() {
		return errors.New("Block watcher must be started before adding orders to orderwatch.Watcher")
	}

	w.eventDecoder.AddKnownExchange(signedOrder.ExchangeAddress)

	err := w.addAssetDataAddressToEventDecoder(signedOrder.MakerAssetData)
	if err != nil {
		return err
	}
	err = w.addAssetDataAddressToEventDecoder(signedOrder.TakerAssetData)
	if err != nil {
		return err
	}

	w.expirationWatcher.Add(signedOrder.ExpirationTimeSeconds.Int64(), orderHash)

	return nil
}

func (w *Watcher) setupExpirationWatcher(expirationPollingInterval time.Duration) {
	go func() {
		for expiredOrders := range w.expirationWatcher.ExpiredOrders {
			for _, expiredOrder := range expiredOrders {
				// TODO(fabio): Handle expired order
				panic(fmt.Sprintf("Handling expired orders is not implemented yet: %+v\n", expiredOrder))
			}
		}
	}()

	w.expirationWatcher.Start(expirationPollingInterval)
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
					if err != nil {
						switch err.(type) {
						case UntrackedTokenError:
							continue
						case UnsupportedEventError:
							// TODO(fabio): Log the event topics
							continue
						default:
							panic(err)
						}
					}
					switch eventType {
					case "ERC20TransferEvent":
						var transferEvent ERC20TransferEvent
						err = w.eventDecoder.Decode(log, &transferEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ERC20ApprovalEvent":
						var approvalEvent ERC20ApprovalEvent
						err = w.eventDecoder.Decode(log, &approvalEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ERC721TransferEvent":
						var transferEvent ERC721TransferEvent
						err = w.eventDecoder.Decode(log, &transferEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ERC721ApprovalEvent":
						var approvalEvent ERC721ApprovalEvent
						err = w.eventDecoder.Decode(log, &approvalEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ERC721ApprovalForAllEvent":
						var approvalForAllEvent ERC721ApprovalForAllEvent
						err = w.eventDecoder.Decode(log, &approvalForAllEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "WethWithdrawalEvent":
						var withdrawalEvent WethWithdrawalEvent
						err = w.eventDecoder.Decode(log, &withdrawalEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "WethDepositEvent":
						var depositEvent WethDepositEvent
						err = w.eventDecoder.Decode(log, &depositEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ExchangeFillEvent":
						var exchangeFillEvent ExchangeFillEvent
						err = w.eventDecoder.Decode(log, &exchangeFillEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ExchangeCancelEvent":
						var exchangeCancelEvent ExchangeCancelEvent
						err = w.eventDecoder.Decode(log, &exchangeCancelEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					case "ExchangeCancelUpToEvent":
						var exchangeCancelUpToEvent ExchangeCancelUpToEvent
						err = w.eventDecoder.Decode(log, &exchangeCancelUpToEvent)
						if err != nil {
							switch err.(type) {
							case UnsupportedEventError:
								// TODO(fabio): Log the event topics
								continue
							default:
								panic(err)
							}
						}
						// TODO(fabio): Handle this event
					default:
						panic(fmt.Sprintf("Did not handle event %s\n", eventType))
					}
				}
			}
		}
	}()
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
