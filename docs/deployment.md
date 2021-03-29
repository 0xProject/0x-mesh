[![Version](https://img.shields.io/badge/version-11.0.3-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

# 0x Mesh Deployment Guide

Welcome to the [0x Mesh](https://github.com/0xProject/0x-mesh) Deployment Guide!
This guide will walk you through how to configure and deploy your own 0x Mesh
node.

## Supported Chains

-   Mainnet
-   Kovan
-   Ropsten
-   Rinkeby
-   [Ganache snapshot](https://cloud.docker.com/u/0xorg/repository/docker/0xorg/ganache-cli)

## Enabling Telemetry

You can optionally help us develop and maintain Mesh by automatically sending your logs, which
requires a few extra steps. If you are interested in enabling telemetry, check out
[this guide](deployment_with_telemetry.md).

## Running Mesh

Make sure you have Docker installed. Then run:

```bash
docker run \
--restart unless-stopped \
-p 60558:60558 \
-p 60559:60559 \
-e ETHEREUM_CHAIN_ID="1" \
-e ETHEREUM_RPC_URL="{your_ethereum_rpc_url}" \
-e VERBOSITY=5 \
-v {local_path_on_host_machine}/0x_mesh:/usr/mesh/0x_mesh \
0xorg/mesh:{version}
```

1. Replace `{your_ethereum_rpc_url}` with the RPC endpoint for an Ethereum node (e.g. Infura or Alchemy).
2. Replace and `{local_path_on_host_machine}` with a directory on your host machine where all Mesh-related data will be stored.
3. Replace `{version}` with [the latest version](https://github.com/0xProject/0x-mesh/releases) of Mesh.

**Notes:**

-   Ports 60558 and 60559 are the default ports used for communicating with other peers in the network.
-   In order to disable P2P order discovery and sharing, set `USE_BOOTSTRAP_LIST` to `false`.
-   If you are running against a POA testnet (e.g., Kovan), you might want to shorten the `BLOCK_POLLING_INTERVAL` since blocks are mined more frequently then on mainnet. If you do this, your node will use more Ethereum RPC calls, so you will also need to adjust the `ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC` upwards (_warning:_ changing this setting can exceed the limits of your Ethereum RPC provider).
-   If you want to run the mesh in "detached" mode, add the `-d` switch to the docker run command so that your console doesn't get blocked.

## Enabling the GraphQL API

In order to enable the GraphQL API, you just need to add these additional arguments

```bash
-p 60557:60557 \
-e ENABLE_GRAPHQL_SERVER=true \
```

Additionally, to enable the GraphQL playground, just add:

```
-e ENABLE_GRAPHQL_PLAYGROUND=true \
```

Note that the GraphQL API is intended to be _private_. If you enable the GraphQL API on
a production server, we recommend using a firewall or VPC to prevent unauthorized access.
See [the GraphQL API page](graphql_api.md) for more information.

## Persisting State

The Docker container is configured to store all Mesh state (e.g. database files,
private key file) in `/usr/mesh/0x_mesh`. If you want the Mesh state to persist
across Docker container re-starts, use the `-v` flag as included in the command
above to mount a local `0x_mesh` directory into your container. This is strongly
recommended.

## Environment Variables

0x Mesh uses environment variables for configuration. Most environment variables
are detailed in the [Config](https://godoc.org/github.com/0xProject/0x-mesh/core#Config)
struct. They are copied here for convenience, although the source code is
authoritative.

```go
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DataDir is the directory to use for persisting all data, including the
	// database and private key files.
	DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
	// P2PTCPPort is the port on which to listen for new TCP connections from
	// peers in the network. Set to 60558 by default.
	P2PTCPPort int `envvar:"P2P_TCP_PORT" default:"60558"`
	// P2PWebSocketsPort is the port on which to listen for new WebSockets
	// connections from peers in the network. Set to 60559 by default.
	P2PWebSocketsPort int `envvar:"P2P_WEBSOCKETS_PORT" default:"60559"`
	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
	// API.
	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL" json:"-"`
	// EthereumChainID is the chain ID specifying which Ethereum chain you wish to
	// run your Mesh node for
	EthereumChainID int `envvar:"ETHEREUM_CHAIN_ID"`
	// UseBootstrapList is whether to bootstrap the DHT by connecting to a
	// specific set of peers.
	UseBootstrapList bool `envvar:"USE_BOOTSTRAP_LIST" default:"true"`
	// BootstrapList is a comma-separated list of multiaddresses to use for
	// bootstrapping the DHT (e.g.,
	// "/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF").
	// If empty, the default bootstrap list will be used.
	BootstrapList string `envvar:"BOOTSTRAP_LIST" default:""`
	// BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
	// that might contain transactions that impact the fillability of orders stored by Mesh. Different
	// chains have different block producing intervals: POW chains are typically slower (e.g., Mainnet)
	// and POA chains faster (e.g., Kovan) so one should adjust the polling interval accordingly.
	BlockPollingInterval time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s"`
	// EthereumRPCMaxContentLength is the maximum request Content-Length accepted by the backing Ethereum RPC
	// endpoint used by Mesh. Geth & Infura both limit a request's content length to 1024 * 512 Bytes. Parity
	// and Alchemy have much higher limits. When batch validating 0x orders, we will fit as many orders into a
	// request without crossing the max content length. The default value is appropriate for operators using Geth
	// or Infura. If using Alchemy or Parity, feel free to double the default max in order to reduce the
	// number of RPC calls made by Mesh.
	EthereumRPCMaxContentLength int `envvar:"ETHEREUM_RPC_MAX_CONTENT_LENGTH" default:"524288"`
	// EnableEthereumRPCRateLimiting determines whether or not Mesh should limit
	// the number of Ethereum RPC requests it sends. It defaults to true.
	// Disabling Ethereum RPC rate limiting can reduce latency for receiving order
	// events in some network conditions, but can also potentially lead to higher
	// costs or other rate limiting issues outside of Mesh, depending on your
	// Ethereum RPC provider. If set to false, ethereumRPCMaxRequestsPer24HrUTC
	// and ethereumRPCMaxRequestsPerSecond will have no effect.
	EnableEthereumRPCRateLimiting bool `envvar:"ENABLE_ETHEREUM_RPC_RATE_LIMITING" default:"true"`
	// EthereumRPCMaxRequestsPer24HrUTC caps the number of Ethereum JSON-RPC requests a Mesh node will make
	// per 24hr UTC time window (time window starts and ends at midnight UTC). It defaults to 200k but
	// can be increased well beyond this limit depending on your infrastructure or Ethereum RPC provider.
	EthereumRPCMaxRequestsPer24HrUTC int `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC" default:"200000"`
	// EthereumRPCMaxRequestsPerSecond caps the number of Ethereum JSON-RPC requests a Mesh node will make per
	// second. This limits the concurrency of these requests and prevents the Mesh node from getting rate-limited.
	// It defaults to the recommended 30 rps for Infura's free tier, and can be increased to 100 rpc for pro users,
	// and potentially higher on alternative infrastructure.
	EthereumRPCMaxRequestsPerSecond float64 `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_SECOND" default:"30"`
	// CustomContractAddresses is a JSON-encoded string representing a set of
	// custom addresses to use for the configured chain ID. The contract
	// addresses for most common chains/networks are already included by default, so this
	// is typically only needed for testing on custom chains/networks. The given
	// addresses are added to the default list of addresses for known chains/networks and
	// overriding any contract addresses for known chains/networks is not allowed. The
	// addresses for exchange, devUtils, erc20Proxy, erc721Proxy and erc1155Proxy are required
	// for each chain/network. For example:
	//
	//    {
	//        "exchange":"0x48bacb9266a570d521063ef5dd96e61686dbe788",
	//        "devUtils": "0x38ef19fdf8e8415f18c307ed71967e19aac28ba1",
	//        "erc20Proxy": "0x1dc4c1cefef38a777b15aa20260a54e584b16c48",
	//        "erc721Proxy": "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401",
	//        "erc1155Proxy": "0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"
	//    }
	//
	CustomContractAddresses string `envvar:"CUSTOM_CONTRACT_ADDRESSES" default:""`
	// MaxOrdersInStorage is the maximum number of orders that Mesh will keep in
	// storage. As the number of orders in storage grows, Mesh will begin
	// enforcing a limit on maximum expiration time for incoming orders and remove
	// any orders with an expiration time too far in the future.
	MaxOrdersInStorage int `envvar:"MAX_ORDERS_IN_STORAGE" default:"100000"`
	// CustomOrderFilter is a stringified JSON Schema which will be used for
	// validating incoming orders. If provided, Mesh will only receive orders from
	// other peers in the network with the same filter.
	//
	// Here is an example filter which will only allow orders with a specific
	// makerAssetData:
	//
	//    {
	//        "properties": {
	//            "makerAssetData": {
	//                "const": "0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"
	//            }
	//        }
	//    }
	//
	// Note that you only need to include the requirements for your specific
	// application in the filter. The default requirements for a valid order (e.g.
	// all the required fields) are automatically included. For more information
	// on JSON Schemas, see https://json-schema.org/
	CustomOrderFilter string `envvar:"CUSTOM_ORDER_FILTER" default:"{}"`
	// MaxBytesPerSecond is the maximum number of bytes per second that a peer is
	// allowed to send before failing the bandwidth check. Defaults to 5 MiB.
	MaxBytesPerSecond float64 `envvar:"MAX_BYTES_PER_SECOND" default:"5242880"`
	// AdditionalPublicIPSources is a list of external public IP source like
	// https://whatismyip.api.0x.org/ which return the IP address in a
	// text/plain format. This list is prepended to the default sources list.
	//
	// It expects a comma delimited list of external sources for example:
	// ADDITIONAL_PUBLIC_IP_SOURCES="https://ifconfig.me/ip,http://192.168.5.10:1337/ip"
	AdditionalPublicIPSources string `envvar:"ADDITIONAL_PUBLIC_IP_SOURCES" default:""`
}
```

There are some additional environment variable in the [main entrypoint for the
Mesh executable](../cmd/mesh/main.go):

```go
type standaloneConfig struct {
	// EnableGraphQLServer determines whether or not to enable the GraphQL server.
	// If enabled, GraphQL queries can be sent to GraphQLServerAddr at the /graphql
	// URL. By default, the GraphQL server is disabled. Please be aware that the GraphQL
	// API is intended to be a *private* API. If you enable the GraphQL server in
	// production we recommend using a firewall/VPC or an authenticated proxy to restrict
	// public access.
	EnableGraphQLServer bool `envvar:"ENABLE_GRAPHQL_SERVER" default:"false"`
	// GraphQLServerAddr is the interface and port to use for the GraphQL API.
	// By default, 0x Mesh will listen on 0.0.0.0 (all available addresses) and
	// port 60557.
	GraphQLServerAddr string `envvar:"GRAPHQL_SERVER_ADDR" default:"0.0.0.0:60557"`
	// EnableGraphQLPlayground determines whether or not to enable GraphiQL, an interactive
	// GraphQL playground which can be accessed by visiting GraphQLServerAddr in a browser.
	// See https://github.com/graphql/graphiql for more information. By default, GraphiQL
	// is disabled.
	EnableGraphQLPlayground bool `envvar:"ENABLE_GRAPHQL_PLAYGROUND" default:"false"`
}
```
