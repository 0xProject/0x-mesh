import { SignedOrder } from '@0x/types';
import { BigNumber } from '@0x/utils';

/**
 * WebSocketClient configs
 * Source: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md#client-config-options
 */
export interface ClientConfig {
    webSocketVersion?: number;
    maxReceivedFrameSize?: number;
    maxReceivedMessageSize?: number;
    fragmentOutgoingMessages?: boolean;
    fragmentationThreshold?: number;
    assembleFragments?: boolean;
    closeTimeout?: number;
    tlsOptions?: any;
}

/**
 * timeout: timeout in milliseconds to enforce on every WS request that expects a response
 * headers: Request headers (e.g., authorization)
 * protocol: requestOptions should be either null or an object specifying additional configuration options to be
 * passed to http.request or https.request. This can be used to pass a custom agent to enable WebSocketClient usage
 * from behind an HTTP or HTTPS proxy server using koichik/node-tunnel or similar.
 * clientConfig: The client configs documented here: https://github.com/theturtle32/WebSocket-Node/blob/master/docs/WebSocketClient.md
 * reconnectAfter: time in milliseconds after which to attempt to reconnect to WS server after an error occurred (default: 5000)
 */
export interface WSOpts {
    timeout?: number;
    headers?: {};
    protocol?: string;
    clientConfig?: ClientConfig;
    reconnectAfter?: number;
}

export interface StringifiedSignedOrder {
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

export enum OrderEventKind {
    Invalid = 'INVALID',
    Added = 'ADDED',
    Filled = 'FILLED',
    FullyFilled = 'FULLY_FILLED',
    Cancelled = 'CANCELLED',
    Expired = 'EXPIRED',
    Unfunded = 'UNFUNDED',
    FillabilityIncreased = 'FILLABILITY_INCREASED',
}

export interface OrderEventPayload {
    subscription: string;
    result: RawOrderEvent[];
}

export interface HeartbeatEventPayload {
    subscription: string;
    result: string;
}

export interface RawOrderEvent {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    kind: OrderEventKind;
    fillableTakerAssetAmount: string;
    txHashes: string[];
}

export interface OrderEvent {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: OrderEventKind;
    fillableTakerAssetAmount: BigNumber;
    txHashes: string[];
}

export interface RawAcceptedOrderInfo {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    fillableTakerAssetAmount: string;
    isNew: boolean;
}

export interface AcceptedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: BigNumber;
    isNew: boolean;
}

export interface RawOrderInfo {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    fillableTakerAssetAmount: string;
}

export interface OrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    fillableTakerAssetAmount: BigNumber;
}

export enum RejectedKind {
    ZeroexValidation = 'ZEROEX_VALIDATION',
    MeshError = 'MESH_ERROR',
    MeshValidation = 'MESH_VALIDATION',
}

export enum RejectedCode {
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

export interface RejectedStatus {
    code: RejectedCode;
    message: string;
}

export interface RawRejectedOrderInfo {
    orderHash: string;
    signedOrder: StringifiedSignedOrder;
    kind: RejectedKind;
    status: RejectedStatus;
}

export interface RejectedOrderInfo {
    orderHash: string;
    signedOrder: SignedOrder;
    kind: RejectedKind;
    status: RejectedStatus;
}

export interface RawValidationResults {
    accepted: RawAcceptedOrderInfo[];
    rejected: RawRejectedOrderInfo[];
}

export interface ValidationResults {
    accepted: AcceptedOrderInfo[];
    rejected: RejectedOrderInfo[];
}

export interface GetOrdersResponse {
    snapshotID: string;
    ordersInfos: RawAcceptedOrderInfo[];
}

export interface WSMessage {
    type: string;
    utf8Data: string;
}

export interface LatestBlock {
    number: number;
    hash: string;
}

export interface GetStatsResponse {
    Version: string;
    PubSubTopic: string;
    Rendezvous: string;
    PeerID: string;
    EthereumNetworkID: number;
    LatestBlock: LatestBlock;
    NumPeers: number;
    NumOrders: number;
}
