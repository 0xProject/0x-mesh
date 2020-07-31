import { BigNumber } from '@0x/utils';
import { ApolloClient, ApolloQueryResult, gql, InMemoryCache, NormalizedCacheObject } from '@apollo/client/core';
import * as R from 'ramda';

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

export interface LatestBlock {
    number: BigNumber;
    hash: string;
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

export class MeshGraphQLClient {
    private readonly _client: ApolloClient<NormalizedCacheObject>;
    constructor(url: string) {
        this._client = new ApolloClient({ uri: url, cache: new InMemoryCache() });
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
}
