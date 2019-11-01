import { Mesh, OrderEvent, BigNumber, Verbosity } from '@0x/mesh-browser';
import { Web3ProviderEngine, RPCSubprovider } from '@0x/subproviders';
import { signatureUtils, Order, orderHashUtils } from '@0x/order-utils';

const ethereumRPCURL = 'http://localhost:8545';

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
        makerAssetAmount: new BigNumber('100000000000000000000') as any,
        makerFee: new BigNumber('0') as any,
        takerAddress: '0x0000000000000000000000000000000000000000',
        takerAssetData: '0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082',
        takerFeeAssetData: '0x',
        takerAssetAmount: new BigNumber('50000000000000000000') as any,
        takerFee: new BigNumber('0') as any,
        senderAddress: '0x0000000000000000000000000000000000000000',
        exchangeAddress: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
        feeRecipientAddress: '0xa258b39954cef5cb142fd567a46cddb31a670124',
        expirationTimeSeconds: new BigNumber(expirationTime) as any,
        salt: new BigNumber('1548619145450') as any,
        chainId: 50,
    };
    const signedOrder = await signatureUtils.ecSignOrderAsync(provider, order, order.makerAddress);
    const orderHash = orderHashUtils.getOrderHashHex(order);
    console.log(
        JSON.stringify({
            message: 'signed order in browser',
            orderHash: orderHash,
        }),
    );

    // Configure Mesh to use our local Ganache instance and local bootstrap
    // node.
    const mesh = new Mesh({
        verbosity: Verbosity.Debug,
        ethereumRPCURL,
        ethereumChainID: 1337,
        bootstrapList: ['/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7'],
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    // This handler will be called whenever an order is added, expired,
    // canceled, or filled. We will check for certain events to be logged in the
    // integration tests.
    mesh.onOrderEvents((events: Array<OrderEvent>) => {
        for (let event of events) {
            console.log(JSON.stringify(event));
        }
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

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

    // This special #jsFinished div is used to signal the headless Chrome driver
    // that the JavaScript code is done running.
    const finishedDiv = document.createElement('div');
    finishedDiv.setAttribute('id', 'jsFinished');
    document.querySelector('body')!.appendChild(finishedDiv);
})().catch(err => {
    if (err instanceof Error) {
        console.error(err.name + ': ' + err.message);
    } else {
        console.error(err.toString());
    }
});
