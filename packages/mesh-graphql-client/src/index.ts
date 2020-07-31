import { BigNumber } from '@0x/utils';
import { GraphQLClient } from 'graphql-request';

import { getSdk, LatestBlock as GeneratedLatestBlock, Sdk, Stats as GeneratedStats } from './generated/graphql';

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

export interface LatestBlock {
    number: BigNumber;
    hash: string;
}

export class MeshGraphQLClient {
    private readonly _sdk: Sdk;
    private readonly _client: GraphQLClient;
    constructor(url: string) {
        this._client = new GraphQLClient(url);
        this._sdk = getSdk(this._client);
    }

    public async getStatsAsync(): Promise<Stats> {
        const resp = await this._sdk.getStats();
        const latestBlock = resp.stats.latestBlock
            ? {
                  ...resp.stats.latestBlock,
                  number: new BigNumber(resp.stats.latestBlock.number),
              }
            : undefined;
        return {
            ...resp.stats,
            maxExpirationTime: new BigNumber(resp.stats.maxExpirationTime),
            latestBlock,
        };
    }
}
