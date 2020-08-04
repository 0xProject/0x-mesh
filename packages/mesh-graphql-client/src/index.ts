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
    OperationVariables,
    QueryOptions,
} from '@apollo/client/core';
import { onError } from '@apollo/client/link/error';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import { SubscriptionClient } from 'subscriptions-transport-ws';
import * as ws from 'ws';
import * as Observable from 'zen-observable';

const defaultOrderQueryLimit = 100;

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

interface OrderResponse {
    order: StringifiedOrderWithMetadata | null;
}

interface OrdersResponse {
    orders: StringifiedOrderWithMetadata[];
}

interface OrderEventResponse {
    orderEvents: StringifiedOrderEvent[];
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

export interface OrderEvent {
    timestampMs: number;
    order: OrderWithMetadata;
    endState: OrderEventEndState;
    contractEvents: ContractEvent[];
}

export interface ContractEvent {
    blockHash: string;
    txHash: string;
    txIndex: number;
    logIndex: number;
    isRemoved: boolean;
    address: string;
    // TODO(albrow): Use an enum type for kind?
    kind: string;
    // TODO(albrow): Use a union type for parameters?
    parameters: any;
}

export enum ContractEventKind {
    ERC20TransferEvent = 'ERC20TransferEvent',
    ERC20ApprovalEvent = 'ERC20ApprovalEvent',
    ERC721TransferEvent = 'ERC721TransferEvent',
    ERC721ApprovalEvent = 'ERC721ApprovalEvent',
    ERC721ApprovalForAllEvent = 'ERC721ApprovalForAllEvent',
    ERC1155ApprovalForAllEvent = 'ERC1155ApprovalForAllEvent',
    ERC1155TransferSingleEvent = 'ERC1155TransferSingleEvent',
    ERC1155TransferBatchEvent = 'ERC1155TransferBatchEvent',
    ExchangeFillEvent = 'ExchangeFillEvent',
    ExchangeCancelEvent = 'ExchangeCancelEvent',
    ExchangeCancelUpToEvent = 'ExchangeCancelUpToEvent',
    WethDepositEvent = 'WethDepositEvent',
    WethWithdrawalEvent = 'WethWithdrawalEvent',
}

export enum OrderEventEndState {
    // The order was successfully validated and added to the Mesh node. The order is now being watched and any changes to
    // the fillability will result in subsequent order events.
    Added = 'ADDED',
    // The order was filled for a partial amount. The order is still fillable up to the fillableTakerAssetAmount.
    Filled = 'FILLED',
    // The order was fully filled and its remaining fillableTakerAssetAmount is 0. The order is no longer fillable.
    FullyFilled = 'FULLY_FILLED',
    // The order was cancelled and is no longer fillable.
    Cancelled = 'CANCELLED',
    // The order expired and is no longer fillable.
    Expired = 'EXPIRED',
    // The order was previously expired, but due to a block re-org it is no longer considered expired (should be rare).
    Unexpired = 'UNEXPIRED',
    // The order has become unfunded and is no longer fillable. This can happen if the maker makes a transfer or changes their allowance.
    Unfunded = 'UNFUNDED',
    // The fillability of the order has increased. This can happen if a previously processed fill event gets reverted due to a block re-org,
    // or if a maker makes a transfer or changes their allowance.
    FillabilityIncreased = 'FILLABILITY_INCREASED',
    // The order is potentially still valid but was removed for a different reason (e.g.
    // the database is full or the peer that sent the order was misbehaving). The order will no longer be watched
    // and no further events for this order will be emitted. In some cases, the order may be re-added in the
    // future.
    StoppedWatching = 'STOPPED_WATCHING',
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

const orderQuery = gql`
    query Order($hash: Hash!) {
        order(hash: $hash) {
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
    }
`;

const ordersQuery = gql`
    query Orders(
        $filters: [OrderFilter!] = []
        $sort: [OrderSort!] = [{ field: hash, direction: ASC }]
        $limit: Int = 100
    ) {
        orders(filters: $filters, sort: $sort, limit: $limit) {
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

export type OrderField = Extract<keyof OrderWithMetadata, string>;

export enum SortDirection {
    Asc = 'ASC',
    Desc = 'DESC',
}

export enum FilterKind {
    Equal = 'EQUAL',
    NotEqual = 'NOT_EQUAL',
    Greater = 'GREATER',
    GreaterOrEqual = 'GREATER_OR_EQUAL',
    Less = 'LESS',
    LessOrEqual = 'LESS_OR_EQUAL',
}

export interface OrderSort {
    field: OrderField;
    direction: SortDirection;
}

export interface OrderFilter {
    field: OrderField;
    kind: FilterKind;
    value: OrderWithMetadata[OrderField];
}

export interface OrderQuery {
    filters?: OrderFilter[];
    sort?: OrderSort[];
    limit?: number;
}

export class MeshGraphQLClient {
    // private readonly _subscriptionClient: SubscriptionClient;
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
                        // console.error(err);
                    } else {
                        // console.log('successfully connected');
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
            if (graphQLErrors != null && graphQLErrors.length > 0) {
                // Throw the first error.
                // TODO(albrow): Is there a clean way to include all the errors?
                const firstErr = graphQLErrors[0];
                throw new Error(
                    `[GraphQL error]: Message: ${firstErr.message}, Location: ${firstErr.locations}, Path: ${firstErr.path}`,
                );
            }
            if (networkError != null) {
                throw new Error(`[Network error]: ${networkError}`);
            }
        });
        const link = from([errorLink, splitLink]);
        // this._subscriptionClient = wsSubClient;
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
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        const results = resp.data.addOrders;
        return fromStringifiedAddOrdersResults(results);
    }

    public async getOrderAsync(hash: string): Promise<OrderWithMetadata | null> {
        const resp: ApolloQueryResult<OrderResponse> = await this._client.query({
            query: orderQuery,
            variables: {
                hash,
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        if (resp.data.order == null) {
            return null;
        }
        return fromStringifiedOrderWithMetadata(resp.data.order);
    }

    public async getOrdersAsync(
        query: OrderQuery = { sort: [], filters: [], limit: defaultOrderQueryLimit },
    ): Promise<OrderWithMetadata[]> {
        const resp: ApolloQueryResult<OrdersResponse> = await this._client.query({
            query: ordersQuery,
            variables: {
                sort: query.sort || [],
                filters: query.filters?.map(convertFilterValue) || [],
                limit: query.limit || defaultOrderQueryLimit,
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        return resp.data.orders.map(fromStringifiedOrderWithMetadata);
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
            return result.data.orderEvents.map(fromStringifiedOrderEvent);
        });
    }

    public async rawQueryAsync<T = any, TVariables = OperationVariables>(
        options: QueryOptions<TVariables>,
    ): Promise<ApolloQueryResult<T>> {
        return this._client.query(options);
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

interface StringifiedOrderEvent {
    timestamp: string;
    order: StringifiedOrderWithMetadata;
    endState: OrderEventEndState;
    fillableTakerAssetAmount: BigNumber;
    contractEvents: ContractEvent[];
}

function fromStringifiedOrderEvent(event: StringifiedOrderEvent): OrderEvent {
    return {
        ...event,
        timestampMs: new Date(event.timestamp).getUTCMilliseconds(),
        order: fromStringifiedOrderWithMetadata(event.order),
    };
}

// converts any filter.value of type BigNumber to string.
function convertFilterValue(filter: OrderFilter): OrderFilter {
    return {
        ...filter,
        value: BigNumber.isBigNumber(filter.value) ? filter.value.toString() : filter.value,
    };
}

// tslint:disable-next-line: max-file-line-count
