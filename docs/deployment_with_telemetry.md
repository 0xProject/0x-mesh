[![Version](https://img.shields.io/badge/version-development-orange.svg)](https://github.com/0xProject/0x-mesh/releases)

## Deploying a Telemetry-Enabled Mesh Node

0x Mesh is completely permissionless and the beta is open to anyone who wants to
participate. You can optionally help improve 0x Mesh by enabling telemetry. Mesh
automatically logs a lot of useful information including the number of orders
processed and details about any errors and warnings that might occur. Sending
this information to us is extraordinarily helpful, but completely optional. If
you don't want to enable telemetry, you can follow the
[Deployment Guide](deployment.md) instead.

This guide will walk you though setting up a telemetry-enabled Mesh node on the
cloud hosting solution of your choice using
[Docker Machine](https://docs.docker.com/machine/). The instructions below will
deploy a Mesh node on [DigitalOcean](https://www.digitalocean.com/), but can be
easily modified to deploy on
[many other cloud providers](https://docs.docker.com/machine/drivers/).

## Prerequisites

-   [Docker](https://www.docker.com/get-started)
-   [Docker Machine](https://docs.docker.com/machine/install-machine/) On Mac/Windows it is installed with Docker Desktop.
-   [Docker Compose](https://docs.docker.com/compose/install/) On Mac, it is installed with Docker Desktop.

## Instructions

### Setting up docker-compose.yml

Let's start by creating a new directory and adding a __docker-compose.yml__ file
with the following contents:

```
version: '3'

services:
    mesh:
        image: 0xorg/mesh:latest
        restart: always
        logging:
            driver: fluentd
            options:
                fluentd-address: localhost:24224
                tag: docker.mesh
        volumes:
            - /root/data:/usr/mesh/0x_mesh
        links:
            - fluentbit
        ports:
            - '60557:60557'
            - '60558:60558'
            - '60559:60559'
        environment:
            - VERBOSITY=5
            - ETHEREUM_CHAIN_ID=1
            - RPC_ADDR=mesh:60557
            # Set your backing Ethereum JSON RPC endpoint below
            - ETHEREUM_RPC_URL=
            - BLOCK_POLLING_INTERVAL=5s

    fluentbit:
        image: 0xorg/mesh-fluent-bit:latest
        links:
            - esproxy
        ports:
            - '24224:24224'
        command: /fluent-bit/bin/fluent-bit -c /fluent-bit/etc/fluent-bit.conf

    esproxy:
        image: overmorrow/auth-es-proxy:latest
        ports:
            - '3333:3333'
        volumes:
            - /root/data/keys:/app/keys
        restart: on-failure:5
        environment:
            - PORT=3333
            - REMOTE_ADDRESS=https://telemetry.mesh.0x.org/_bulk
            - INPUT_VALIDATION=false
            - OUTPUT_SIGNING=true
            - PRIVATE_KEY_PATH=/app/keys/privkey
```

In most cases, the only change you need to make to the __docker-compose.yml__
file is to set `ETHEREUM_RPC_URL` to your own Ethereum JSON RPC endpoint.

### Deploying with Docker Machine

We are now ready to deploy our instance. Before we can continue, you will need
to set up an account with the cloud hosting provider of your choice, and
retrieve your access token/key/secret. We will use them to create a new machine
with name `mesh-beta`. Docker has great documentation on doing all of that for
[DigitalOcean](https://docs.docker.com/machine/examples/ocean/) and
[AWS](https://docs.docker.com/machine/examples/aws/). Instead of naming the
machine `docker-sandbox` as in those examples, let's name ours `mesh-beta` as
shown below.

```bash
docker-machine create --driver digitalocean --digitalocean-access-token xxxxx mesh-beta
```

Make sure you replaced `xxxxx` with your access token. This command will spin up
a new instance on your cloud provider, pre-installed with Docker so that it's
ready-to-use with the `docker-machine` command. Once the command completes,
let's make sure the machine exists:

```bash
docker-machine ls
```

You should see something like:

```
mesh-beta     -        digitalocean   Running   tcp://162.31.121.332:2376            v18.09.7
```

Now comes the Docker Machine magic. By running the following commands, we can
ask Docker Machine to let us execute any Docker command in our local shell AS IF
we were executing them directly on the `mesh-beta` machine:

```bash
eval $(docker-machine env mesh-beta)
```

Presto! We are now ready to spin up our telemetry-enabled Mesh node! We do this
using the Docker Compose command `up`:

```bash
docker-compose up -d
```

Houston, we have lift-off! All the logs from the Mesh node are being piped to
the [Fluentbit](https://fluentbit.io/) instance that got deployed alongside
Mesh. So if you want to inspect the logs, you need to do:

```bash
docker logs <fluent-bit-container-id> -f
```

Instead of reading them from the `0xorg/mesh` container.

Finally, in order to prevent our log aggregation stack from getting overloaded,
we whitelist the peers that are allowed to send us logs. Look for a field in the
logs called `myPeerID`:

```json
{
    "myPeerID": "QmbKkHnmkmFxKbPWbBNz3inKizDuqjTsWsVyutnshYULLp",
}
```

Ping us in [Discord](https://discord.gg/HF7fHwk) and let us know your peer ID.
You can DM `fabio#1058`, `Alex Browne | 0x#2975` or `ovrmrrw#0454` and we'll
whitelist your node :)

I hope that was easy enough! If you ran into any issues, please ping us in the
`#mesh` channel on [Discord](https://discord.gg/HF7fHwk). To learn more about
connecting to your Mesh node's JSON RPC interface, check out the
[JSON-RPC API Documentation](rpc_api.md). Your node's JSON RPC endpoint
should be available at `ws://<your-ip-address>:60557` and you can discover your
machine's IP address by running:

```
docker-machine ip mesh-beta
```
