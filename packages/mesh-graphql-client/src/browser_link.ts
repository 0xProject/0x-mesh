import { Mesh } from '@0x/mesh-browser-lite';
import { ApolloLink, FetchResult, Operation } from '@apollo/client/link/core';
import * as Observable from 'zen-observable';

import { AddOrdersResponse, OrderResponse, OrdersResponse, StatsResponse } from './types';

export class BrowserLink extends ApolloLink {
    constructor(private readonly _mesh: Mesh) {
        super();
    }

    public request(operation: Operation): Observable<FetchResult> {
        const wrapper = this._mesh.wrapper;
        if (wrapper === undefined) {
            throw new Error('mesh-graphql-client: Mesh node is not ready to receive requests');
        }
        switch (operation.operationName) {
            case 'AddOrders':
                if (
                    operation.variables.opts.keepCancelled ||
                    operation.variables.opts.keepExpired ||
                    operation.variables.opts.keepFullyFilled ||
                    operation.variables.opts.keepUnfunded
                ) {
                    throw new Error('mesh-graphql-client: Browser nodes do not support true values in AddOrdersOpts');
                }
                return new Observable<{ data: AddOrdersResponse }>((observer) => {
                    wrapper
                        .gqlAddOrdersAsync(operation.variables.orders, operation.variables.pinned)
                        .then((addOrders) => {
                            observer.next({ data: { addOrders } });
                            observer.complete();
                            return { data: { addOrders } };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            case 'Order':
                return new Observable<{ data: OrderResponse }>((observer) => {
                    wrapper
                        .gqlGetOrderAsync(operation.variables.hash)
                        .then((order) => {
                            observer.next({ data: { order } });
                            observer.complete();
                            return { data: { order } };
                        })
                        .catch((err: Error) => {
                            throw err;
                        });
                });
            case 'Orders':
                return new Observable<{ data: OrdersResponse }>((observer) => {
                    wrapper
                        .gqlFindOrdersAsync(
                            operation.variables.sort,
                            operation.variables.filters,
                            operation.variables.limit,
                        )
                        .then((orders) => {
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
                return new Observable<{ data: StatsResponse }>((observer) => {
                    wrapper
                        .gqlGetStatsAsync()
                        .then((stats) => {
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
