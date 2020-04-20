- [Docker image](https://hub.docker.com/r/0xorg/mesh/tags)
- [README](https://github.com/0xProject/0x-mesh/blob/v9.3.0/README.md)

## Summary

### Features ✅

- Mesh now ensures on startup that the chain ID of your Ethereum RPC endpoint matches config.EthereumChainID [#733](https://github.com/0xProject/0x-mesh/pull/733).

### Bug fixes 🐞

- Fixed a compatibility issue in `@0x/mesh-browser-lite` for Safari and some other browsers [#770](https://github.com/0xProject/0x-mesh/pull/770).
- Fixes an issue that would allow expired orders to be returned in `GetOrders` [773](http://github.com/0xProject/0x-mesh/pull/773)
- Fixed a rare bug where ERC721 approval events could be missed [#782](https://github.com/0xProject/0x-mesh/pull/782)



