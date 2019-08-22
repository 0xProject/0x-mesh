import { Mesh, OrderEvent, SignedOrder } from '@0x/mesh-browser';

(async () => {
    const mesh = new Mesh({
        ethereumRPCURL: 'https://mainnet.infura.io/v3/af2e590be00f463fbfd0b546784065ad',
        ethereumNetworkID: 1,
    });
    mesh.setOrderEventsHandler((events: Array<OrderEvent>) => {
        for (let event of events) {
            console.log(event);
        }
    });
    await mesh.startAsync();

    // This order is for demonstration purposes only and is invalid. It will be
    // rejected by Mesh. You can replace it with a valid order.
    const order: SignedOrder = {
        makerAddress: '0xa3eCE5D5B6319Fa785EfC10D3112769a46C6E149',
        makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
        makerAssetAmount: '1000000000000000000',
        makerFee: '0',
        takerAddress: '0x0000000000000000000000000000000000000000',
        takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
        takerAssetAmount: '10000000000000000000000',
        takerFee: '0',
        senderAddress: '0x0000000000000000000000000000000000000000',
        exchangeAddress: '0x080bf510FCbF18b91105470639e9561022937712',
        feeRecipientAddress: '0x0000000000000000000000000000000000000000',
        expirationTimeSeconds: '1586340602',
        salt: '41253767178111694375645046549067933145709740457131351457334397888365956743955',
        signature:
            '0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec03',
    };
    const result = await mesh.addOrdersAsync([order]);
    console.log(result);
})().catch(err => {
    console.error(err);
});
