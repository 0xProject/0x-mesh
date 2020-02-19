import { Mesh, OrderEvent, Verbosity } from '@0x/mesh-browser';
import { BigNumber } from '@0x/utils';
import { Web3ProviderEngine, RPCSubprovider } from '@0x/subproviders';
import { signatureUtils, Order, orderHashUtils } from '@0x/order-utils';

const ethereumRPCURL = 'http://localhost:8545';

// Set up a Web3 Provider that uses the RPC endpoint
const provider = new Web3ProviderEngine();
provider.addProvider(new RPCSubprovider(ethereumRPCURL));
provider.start();

(async () => {
    // Sign an order and log the order hash so that we can use it in the
    // integration tests.
    console.log('signing order...');
    const currentTime = Math.floor(Date.now() / 1000);
    const expirationTime = currentTime + 24 * 60 * 60;
    const order: Order = {
        makerAddress: '0x6ecbe1db9ef729cbe972c83fb886247691fb6beb',
        makerAssetData: '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
        makerFeeAssetData: '0x',
        makerAssetAmount: new BigNumber('100000000000000000000'),
        makerFee: new BigNumber('0'),
        takerAddress: '0x0000000000000000000000000000000000000000',
        takerAssetData: '0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082',
        takerFeeAssetData: '0x',
        takerAssetAmount: new BigNumber('50000000000000000000'),
        takerFee: new BigNumber('0'),
        senderAddress: '0x0000000000000000000000000000000000000000',
        exchangeAddress: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
        feeRecipientAddress: '0xa258b39954cef5cb142fd567a46cddb31a670124',
        expirationTimeSeconds: new BigNumber(expirationTime),
        salt: new BigNumber('1548619145450'),
        chainId: 1337,
    };
    const signedOrder = await signatureUtils.ecSignOrderAsync(provider, order, order.makerAddress);
    const orderHash = await orderHashUtils.getOrderHashAsync(order);
    console.log(
        JSON.stringify({
            message: 'signed order in browser',
            orderHash: orderHash,
        }),
    );

    // Configure Mesh to use our local Ganache instance and local bootstrap
    // node.
    const mesh = new Mesh({
        verbosity: Verbosity.Trace,
        ethereumChainID: 1337,
        bootstrapList: ['/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7'],
        customOrderFilter: {
            properties: { makerAddress: { const: '0x6ecbe1db9ef729cbe972c83fb886247691fb6beb' } },
        },
        web3Provider: provider,
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    // This handler will be called whenever an order is added, expired,
    // canceled, or filled. We will check for certain events to be logged in the
    // integration tests.
    mesh.onOrderEvents((events: Array<OrderEvent>) => {
        (async () => {
            for (let event of events) {
                // Check the happy path for getOrdersForPageAsync. There should
                // be two orders. (just make sure it doesn't throw/reject).
                const firstOrdersResponse = await mesh.getOrdersForPageAsync(0, 1, '');
                console.log(JSON.stringify(firstOrdersResponse));
                const secondOrdersResponse = await mesh.getOrdersForPageAsync(1, 1, firstOrdersResponse.snapshotID);
                console.log(JSON.stringify(secondOrdersResponse));

                // Check the happy path for getOrders (just make sure it
                // doesn't throw/reject).
                await mesh.getOrdersAsync();

                // Log the event. The Go code will be watching the logs for
                // this.
                console.log(JSON.stringify(event));
            }
        })().catch(err => console.error(err));
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // HACK(albrow): Wait for GossipSub to initialize. We could remove this if we adjust
    // how we are waiting for the order (what log message we look for). As the test is
    // currently written it only passes when the order is received through GossipSub and
    // fails if it was received through ordersync.
    await sleepAsync(5000);

    // Send an order to the network. In the integration tests we will check that
    // the order was received.
    const result = await mesh.addOrdersAsync([signedOrder as any]);
    if (result.accepted.length !== 1) {
        console.log(JSON.stringify(result));
        throw new Error('Expected exactly one order to be accepted but got: ' + result.accepted.length);
    }
    if (result.rejected.length !== 0) {
        console.log(JSON.stringify(result));
        throw new Error('Expected no orders to be rejected but got: ' + result.rejected.length);
    }

    // Call getStatsAsync and make sure it works.
    const stats = await mesh.getStatsAsync();
    console.log(JSON.stringify(stats));

    // This special #jsFinished div is used to signal the headless Chrome driver
    // that the JavaScript code is done running. This is not a native Javascript
    // concept. Rather, it is our way of letting the Go program that serves this
    // Javascript know whether or not the test has completed.
    const finishedDiv = document.createElement('div');
    finishedDiv.setAttribute('id', 'jsFinished');
    document.body.appendChild(finishedDiv);
})().catch(err => {
    if (err instanceof Error) {
        console.error(err.name + ': ' + err.message);
    } else {
        console.error(err.toString());
    }
});

async function sleepAsync(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
}
