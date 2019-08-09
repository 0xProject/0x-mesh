# CHANGELOG

## Upcoming release

### Breaking changes 🛠 

- Modified Mesh's validation logic to reject and consider invalid any _partially fillable_ orders. While this is
  technically a breaking change, partially fillable orders are rare in the wild and we don't expect this will
  affect many users. (#333)

### Bug fixes 🐞 

- De-dup order submitted via the JSON-RPC method `mesh_addOrders` before performing validation (#331)
- Added `"declaration": true,` to TS client's `tsconfig.json` so that downstream projects can use it's TS typings. (#325)


## v2.0.0-beta

### Breaking changes 🛠 

- Modified how `mesh_addOrders` treats orders that are already stored on the Mesh node. Previously, they would be rejected with code `OrderAlreadyStored`. Now, if the order is stored and fillable, it will be accepted. If it is stored but unfillable, it will be rejected with `OrderAlreadyStoredAndUnfillable`. We additionally added a `isNew` property to the accepted orderInfos returned, so that callers can discern which orders Mesh already knew about. (#316)

### Features ✅ 

- Added backup bootstrap nodes provided by the libp2p community
- Improved log formatting and reduced verbosity in a few cases (#314, #287)
- Reduced AdvertiseBootDelay for bootstrap nodes
- Implemented a check that will alerts you when switching to a different Ethereum network ID. (#301) -- special thanks to @hrharder!
- Made environment variable parsing more generous by automatically removing quotes if needed (#306)
- Improved tests by adding timeouts and closing resources where appropriate (#310, #309, #308)
- Increased robustness by removing panics and failing more gracefully (#312)
- RPC server is now started while block event backfilling is happening under the hood instead of waiting for it to complete (#318)
- Added a `mesh_getStats` endpoint which returns a host of useful information about the state of the Mesh node (e.g., number of fillable order stored, number of peers, peerID, etc...) (#322)

### Bug fixes 🐞 

- Log messages are no longer incorrectly fired when receiving orders which have already been seen (#286)
- Fixed a bug where Mesh was still running after the database was closed (#300)
- Handled Parity "unknown block" error gracefully like we do Geth's (#285)

## v1.0.6-beta

This release fixes several bugs:

- Uninitialized TxHashes map & accidental inclusion of null address txHash in order events (#280)
- Concurrent read/write issue in OrderWatcher's EventDecoder (#278)
- Non-unique logging keys causing Elastic Search indexing issues (#275)

It also includes a reduction in the delay before which bootstrap nodes advertise themselves as relays from 15mins to 30sec.

## v1.0.5-beta

This version introduces a temporary hack in order to help some Mesh nodes find their public IP address under certain circumstances.

## v1.0.4-beta

This release fixes a bug in how AutoNAT and AutoRelay services were discovered.

## v1.0.3-beta

This release fixes some networking related issues with our bootstrap nodes.

- Bootstrap nodes now advertise the correct public IP address.
- Bootstrap nodes now also operate as a relay.

These fixes will help smooth out any issues with peer discovery.

## v1.0.2-beta

This release fixes a few minor bugs and includes additional documentation.

1. Set a custom protocol ID for our DHT in order to separate it from the default IPFS DHT.
2. Fixed a bug in the getOrders JSON-RPC endpoint where fillableTakerAssetAmount was sometimes being encoded as numbers instead of strings.
3. Converted addresses to all lowercase (non-checksummed) in the JSON-RPC API.
4. Improved logging.
5. Added a guide for running a Mesh node with telemetry-enabled.


## v1.0.1-beta

This release adds AutoNAT support for our bootstrap nodes. This enables peers who are behind NATs to find each and connect to each other (in most cases).

In addition, this release includes some changes to default network timeouts and documentation improvements.

## v1.0.0-beta

This is the initial beta release for 0x Mesh!

This release supports the following features:

- Automatic peer discovery via a DHT.
- A JSON-RPC endpoint for interacting with your Mesh node. It includes support for adding new orders and subscribing to order updates.
- Efficient order validation and order watching under the hood. You will get notified quickly when orders are expired, canceled, or filled.
- Basic limitations on the size and types of messages sent by peers. Peers that send malformed messages or messages that are too big will be dropped.

In addition to improving stability and scalability we plan to release many more important features in the near future. Check out https://github.com/0xProject/0x-mesh/issues for more information about what we are working on.
