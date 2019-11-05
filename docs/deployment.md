[![Version](https://img.shields.io/badge/version-5.1.0--beta-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

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

## Running Mesh

If you would like to participate in the Mesh Beta, check out [this guide](deployment_with_telemetry.md) to deploying a telemetry-enabled Mesh node.

Make sure you have Docker installed. Then run:

```bash
docker run \
-it \
--rm \
-p 60557:60557 \
-p 60558:60558 \
-p 60559:60559 \
-e ETHEREUM_CHAIN_ID="1" \
-e ETHEREUM_RPC_URL="{your_ethereum_rpc_url}" \
-e VERBOSITY=5 \
-v {local_path_on_host_machine}/0x_mesh:/usr/mesh/0x_mesh
0xorg/mesh:latest
```

You should replace `{your_ethereum_rpc_url}` with the RPC endpoint for an
Ethereum node and `{local_path_on_host_machine}` with a directory on your host
machine where all Mesh-related data will be stored.

**Notes:**

-   Ports 60557, 60558, and 60559 are the default ports used for the JSON RPC endpoint, communicating with peers over TCP, and communicating with peers over WebSockets, respectively.
-   In order to disable P2P order discovery and sharing, set `USE_BOOTSTRAP_LIST` to `false`.
-   Running a VPN may interfere with Mesh. If you are having difficulty connecting to peers, disable your VPN.
-   If you are running against a POA testnet (e.g., Kovan), you might want to shorten the `BLOCK_POLLING_INTERVAL` since blocks are mined more frequently then on mainnet.
-   If you want to run the mesh in "detached" mode, add the `-d` switch to the docker run command so that your console doesn't get blocked.

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
	// P2PListenPort is the port on which to listen for new peer connections.
	P2PListenPort int `envvar:"P2P_LISTEN_PORT"`
	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
	// API.
	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL"`
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
}
```

There is one additional environment variable in the [main entrypoint for the
Mesh executable](../cmd/mesh/main.go):

```go
type standaloneConfig struct {
	// RPCAddr is the interface and port to use for the JSON-RPC API over
	// WebSockets. By default, 0x Mesh will listen on localhost and will let the
	// OS select a randomly available port.
	RPCAddr string `envvar:"RPC_ADDR" default:"localhost:0"`
}
```
