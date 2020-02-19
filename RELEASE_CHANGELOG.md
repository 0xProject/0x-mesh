- [Docker image](https://hub.docker.com/r/0xorg/mesh/tags)
- [README](https://github.com/0xProject/0x-mesh/blob/v9.0.0/README.md)

## Summary

### Breaking changes 🛠

- As a result of implementing custom order filters, some of the code Mesh uses under the hood to share orders with peers has changed. As a result this version of Mesh cannot share orders with any older versions and vice versa ([#630](https://github.com/0xProject/0x-mesh/pull/630)).
- Implemented a new protocol for sharing existing orders with peers. This will drastically reduce bandwidth and CPU usage and increase the speed at which _new_ orders are propagated. ([#692](https://github.com/0xProject/0x-mesh/pull/692)).
- Rename `RPC_ADDR` to `WS_RPC_ADDR` since we now support both WS and HTTP JSON-RPC endpoints. ([#658](https://github.com/0xProject/0x-mesh/pull/658))

### Features ✅

- Implemented custom order filters, which allow users to filter out all but the orders they care about. When a custom order filter is specified, Mesh will only send and receive orders that pass the filter. ([#630](https://github.com/0xProject/0x-mesh/pull/630)).
- Developers can now override the contract addresses for any testnet using the `CUSTOM_CONTRACT_ADDRESSES` env config ([#640](https://github.com/0xProject/0x-mesh/pull/640)).
- Added `getOrdersForPageAsync` method to `@0x/mesh-rpc-client` WS client interface so that clients can paginate through the retrieved orders themselves ([#642](https://github.com/0xProject/0x-mesh/pull/642)).
- Added support for passing in your own Web3 provider when using the `@0x/mesh-browser` package. ([#665](https://github.com/0xProject/0x-mesh/pull/665)).
- Add support for orders involving Chai ERC20Bridge assetData ([#663](https://github.com/0xProject/0x-mesh/pull/663))
- Add support for calling JSON-RPC methods over HTTP (env config `HTTP_RPC_ADDR` defaults to `localhost:60556`). ([#658](https://github.com/0xProject/0x-mesh/pull/658))

### Bug fixes 🐞

- Fixed some of the browser typescript bindings to be consistent with the Go and smart contract implementations ([#697](https://github.com/0xProject/0x-mesh/pull/697)).


