import { Config, Mesh, OrderEvent } from '@0x/mesh-browser';
import { RPCSubprovider, Web3ProviderEngine } from '@0x/subproviders';

function importConfig(config: any): Config {
    config.customOrderFilter = config.customOrderFilter === '' ? null : JSON.parse(config.customOrderFilter);
    config.customContractAddresses =
        config.customContractAddresses === '' ? null : JSON.parse(config.customContractAddresses);
    config.bootstrapList = config.bootstrapList === '' ? [] : config.bootstrapList.split(',');

    return config;
}

// addDivMark adds a specific div to be used by as a marker for chrome dev
// protocol usage
function addDivMark(idAttribute: string): void {
    const newDiv = document.createElement('div');
    newDiv.setAttribute('id', idAttribute);
    document.body.appendChild(newDiv);
}

// tslint:disable:no-console
// NOTE: meshConfig is passed in a marshalled core.Config JSON string.
async function startMeshAsync(meshConfig: any): Promise<string> {
    const config = importConfig(meshConfig);
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

    mesh.onError((err: Error) => {
        console.error(err);
    });

    mesh.onOrderEvents((events: OrderEvent[]) => {
        for (const event of events) {
            console.log(event);
        }
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // Call getStatsAsync and make sure it works.
    const stats = await mesh.getStatsAsync();
    console.log(JSON.stringify(stats));
    addDivMark('jsFinished');
    return JSON.stringify(stats);
}

(async () => {
    (window as any).startMesh = startMeshAsync;
    addDivMark('jsFinishedLoading');
})().catch(err => {
    if (err instanceof Error) {
        console.error(`${err.name}: ${err.message}`);
    } else {
        console.error(err.toString());
    }
});
// tslint:enable:no-console
