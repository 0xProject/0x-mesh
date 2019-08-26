import { Mesh, OrderEvent, SignedOrder, BigNumber } from '@0x/mesh-browser';

(async () => {
    // Configure Mesh to use our local Ganache instance and network ID 51
    const mesh = new Mesh({
        verbosity: 5,
        ethereumRPCURL: 'http://localhost:8545',
        ethereumNetworkID: 50,
        bootstrapList: ['/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7'],
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    // This handler will be called whenever an order is added, expired,
    // canceled, or filled.
    mesh.onOrderEvents((events: Array<OrderEvent>) => {
        for (let event of events) {
            // console.log('received order event: ' + JSON.stringify(event));
        }
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // This order is for demonstration purposes only and is invalid. It will be
    // rejected by Mesh. You can replace it with a valid order.
    const order: SignedOrder = {
        makerAddress: '0x5409ed021d9299bf6814279a6a1411a7e866a631',
        makerAssetData: '0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c',
        makerAssetAmount: new BigNumber('1000'),
        makerFee: new BigNumber('0'),
        takerAddress: '0x0000000000000000000000000000000000000000',
        takerAssetData: '0xf47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082',
        takerAssetAmount: new BigNumber('2000'),
        takerFee: new BigNumber('0'),
        senderAddress: '0x0000000000000000000000000000000000000000',
        exchangeAddress: '0x48bacb9266a570d521063ef5dd96e61686dbe788',
        feeRecipientAddress: '0xa258b39954cef5cb142fd567a46cddb31a670124',
        expirationTimeSeconds: new BigNumber('1567120395'),
        salt: new BigNumber('1548619145450'),
        signature:
            '0x1be53f1a3cc8508b51995f03e5c22948c4988d113c2085d10135d9bbc20dd272275d31c59ec09790ea4691b35aec4a0ed0558b18add3e778dda20eb05d8a07097303',
    };

    // Add the order and log the result.
    const result = await mesh.addOrdersAsync([order]);
    if (result.accepted.length !== 1) {
        throw new Error('Expected exactly one order to be accepted but got: ' + result.accepted.length);
    }
    if (result.rejected.length !== 0) {
        throw new Error('Expected no orders to be rejected but got: ' + result.rejected.length);
    }

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
