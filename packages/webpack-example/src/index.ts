import { BigNumber, Mesh, OrderEvent, SignedOrder, SupportedProvider } from '@0x/mesh-browser';

// tslint:disable:no-console
(async () => {
    // Configure Mesh to use web3.currentProvider (e.g. provided by MetaMask).
    const mesh = new Mesh({
        verbosity: 4,
        ethereumChainID: 1,
        web3Provider: (window as any).web3.currentProvider as SupportedProvider,
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    // This handler will be called whenever an order is added, expired,
    // canceled, or filled.
    mesh.onOrderEvents((events: OrderEvent[]) => {
        for (const event of events) {
            console.log(event);
        }
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // This order is for demonstration purposes only and will likely be expired
    // by the time you run this example. If so, it will be rejected by Mesh. You
    // can replace it with a valid order.
    const order: SignedOrder = {
        signature:
            '0x1c68eb1e2577e9f51776bdb06ec51fcec9aec0ea1565eca5e243917cecaafaa46b3b9590ff6575bf1c048d0b4ec5773a2e3a8df3bf117e1613e2a7b57d6f95c95a02',
        senderAddress: '0x0000000000000000000000000000000000000000',
        makerAddress: '0x4418755f710468e223797a006603e29937e825bc',
        takerAddress: '0x0000000000000000000000000000000000000000',
        makerFee: new BigNumber('0'),
        takerFee: new BigNumber('0'),
        makerAssetAmount: new BigNumber('3000000000'),
        takerAssetAmount: new BigNumber('19500000000000000000'),
        makerAssetData: '0xf47261b0000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48',
        takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
        salt: new BigNumber('1579725034907'),
        exchangeAddress: '0x61935cbdd02287b511119ddb11aeb42f1593b7ef',
        feeRecipientAddress: '0xa258b39954cef5cb142fd567a46cddb31a670124',
        expirationTimeSeconds: new BigNumber('1580329834'),
        makerFeeAssetData: '0x',
        chainId: 1,
        takerFeeAssetData: '0x',
    };

    // Add the order and log the result.
    const result = await mesh.addOrdersAsync([order]);
    console.log(result);
})().catch(err => {
    console.error(err);
});
// tslint:enable:no-console
