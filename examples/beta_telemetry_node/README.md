## Setting up a Telemetry-Enabled Mesh Deployment for the Beta Period

This guide will walk you though setting up a Mesh node on the cloud hosting solution of your choice using [Docker Machine](https://docs.docker.com/machine/). The instructions below will deploy to Mesh node on [DigitalOcean](https://www.digitalocean.com/), but can be easily modified to deploy on [many other cloud providers](https://docs.docker.com/machine/drivers/).

## Prerequisites

-   [Docker](https://www.docker.com/get-started)
-   [Docker Machine](https://docs.docker.com/machine/install-machine/) On Mac/Windows it is installed with Docker Desktop.
-   [Docker Compose](https://docs.docker.com/compose/install/) On Mac, it is installed with Docker Desktop.

## Instructions

Let's start by cloning the Mesh repo and navigating to the `beta_telemetry_node` directory:

```
git clone https://github.com/0xProject/0x-mesh.git
cd 0x-mesh/examples/beta_telemetry_node
code . # open this directory in your text editor
```

Next, let's customize this setup for our deployment. Let's open up `fluent-bit.conf` and modify the `Logstash_Prefix` config so that it has a name that uniquely identifies logs originating from your Mesh node as belonging to you or your organization. Let's edit it from `mesh_beta` to `mesh_{NAME}` where `{name}` is your personal/organization name.

Now, let's open `data/keys/privkey` and paste in the private key given to you by the 0x developers. This private key will be used to generate both your `peerID` on the network, and allow your node's logs to be sent to 0x's log aggregation service. Do ahead and replace `<PASTE_YOUR_PRIVATE_KEY_HERE>` with your private key.

Please open up `docker-compose.yml` and set `ETHEREUM_RPC_URL` to your own Ethereum JSON RPC endpoint (if everyone uses our Infura API key, it'll get rate-limited!).

We are now ready to deploy our instance. Before we can continue, you will need to set up an account with the cloud hosting provider of your choice, and retrieve your access token/key/secret. We will use them to create a new machine with name `mesh-beta`. Docker has great documentation on doing all of that for [DigitalOcean](https://docs.docker.com/machine/examples/ocean/) and [AWS](https://docs.docker.com/machine/examples/aws/). Instead of naming the machine `docker-sandbox` as in those examples, let's name ours `mesh-beta` as shown below.

```bash
docker-machine create --driver digitalocean --digitalocean-access-token xxxxx mesh-beta
```

Make sure you replaced `xxxxx` with your access token. This command will spin up a new instance on your cloud provider, pre-installed with Docker so that it's ready-to-use with the `docker-machine` command. Once the command completes, let's make sure the machine exists:

```
docker-machine ls
```

You should see something like:

```
mesh-beta     -        digitalocean   Running   tcp://162.31.121.332:2376            v18.09.7
```

Our remote machine is alive! The next step is to copy over our config files and private key to this remote machine. We can use [docker-machine scp](https://docs.docker.com/machine/reference/scp/). The following commands ask `docker-machine` to copy over the `data` dir and `fluent-bit.conf` and `parsers.conf` files from our local computer to a directory called `root` on the `mesh-beta` machine.

```
docker-machine scp -r -d data/ mesh-beta:/root/data/
docker-machine scp fluent-bit.conf mesh-beta:/root/fluent-bit.conf
docker-machine scp parsers.conf mesh-beta:/root/parsers.conf
```

Now comes the Docker Machine magic. By running the following commands, we can ask Docker Machine to let us execute any Docker command in our local shell AS IF we were executing them on the `mesh-beta` machine:

```
docker-machine env mesh-beta
eval $(docker-machine env mesh-beta)
```

Presto! We are now ready to spin up our telemetry-enabled Mesh node! We do this using Docker Compose:

```
docker-compose up -d
```

Houston, we have lift-off! All the logs from the Mesh node are being piped to the [Fluentbit](https://fluentbit.io/) instance that got deployed alongside Mesh. So if you want to inspect the logs, you need to do:

```
docker logs <fluent-bit-container-id> -f
```

Instead of reading them from the `0xorg/mesh` container.

I hope that was easy enough! If you ran into any issues, please ping us in the #mesh channel on [Discord](https://discord.gg/HF7fHwk). To learn more about connecting to your Mesh node's JSON RPC interface, check out our [Usage docs](USAGE.md).
