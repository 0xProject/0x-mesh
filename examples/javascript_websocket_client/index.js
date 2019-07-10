const Web3Providers = require('web3-providers');

const MESH_WS_PORT = 60557;
const MESH_WS_ENDPOINT = `ws://localhost:${MESH_WS_PORT}`;
console.log('Mesh WebSocket endpoint: ', MESH_WS_ENDPOINT);

// Instantiate the WebSocket provider/client
const websocketProvider = new Web3Providers.WebsocketProvider(MESH_WS_ENDPOINT);

// Listen for the close event which will fire if Mesh goes down
websocketProvider.connection.addEventListener('close', function() {
    console.log('close event received');
    process.exit(1);
});

console.log('Subscribing to heartbeat...');
websocketProvider.subscribe('mesh_subscribe', 'heartbeat', []).then(function(subscriptionId) {
    console.log('Heartbeat subscriptionId', subscriptionId);

    // Listen to event on the subscription (topic is the subscriptionId)
    websocketProvider.on(subscriptionId, function(eventPayload) {
        console.log('Received:', JSON.stringify(eventPayload, null, '\t'));
    });
});

console.log('About to subscribe to order events...');
websocketProvider
    .subscribe('mesh_subscribe', 'orders', [])
    .then(function(subscriptionId) {
        console.log('Order events subscriptionId', subscriptionId);
        // Listen to event on the subscription (topic is the subscriptionId)
        websocketProvider.on(subscriptionId, function(eventPayload) {
            console.log('Received:', JSON.stringify(eventPayload, null, '\t'));
        });

        // Submit an order to the Mesh node
        console.log('About to send order...');
        var order = {
            makerAddress: '0xa3ece5d5b6319fa785efc10d3112769a46c6e149',
            takerAddress: '0x0000000000000000000000000000000000000000',
            makerAssetAmount: '100000000000000000000',
            takerAssetAmount: '100000000000000000000000',
            expirationTimeSeconds: '1559856615025',
            makerFee: '0',
            takerFee: '0',
            feeRecipientAddress: '0x0000000000000000000000000000000000000000',
            senderAddress: '0x0000000000000000000000000000000000000000',
            salt: '46108882540880341679561755865076495033942060608820537332859096815711589201849',
            makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
            takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
            exchangeAddress: '0x4f833a24e1f95d70f028921e27040ca56e09ab0b',
            signature:
                '0x1c52f75daa4bd2ad9e6e8a7c35adbd089d709e48ae86463f2abfafa3578747fafc264a04d02fa26227e90476d57bca94e24af32f1cc8da444bba21092ca56cd85603',
        };
        const orders = [order];
        websocketProvider.send('mesh_addOrders', [orders]).then(function(validationResults) {
            console.log('mesh_addOrders Response: ', JSON.stringify(validationResults, null, '\t'));
        });
    })
    .catch(function(err) {
        console.log('Error:', err);
    });

(async () => {
    // Get all orders stored in Mesh at a snapshot in time
    const perPage = 200;
    let snapshotID = ''; // New snapshot
    let ordersInfosLen = 1;

    let i = 0;
    const allOrdersInfos = [];
    while (ordersInfosLen !== 0) {
        const page = i;
        const getOrdersResponse = await websocketProvider.send('mesh_getOrders', [page, perPage, snapshotID]);
        console.log('mesh_getOrders Response:', JSON.stringify(getOrdersResponse, null, '\t'));
        snapshotID = getOrdersResponse.snapshotID;
        allOrdersInfos.push(...getOrdersResponse.ordersInfos);
        ordersInfosLen = getOrdersResponse.ordersInfos.length;
        i++;
    }
    console.log('Got ', allOrdersInfos.length, 'orders from snapshot ', snapshotID);
})();
