import { SignedOrder } from '@0x/types';
import { from, HttpLink, split } from '@apollo/client';
import {
    ApolloClient,
    ApolloQueryResult,
    FetchResult,
    InMemoryCache,
    NormalizedCacheObject,
    OperationVariables,
    QueryOptions,
} from '@apollo/client/core';
import { ApolloLink } from '@apollo/client/link/core';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import { SubscriptionClient } from 'subscriptions-transport-ws';
import * as ws from 'ws';
import * as Observable from 'zen-observable';
import {
    addOrdersMutation,
    addOrdersMutationV4,
    orderEventsSubscription,
    orderQuery,
    orderQueryV4,
    ordersQuery,
    ordersQueryV4,
    statsQuery,
} from './queries';
import {
    AddOrdersOpts,
    AddOrdersResponse,
    AddOrdersResponseV4,
    AddOrdersResults,
    convertFilterValue,
    fromStringifiedAddOrdersResults,
    fromStringifiedAddOrdersResultsV4,
    fromStringifiedOrderEvent,
    fromStringifiedOrderWithMetadata,
    fromStringifiedOrderWithMetadataV4,
    fromStringifiedStats,
    OrderEvent,
    OrderEventResponse,
    OrderQuery,
    OrderResponse,
    OrderResponseV4,
    OrdersResponse,
    OrdersResponseV4,
    OrderWithMetadata,
    OrderWithMetadataV4,
    SignedOrderV4,
    Stats,
    StatsResponse,
    StringifiedOrderWithMetadata,
    StringifiedOrderWithMetadataV4,
    StringifiedSignedOrder,
    StringifiedSignedOrderV4,
    toStringifiedSignedOrder,
    toStringifiedSignedOrderV4,
} from './types';

export { SignedOrder } from '@0x/types';
export { ApolloQueryResult, QueryOptions } from '@apollo/client/core';
export {
    AcceptedOrderResult,
    AddOrdersResults,
    FilterKind,
    OrderEvent,
    OrderEventEndState,
    OrderField,
    OrderFilter,
    OrderQuery,
    OrderSort,
    OrderWithMetadata,
    OrderWithMetadataV4,
    RejectedOrderCode,
    RejectedOrderResult,
    SortDirection,
    Stats,
} from './types';
export { Observable };

const defaultOrderQueryLimit = 100;
export interface LinkConfig {
    httpUrl?: string;
    webSocketUrl?: string;
}

export class MeshGraphQLClient {
    // NOTE(jalextowle): BrowserLink doesn't support subscriptions at this time.
    private readonly _subscriptionClient?: SubscriptionClient;
    private readonly _client: ApolloClient<NormalizedCacheObject>;
    private readonly _onReconnectedCallbacks: (() => void)[] = [];
    constructor(linkConfig: LinkConfig) {
        let link: ApolloLink;
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
                reconnect: true,
            },
            // Use ws in Node.js and native WebSocket in browsers.
            (process as any).browser ? undefined : ws,
        );
        const wsLink = new WebSocketLink(wsSubClient);

        // HACK(kimpers): See https://github.com/apollographql/apollo-client/issues/5115#issuecomment-572318778
        // @ts-ignore at the time of writing the field is private and untyped
        const subscriptionClient = wsLink.subscriptionClient as SubscriptionClient;

        subscriptionClient.onReconnected(() => {
            for (const cb of this._onReconnectedCallbacks) {
                cb();
            }
        });

        const splitLink = split(
            ({ query }) => {
                const definition = getMainDefinition(query);
                return definition.kind === 'OperationDefinition' && definition.operation === 'subscription';
            },
            wsLink,
            httpLink,
        );
        link = from([splitLink]);
        this._subscriptionClient = wsSubClient;

        this._client = new ApolloClient({
            cache: new InMemoryCache({
                resultCaching: false,
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
            fetchPolicy: 'no-cache',
            query: statsQuery,
        });
        if (resp.data === undefined) {
            throw new Error('received no data');
        }
        const stats = resp.data.stats;
        return fromStringifiedStats(stats);
    }

    public async addOrdersAsync(
        orders: SignedOrder[],
        pinned: boolean = true,
        opts?: AddOrdersOpts,
    ): Promise<AddOrdersResults<OrderWithMetadata, SignedOrder>> {
        const resp: FetchResult<AddOrdersResponse<
            StringifiedOrderWithMetadata,
            StringifiedSignedOrder
        >> = await this._client.mutate({
            mutation: addOrdersMutation,
            variables: {
                orders: orders.map(toStringifiedSignedOrder),
                pinned,
                opts: {
                    keepCancelled: false,
                    keepExpired: false,
                    keepFullyFilled: false,
                    keepUnfunded: false,
                    ...opts,
                },
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        const results = resp.data.addOrders;
        return fromStringifiedAddOrdersResults(results);
    }

    public async addOrdersV4Async(
        orders: SignedOrderV4[],
        pinned: boolean = true,
        opts?: AddOrdersOpts,
    ): Promise<AddOrdersResults<OrderWithMetadataV4, SignedOrderV4>> {
        const resp: FetchResult<AddOrdersResponseV4<
            StringifiedOrderWithMetadataV4,
            StringifiedSignedOrderV4
        >> = await this._client.mutate({
            mutation: addOrdersMutationV4,
            variables: {
                orders: orders.map(toStringifiedSignedOrderV4),
                pinned,
                opts: {
                    keepCancelled: false,
                    keepExpired: false,
                    keepFullyFilled: false,
                    keepUnfunded: false,
                    ...opts,
                },
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }

        const results = resp.data.addOrdersV4;

        return fromStringifiedAddOrdersResultsV4(results);
    }

    public async getOrderAsync(hash: string): Promise<OrderWithMetadata | null> {
        const resp: ApolloQueryResult<OrderResponse> = await this._client.query({
            query: orderQuery,
            fetchPolicy: 'no-cache',
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

    public async getOrderV4Async(hash: string): Promise<OrderWithMetadataV4 | null> {
        const resp: ApolloQueryResult<OrderResponseV4> = await this._client.query({
            query: orderQueryV4,
            fetchPolicy: 'no-cache',
            variables: {
                hash,
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        if (resp.data.orderv4 == null) {
            return null;
        }
        return fromStringifiedOrderWithMetadataV4(resp.data.orderv4);
    }

    public async findOrdersAsync(
        query: OrderQuery = { sort: [], filters: [], limit: defaultOrderQueryLimit },
    ): Promise<OrderWithMetadata[]> {
        const resp: ApolloQueryResult<OrdersResponse> = await this._client.query({
            query: ordersQuery,
            fetchPolicy: 'no-cache',
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

    public async findOrdersV4Async(
        query: OrderQuery = { sort: [], filters: [], limit: defaultOrderQueryLimit },
    ): Promise<OrderWithMetadataV4[]> {
        const resp: ApolloQueryResult<OrdersResponseV4> = await this._client.query({
            query: ordersQueryV4,
            fetchPolicy: 'no-cache',
            variables: {
                sort: query.sort || [],
                filters: query.filters?.map(convertFilterValue) || [],
                limit: query.limit || defaultOrderQueryLimit,
            },
        });
        if (resp.data == null) {
            throw new Error('received no data');
        }
        return resp.data.ordersv4.map(fromStringifiedOrderWithMetadataV4);
    }

    public onReconnected(cb: () => void): void {
        this._onReconnectedCallbacks.push(cb);
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
                fetchPolicy: 'no-cache',
                query: orderEventsSubscription,
            }) as Observable<FetchResult<OrderEventResponse>>;
            const outgoingObservable = new Observable<OrderEvent[]>((observer) => {
                subscriptionClient.onError((err: ErrorEvent) => {
                    observer.error(new Error(err.message));
                });
                subscriptionClient.onDisconnected((event: Event) => {
                    observer.error(new Error('WebSocket connection lost'));
                });
                incomingObservable.subscribe({
                    next: (result: FetchResult<OrderEventResponse>) => {
                        if (result.errors != null && result.errors.length > 0) {
                            result.errors.forEach((err) => observer.error(err));
                        } else if (result.data == null) {
                            observer.error(new Error('received no data'));
                        } else {
                            observer.next(result.data.orderEvents.map(fromStringifiedOrderEvent));
                        }
                    },
                    error: (err) => observer.error(err),
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
