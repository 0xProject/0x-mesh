import { Mesh, OrderEvent, Verbosity, Config } from '@0x/mesh-browser';
import { Order, orderHashUtils, signatureUtils } from '@0x/order-utils';
import { RPCSubprovider, Web3ProviderEngine } from '@0x/subproviders';
import { BigNumber } from '@0x/utils';

const ethereumRPCURL = 'http://localhost:8545';

// tslint:disable:no-console
(async () => {
    (window as any).startMesh = startMesh;
    const finishedDiv = document.createElement('div');
    finishedDiv.setAttribute('id', 'jsFinishedLoading');
    document.body.appendChild(finishedDiv);
})().catch(err => {
    if (err instanceof Error) {
        console.error(`${err.name}: ${err.message}`);
    } else {
        console.error(err.toString());
    }
});

function importConfig(config: any): Config {
    config.customOrderFilter = config.customOrderFilter === '' ? null : JSON.parse(config.customOrderFilter);
    config.customContractAddresses =
        config.customContractAddresses === '' ? null : JSON.parse(config.customContractAddresses);
    config.bootstrapList = config.bootstrapList === '' ? [] : config.bootstrapList.split(',');

    return config;
}

// NOTE: meshConfig is passed in a marshalled core.Config JSON string.
async function startMesh(meshConfig: any) {
    const config = importConfig(meshConfig);
    // Set up a Web3 Provider that uses the RPC endpoint
    const provider = new Web3ProviderEngine();
    if (!config.ethereumRPCURL) {
        throw new Error('missing ethereumRPCUrl');
    }
    provider.addProvider(new RPCSubprovider(config.ethereumRPCURL));
    provider.start();

    const mesh = new Mesh({
        ...config,
        web3Provider: provider,
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    mesh.onOrderEvents((events: OrderEvent[]) => {
        (async () => {
            for (const event of events) {
                console.log(event);
            }
        })().catch(err => console.error(err));
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // HACK(albrow): Wait for GossipSub to initialize. We could remove this if we adjust
    // how we are waiting for the order (what log message we look for). As the test is
    // currently written it only passes when the order is received through GossipSub and
    // fails if it was received through ordersync.
    const fiveSeconds = 5000;
    await sleepAsync(fiveSeconds);

    // Call getStatsAsync and make sure it works.
    const stats = await mesh.getStatsAsync();
    console.log(JSON.stringify(stats));
    const finishedDiv = document.createElement('div');
    finishedDiv.setAttribute('id', 'jsFinished');
    document.body.appendChild(finishedDiv);
    return JSON.stringify(stats);
}

async function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
// tslint:enable:no-console
