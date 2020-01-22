import { Mesh, OrderEvent, SignedOrder, BigNumber, Provider } from '@0x/mesh-browser';

(async () => {
    // Configure Mesh to use mainnet and Infura.
    const mesh = new Mesh({
        verbosity: 4,
        ethereumChainID: 1,
        web3Provider: (window as any).web3.currentProvider as Provider,
    });

    // This handler will be called whenver there is a critical error.
    mesh.onError((err: Error) => {
        console.error(err);
    });

    // This handler will be called whenever an order is added, expired,
    // canceled, or filled.
    mesh.onOrderEvents((events: Array<OrderEvent>) => {
        for (let event of events) {
            console.log(event);
        }
    });

    // Start Mesh *after* we set up the handlers.
    await mesh.startAsync();

    // This order is for demonstration purposes only and is invalid. It will be
    // rejected by Mesh. You can replace it with a valid order.
    const order: SignedOrder = {
        chainId: 1,
        makerAddress: '0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149',
        makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
        makerAssetAmount: new BigNumber('1000000000000000000'),
        makerFee: new BigNumber('0'),
        makerFeeAssetData: '0x',
        takerAddress: '0x0000000000000000000000000000000000000000',
        takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
        takerAssetAmount: new BigNumber('10000000000000000000000'),
        takerFee: new BigNumber('0'),
        takerFeeAssetData: '0x',
        senderAddress: '0x0000000000000000000000000000000000000000',
        exchangeAddress: '0x080bf510FCbF18b91105470639e9561022937712',
        feeRecipientAddress: '0x0000000000000000000000000000000000000000',
        expirationTimeSeconds: new BigNumber('1586340602'),
        salt: new BigNumber('41253767178111694375645046549067933145709740457131351457334397888365956743955'),
        signature:
            '0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec03',
    };

    // Add the order and log the result.
    const result = await mesh.addOrdersAsync([order]);
    console.log(result);
})().catch(err => {
    console.error(err);
});
