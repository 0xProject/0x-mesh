# CHANGELOG

This changelog is a work in progress and may contain notes for versions which have not actually been released. Check the [Releases](https://github.com/0xProject/0x-mesh/releases) page to see full release notes and more information about the latest released versions.

## v10.0.0

### Features ✅

-   Upgraded to Go 1.14, which contains several WebAssembly performance improvements [#815](https://github.com/0xProject/0x-mesh/pull/815).
-   Reject orders that have taker addresses that are non-zero and are not whitelisted. [#896](https://github.com/0xProject/0x-mesh/pull/896).

## v9.4.1

### Bug fixes 🐞

-   Fixed a problem in the filtered pagination subprotocols of ordersync that caused the nodes to use the wrong orderfilter [#882](https://github.com/0xProject/0x-mesh/pull/882)

## v9.4.0

### Features ✅

-   Improved the performance of validating order filters in the browser [#809](https://github.com/0xProject/0x-mesh/pull/809).

## v9.3.0

### Features ✅

-   Mesh now ensures on startup that the chain ID of your Ethereum RPC endpoint matches config.EthereumChainID [#733](https://github.com/0xProject/0x-mesh/pull/733).

### Bug fixes 🐞

-   Fixed a compatibility issue in `@0x/mesh-browser-lite` for Safari and some other browsers [#770](https://github.com/0xProject/0x-mesh/pull/770).
-   Fixes an issue that would allow expired orders to be returned in `GetOrders` [773](http://github.com/0xProject/0x-mesh/pull/773)
-   Fixed a rare bug where ERC721 approval events could be missed [#782](https://github.com/0xProject/0x-mesh/pull/782)

## v9.2.1

### Bug fixes 🐞

-   Fixed a critical bug in the ordersync protocol which resulted in only 50% of existing orders being shared when a new peer joins the network. New orders are shared separately and were unaffected. [#760](https://github.com/0xProject/0x-mesh/pull/760).

## v9.2.0

### Features ✅

-   Greatly reduced latency for propagating orders, especially for browser nodes [#756](https://github.com/0xProject/0x-mesh/pull/756).
-   Added support for `checkGasPrice` StaticCall asset data [#744](https://github.com/0xProject/0x-mesh/pull/744)

## v9.1.0

### Features ✅

-   Improved speed and efficiency of peer discovery, especially when using custom order filters [#729](https://github.com/0xProject/0x-mesh/pull/729).
-   Added a lightweight package to use for loading Mesh's Wasm binary in a streaming manner [#707](https://github.com/0xProject/0x-mesh/pull/707).

### Bug fixes 🐞

-   Fixed an issue where incoming orders could sometimes be dropped by peers [#732](https://github.com/0xProject/0x-mesh/pull/732).

## v9.0.1

### Bug fixes 🐞

-   Fix bug where we weren't enforcing that we never store more than `miniHeaderRetentionLimit` block headers in the DB. This caused [issue #667](https://github.com/0xProject/0x-mesh/issues/667) and also caused the Mesh node's DB storage to continuously grow over time. ([#716](https://github.com/0xProject/0x-mesh/pull/716))

## v9.0.0

### Breaking changes 🛠

-   As a result of implementing custom order filters, some of the code Mesh uses under the hood to share orders with peers has changed. As a result this version of Mesh cannot share orders with any older versions and vice versa ([#630](https://github.com/0xProject/0x-mesh/pull/630)).
-   Implemented a new protocol for sharing existing orders with peers. This will drastically reduce bandwidth and CPU usage and increase the speed at which _new_ orders are propagated. ([#692](https://github.com/0xProject/0x-mesh/pull/692)).
-   Rename `RPC_ADDR` to `WS_RPC_ADDR` since we now support both WS and HTTP JSON-RPC endpoints. ([#658](https://github.com/0xProject/0x-mesh/pull/658))

### Features ✅

-   Implemented custom order filters, which allow users to filter out all but the orders they care about. When a custom order filter is specified, Mesh will only send and receive orders that pass the filter. ([#630](https://github.com/0xProject/0x-mesh/pull/630)).
-   Developers can now override the contract addresses for any testnet using the `CUSTOM_CONTRACT_ADDRESSES` env config ([#640](https://github.com/0xProject/0x-mesh/pull/640)).
-   Added `getOrdersForPageAsync` method to `@0x/mesh-rpc-client` WS client interface so that clients can paginate through the retrieved orders themselves ([#642](https://github.com/0xProject/0x-mesh/pull/642)).
-   Added support for passing in your own Web3 provider when using the `@0x/mesh-browser` package. ([#665](https://github.com/0xProject/0x-mesh/pull/665)).
-   Add support for orders involving Chai ERC20Bridge assetData ([#663](https://github.com/0xProject/0x-mesh/pull/663))
-   Add support for calling JSON-RPC methods over HTTP (env config `HTTP_RPC_ADDR` defaults to `localhost:60556`). ([#658](https://github.com/0xProject/0x-mesh/pull/658))

### Bug fixes 🐞

-   Fixed some of the browser typescript bindings to be consistent with the Go and smart contract implementations ([#697](https://github.com/0xProject/0x-mesh/pull/697)).

## v8.2.0

### Features ✅

-   Added `getStatsAsync` to the `@0x/mesh-browser` package. ([#654](https://github.com/0xProject/0x-mesh/pull/654)).
-   Added `getOrdersAsync` and `getOrdersForPageAsync` to the `@0x/mesh-browser` package. ([#655](https://github.com/0xProject/0x-mesh/pull/655)).

### Bug fixes 🐞

-   Update DevUtils contract address to fix intermittent revert issues. ([#671](https://github.com/0xProject/0x-mesh/pull/671)).

## v8.1.2

### Bug fixes 🐞

-   Update DevUtils contract to version that removed maker transfer simulation. ([#662](https://github.com/0xProject/0x-mesh/pull/662)).
-   Fix faulty Go to Javascript conversion logic. ([#659](https://github.com/0xProject/0x-mesh/pull/659)).
-   Updated dockerfiles to work with Go modules. ([#646](https://github.com/0xProject/0x-mesh/pull/646)).
-   Update DevUtils mainnet contract address to version that fixes MAP order validation issue ([#644](https://github.com/0xProject/0x-mesh/pull/644)).

## v8.1.1

### Bug fixes 🐞

-   Fixed a regression which can result in memory leaks. ([#650](https://github.com/0xProject/0x-mesh/pull/650)).

## v8.1.0

### Features ✅

-   Reduced startup time for Mesh node by only waiting for a block to be processed if Mesh isn't already sync'ed up to the latest block. ([#622](https://github.com/0xProject/0x-mesh/pull/622))
-   Increased the maximum size for encoded orders from ~8kB to 16kB ([#631](https://github.com/0xProject/0x-mesh/pull/631)).

### Bug fixes 🐞

-   Fixed a typo ("rendervouz" --> "rendezvous") in GetStatsResponse. ([#611](https://github.com/0x-mesh/pull/611)).
-   Fixed a bug where we attempted to update the same order multiple times in a single DB txn, causing the later update to noop. ([#623](https://github.com/0xProject/0x-mesh/pull/623)).
-   Fixed a bug which could cause Mesh to exit if a re-org condition occurs causing a block to be added and removed within the same block sync operation. ([#614](https://github.com/0xProject/0x-mesh/pull/614)).

## v8.0.0-beta-0xv3

### Breaking changes 🛠

-   Changed the response from `@0x/mesh-ts-client`'s `getOrdersAsync` endpoint to include the `snapshotID` and `snapshotTimestamp` at which the Mesh DB was queried along with the orders found. ([#591](https://github.com/0xProject/0x-mesh/pull/591))
-   Increased the default `ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC` from 100k to 200k ([#596](https://github.com/0xProject/0x-mesh/pull/596)).

### Features ✅

-   Instead of progressing Mesh forward by a single block on every invocation of the `BLOCK_POLLING_INTERVAL`, we now attempt to sync as many blocks as necessary to reach the latest block available. This will reduce the chances of Mesh becoming out-of-sync with it's backing Ethereum node ([#564](https://github.com/0xProject/0x-mesh/pull/564))
-   Added a new environment variable `ENABLE_ETHEREUM_RPC_RATE_LIMITING` and config option `enableEthereumRPCRateLimiting` which can be used to completely disable Mesh's internal Ethereum RPC rate limiting features. By default it is enabled, and disabling can have some consequences depending on your RPC provider. ([#584](https://github.com/0xProject/0x-mesh/pull/584))
-   Added a `SnapshotTimestamp` field to `GetOrdersResponse`, the return type of the `mesh_getOrders` RPC method. This way, the caller can know at what point in time the snapshot had been created. ([#591](https://github.com/0xProject/0x-mesh/pull/591))
-   Improved batching of events emitted from order events subscriptions ([#566](https://github.com/0xProject/0x-mesh/pull/566))
-   Added timestamp to order events ([#602](https://github.com/0xProject/0x-mesh/pull/602))

### Bug fixes 🐞

-   Fixed an issue where order updates could have been missed if Mesh discovered blocks but didn't have time to process them before getting shut down. Now, blocks are only persisted to the DB once any order updates resulting from it have been processed. ([#566](https://github.com/0xProject/0x-mesh/pull/566)).
-   Fixed a race-condition when adding new orders to Mesh which could result in order-relevant events being missed if they occured very soon after the order was submitted and the order validation RPC call took a long time ([#566](https://github.com/0xProject/0x-mesh/pull/566)).
-   Upgraded the `web3-provider` dependency used by `@0x/mesh-rpc-client` in order to fix a bug where it was requiring either `process` OR `window` to exist in the global scope ([#601](https://github.com/0xProject/0x-mesh/pull/601)).
-   Fixed an issue where the internal Ethereum RPC rate limiter could be too aggressive in certain scenarios ([#596](https://github.com/0xProject/0x-mesh/pull/596)).
-   Add a default RPC request timeout of 30sec to all non-subscription requests sent by `@0x/mesh-rpc-client` to avoid the client from hanging endlessly if it misses a response due to network disruption ([#603](https://github.com/0xProject/0x-mesh/pull/603)).

## v7.2.2-beta-0xv3

### Bug fixes 🐞

-   Fixed a bug where `makerAssetFeeData` or `takerAssetFeeData` were sometimes being set to an empty string instead of `0x` in the browser ([#579](https://github.com/0xProject/0x-mesh/pull/579)).

## v7.2.1-beta-0xv3

### Bug fixes 🐞

-   Fixed a bug which could cause Mesh to crash with a nil pointer exception if RPC requests are sent too quickly during/immediately after start up ([#560](https://github.com/0xProject/0x-mesh/pull/560)).

## v7.2.0-beta-0xv3

-   Update V3 contract addresses for Ganache snapshot. ([#569](https://github.com/0xProject/0x-mesh/pull/569))

## v7.1.1-beta-0xv3

### Bug fixes 🐞

-   Fixed a bug where the internal order event feed could be come blocked, rendering Mesh unable to receive any new orders or update existing ones ([#552](https://github.com/0xProject/0x-mesh/pull/552)).

## v7.1.0-beta-0xV3

-   Update V3 contract addresses for mainnet/testnets. ([#547](https://github.com/0xProject/0x-mesh/pull/547))

### Features ✅

-   Added support for persistence in the browser. Users of the `@0x/mesh-browser` package will now be able to retain orders and other relevant parts of the state when refreshing the page or closing and re-opening the browser. ([#533](https://github.com/0xProject/0x-mesh/pull/533)).

### Bug fixes 🐞

-   Fix bug where Mesh nodes were logging receipt and re-sharing with peers duplicate orders already stored in it's DB, if the duplicate order was submitted via JSON-RPC. ([#529](https://github.com/0xProject/0x-mesh/pull/529))
-   Add missing `UNEXPIRED` `OrderEventEndState` enum value to both `@0x/mesh-rpc-client` and `@0x/mesh-browser` and missing `STOPPED_WATCHING` value from `@0x/mesh-rpc-client`.
-   Fixed a potential memory leak by using the latest version of `github.com/libp2p/go-libp2p-kad-dht` ([#539](https://github.com/0xProject/0x-mesh/pull/539)).
-   Changed the default port for `RPC_ADDR` from a random available port to `60557`. _Some_ documentation already assumed `60557` was the default port. Now all documentation has been updated for consistency with this change. ([#542](https://github.com/0xProject/0x-mesh/pull/542)).
-   Fixed a potential nil pointer exception in log hooks ([#543](https://github.com/0xProject/0x-mesh/pull/543)).
-   Fixed a bug where successful closes of an rpc subscription were being reported as errors ([#544](https://github.com/0xProject/0x-mesh/pull/544)).
-   We now log the error and stack trace if an RPC method panics. Before, these errors were swallowed by the panic recovery logic in `go-ethereum`'s `rpc` package. ([#545](https://github.com/0xProject/0x-mesh/pull/545))
-   Previously, we used to fast-sync block events missed since a Mesh node was last online. If this was more than 128 blocks ago, the fast-sync would fail if Mesh was not connected to an Ethereum node with the `--archive` flag enabled. We now fast-sync only if less than 128 blocks have elapsed. Otherwise, we simply re-validate all orders and continue processing block events from the latest block. ([#407](https://github.com/0xProject/0x-mesh/issues/407))

## v7.0.4-beta-0xv3

-   Upgraded `@0x` deps in `@0x/mesh-rpc-client`

## v7.0.3-beta-0xv3

-   Upgrade BigNumber dep used by `@0x/mesh-rpc-client` to `~9.0.0` ([#527](https://github.com/0xProject/0x-mesh/pull/527))

## v7.0.2-beta-0xv3

### Bug fixes 🐞

-   Upgrade Mesh's `libp2p/go-flow-metrics` dep in an attempt to fix a bug in the flow metrics module. ([#521](https://github.com/0xProject/0x-mesh/pull/521))
-   Updates `DevUtils.sol` addresses on all networks to the latest version which fixed a deploy issue. ([#520](https://github.com/0xProject/0x-mesh/pull/520))

## v7.0.1-beta-0xv3

### Bug fixes 🐞

-   Fixed an oversight which granted immunity from bandwidth banning for any peer using a relayed connection ([#509](https://github.com/0xProject/0x-mesh/pull/509)).
-   Fixed a typo in the `@0x/mesh-browser` package that resulted in some config options not being passed through correctly ([#502](https://github.com/0xProject/0x-mesh/pull/502)).
-   Fixed a bug in ETH JSON-RPC rate limiter where not all dates were being properly converted to UTC, causing Mesh to malfunction if the local time was a day earlier or later than UTC. ([#505](https://github.com/0xProject/0x-mesh/pull/505))
-   Fixed a bug in the TypeScript RPC client that prevented orders from being added ([#514](https://github.com/0xProject/0x-mesh/pull/514)).

## v7.0.0-beta-0xv3

### Breaking changes 🛠

_Note:_ This release will require wiping your Mesh's DB before upgrading. The DB location defaults to `./0x_mesh/db`.

-   Removed `RPC_PORT` environment variable. The new `RPC_ADDR` environment variable allows specifying both the interface and port ([#487](https://github.com/0xProject/0x-mesh/pull/487)).
-   Due to security concerns and new rate limiting mechanisms, the default bind address for the RPC API has changed from `0.0.0.0` to `localhost`. Users who previously did not set `RPC_PORT` may need to now manually set `RPC_ADDR` to enable other applications to access the RPC API. If you are using Docker Compose, we recommend using [links](https://docs.docker.com/compose/networking/#links). If you do need to set `RPC_ADDR` to bind on `0.0.0.0`, please be aware of the security implications and consider protecting access to Mesh via a third-party firewall. (See [#444](https://github.com/0xProject/0x-mesh/pull/444) and [#487](https://github.com/0xProject/0x-mesh/pull/487) for more details).
-   Changed the `EXPIRED` event such that it is emitted when an order is expired according to the latest block timestamp, not anymore based on UTC time. ([#490](https://github.com/0xProject/0x-mesh/pull/490))
-   Removed the `EXPIRATION_BUFFER_SECONDS` env config since we no longer compute order expiration using UTC time. ([#490](https://github.com/0xProject/0x-mesh/pull/490))

### Features ✅

-   Added an `UNEXPIRED` order event kind which is emitted for orders that were previously considered expired but due to a block-reorg causing the latest block timestamp to be earlier than the previous latest block timestamp, are no longer expired. ([#490](https://github.com/0xProject/0x-mesh/pull/490))
-   Added support for decoding Axie Infinity `Transfer` and `Approve` ERC721 events which differ from the ERC721 standard. ([#494](https://github.com/0xProject/0x-mesh/pull/494))

### Bug fixes 🐞

-   Fixed a bug in the Go RPC client which resulted in errors when receving order events with at least one contract event ([#496](https://github.com/0xProject/0x-mesh/pull/496)).

## v6.0.0-beta-0xv3

### Breaking changes 🛠

-   Renamed env config from `ETHEREUM_NETWORK_ID` to `ETHEREUM_CHAIN_ID` since `network` is a misnomer here and what we actually care about is the `chainID`. Most chains have the same id for their p2p network and chain. From the ones we support, the only outlier is Ganache, for which you will now need to supply `1337` instead of `50` (Learn more: https://medium.com/@pedrouid/chainid-vs-networkid-how-do-they-differ-on-ethereum-eec2ed41635b) ([#485](https://github.com/0xProject/0x-mesh/pull/485))
-   Rejected order code `OrderForIncorrectNetwork` has been changed to `OrderForIncorrectChain` ([#485](https://github.com/0xProject/0x-mesh/pull/485))

### Features ✅

-   Implemented a new strategy for limiting the amount of database storage used by Mesh and removing orders when the database is full. This strategy involves a dynamically adjusting maximum expiration time. When the database is full, Mesh will enforce a maximum expiration time for all incoming orders and remove any existing orders with an expiration time too far in the future. If conditions change and there is enough space in the database again, the max expiration time will slowly increase. This is a short term solution which solves the immediate issue of finite storage capacities and does a decent job of protecting against spam. We expect to improve and possibly replace it in the future. See [#450](https://github.com/0xProject/0x-mesh/pull/450) for more details.
-   Added support for a new feature called "order pinning" ([#474](https://github.com/0xProject/0x-mesh/pull/474)). Pinned orders will not be affected by any DDoS prevention or incentive mechanisms (including the new dynamic max expiration time feature) and will always stay in storage until they are no longer fillable. By default, all orders which are submitted via either the JSON-RPC API or the `addOrdersAsync` function in the TypeScript bindings will be pinned.
-   Re-enabled bandwidth-based peer banning with a workaround to deal with erroneous spikes [#478](https://github.com/0xProject/0x-mesh/pull/478).

### Bug fixes 🐞

-   Improved the aggressiveness at which we permanently delete orders that have been flagged for removal. Previously we would wait for the cleanup job to handle this (once an hour), but that meant many removed orders would accumulate. We now prune them every 5 minutes. ([#471](https://github.com/0xProject/0x-mesh/pull/471))

## v5.1.0-beta

### Features ✅

-   The `getStats` RPC endpoint now includes a new field which accounts for the number of orders that have been marked as "removed" but not yet permanently deleted ([#461](https://github.com/0xProject/0x-mesh/pull/461)).
-   Improved historical order sharing using round-robin algorithm instead of random selection ([#454](https://github.com/0xProject/0x-mesh/pull/454)). This will reduce the warm-up time for receiving existing orders when first joining the network.
-   Added ERC1155 assetData support ([#453](https://github.com/0xProject/0x-mesh/pull/453)). This includes order watching and order events for orders involving ERC1155 tokens.
-   Added Ability to specify custom contract addresses via the `CUSTOM_ADDRESSES` environment variable or the `customAddresses` field in the TypeScript bindings ([#445](https://github.com/0xProject/0x-mesh/pull/445)).

### Bug fixes 🐞

-   Temporarily disabled bandwidth-based peer banning ([#448](https://github.com/0xProject/0x-mesh/pull/448)). A [bug in libp2p](https://github.com/libp2p/go-libp2p-core/issues/65) was occasionally causing observed bandwidth usage to spike to unrealistic levels, which can result in peers being erroneously banned. We decided to temporarily stop banning peers while we're working with the libp2p team to resolve the issue.

## v5.0.2-beta-0xv3

### Bug fixes 🐞

-   Updated DevUtils.sol contract address for the Kovan network to one including a bug fix for validating orders with nulled out `feeAssetData` fields. ([#446](https://github.com/0xProject/0x-mesh/pull/446)])
-   Added back Ropsten and Rinkeby support and fixed `exchange` address on Kovan ([#451](https://github.com/0xProject/0x-mesh/pull/451))
-   Segregate V3 Mesh network from V2 network ([#455](https://github.com/0xProject/0x-mesh/pull/455))

## v5.0.0-beta

### Breaking changes 🛠

-   Removes the `txHashes` key in the `OrderEvent`s emitted from the `orders` JSON-RPC subscription and replaced it with `contractEvents`, an array of decoded order-relevant contract events. Parsing these events allows callers to find every discrete order fill/cancel event. ([#420](https://github.com/0xProject/0x-mesh/pull/420))
-   Renames the `Kind` key in `OrderEvent` to `EndState` to better elucidate that it represents the aggregate change to the orders state since it was last re-validated. As an end state, it does not capture any possible intermediate states the order might have been in since the last re-validation. Intermediate states can be inferred from the `contractEvents` included ([#420](https://github.com/0xProject/0x-mesh/pull/420))

### Features ✅

-   Removed the max expiration limit for orders. The only remaining expiration constraint is that the unix timestamp does not overflow int64 (i.e., is not larger than 9223372036854775807). ([#400](https://github.com/0xProject/0x-mesh/pull/400))

### Bug fixes 🐞

-   Fixed bug where we weren't updating an orders `fillableTakerAssetAmount` in the DB when orders were being partially filled or when their fillability increased due to a block re-org. ([#439](https://github.com/0xProject/0x-mesh/pull/439))
-   Made `verbosity` field optional in the TypeScript `Config` type. ([#410](https://github.com/0xProject/0x-mesh/pull/410))
-   Fixed issue where we weren't re-validating orders potentially impacted by the balance increase of the recipient of an ERC20 or ERC721 transfer. ([#416](https://github.com/0xProject/0x-mesh/pull/416))

## v4.0.1-beta

### Bug fixes 🐞

-   Fixed a DB transaction deadlock accidentally introduced in the v4.0.0-beta release. ([#403](https://github.com/0xProject/0x-mesh/pull/403))

## v4.0.0-beta

### Breaking changes 🛠

-   Renamed the environment variable `P2P_LISTEN_PORT` to `P2P_TCP_PORT` ([#366](https://github.com/0xProject/0x-mesh/pull/366)). This makes it possible to configure Mesh to use both the TCP and Websocket transports by listening on different ports.

### Features ✅

-   Enabled WebSocket transport for bootstrap nodes and all other nodes ([#361](https://github.com/0xProject/0x-mesh/pull/361), [#363](https://github.com/0xProject/0x-mesh/pull/363), [#366](https://github.com/0xProject/0x-mesh/pull/366)). By default the WebSocket transport listens on port `60559` but this can be changed via the `P2P_WEBSOCKETS_PORT` environment variable.
-   Created TypeScript bindings and an NPM package for running Mesh directly in the browser ([#369](https://github.com/0xProject/0x-mesh/pull/369)). Documentation for the NPM package and a guide for running Mesh in the browser can be found at [https://0x-org.gitbook.io/mesh/](https://0x-org.gitbook.io/mesh/).
-   Added ability to use custom bootstrap list via the `BOOTSTRAP_LIST` environment variable ([#374](https://github.com/0xProject/0x-mesh/pull/374)). Typically this should only be used for testing/debugging.
-   Added WebAssembly/Browser support to packages that previously did not support it ([#358](https://github.com/0xProject/0x-mesh/pull/358), [#359](https://github.com/0xProject/0x-mesh/pull/359), [#366](https://github.com/0xProject/0x-mesh/pull/366)).
-   Order hash calculations are now cached, which slightly improves performance ([#365](https://github.com/0xProject/0x-mesh/pull/365)).
-   Refactored `BlockWatch` so that it can be used without using `LevelDB` for Ethereum block storage. ([#355](https://github.com/0xProject/0x-mesh/pull/355))

### Bug fixes 🐞

-   Fixed two related bugs: One where order expiration events would be emitted multiple times and another that meant subsequent fill/cancel events for orders deemed expired were not emitted. Fills/cancels for expired orders will continue to be emitted if they occur within ~4 mins (i.e. 20 blocks) of the expiration ([#385](https://github.com/0xProject/0x-mesh/pull/385)).
-   Fixed a data race-condition in OrderWatcher that could have caused order collection updates to be overwritten in the DB. ([#386](https://github.com/0xProject/0x-mesh/pull/386))
-   Fixed a bug where `fillableTakerAssetAmount` and `lastUpdated` were not always being properly updated in the DB. ([#386](https://github.com/0xProject/0x-mesh/pull/386))
-   Fixed some issues with key prefixes for certain types not being applied correctly to logs ([#375](https://github.com/0xProject/0x-mesh/pull/375)).
-   Fixed an issue where order hashes were not being correctly logged ([#368](https://github.com/0xProject/0x-mesh/pull/368)).
-   Mesh will now properly shut down if the database is unexpectedly closed ([#370](https://github.com/0xProject/0x-mesh/pull/370)).

## v3.0.1-beta

### Bug fixes 🐞

-   Fixed bug where block number would sometimes be converted to hex with a leading zero, an invalid hex value per the [Ethereum JSON-RPC specification](https://github.com/ethereum/wiki/wiki/JSON-RPC#hex-value-encoding). ([#353](https://github.com/0xProject/0x-mesh/pull/353))
-   Fixed bug which resulted in orders that were close to expiring being re-added and removed multiple times, resulting in multiple ADDED and EXPIRED events for the same order ([#352](https://github.com/0xProject/0x-mesh/pull/352)).

## v3.0.0-beta

### Breaking changes 🛠

-   Modified Mesh's validation logic to reject and consider invalid any _partially fillable_ orders. While this is
    technically a breaking change, partially fillable orders are rare in the wild and we don't expect this will
    affect many users. ([#333](https://github.com/0xProject/0x-mesh/pull/333))
-   Lowercased `GetStatsAsync` method to `getStatsAsync` in TS client

### Bug fixes 🐞

-   De-dup order submitted via the JSON-RPC method `mesh_addOrders` before performing validation ([#331](https://github.com/0xProject/0x-mesh/pull/331))
-   Added `"declaration": true,` to TS client's `tsconfig.json` so that downstream projects can use it's TS typings. ([#325](https://github.com/0xProject/0x-mesh/pull/325))

## v2.0.0-beta

### Breaking changes 🛠

-   Modified how `mesh_addOrders` treats orders that are already stored on the Mesh node. Previously, they would be rejected with code `OrderAlreadyStored`. Now, if the order is stored and fillable, it will be accepted. If it is stored but unfillable, it will be rejected with `OrderAlreadyStoredAndUnfillable`. We additionally added a `isNew` property to the accepted orderInfos returned, so that callers can discern which orders Mesh already knew about. ([#316](https://github.com/0xProject/0x-mesh/pull/316))

### Features ✅

-   Added backup bootstrap nodes provided by the libp2p community
-   Improved log formatting and reduced verbosity in a few cases ([#314](https://github.com/0xProject/0x-mesh/pull/314), [#287](https://github.com/0xProject/0x-mesh/pull/287))
-   Reduced AdvertiseBootDelay for bootstrap nodes
-   Implemented a check that will alerts you when switching to a different Ethereum network ID. ([#301](https://github.com/0xProject/0x-mesh/pull/301)) -- special thanks to @hrharder!
-   Made environment variable parsing more generous by automatically removing quotes if needed ([#306](https://github.com/0xProject/0x-mesh/pull/306))
-   Improved tests by adding timeouts and closing resources where appropriate ([#310](https://github.com/0xProject/0x-mesh/pull/310), [#309](https://github.com/0xProject/0x-mesh/pull/309), [#308](https://github.com/0xProject/0x-mesh/pull/308))
-   Increased robustness by removing panics and failing more gracefully ([#312](https://github.com/0xProject/0x-mesh/pull/312))
-   RPC server is now started while block event backfilling is happening under the hood instead of waiting for it to complete ([#318](https://github.com/0xProject/0x-mesh/pull/318))
-   Added a `mesh_getStats` endpoint which returns a host of useful information about the state of the Mesh node (e.g., number of fillable order stored, number of peers, peerID, etc...) ([#322](https://github.com/0xProject/0x-mesh/pull/322))

### Bug fixes 🐞

-   Log messages are no longer incorrectly fired when receiving orders which have already been seen ([#286](https://github.com/0xProject/0x-mesh/pull/286))
-   Fixed a bug where Mesh was still running after the database was closed ([#300](https://github.com/0xProject/0x-mesh/pull/300))
-   Handled Parity "unknown block" error gracefully like we do Geth's ([#285](https://github.com/0xProject/0x-mesh/pull/285))

## v1.0.6-beta

This release fixes several bugs:

-   Uninitialized TxHashes map & accidental inclusion of null address txHash in order events ([#280](https://github.com/0xProject/0x-mesh/pull/280))
-   Concurrent read/write issue in OrderWatcher's EventDecoder ([#278](https://github.com/0xProject/0x-mesh/issues/278))
-   Non-unique logging keys causing Elastic Search indexing issues ([#275](https://github.com/0xProject/0x-mesh/pull/275))

It also includes a reduction in the delay before which bootstrap nodes advertise themselves as relays from 15mins to 30sec.

## v1.0.5-beta

This version introduces a temporary hack in order to help some Mesh nodes find their public IP address under certain circumstances.

## v1.0.4-beta

This release fixes a bug in how AutoNAT and AutoRelay services were discovered.

## v1.0.3-beta

This release fixes some networking related issues with our bootstrap nodes.

-   Bootstrap nodes now advertise the correct public IP address.
-   Bootstrap nodes now also operate as a relay.

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

-   Automatic peer discovery via a DHT.
-   A JSON-RPC endpoint for interacting with your Mesh node. It includes support for adding new orders and subscribing to order updates.
-   Efficient order validation and order watching under the hood. You will get notified quickly when orders are expired, canceled, or filled.
-   Basic limitations on the size and types of messages sent by peers. Peers that send malformed messages or messages that are too big will be dropped.

In addition to improving stability and scalability we plan to release many more important features in the near future. Check out https://github.com/0xProject/0x-mesh/issues for more information about what we are working on.
