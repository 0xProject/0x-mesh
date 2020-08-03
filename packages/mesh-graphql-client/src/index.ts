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
    peerID: string;
    ethereumChainID: number;
    latestBlock?: LatestBlock;
    numPeers: number;
    numOrders: number;
    numOrdersIncludingRemoved: number;
    startOfCurrentUTCDay: Date;
    ethRPCRequestsSentInCurrentUTCDay: number;
    ethRPCRateLimitExpiredRequests: number;
    maxExpirationTime: BigNumber;
}

interface StatsResponse {
    stats: Stats;
}

interface OrderEventResponse {
    orderEvents: OrderEvent[];
}

export interface LatestBlock {
    number: BigNumber;
    hash: string;
}

export interface OrderEvent {
    timestamp: Date;
    endState: string;
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
        wsSubClient.onDisconnected(() => console.error('detected disconnect'));
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
        this._client = new ApolloClient({ cache: new InMemoryCache(), link });
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
        const latestBlock = stats.latestBlock
            ? {
                  ...stats.latestBlock,
                  number: new BigNumber(stats.latestBlock.number),
              }
            : undefined;
        return {
            ...R.omit(['__typename'], stats),
            maxExpirationTime: new BigNumber(stats.maxExpirationTime),
            latestBlock,
        };
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
