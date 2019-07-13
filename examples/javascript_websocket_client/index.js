const Web3Providers = require('web3-providers');

const MESH_WS_PORT = 60557;
const MESH_WS_ENDPOINT = `ws://localhost:${MESH_WS_PORT}`;
console.log('Mesh WebSocket endpoint: ', MESH_WS_ENDPOINT);

// Instantiate the WebSocket provider/client
const websocketProvider = new Web3Providers.WebsocketProvider(MESH_WS_ENDPOINT, {
    clientConfig: {
        // For some reason fragmenting the payloads causes the connection to close
        // Source: https://github.com/theturtle32/WebSocket-Node/issues/359
        fragmentOutgoingMessages: false,
    },
});

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
            makerAssetAmount: '1000000000000000000',
            takerAssetAmount: '10000000000000000000000',
            expirationTimeSeconds: '1586340602',
            makerFee: '0',
            takerFee: '0',
            feeRecipientAddress: '0x0000000000000000000000000000000000000000',
            senderAddress: '0x0000000000000000000000000000000000000000',
            salt: '41253767178111694375645046549067933145709740457131351457334397888365956743955',
            makerAssetData: '0xf47261b0000000000000000000000000e41d2489571d322189246dafa5ebde1f4699f498',
            takerAssetData: '0xf47261b0000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2',
            exchangeAddress: '0x080bf510fcbf18b91105470639e9561022937712',
            signature:
                '0x1c0827552a3bde2c72560362950a69f581ae7a1e6fa8c160bb437f3a61002bb96c22b646edd3b103b976db4aa4840a11c13306b2a02a0bb6ce647806c858c238ec02',
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
