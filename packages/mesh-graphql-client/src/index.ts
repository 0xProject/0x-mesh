import { Mesh } from '@0x/mesh-browser-lite';
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
import { ApolloLink } from '@apollo/client/link/core';
import { onError } from '@apollo/client/link/error';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import { SubscriptionClient } from 'subscriptions-transport-ws';
import * as ws from 'ws';
import * as Observable from 'zen-observable';

import { BrowserLink } from './browser_link';
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
    query Order($hash: String!) {
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

export interface LinkConfig {
    httpUrl?: string;
    webSocketUrl?: string;
    mesh?: Mesh;
}

export class MeshGraphQLClient {
    // NOTE(jalextowle): BrowserLink doesn't support subscriptions at this time.
    private readonly _subscriptionClient?: SubscriptionClient;
    private readonly _client: ApolloClient<NormalizedCacheObject>;
    constructor(linkConfig: LinkConfig) {
        let link: ApolloLink;
        if (linkConfig.httpUrl && linkConfig.webSocketUrl) {
            if (!linkConfig.httpUrl || !linkConfig.webSocketUrl) {
                throw new Error(
                    'mesh-graphql-client: Both "httpUrl" and "webSocketUrl" must be provided in "linkConfig" if a network link is used',
                );
            }

            // Set up an apollo client with WebSocket and HTTP links. This allows
            // us to use the appropriate transport based on the type of the query.
            const httpLink = new HttpLink({
                uri: linkConfig.httpUrl,
            });
            const wsSubClient = new SubscriptionClient(
                linkConfig.webSocketUrl,
                {
                    reconnect: false,
                },
                // Use ws in Node.js and native WebSocket in browsers.
                (process as any).browser ? undefined : ws,
            );
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
                    const allMessages = graphQLErrors.map(err => err.message).join('\n');
                    throw new Error(`GraphQL error(s): ${allMessages}`);
                }
                if (networkError != null) {
                    throw new Error(`Network error: ${networkError.message}`);
                }
            });
            link = from([errorLink, splitLink]);
            this._subscriptionClient = wsSubClient;
        } else {
            if (!linkConfig.mesh) {
                throw new Error(
                    'mesh-graphql-client: "httpUrl" and "webSocketUrl" cannot be provided if a browser link is used',
                );
            }

            link = new BrowserLink(linkConfig.mesh);
        }
        this._client = new ApolloClient({
            cache: new InMemoryCache({
                // This custom merge function is required for our orderEvents subscription.
                // See https://www.apollographql.com/docs/react/caching/cache-field-behavior/#the-merge-function
                typePolicies: {
                    Subscription: {
                        fields: {
                            orderEvents: {
                                merge(existing: OrderEvent[] = [], incoming: OrderEvent[]): OrderEvent[] {
                                    return [...existing, ...incoming];
                                },
                            },
                        },
                    },
                },
                // Stop apollo client from injecting `__typename` fields. These extra fields mess up our tests.
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

    public async findOrdersAsync(
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
        if (this._subscriptionClient !== undefined) {
            // NOTE(jalextowle): We must use a variable here because Typescript
            // thinks that this._subscriptionClient can become undefined between
            // Observable events.
            const subscriptionClient = this._subscriptionClient;

            // We handle incomingObservable and return a new outgoingObservable. This
            // can be thought of as "wrapping" the observable and we do it for two reasons:
            //
            // 1. Convert FetchResult<OrderEventResponse> to OrderEvent[]
            // 2. Handle errors and disconnects from the underlying websocket transport. If we don't
            //    do this, Apollo Client just ignores them completely and acts like everything is fine :(
            //
            const incomingObservable = this._client.subscribe({
                query: orderEventsSubscription,
            }) as Observable<FetchResult<OrderEventResponse>>;
            const outgoingObservable = new Observable<OrderEvent[]>(observer => {
                subscriptionClient.onError((err: ErrorEvent) => {
                    observer.error(new Error(err.message));
                });
                subscriptionClient.onDisconnected((event: Event) => {
                    observer.error(new Error('WebSocket connection lost'));
                });
                incomingObservable.subscribe({
                    next: (result: FetchResult<OrderEventResponse>) => {
                        if (result.errors != null && result.errors.length > 0) {
                            result.errors.forEach(err => observer.error(err));
                        } else if (result.data == null) {
                            observer.error(new Error('received no data'));
                        } else {
                            observer.next(result.data.orderEvents.map(fromStringifiedOrderEvent));
                        }
                    },
                    error: err => observer.error(err),
                    complete: () => observer.complete(),
                });
            });
            return outgoingObservable;
        } else {
            throw new Error(
                'mesh-graphql-client: Browser GraphQl API does not support subscriptions. Please use the legacy API to listen to events and errors',
            );
        }
    }

    public async rawQueryAsync<T = any, TVariables = OperationVariables>(
        options: QueryOptions<TVariables>,
    ): Promise<ApolloQueryResult<T>> {
        if (!this._subscriptionClient) {
            throw new Error('mesh-graphql-client: Raw queries are not currently supported by browser nodes');
        }
        return this._client.query<T>(options);
    }
}
