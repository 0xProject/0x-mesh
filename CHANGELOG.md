# CHANGELOG

This changelog is a work in progress and may contain notes for versions which have not actually been released. Check the [Releases](https://github.com/0xProject/0x-mesh/releases) page to see full release notes and more information about the latest released versions.

## v6.1.0-beta

### Features ✅

- Added support for persistence in the browser. Users of the `@0x/mesh-browser` package will now be able to retain orders and other relevant parts of the state when refreshing the page or closing and re-opening the browser. ([#533](https://github.com/0xProject/0x-mesh/pull/533)).

### Bug fixes 🐞

- Fix bug where Mesh nodes were logging receipt and re-sharing with peers duplicate orders already stored in it's DB, if the duplicate order was submitted via JSON-RPC. ([#529](https://github.com/0xProject/0x-mesh/pull/529))
- Add missing `UNEXPIRED` `OrderEventEndState` enum value to both `@0x/mesh-rpc-client` and `@0x/mesh-browser` and missing `STOPPED_WATCHING` value from `@0x/mesh-rpc-client`.
- Fixed a potential memory leak by using the latest version of `github.com/libp2p/go-libp2p-kad-dht` ([#539](https://github.com/0xProject/0x-mesh/pull/539)).
- Changed the default port for `RPC_ADDR` from a random available port to `60557`. _Some_ documentation already assumed `60557` was the default port. Now all documentation has been updated for consistency with this change. ([#542](https://github.com/0xProject/0x-mesh/pull/542)). 
- Fixed a potential nil pointer exception in log hooks ([#543](https://github.com/0xProject/0x-mesh/pull/543)).
- Fixed a bug where successful closes of an rpc subscription were being reported as errors ([#544](https://github.com/0xProject/0x-mesh/pull/544)).
- We now log the error and stack trace if an RPC method panics. Before, these errors were swallowed by the panic recovery logic in `go-ethereum`'s `rpc` package. ([#545](https://github.com/0xProject/0x-mesh/pull/545))

## v6.0.1-beta

### Bug fixes 🐞 

- Fixed an oversight which granted immunity from bandwidth banning for any peer using a relayed connection ([#509](https://github.com/0xProject/0x-mesh/pull/509)).
- Fixed a typo in the `@0x/mesh-browser` package that resulted in some config options not being passed through correctly ([#502](https://github.com/0xProject/0x-mesh/pull/502)).
- Fixed a bug in ETH JSON-RPC rate limiter where not all dates were being properly converted to UTC, causing Mesh to malfunction if the local time was a day earlier or later than UTC. ([#505](https://github.com/0xProject/0x-mesh/pull/505))
- Fixed a bug in the TypeScript RPC client that prevented orders from being added ([#514](https://github.com/0xProject/0x-mesh/pull/514)).


## v6.0.0-beta

### Breaking changes 🛠 

*Note:* This release will require wiping your Mesh's DB before upgrading. The DB location defaults to `./0x_mesh/db`.

- Renamed env config from `ETHEREUM_NETWORK_ID` to `ETHEREUM_CHAIN_ID` since `network` is a misnomer here and what we actually care about is the `chainID`. Most chains have the same id for their p2p network and chain. From the ones we support, the only outlier is Ganache, for which you will now need to supply `1337` instead of `50` (Learn more: https://medium.com/@pedrouid/chainid-vs-networkid-how-do-they-differ-on-ethereum-eec2ed41635b) ([#485](https://github.com/0xProject/0x-mesh/pull/485))
- Rejected order code `OrderForIncorrectNetwork` has been changed to `OrderForIncorrectChain` ([#485](https://github.com/0xProject/0x-mesh/pull/485))
- Removed `RPC_PORT` environment variable. The new `RPC_ADDR` environment variable allows specifying both the interface and port ([#487](https://github.com/0xProject/0x-mesh/pull/487)).
- Due to security concerns and new rate limiting mechanisms, the default bind address for the RPC API has changed from `0.0.0.0` to `localhost`. Users who previously did not set `RPC_PORT` may need to now manually set `RPC_ADDR` to enable other applications to access the RPC API. If you are using Docker Compose, we recommend using [links](https://docs.docker.com/compose/networking/#links). If you do need to set `RPC_ADDR` to bind on `0.0.0.0`, please be aware of the security implications and consider protecting access to Mesh via a third-party firewall. (See [#444](https://github.com/0xProject/0x-mesh/pull/444) and [#487](https://github.com/0xProject/0x-mesh/pull/487) for more details).
- Changed the `EXPIRED` event such that it is emitted when an order is expired according to the latest block timestamp, not anymore based on UTC time. ([#490](https://github.com/0xProject/0x-mesh/pull/490))
- Removed the `EXPIRATION_BUFFER_SECONDS` env config since we no longer compute order expiration using UTC time. ([#490](https://github.com/0xProject/0x-mesh/pull/490))

### Features ✅ 

- Implemented a new strategy for limiting the amount of database storage used by Mesh and removing orders when the database is full. This strategy involves a dynamically adjusting maximum expiration time. When the database is full, Mesh will enforce a maximum expiration time for all incoming orders and remove any existing orders with an expiration time too far in the future. If conditions change and there is enough space in the database again, the max expiration time will slowly increase. This is a short term solution which solves the immediate issue of finite storage capacities and does a decent job of protecting against spam. We expect to improve and possibly replace it in the future. See [#450](https://github.com/0xProject/0x-mesh/pull/450) for more details.
- Added support for a new feature called "order pinning" ([#474](https://github.com/0xProject/0x-mesh/pull/474)). Pinned orders will not be affected by any DDoS prevention or incentive mechanisms (including the new dynamic max expiration time feature) and will always stay in storage until they are no longer fillable. By default, all orders which are submitted via either the JSON-RPC API or the `addOrdersAsync` function in the TypeScript bindings will be pinned.
- Re-enabled bandwidth-based peer banning with a workaround to deal with erroneous spikes [#478](https://github.com/0xProject/0x-mesh/pull/478).
- Added an `UNEXPIRED` order event kind which is emitted for orders that were previously considered expired but due to a block-reorg causing the latest block timestamp to be earlier than the previous latest block timestamp, are no longer expired. ([#490](https://github.com/0xProject/0x-mesh/pull/490))
- Added support for decoding Axie Infinity `Transfer` and `Approve` ERC721 events which differ from the ERC721 standard. ([#494](https://github.com/0xProject/0x-mesh/pull/494))

### Bug fixes 🐞 

- Improved the aggressiveness at which we permanently delete orders that have been flagged for removal. Previously we would wait for the cleanup job to handle this (once an hour), but that meant many removed orders would accumulate. We now prune them every 5 minutes. ([#471](https://github.com/0xProject/0x-mesh/pull/471))
- Fixed a bug in the Go RPC client which resulted in errors when receving order events with at least one contract event ([#496](https://github.com/0xProject/0x-mesh/pull/496)).

## v5.1.0-beta

### Features ✅ 

- The `getStats` RPC endpoint now includes a new field which accounts for the number of orders that have been marked as "removed" but not yet permanently deleted ([#461](https://github.com/0xProject/0x-mesh/pull/461)).
- Improved historical order sharing using round-robin algorithm instead of random selection ([#454](https://github.com/0xProject/0x-mesh/pull/454)). This will reduce the warm-up time for receiving existing orders when first joining the network.
- Added ERC1155 assetData support ([#453](https://github.com/0xProject/0x-mesh/pull/453)). This includes order watching and order events for orders involving ERC1155 tokens.
- Added Ability to specify custom contract addresses via the `CUSTOM_ADDRESSES` environment variable or the `customAddresses` field in the TypeScript bindings ([#445](https://github.com/0xProject/0x-mesh/pull/445)).

### Bug fixes 🐞 

- Temporarily disabled bandwidth-based peer banning ([#448](https://github.com/0xProject/0x-mesh/pull/448)). A [bug in libp2p](https://github.com/libp2p/go-libp2p-core/issues/65) was occasionally causing observed bandwidth usage to spike to unrealistic levels, which can result in peers being erroneously banned. We decided to temporarily stop banning peers while we're working with the libp2p team to resolve the issue.

## v5.0.0-beta

### Breaking changes 🛠 

- Removes the `txHashes` key in the `OrderEvent`s emitted from the `orders` JSON-RPC subscription and replaced it with `contractEvents`, an array of decoded order-relevant contract events. Parsing these events allows callers to find every discrete order fill/cancel event. ([#420](https://github.com/0xProject/0x-mesh/pull/420))
- Renames the `Kind` key in `OrderEvent` to `EndState` to better elucidate that it represents the aggregate change to the orders state since it was last re-validated. As an end state, it does not capture any possible intermediate states the order might have been in since the last re-validation. Intermediate states can be inferred from the `contractEvents` included ([#420](https://github.com/0xProject/0x-mesh/pull/420))

### Features ✅ 

- Removed the max expiration limit for orders. The only remaining expiration constraint is that the unix timestamp does not overflow int64 (i.e., is not larger than 9223372036854775807). ([#400](https://github.com/0xProject/0x-mesh/pull/400))

### Bug fixes 🐞 

- Fixed bug where we weren't updating an orders `fillableTakerAssetAmount` in the DB when orders were being partially filled or when their fillability increased due to a block re-org. ([#439](https://github.com/0xProject/0x-mesh/pull/439))
- Made `verbosity` field optional in the TypeScript `Config` type. ([#410](https://github.com/0xProject/0x-mesh/pull/410))
- Fixed issue where we weren't re-validating orders potentially impacted by the balance increase of the recipient of an ERC20 or ERC721 transfer. ([#416](https://github.com/0xProject/0x-mesh/pull/416))

## v4.0.1-beta

### Bug fixes 🐞 

- Fixed a DB transaction deadlock accidentally introduced in the v4.0.0-beta release. ([#403](https://github.com/0xProject/0x-mesh/pull/403))

## v4.0.0-beta

### Breaking changes 🛠 

- Renamed the environment variable `P2P_LISTEN_PORT` to `P2P_TCP_PORT` ([#366](https://github.com/0xProject/0x-mesh/pull/366)). This makes it possible to configure Mesh to use both the TCP and Websocket transports by listening on different ports.

### Features ✅ 

- Enabled WebSocket transport for bootstrap nodes and all other nodes ([#361](https://github.com/0xProject/0x-mesh/pull/361), [#363](https://github.com/0xProject/0x-mesh/pull/363), [#366](https://github.com/0xProject/0x-mesh/pull/366)). By default the WebSocket transport listens on port `60559` but this can be changed via the `P2P_WEBSOCKETS_PORT` environment variable.
- Created TypeScript bindings and an NPM package for running Mesh directly in the browser ([#369](https://github.com/0xProject/0x-mesh/pull/369)). Documentation for the NPM package and a guide for running Mesh in the browser can be found at [https://0x-org.gitbook.io/mesh/](https://0x-org.gitbook.io/mesh/).
- Added ability to use custom bootstrap list via the `BOOTSTRAP_LIST` environment variable ([#374](https://github.com/0xProject/0x-mesh/pull/374)). Typically this should only be used for testing/debugging.
- Added WebAssembly/Browser support to packages that previously did not support it ([#358](https://github.com/0xProject/0x-mesh/pull/358), [#359](https://github.com/0xProject/0x-mesh/pull/359), [#366](https://github.com/0xProject/0x-mesh/pull/366)).
- Order hash calculations are now cached, which slightly improves performance ([#365](https://github.com/0xProject/0x-mesh/pull/365)).
- Refactored `BlockWatch` so that it can be used without using `LevelDB` for Ethereum block storage. ([#355](https://github.com/0xProject/0x-mesh/pull/355))

### Bug fixes 🐞 

- Fixed two related bugs: One where order expiration events would be emitted multiple times and another that meant subsequent fill/cancel events for orders deemed expired were not emitted. Fills/cancels for expired orders will continue to be emitted if they occur within ~4 mins (i.e. 20 blocks) of the expiration ([#385](https://github.com/0xProject/0x-mesh/pull/385)).
- Fixed a data race-condition in OrderWatcher that could have caused order collection updates to be overwritten in the DB. ([#386](https://github.com/0xProject/0x-mesh/pull/386))
- Fixed a bug where `fillableTakerAssetAmount` and `lastUpdated` were not always being properly updated in the DB. ([#386](https://github.com/0xProject/0x-mesh/pull/386))
- Fixed some issues with key prefixes for certain types not being applied correctly to logs ([#375](https://github.com/0xProject/0x-mesh/pull/375)).
- Fixed an issue where order hashes were not being correctly logged ([#368](https://github.com/0xProject/0x-mesh/pull/368)).
- Mesh will now properly shut down if the database is unexpectedly closed ([#370](https://github.com/0xProject/0x-mesh/pull/370)).

## v3.0.1-beta

### Bug fixes 🐞 

- Fixed bug where block number would sometimes be converted to hex with a leading zero, an invalid hex value per the [Ethereum JSON-RPC specification](https://github.com/ethereum/wiki/wiki/JSON-RPC#hex-value-encoding). ([#353](https://github.com/0xProject/0x-mesh/pull/353))
- Fixed bug which resulted in orders that were close to expiring being re-added and removed multiple times, resulting in multiple ADDED and EXPIRED events for the same order ([#352](https://github.com/0xProject/0x-mesh/pull/352)).

## v3.0.0-beta

### Breaking changes 🛠 

- Modified Mesh's validation logic to reject and consider invalid any _partially fillable_ orders. While this is
  technically a breaking change, partially fillable orders are rare in the wild and we don't expect this will
  affect many users. ([#333](https://github.com/0xProject/0x-mesh/pull/333))
- Lowercased `GetStatsAsync` method to `getStatsAsync` in TS client

### Bug fixes 🐞 

- De-dup order submitted via the JSON-RPC method `mesh_addOrders` before performing validation ([#331](https://github.com/0xProject/0x-mesh/pull/331))
- Added `"declaration": true,` to TS client's `tsconfig.json` so that downstream projects can use it's TS typings. ([#325](https://github.com/0xProject/0x-mesh/pull/325))


## v2.0.0-beta

### Breaking changes 🛠 

- Modified how `mesh_addOrders` treats orders that are already stored on the Mesh node. Previously, they would be rejected with code `OrderAlreadyStored`. Now, if the order is stored and fillable, it will be accepted. If it is stored but unfillable, it will be rejected with `OrderAlreadyStoredAndUnfillable`. We additionally added a `isNew` property to the accepted orderInfos returned, so that callers can discern which orders Mesh already knew about. ([#316](https://github.com/0xProject/0x-mesh/pull/316))

### Features ✅ 

- Added backup bootstrap nodes provided by the libp2p community
- Improved log formatting and reduced verbosity in a few cases ([#314](https://github.com/0xProject/0x-mesh/pull/314), [#287](https://github.com/0xProject/0x-mesh/pull/287))
- Reduced AdvertiseBootDelay for bootstrap nodes
- Implemented a check that will alerts you when switching to a different Ethereum network ID. ([#301](https://github.com/0xProject/0x-mesh/pull/301)) -- special thanks to @hrharder!
- Made environment variable parsing more generous by automatically removing quotes if needed ([#306](https://github.com/0xProject/0x-mesh/pull/306))
- Improved tests by adding timeouts and closing resources where appropriate ([#310](https://github.com/0xProject/0x-mesh/pull/310), [#309](https://github.com/0xProject/0x-mesh/pull/309), [#308](https://github.com/0xProject/0x-mesh/pull/308))
- Increased robustness by removing panics and failing more gracefully ([#312](https://github.com/0xProject/0x-mesh/pull/312))
- RPC server is now started while block event backfilling is happening under the hood instead of waiting for it to complete ([#318](https://github.com/0xProject/0x-mesh/pull/318))
- Added a `mesh_getStats` endpoint which returns a host of useful information about the state of the Mesh node (e.g., number of fillable order stored, number of peers, peerID, etc...) ([#322](https://github.com/0xProject/0x-mesh/pull/322))

### Bug fixes 🐞 

- Log messages are no longer incorrectly fired when receiving orders which have already been seen ([#286](https://github.com/0xProject/0x-mesh/pull/286))
- Fixed a bug where Mesh was still running after the database was closed ([#300](https://github.com/0xProject/0x-mesh/pull/300))
- Handled Parity "unknown block" error gracefully like we do Geth's ([#285](https://github.com/0xProject/0x-mesh/pull/285))

## v1.0.6-beta

This release fixes several bugs:

- Uninitialized TxHashes map & accidental inclusion of null address txHash in order events ([#280](https://github.com/0xProject/0x-mesh/pull/280))
- Concurrent read/write issue in OrderWatcher's EventDecoder ([#278](https://github.com/0xProject/0x-mesh/issues/278))
- Non-unique logging keys causing Elastic Search indexing issues ([#275](https://github.com/0xProject/0x-mesh/pull/275))

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
