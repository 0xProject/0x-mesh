import { MeshWrapper } from '@0x/mesh-browser-lite/lib/types';
import { Observable } from '@apollo/client';
import { ApolloLink, FetchResult, Operation } from '@apollo/client/link/core';

export class BrowserLink extends ApolloLink {
    // FIXME - This should be a required parameter
    constructor(private readonly _meshWrapper?: MeshWrapper) {
        super();
    }

    public request(operation: Operation): Observable<FetchResult> | null {
        // FIXME: We probably want to do some basic sanity checks on these operations.
        // This could include checking that `AddOrders` is a mutation and a few
        // other things.
        switch (operation.operationName) {
            case 'AddOrders':
                // FIXME - Remove `if-statement` once testing is done
                if (this._meshWrapper) {
                    // FIXME - Is this how observables should be used?
                    return new Observable(async () => {
                        return {
                            // tslint:disable-next-line:no-non-null-assertion
                            data: await this._meshWrapper!.gqlAddOrdersAsync(
                                operation.variables.orders,
                                operation.variables.pinned,
                            ),
                        };
                    });
                } else {
                    console.log(operation.variables);
                }
                break;
            case 'Order':
                if (this._meshWrapper) {
                    // FIXME - Is this how observables should be used?
                    return new Observable(async () => {
                        return {
                            // tslint:disable-next-line:no-non-null-assertion
                            data: await this._meshWrapper!.gqlGetOrderAsync(operation.variables.hash),
                        };
                    });
                } else {
                    console.log(operation.variables);
                }
                break;
            case 'Orders':
                if (this._meshWrapper) {
                    // FIXME - Is this how observables should be used?
                    return new Observable(async () => {
                        return {
                            // tslint:disable-next-line:no-non-null-assertion
                            data: await this._meshWrapper!.gqlFindOrdersAsync(
                                operation.variables.sort,
                                operation.variables.filters,
                                operation.variables.limit,
                            ),
                        };
                    });
                } else {
                    console.log(operation.variables);
                }
                break;
            case 'Stats':
                if (this._meshWrapper) {
                    // FIXME - Is this how observables should be used?
                    return new Observable(async () => {
                        return {
                            // tslint:disable-next-line:no-non-null-assertion
                            data: await this._meshWrapper!.gqlGetStatsAsync(),
                        };
                    });
                } else {
                    console.log(operation);
                }
                break;
            // FIXME - A few of the operations do not populate the operationName
            // field
            default:
                throw new Error('browser link: unrecognized operation name');
        }
        return null;
    }
}
