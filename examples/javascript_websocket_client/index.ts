import * as Web3Providers from 'web3-providers';

const MESH_WS_PORT = 60557;
const MESH_WS_ENDPOINT = `ws://localhost:${MESH_WS_PORT}`;

interface StringifiedSignedOrder {
    senderAddress: string;
    makerAddress: string;
    takerAddress: string;
    makerFee: string;
    takerFee: string;
    makerAssetAmount: string;
    takerAssetAmount: string;
    makerAssetData: string;
    takerAssetData: string;
    salt: string;
    exchangeAddress: string;
    feeRecipientAddress: string;
    expirationTimeSeconds: string;
    signature: string;
}

enum OrderEventKind {
    Invalid = 'INVALID',
    Added = 'ADDED',
    Filled = 'FILLED',
    Fully_filled = 'FULLY_FILLED',
    Cancelled = 'CANCELLED',
    Expired = 'EXPIRED',
    Unfunded = 'UNFUNDED',
    Fillability_increased = 'FILLABILITY_INCREASED',
}

interface OrderEventPayload {
    subscription: string;
    result: OrderEvent[];
}

interface OrderEvent {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    kind: OrderEventKind;
    fillableTakerAssetAmount: string;
    txHash: string;
}

interface AcceptedOrderInfo {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    fillableTakerAssetAmount: string;
}

enum RejectedKind {
    Zeroex_validation = 'ZEROEX_VALIDATION',
    Mesh_error = 'MESH_ERROR',
    Mesh_validation = 'MESH_VALIDATION',
}

enum RejectedCode {
    InternalError = 'InternalError',
    MaxOrderSizeExceeded = 'MaxOrderSizeExceeded',
    OrderAlreadyStored = 'OrderAlreadyStored',
    OrderForIncorrectNetwork = 'OrderForIncorrectNetwork',
    NetworkRequestFailed = 'NetworkRequestFailed',
    OrderHasInvalidMakerAssetAmount = 'OrderHasInvalidMakerAssetAmount',
    OrderHasInvalidTakerAssetAmount = 'OrderHasInvalidTakerAssetAmount',
    OrderExpired = 'OrderExpired',
    OrderFullyFilled = 'OrderFullyFilled',
    OrderCancelled = 'OrderCancelled',
    OrderUnfunded = 'OrderUnfunded',
    OrderHasInvalidMakerAssetData = 'OrderHasInvalidMakerAssetData',
    OrderHasInvalidTakerAssetData = 'OrderHasInvalidTakerAssetData',
    OrderHasInvalidSignature = 'OrderHasInvalidSignature',
}

interface RejectedStatus {
    code: RejectedCode;
    message: string;
}

interface RejectedOrderInfo {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    kind: RejectedKind;
    status: RejectedStatus;
}

interface ValidationResults {
    accepted: AcceptedOrderInfo[];
    rejected: RejectedOrderInfo[];
}

interface GetOrdersResponse {
    snapshotID: string;
    ordersInfos: AcceptedOrderInfo[];
}

console.log('Mesh WebSocket endpoint: ', MESH_WS_ENDPOINT);

(async () => {
    // Instantiate the WebSocket provider/client
    const websocketProvider = new Web3Providers.WebsocketProvider(MESH_WS_ENDPOINT, {
        clientConfig: {
            // For some reason fragmenting the payloads causes the connection to close
            // Source: https://github.com/theturtle32/WebSocket-Node/issues/359
            fragmentOutgoingMessages: false,
        },
    } as any);

    // Listen for the close event which will fire if Mesh goes down
    (websocketProvider as any).connection.addEventListener('close', () => {
        console.log('close event received');
        process.exit(1);
    });

    console.log('Subscribing to heartbeat...');
    const heartbeatSubscriptionId = await websocketProvider.subscribe('mesh_subscribe', 'heartbeat', []);
    console.log('Heartbeat subscriptionId', heartbeatSubscriptionId);
    // Listen to event on the subscription (topic is the subscriptionId)
    const heartbeatCallback = (eventPayload: OrderEventPayload) => {
        console.log('Received:', JSON.stringify(eventPayload, null, '\t'));
    };
    websocketProvider.on(heartbeatSubscriptionId, heartbeatCallback as any);

    console.log('Subscribing to order events...');
    const orderEventsSubscriptionId = await websocketProvider.subscribe('mesh_subscribe', 'orders', []);
    console.log('Order events subscriptionId', orderEventsSubscriptionId);
    // Listen to event on the subscription (topic is the subscriptionId)
    const orderEventsCallback = (eventPayload: OrderEventPayload) => {
        console.log('Received:', JSON.stringify(eventPayload, null, '\t'));
    };
    websocketProvider.on(orderEventsSubscriptionId, orderEventsCallback as any);

    // Submit an order to the Mesh node
    console.log('Sending order...');
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
    const validationResults: ValidationResults = await (websocketProvider as any).send('mesh_addOrders', [orders]);
    console.log('mesh_addOrders Response: ', JSON.stringify(validationResults, null, '\t'));

    // Get all orders stored in Mesh at a snapshot in time
    const perPage = 200;
    let snapshotID = ''; // New snapshot
    let ordersInfosLen = 1;

    let i = 0;
    const allOrdersInfos = [];
    while (ordersInfosLen !== 0) {
        const page = i;
        const getOrdersResponse: GetOrdersResponse = await (websocketProvider as any).send('mesh_getOrders', [
            page,
            perPage,
            snapshotID,
        ]);
        console.log('mesh_getOrders Response:', JSON.stringify(getOrdersResponse, null, '\t'));
        snapshotID = getOrdersResponse.snapshotID;
        allOrdersInfos.push(...getOrdersResponse.ordersInfos);
        ordersInfosLen = getOrdersResponse.ordersInfos.length;
        i++;
    }
    console.log('Got ', allOrdersInfos.length, 'orders from snapshot ', snapshotID);
})();
