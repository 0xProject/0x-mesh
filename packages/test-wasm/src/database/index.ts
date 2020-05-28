import Dexie, { Transaction } from 'dexie';

export interface Options {
    dataSourceName: string;
    maxOrders: number;
    maxMiniHeaders: number;
}

export enum SortDirection {
    Asc = 'ASC',
    Desc = 'DESC',
}

export enum FilterKind {
    Equal = '=',
    NotEqual = '!=',
    Less = '<',
    Greater = '>',
    LessOrEqual = '<=',
    GreaterOrEqual = '>=',
    Contains = 'CONTAINS',
}

export interface Order {
    hash: string;
    chainId: number;
    makerAddress: string;
    makerAssetData: string;
    makerAssetAmount: string;
    makerFee: string;
    makerFeeAssetData: string;
    takerAddress: string;
    takerAssetData: string;
    takerFeeAssetData: string;
    takerAssetAmount: string;
    takerFee: string;
    senderAddress: string;
    feeRecipientAddress: string;
    expirationTimeSeconds: string;
    salt: string;
    signature: string;
    exchangeAddress: string;
    fillableTakerAssetAmount: string;
    lastUpdated: string;
    isRemoved: boolean;
    isPinned: boolean;
    parsedMakerAssetData: string;
    parsedMakerFeeAssetData: string;
}

export type OrderField = keyof Order;

export interface OrderQuery {
    filters?: OrderFilter[];
    sort?: OrderSort[];
    limit?: number;
    offset?: number;
}

export interface OrderSort {
    field: OrderField;
    direction: SortDirection;
}

export interface OrderFilter {
    field: OrderField;
    kind: FilterKind;
    value: any;
}

export interface AddOrdersResult {
    added: Order[];
    removed: Order[];
}

export interface MiniHeader {
    // TODO(albrow): define this
}

export interface MiniHeaderQuery {
    // TODO(albrow): define this
}

export interface AddMiniHeadersResult {
    added: MiniHeader[];
    removed: MiniHeader[];
}

export interface Metadata {
    // TODO(albrow): define this
}

function newNotFoundError(): Error {
    return new Error('could not find existing model or row in database');
}

// TODO(albrow): Implement remaining methods.
export class Database {
    private db: Dexie;
    private maxOrders: number;
    private maxMiniHeaders: number;
    private orders: Dexie.Table<Order, string>;

    constructor(opts: Options) {
        this.db = new Dexie(opts.dataSourceName);
        this.maxOrders = opts.maxOrders;
        this.maxMiniHeaders = opts.maxMiniHeaders;

        this.db.version(1).stores({
            // TODO(albrow): Add more indexes. https://dexie.org/docs/Version/Version.stores()
            orders:
                '&hash,chainId,makerAddress,makerAssetData,makerAssetAmount,makerFee,makerFeeAssetData,takerAddress,takerAssetData,takerFeeAssetData,takerAssetAmount,takerFee,senderAddress,feeRecipientAddress,expirationTimeSeconds,salt,signature,exchangeAddress,fillableTakerAssetAmount,lastUpdated,isRemoved,isPinned,parsedMakerAssetData,parsedMakerFeeAssetData',
            miniheaders: '&hash',
            metadata: '',
        });

        this.orders = this.db.table('orders');
    }

    // AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
    public async addOrdersAsync(orders: Order[]): Promise<AddOrdersResult> {
        // TODO(albrow): Remove orders with max expiration time.
        const added: Order[] = [];
        await this.db.transaction('rw', this.orders, async (txn: Transaction) => {
            for (const order of orders) {
                try {
                    await this.orders.add(order);
                } catch (e) {
                    if (e.name === 'ConstraintError') {
                        // An order with this hash already exists. This is fine based on the semantics of
                        // addOrders.
                        continue;
                    }
                    throw e;
                }
                added.push(order);
            }
        });
        return {
            added,
            removed: [],
        };
    }

    // GetOrder(hash common.Hash) (*types.OrderWithMetadata, error)
    public async getOrderAsync(hash: string): Promise<Order> {
        const order = await this.orders.get(hash);
        if (order === undefined) {
            throw newNotFoundError();
        }
        return order;
    }

    // FindOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async findOrdersAsync(query?: OrderQuery): Promise<Order[]> {
        const collection = this.prepareOrderQuery(query);
        return this.runQueryAsync(query, collection);
    }

    // CountOrders(opts *OrderQuery) (int, error)
    public async countOrdersAsync(query: OrderQuery): Promise<number> {
        const collection = this.prepareOrderQuery(query);
        return collection.count();
    }

    // DeleteOrder(hash common.Hash) error
    public async deleteOrderAsync(hash: string): Promise<void> {
        return this.orders.delete(hash);
    }

    // DeleteOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async deleteOrdersAsync(query: OrderQuery | undefined): Promise<Order[]> {
        const deletedOrders: Order[] = [];
        await this.db.transaction('rw', this.orders, async () => {
            // TODO(albrow): Pay special attention to this code and make sure it works.
            // Behavior of Dexie.js regarding transactions and function scope is a little
            // too complicated/magical.
            const orders = await this.findOrdersAsync(query);
            for (const order of orders) {
                await this.orders.delete(order.hash);
                deletedOrders.push(order);
            }
        });
        return deletedOrders;
    }

    // UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error
    public async updateOrderAsync(hash: string, updateFunc: (existingOrder: Order) => Order): Promise<void> {
        await this.db.transaction('rw', this.orders, async () => {
            // TODO(albrow): Pay special attention to this code and make sure it works.
            // Behavior of Dexie.js regarding transactions and function scope is a little
            // too complicated/magical.
            const existingOrder = await this.getOrderAsync(hash);
            const updatedOrder = updateFunc(existingOrder);
            await this.orders.put(updatedOrder, hash);
        });
    }

    // AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error)
    public async addMiniHeadersAsync(miniHeaders: MiniHeader[]): Promise<AddMiniHeadersResult> {
        return Promise.reject('not yet implemented');
    }

    // GetMiniHeader(hash common.Hash) (*types.MiniHeader, error)
    public async getMiniHeaderAsync(hash: string): Promise<MiniHeader> {
        return Promise.reject('not yet implemented');
    }

    // FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async findMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        return Promise.reject('not yet implemented');
    }

    // DeleteMiniHeader(hash common.Hash) error
    public async deleteMiniHeaderAsync(hash: string): Promise<void> {
        return Promise.reject('not yet implemented');
    }

    // DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async deleteMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        return Promise.reject('not yet implemented');
    }

    // GetMetadata() (*types.Metadata, error)
    public async getMetadataAsync(): Promise<Metadata> {
        return Promise.reject('not yet implemented');
    }

    // SaveMetadata(metadata *types.Metadata) error
    public async saveMetadataAsync(metadata: Metadata): Promise<void> {
        return Promise.reject('not yet implemented');
    }

    // UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error
    public async updateMetadatasAsync(updateFunc: (existingMetadata: Metadata) => Metadata): Promise<void> {
        return Promise.reject('not yet implemented');
    }

    prepareOrderQuery(query: OrderQuery | undefined): Dexie.Collection<Order, string> {
        if (query === null || query === undefined) {
            return this.orders.toCollection();
        }
        var col: Dexie.Collection;
        if (query.filters !== undefined && query.filters !== null && query.filters.length > 0) {
            const firstFilter = query.filters[0];
            switch (query.filters[0].kind) {
                case FilterKind.Equal:
                    col = this.orders.where(firstFilter.field).equals(firstFilter.value);
                    break;
                case FilterKind.NotEqual:
                    col = this.orders.where(firstFilter.field).notEqual(firstFilter.value);
                    break;
                case FilterKind.Greater:
                    col = this.orders.where(firstFilter.field).above(firstFilter.value);
                    break;
                case FilterKind.GreaterOrEqual:
                    col = this.orders.where(firstFilter.field).aboveOrEqual(firstFilter.value);
                    break;
                case FilterKind.Less:
                    col = this.orders.where(firstFilter.field).below(firstFilter.value);
                    break;
                case FilterKind.LessOrEqual:
                    col = this.orders.where(firstFilter.field).belowOrEqual(firstFilter.value);
                    break;
                case FilterKind.Contains:
                    // TODO(albrow): This iterates through all orders and is very inefficient.
                    // Is there a way to optimize this?
                    col = this.orders.filter(order => {
                        return order[firstFilter.field].toString().includes(firstFilter.value);
                    });
                    break;
            }
            if (query.filters.length > 1) {
                // TODO(albrow): Dexie.js does not support multiple where conditions. We have to
                // use Collection.and which iterates through all orders in the collection so far
                // and is very inefficient. Is there a way to optimize this?
                query.filters.slice(1).forEach(filter => {
                    switch (filter.kind) {
                        case FilterKind.Equal:
                            col.and(order => order[filter.field] === filter.value);
                            break;
                        case FilterKind.NotEqual:
                            col.and(order => order[filter.field] !== filter.value);
                            break;
                        case FilterKind.Greater:
                            col.and(order => order[filter.field] > filter.value);
                            break;
                        case FilterKind.GreaterOrEqual:
                            col.and(order => order[filter.field] >= filter.value);
                            break;
                        case FilterKind.Less:
                            col.and(order => order[filter.field] < filter.value);
                            break;
                        case FilterKind.LessOrEqual:
                            col.and(order => order[filter.field] <= filter.value);
                            break;
                        case FilterKind.Contains:
                            col.and(order => {
                                return order[filter.field].toString().includes(filter.value);
                            });
                            break;
                    }
                });
            }
        } else {
            col = this.orders.toCollection();
        }
        if (query.offset !== undefined && query.offset !== 0) {
            col.offset(query.offset);
        }
        if (query.limit !== undefined && query.limit !== 0) {
            col.limit(query.limit);
        }
        return col;
    }

    async runQueryAsync(query: OrderQuery | undefined, col: Dexie.Collection<Order, string>): Promise<Order[]> {
        if (
            query === null ||
            query === undefined ||
            query.sort === null ||
            query.sort === undefined ||
            query.sort.length === 0
        ) {
            return col.toArray();
        } else {
            if (
                (query.offset !== null && query.offset !== undefined && query.offset !== 0) ||
                (query.limit !== null && query.limit !== undefined && query.limit !== 0)
            ) {
                if (
                    query.filters === null ||
                    query.filters === undefined ||
                    (query.filters.length === 1 && query.filters[0].field === 'hash')
                ) {
                    // This is okay.
                } else {
                    // TODO(albrow): Technically this is allowed if and only if
                    // there is exactly one filter, exactly one sort, and the sort
                    // field is equal to the filter field.
                    throw new Error('sorting by arbitrary fields with limit and offset is not supported by Dexie.js');
                }
            }

            // Note(albrow): Dexie.js can't sort by more than one field. Looks like
            // we have no choice but to manually sort here. This is not fast or
            // efficient.
            return (await col.toArray()).sort((a: Order, b: Order) => {
                for (const s of query.sort!) {
                    switch (s.direction) {
                        case SortDirection.Asc:
                            if (a[s.field] < b[s.field]) {
                                return -1;
                            } else if (a[s.field] > b[s.field]) {
                                return 1;
                            }
                            break;
                        case SortDirection.Desc:
                            if (a[s.field] > b[s.field]) {
                                return -1;
                            } else if (a[s.field] < b[s.field]) {
                                return 1;
                            }
                            break;
                    }
                }
                return 0;
            });
        }
    }
}
