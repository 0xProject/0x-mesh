import { SignedOrder } from '@0x/types';
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

import {
    AddOrdersResponse,
    AddOrdersResults,
    convertFilterValue,
    fromStringifiedAddOrdersResults,
    fromStringifiedOrderEvent,
    fromStringifiedOrderWithMetadata,
    fromStringifiedStats,
    OrderEvent,
    OrderEventResponse,
    OrderQuery,
    OrderResponse,
    OrdersResponse,
    OrderWithMetadata,
    Stats,
    StatsResponse,
    toStringifiedSignedOrder,
} from './types';

export {
    AddOrdersResults,
    OrderEvent,
    OrderQuery,
    OrderWithMetadata,
    Stats,
    OrderFilter,
    FilterKind,
    OrderField,
    OrderSort,
    SortDirection,
    OrderEventEndState,
    RejectedOrderCode,
} from './types';
export { SignedOrder } from '@0x/types';
export { ApolloQueryResult, QueryOptions } from '@apollo/client/core';
export { Observable };

const defaultOrderQueryLimit = 100;

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
        return this._client.query<T>(options);
    }
}
