[![Version](https://img.shields.io/badge/version-1.0.4--beta-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

# 0x Mesh Deployment Guide

Welcome to the [0x Mesh](https://github.com/0xProject/0x-mesh) Deployment Guide!
This guide will walk you through how to configure and deploy your own 0x Mesh
node.

## Supported Networks

-   Mainnet
-   Kovan
-   Ropsten
-   Rinkeby
-   [Ganache snapshot](https://cloud.docker.com/u/0xorg/repository/docker/0xorg/mesh-ganache-cli)

## Running Mesh

Make sure you have Docker installed. Then run:

```bash
docker run \
-it \
--rm \
-p 60557:60557 \
-p 60558:60558 \
-e ETHEREUM_NETWORK_ID="1" \
-e ETHEREUM_RPC_URL="https://mainnet.infura.io/v3/a9a23d2566e542629179d6372ace13c9" \
-e VERBOSITY=5 \
0xorg/mesh:latest
```

**Notes:**

-   `60557` is the `RPC_PORT` and `60558` is the `P2P_LISTEN_PORT`
-   In order to disable P2P order discovery and sharing, set `USE_BOOTSTRAP_LIST` to `false`.
-   Running a VPN may interfere with Mesh. If you are having difficulty connecting to peers, disable your VPN.
-   If you are running against a POA testnet (e.g., Kovan), you might want to shorten the `BLOCK_POLLING_INTERVAL` since blocks are mined more frequently then on mainnet.

## Persisting State

By default in the Docker container, Mesh stores all state (e.g. database files,
private key file) in `/usr/mesh/0x_mesh`. If you want the Mesh state to persist
across Docker container re-starts, mount a local `0x_mesh` directory to your
container. Add the following to the `docker run` command above:

```
-v {abs_local_path}/0x_mesh:/usr/mesh/0x_mesh
```

**Note:** Replace `{abs_local_path}` with the absolute path to the desired `0x_mesh` directory on the host machine.

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
    // EthereumNetworkID is the network ID to use when communicating with
    // Ethereum.
    EthereumNetworkID int `envvar:"ETHEREUM_NETWORK_ID"`
    // UseBootstrapList is whether to use the predetermined list of peers to
    // bootstrap the DHT and peer discovery.
    UseBootstrapList bool `envvar:"USE_BOOTSTRAP_LIST" default:"true"`
    // OrderExpirationBuffer is the amount of time before the order's stipulated expiration time
    // that you'd want it pruned from the Mesh node.
    OrderExpirationBuffer time.Duration `envvar:"ORDER_EXPIRATION_BUFFER" default:"10s"`
    // BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
    // that might contain transactions that impact the fillability of orders stored by Mesh. Different
    // networks have different block producing intervals: POW networks are typically slower (e.g., Mainnet)
    // and POA networks faster (e.g., Kovan) so one should adjust the polling interval accordingly.
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
Mesh executable](cmd/mesh/main.go):

```go
type standaloneConfig struct {
    // RPCPort is the port to use for the JSON RPC API over WebSockets. By
    // default, 0x Mesh will let the OS select a randomly available port.
    RPCPort int `envvar:"RPC_PORT" default:"0"`
}
```
