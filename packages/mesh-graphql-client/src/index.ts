import { SignedOrder } from '@0x/types';
import { BigNumber } from '@0x/utils';
import { from, HttpLink, split } from '@apollo/client';
import {
    ApolloClient,
    ApolloQueryResult,
    FetchResult,
    gql,
    InMemoryCache,
    NormalizedCacheObject,
} from '@apollo/client/core';
import { onError } from '@apollo/client/link/error';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import * as R from 'ramda';
import { SubscriptionClient } from 'subscriptions-transport-ws';
import * as ws from 'ws';
import * as Observable from 'zen-observable';

export interface Stats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock: LatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: BigNumber;
    startOfCurrentUTCDay: Date;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}

export interface LatestBlock {
    number: BigNumber;
    hash: string;
}

interface StringifiedLatestBlock {
    number: string;
    hash: string;
}

interface StatsResponse {
    stats: StringifiedStats;
}

interface AddOrdersResponse {
    addOrders: StringifiedAddOrdersResults;
}

interface OrderEventResponse {
    orderEvents: OrderEvent[];
}

export interface OrderEvent {
    timestamp: Date;
    endState: string;
}

export interface OrderWithMetadata extends SignedOrder {
    hash: string;
    fillableTakerAssetAmount: BigNumber;
}

export interface AddOrdersResults {
    // The set of orders that were accepted. Accepted orders will be watched and order events will be emitted if
    // their status changes.
    accepted: AcceptedOrderResult[];
    // The set of orders that were rejected, including the reason they were rejected. Rejected orders will not be
    // watched.
    rejected: RejectedOrderResult[];
}

export interface AcceptedOrderResult {
    // The order that was accepted, including metadata.
    order: OrderWithMetadata;
    // Whether or not the order is new. Set to true if this is the first time this Mesh node has accepted the order
    // and false otherwise.
    isNew: boolean;
}

export interface RejectedOrderResult {
    // The hash of the order. May be null if the hash could not be computed.
    hash?: string;
    // The order that was rejected.
    order: SignedOrder;
    // A machine-readable code indicating why the order was rejected. This code is designed to
    // be used by programs and applications and will never change without breaking backwards-compatibility.
    code: RejectedOrderCode;
    // A human-readable message indicating why the order was rejected. This message may change
    // in future releases and is not covered by backwards-compatibility guarantees.
    message: string;
}

export enum RejectedOrderCode {
    EthRpcRequestFailed = 'ETH_RPC_REQUEST_FAILED',
    CoordinatorRequestFailed = 'COORDINATOR_REQUEST_FAILED',
    CoordinatorSoftCancelled = 'COORDINATOR_SOFT_CANCELLED',
    CoordinatorEndpointNotFound = 'COORDINATOR_ENDPOINT_NOT_FOUND',
    OrderHasInvalidMakerAssetAmount = 'ORDER_HAS_INVALID_MAKER_ASSET_AMOUNT',
    OrderHasInvalidTakerAssetAmount = 'ORDER_HAS_INVALID_TAKER_ASSET_AMOUNT',
    OrderExpired = 'ORDER_EXPIRED',
    OrderFullyFilled = 'ORDER_FULLY_FILLED',
    OrderCancelled = 'ORDER_CANCELLED',
    OrderUnfunded = 'ORDER_UNFUNDED',
    OrderHasInvalidMakerAssetData = 'ORDER_HAS_INVALID_MAKER_ASSET_DATA',
    OrderHasInvalidMakerFeeAssetData = 'ORDER_HAS_INVALID_MAKER_FEE_ASSET_DATA',
    OrderHasInvalidTakerAssetData = 'ORDER_HAS_INVALID_TAKER_ASSET_DATA',
    OrderHasInvalidTakerFeeAssetData = 'ORDER_HAS_INVALID_TAKER_FEE_ASSET_DATA',
    OrderHasInvalidSignature = 'ORDER_HAS_INVALID_SIGNATURE',
    OrderMaxExpirationExceeded = 'ORDER_MAX_EXPIRATION_EXCEEDED',
    InternalError = 'INTERNAL_ERROR',
    MaxOrderSizeExceeded = 'MAX_ORDER_SIZE_EXCEEDED',
    OrderAlreadyStoredAndUnfillable = 'ORDER_ALREADY_STORED_AND_UNFILLABLE',
    OrderForIncorrectChain = 'ORDER_FOR_INCORRECT_CHAIN',
    IncorrectExchangeAddress = 'INCORRECT_EXCHANGE_ADDRESS',
    SenderAddressNotAllowed = 'SENDER_ADDRESS_NOT_ALLOWED',
    DatabaseFullOfOrders = 'DATABASE_FULL_OF_ORDERS',
}

const statsQuery = gql`
    query Stats {
        stats {
            version
            pubSubTopic
            rendezvous
            peerID
            ethereumChainID
            latestBlock {
                number
                hash
            }
            numPeers
            numOrders
            numOrdersIncludingRemoved
            startOfCurrentUTCDay
            ethRPCRequestsSentInCurrentUTCDay
            ethRPCRateLimitExpiredRequests
            maxExpirationTime
        }
    }
`;

const addOrdersMutation = gql`
    mutation AddOrders($orders: [NewOrder!]!, $pinned: Boolean = true) {
        addOrders(orders: $orders, pinned: $pinned) {
            accepted {
                order {
                    hash
                    chainId
                    exchangeAddress
                    makerAddress
                    makerAssetData
                    makerAssetAmount
                    makerFeeAssetData
                    makerFee
                    takerAddress
                    takerAssetData
                    takerAssetAmount
                    takerFeeAssetData
                    takerFee
                    senderAddress
                    feeRecipientAddress
                    expirationTimeSeconds
                    salt
                    signature
                    fillableTakerAssetAmount
                }
                isNew
            }
            rejected {
                hash
                code
                message
                order {
                    chainId
                    exchangeAddress
                    makerAddress
                    makerAssetData
                    makerAssetAmount
                    makerFeeAssetData
                    makerFee
                    takerAddress
                    takerAssetData
                    takerAssetAmount
                    takerFeeAssetData
                    takerFee
                    senderAddress
                    feeRecipientAddress
                    expirationTimeSeconds
                    salt
                    signature
                }
            }
        }
    }
`;

const orderEventsSubscription = gql`
    subscription {
        orderEvents {
            timestamp
            endState
            order {
                hash
                chainId
                exchangeAddress
                makerAddress
                makerAssetData
                makerAssetAmount
                makerFeeAssetData
                makerFee
                takerAddress
                takerAssetData
                takerAssetAmount
                takerFeeAssetData
                takerFee
                senderAddress
                feeRecipientAddress
                expirationTimeSeconds
                salt
                signature
                fillableTakerAssetAmount
            }
            contractEvents {
                blockHash
                txHash
                txIndex
                logIndex
                isRemoved
                address
                kind
                parameters
            }
        }
    }
`;

export class MeshGraphQLClient {
    private readonly _subscriptionClient: SubscriptionClient;
    private readonly _client: ApolloClient<NormalizedCacheObject>;
    constructor(httpUrl: string, webSocketUrl: string) {
        // Set up an apollo client with WebSocket and HTTP links. This allows
        // us to use the appropriate transport based on the type of the query.
        const httpLink = new HttpLink({
            uri: httpUrl,
        });
        const wsSubClient = new SubscriptionClient(
            webSocketUrl,
            {
                reconnectionAttempts: 3,
                reconnect: true,
                connectionCallback: err => {
                    if (err) {
                        console.error(err);
                    } else {
                        console.log('successfully connected');
                    }
                },
            },
            // Use ws in Node.js and native WebSocket in browsers.
            (process as any).browser ? undefined : ws,
        );
        // wsSubClient.onError(err => console.error(err));
        // wsSubClient.onDisconnected(() => console.error('detected disconnect'));
        const wsLink = new WebSocketLink(wsSubClient);
        const splitLink = split(
            ({ query }) => {
                const definition = getMainDefinition(query);
                return definition.kind === 'OperationDefinition' && definition.operation === 'subscription';
            },
            wsLink,
            httpLink,
        );
        const errorLink = onError(({ graphQLErrors, networkError }) => {
            if (graphQLErrors) {
                graphQLErrors.map(({ message, locations, path }) =>
                    console.log(`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`),
                );
            }
            if (networkError) {
                console.log(`[Network error]: ${networkError}`);
            }
        });
        const link = from([errorLink, splitLink]);
        this._subscriptionClient = wsSubClient;
        this._client = new ApolloClient({
            cache: new InMemoryCache({
                // Stop apollo client from injecting `__typenme` fields.
                addTypename: false,
            }),
            link,
        });
    }

    public async getStatsAsync(): Promise<Stats> {
        const resp: ApolloQueryResult<StatsResponse> = await this._client.query({
            query: statsQuery,
            errorPolicy: 'none',
        });
        if (resp.data === undefined) {
            throw new Error('received no data');
        }
        const stats = resp.data.stats;
        return fromStringifiedStats(stats);
    }

    public async addOrdersAsync(orders: SignedOrder[], pinned: boolean = true): Promise<AddOrdersResults> {
        const resp: FetchResult<AddOrdersResponse> = await this._client.mutate({
            mutation: addOrdersMutation,
            variables: {
                orders: orders.map(toStringifiedSignedOrder),
                pinned,
            },
            errorPolicy: 'none',
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        const results = resp.data.addOrders;
        // TODO(albrow): Convert response type.
        return fromStringifiedAddOrdersResults(results);
    }

    public onOrderEvents(): Observable<OrderEvent[]> {
        const observable = this._client.subscribe({
            query: orderEventsSubscription,
        }) as Observable<FetchResult<OrderEventResponse>>;
        return observable.map(result => {
            if (result.errors != null && result.errors.length > 0) {
                throw new Error(result.errors[0].message);
            }
            if (result.data === undefined || result.data === null) {
                throw new Error('no data received');
            }
            return result.data.orderEvents;
        });
    }
}

interface StringifiedStats {
    version: string;
    pubSubTopic: string;
    rendezvous: string;
    secondaryRendezvous: string[];
    peerID: string;
    ethereumChainID: number;
    latestBlock: StringifiedLatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    numPinnedOrders: number;
    maxExpirationTime: string;
    startOfCurrentUTCDay: string;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
}

function fromStringifiedStats(stats: StringifiedStats): Stats {
    return {
        ...stats,
        latestBlock: fromStringifiedLatestBlock(stats.latestBlock),
        maxExpirationTime: new BigNumber(stats.maxExpirationTime),
        startOfCurrentUTCDay: new Date(stats.startOfCurrentUTCDay),
    };
}

function fromStringifiedLatestBlock(latestBlock: StringifiedLatestBlock): LatestBlock {
    return {
        ...latestBlock,
        number: new BigNumber(latestBlock.number),
    };
}

interface StringifiedSignedOrder {
    chainId: string;
    exchangeAddress: string;
    makerAddress: string;
    takerAddress: string;
    feeRecipientAddress: string;
    senderAddress: string;
    makerAssetAmount: string;
    takerAssetAmount: string;
    makerFee: string;
    takerFee: string;
    expirationTimeSeconds: string;
    salt: string;
    makerAssetData: string;
    takerAssetData: string;
    makerFeeAssetData: string;
    takerFeeAssetData: string;
    signature: string;
}

function toStringifiedSignedOrder(order: SignedOrder): StringifiedSignedOrder {
    return {
        ...order,
        chainId: order.chainId.toString(),
        makerAssetAmount: order.makerAssetAmount.toString(),
        takerAssetAmount: order.takerAssetAmount.toString(),
        makerFee: order.makerFee.toString(),
        takerFee: order.takerFee.toString(),
        expirationTimeSeconds: order.expirationTimeSeconds.toString(),
        salt: order.salt.toString(),
    };
}

interface StringifiedOrderWithMetadata extends StringifiedSignedOrder {
    hash: string;
    fillableTakerAssetAmount: string;
}

function fromStringifiedOrderWithMetadata(order: StringifiedOrderWithMetadata): OrderWithMetadata {
    return {
        ...order,
        // tslint:disable-next-line: custom-no-magic-numbers
        chainId: Number.parseInt(order.chainId, 10),
        makerAssetAmount: new BigNumber(order.makerAssetAmount),
        takerAssetAmount: new BigNumber(order.takerAssetAmount),
        makerFee: new BigNumber(order.makerFee),
        takerFee: new BigNumber(order.takerFee),
        expirationTimeSeconds: new BigNumber(order.expirationTimeSeconds),
        salt: new BigNumber(order.salt),
        fillableTakerAssetAmount: new BigNumber(order.fillableTakerAssetAmount),
    };
}

function fromStringifiedSignedOrder(order: StringifiedSignedOrder): SignedOrder {
    return {
        ...order,
        // tslint:disable-next-line: custom-no-magic-numbers
        chainId: Number.parseInt(order.chainId, 10),
        makerAssetAmount: new BigNumber(order.makerAssetAmount),
        takerAssetAmount: new BigNumber(order.takerAssetAmount),
        makerFee: new BigNumber(order.makerFee),
        takerFee: new BigNumber(order.takerFee),
        expirationTimeSeconds: new BigNumber(order.expirationTimeSeconds),
        salt: new BigNumber(order.salt),
    };
}

export interface StringifiedAddOrdersResults {
    accepted: StringifiedAcceptedOrderResult[];
    rejected: StringifiedRejectedOrderResult[];
}

export interface StringifiedAcceptedOrderResult {
    order: StringifiedOrderWithMetadata;
    isNew: boolean;
}

export interface StringifiedRejectedOrderResult {
    hash?: string;
    order: StringifiedSignedOrder;
    code: RejectedOrderCode;
    message: string;
}

function fromStringifiedAddOrdersResults(results: StringifiedAddOrdersResults): AddOrdersResults {
    return {
        accepted: results.accepted.map(fromStringifiedAcceptedOrderResult),
        rejected: results.rejected.map(fromStringifiedRejectedOrderResult),
    };
}

function fromStringifiedAcceptedOrderResult(acceptedResult: StringifiedAcceptedOrderResult): AcceptedOrderResult {
    return {
        ...acceptedResult,
        order: fromStringifiedOrderWithMetadata(acceptedResult.order),
    };
}

function fromStringifiedRejectedOrderResult(rejectedResult: StringifiedRejectedOrderResult): RejectedOrderResult {
    return {
        ...rejectedResult,
        order: fromStringifiedSignedOrder(rejectedResult.order),
    };
}
