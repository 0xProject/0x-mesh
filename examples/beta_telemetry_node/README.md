## Deploying a Telemetry-Enabled Mesh Node

0x Mesh is completely permissionless and the beta is open to anyone who wants to
participate. You can optionally help improve 0x Mesh by enabling telemetry. Mesh
automatically logs a lot of useful information including the number of orders
processed and details about any errors and warnings that might occur. Sending
this information to us is extraordinarily helpful, but completely optional. If
you don't want to enable telemetry, you can follow the
[Deployment Guide](../../DEPLOYMENT.md) instead.

This guide will walk you though setting up a telemetry-enabled Mesh node on the cloud hosting solution of your choice using [Docker Machine](https://docs.docker.com/machine/). The instructions below will deploy a Mesh node on [DigitalOcean](https://www.digitalocean.com/), but can be easily modified to deploy on [many other cloud providers](https://docs.docker.com/machine/drivers/).

## Prerequisites

-   [Docker](https://www.docker.com/get-started)
-   [Docker Machine](https://docs.docker.com/machine/install-machine/) On Mac/Windows it is installed with Docker Desktop.
-   [Docker Compose](https://docs.docker.com/compose/install/) On Mac, it is installed with Docker Desktop.

## Instructions

Let's start by cloning the Mesh repo and navigating to the `beta_telemetry_node` directory:

```bash
git clone https://github.com/0xProject/0x-mesh.git
cd 0x-mesh/examples/beta_telemetry_node
code . # open this directory in your text editor
```

Please open up `docker-compose.yml` and set `ETHEREUM_RPC_URL` to your own Ethereum JSON RPC endpoint.

We are now ready to deploy our instance. Before we can continue, you will need to set up an account with the cloud hosting provider of your choice, and retrieve your access token/key/secret. We will use them to create a new machine with name `mesh-beta`. Docker has great documentation on doing all of that for [DigitalOcean](https://docs.docker.com/machine/examples/ocean/) and [AWS](https://docs.docker.com/machine/examples/aws/). Instead of naming the machine `docker-sandbox` as in those examples, let's name ours `mesh-beta` as shown below.

```bash
docker-machine create --driver digitalocean --digitalocean-access-token xxxxx mesh-beta
```

Make sure you replaced `xxxxx` with your access token. This command will spin up a new instance on your cloud provider, pre-installed with Docker so that it's ready-to-use with the `docker-machine` command. Once the command completes, let's make sure the machine exists:

```bash
docker-machine ls
```

You should see something like:

```
mesh-beta     -        digitalocean   Running   tcp://162.31.121.332:2376            v18.09.7
```

Our remote machine is alive! The next step is to copy over two config files to this remote machine. We can use [docker-machine scp](https://docs.docker.com/machine/reference/scp/) to do this. The following commands ask `docker-machine` to copy over the `fluent-bit.conf` and `parsers.conf` files from our local computer to a directory called `root` on the `mesh-beta` machine.

```bash
docker-machine scp fluent-bit.conf mesh-beta:/root/fluent-bit.conf
docker-machine scp parsers.conf mesh-beta:/root/parsers.conf
```

Now comes the Docker Machine magic. By running the following commands, we can ask Docker Machine to let us execute any Docker command in our local shell AS IF we were executing them directly on the `mesh-beta` machine:

```bash
docker-machine env mesh-beta
eval $(docker-machine env mesh-beta)
```

Presto! We are now ready to spin up our telemetry-enabled Mesh node! We do this using the Docker Compose command `up`:

```bash
docker-compose up -d
```

Houston, we have lift-off! All the logs from the Mesh node are being piped to the [Fluentbit](https://fluentbit.io/) instance that got deployed alongside Mesh. So if you want to inspect the logs, you need to do:

```bash
docker logs <fluent-bit-container-id> -f
```

Instead of reading them from the `0xorg/mesh` container.

Finally, in order to prevent our log aggregation stack from getting overloaded,
we whitelist the peers that are allowed to send us logs. Look for a log message
that looks like this:

```json
{
    "addresses": ["/ip4/127.0.0.1/tcp/60557", "/ip4/172.17.0.2/tcp/60557"],
    "level": "info",
    "msg": "started p2p node",
    "peerID": "QmbKkHnmkmFxKbPWbBNz3inKizDuqjTsWsVyutnshYULLp",
    "time": "2019-07-15T17:36:46-07:00"
}
```

Ping us in [Discord](https://discord.gg/HF7fHwk) and let us know your peer ID. You can DM `fabio#1058`, `Alex Browne | 0x#2975` or `ovrmrrw#0454` and we'll whitelist your node :)

I hope that was easy enough! If you ran into any issues, please ping us in the #mesh channel on [Discord](https://discord.gg/HF7fHwk). To learn more about connecting to your Mesh node's JSON RPC interface, check out our [Usage docs](../../USAGE.md). Your node's JSON RPC endpoint should be available at `ws://<your-ip-address>:60557` and you can discover your machine's IP address by running:

```
docker-machine ip mesh-beta
```

Cheers!
