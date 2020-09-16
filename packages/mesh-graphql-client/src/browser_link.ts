import { MeshWrapper } from '@0x/mesh-browser-lite/lib/types';
import { SignedOrder } from '@0x/order-utils';
import { ApolloLink, FetchResult, Operation } from '@apollo/client/link/core';
import * as Observable from 'zen-observable';

import { AddOrdersResponse, OrderResponse, OrdersResponse, StatsResponse } from './types';

export class BrowserLink extends ApolloLink {
    constructor(private readonly _meshWrapper: MeshWrapper) {
        super();
    }

    public request(operation: Operation): Observable<FetchResult> {
        switch (operation.operationName) {
            case 'AddOrders':
                return new Observable<{ data: AddOrdersResponse }>(observer => {
                    this._meshWrapper
                        .gqlAddOrdersAsync(operation.variables.orders, operation.variables.pinned)
                        .then(addOrders => {
                            observer.next({ data: { addOrders } });
                            observer.complete();
                            return { data: { addOrders } };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            case 'Order':
                return new Observable<{ data: OrderResponse }>(observer => {
                    this._meshWrapper
                        .gqlGetOrderAsync(operation.variables.hash)
                        .then(order => {
                            observer.next({ data: { order } });
                            observer.complete();
                            return { data: { order } };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            case 'Orders':
                return new Observable<{ data: OrdersResponse }>(observer => {
                    this._meshWrapper
                        .gqlFindOrdersAsync(
                            operation.variables.sort,
                            operation.variables.filters,
                            operation.variables.limit,
                        )
                        .then(orders => {
                            observer.next({
                                data: {
                                    orders,
                                },
                            });
                            observer.complete();
                            return {
                                data: {
                                    orders,
                                },
                            };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            case 'Stats':
                return new Observable<{ data: StatsResponse }>(observer => {
                    this._meshWrapper
                        .gqlGetStatsAsync()
                        .then(stats => {
                            observer.next({
                                data: {
                                    stats,
                                },
                            });
                            observer.complete();
                            return {
                                data: {
                                    stats,
                                },
                            };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            default:
                throw new Error('browser link: unrecognized operation name');
        }
    }
}
