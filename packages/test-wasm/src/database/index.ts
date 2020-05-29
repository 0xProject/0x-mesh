import Dexie, { Transaction } from 'dexie';

export type Record = Order | MiniHeader | Metadata;

export interface Options {
    dataSourceName: string;
    maxOrders: number;
    maxMiniHeaders: number;
}

export interface Query<T extends Record> {
    filters?: FilterOption<T>[];
    sort?: SortOption<T>[];
    limit?: number;
    offset?: number;
}

export interface SortOption<T extends Record> {
    field: Extract<keyof T, string>;
    direction: SortDirection;
}

export interface FilterOption<T extends Record> {
    field: Extract<keyof T, string>;
    kind: FilterKind;
    value: any;
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

export type OrderQuery = Query<Order>;

export type OrderSort = SortOption<Order>;

export type OrderFilter = FilterOption<Order>;

export interface AddOrdersResult {
    added: Order[];
    removed: Order[];
}

export interface MiniHeader {
    hash: string;
    parent: string;
    number: string;
    timestamp: string;
    logs: string;
}

export type MiniHeaderField = keyof MiniHeader;

export type MiniHeaderQuery = Query<MiniHeader>;

export type MiniHeaderSort = SortOption<MiniHeader>;

export type MiniHeaderFilter = FilterOption<MiniHeader>;

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
    private miniHeaders: Dexie.Table<MiniHeader, string>;

    constructor(opts: Options) {
        this.db = new Dexie(opts.dataSourceName);
        this.maxOrders = opts.maxOrders;
        this.maxMiniHeaders = opts.maxMiniHeaders;

        this.db.version(1).stores({
            // TODO(albrow): Add more indexes. https://dexie.org/docs/Version/Version.stores()
            orders:
                '&hash,chainId,makerAddress,makerAssetData,makerAssetAmount,makerFee,makerFeeAssetData,takerAddress,takerAssetData,takerFeeAssetData,takerAssetAmount,takerFee,senderAddress,feeRecipientAddress,expirationTimeSeconds,salt,signature,exchangeAddress,fillableTakerAssetAmount,lastUpdated,isRemoved,isPinned,parsedMakerAssetData,parsedMakerFeeAssetData',
            miniHeaders: '&hash,parent,number,timestamp,logs',
            metadata: '',
        });

        this.orders = this.db.table('orders');
        this.miniHeaders = this.db.table('miniHeaders');
    }

    // AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
    public async addOrdersAsync(orders: Order[]): Promise<AddOrdersResult> {
        // TODO(albrow): Remove orders with max expiration time.
        const added: Order[] = [];
        await this.db.transaction('rw', this.orders, async () => {
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
        const collection = this.prepareQuery(this.orders, query);
        return this.runQueryAsync(this.orders, collection, query);
    }

    // CountOrders(opts *OrderQuery) (int, error)
    public async countOrdersAsync(query: OrderQuery): Promise<number> {
        const collection = this.prepareQuery(this.orders, query);
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
        const added: MiniHeader[] = [];
        const removed: MiniHeader[] = [];
        await this.db.transaction('rw', this.miniHeaders, async () => {
            for (const miniHeader of miniHeaders) {
                try {
                    await this.miniHeaders.add(miniHeader);
                } catch (e) {
                    if (e.name === 'ConstraintError') {
                        // A miniHeader with this hash already exists. This is fine based on the semantics of
                        // addMiniHeaders.
                        continue;
                    }
                    throw e;
                }
                added.push(miniHeader);
                const outdatedMiniHeaders = await this.miniHeaders
                    .orderBy('number')
                    .offset(this.maxMiniHeaders)
                    .reverse()
                    .toArray();
                for (const outdated of outdatedMiniHeaders) {
                    await this.miniHeaders.delete(outdated.hash);
                    removed.push(outdated);
                }
            }
        });
        return {
            added,
            removed,
        };
    }

    // GetMiniHeader(hash common.Hash) (*types.MiniHeader, error)
    public async getMiniHeaderAsync(hash: string): Promise<MiniHeader> {
        const miniHeader = await this.miniHeaders.get(hash);
        if (miniHeader === undefined) {
            throw newNotFoundError();
        }
        return miniHeader;
    }

    // FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async findMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        const collection = this.prepareQuery(this.miniHeaders, query);
        console.log(collection === undefined);
        return this.runQueryAsync(this.miniHeaders, collection, query);
    }

    // DeleteMiniHeader(hash common.Hash) error
    public async deleteMiniHeaderAsync(hash: string): Promise<void> {
        return this.miniHeaders.delete(hash);
    }

    // DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async deleteMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        const deletedMiniHeaders: MiniHeader[] = [];
        await this.db.transaction('rw', this.miniHeaders, async () => {
            // TODO(albrow): Pay special attention to this code and make sure it works.
            // Behavior of Dexie.js regarding transactions and function scope is a little
            // too complicated/magical.
            const miniHeaders = await this.findMiniHeadersAsync(query);
            for (const miniHeader of miniHeaders) {
                await this.miniHeaders.delete(miniHeader.hash);
                deletedMiniHeaders.push(miniHeader);
            }
        });
        return deletedMiniHeaders;
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

    prepareQuery<T extends Record, Key>(table: Dexie.Table<T, Key>, query?: Query<T>): Dexie.Collection<T, Key> {
        if (query === null || query === undefined) {
            return table.toCollection();
        }
        var col: Dexie.Collection<T, Key>;
        if (query.filters !== undefined && query.filters !== null && query.filters.length > 0) {
            console.log('at least one filter');
            const firstFilter = query.filters[0];
            switch (firstFilter.kind) {
                case FilterKind.Equal:
                    col = table.where(firstFilter.field).equals(firstFilter.value);
                    break;
                case FilterKind.NotEqual:
                    col = table.where(firstFilter.field).notEqual(firstFilter.value);
                    break;
                case FilterKind.Greater:
                    col = table.where(firstFilter.field).above(firstFilter.value);
                    break;
                case FilterKind.GreaterOrEqual:
                    col = table.where(firstFilter.field).aboveOrEqual(firstFilter.value);
                    break;
                case FilterKind.Less:
                    col = table.where(firstFilter.field).below(firstFilter.value);
                    break;
                case FilterKind.LessOrEqual:
                    col = table.where(firstFilter.field).belowOrEqual(firstFilter.value);
                    break;
                case FilterKind.Contains:
                    // TODO(albrow): This iterates through all orders and is very inefficient.
                    // Is there a way to optimize this?)
                    col = table.filter(containsFilterFunc(firstFilter));
                    break;
                default:
                    throw new Error(`unexpected filter kind: ${firstFilter.kind}`);
            }
            if (query.filters.length > 1) {
                // TODO(albrow): Dexie.js does not support multiple where conditions. We have to
                // use Collection.and which iterates through all orders in the collection so far
                // and is very inefficient. Is there a way to optimize this?
                query.filters.slice(1).forEach(filter => {
                    switch (filter.kind) {
                        case FilterKind.Equal:
                            col.and(record => record[filter.field] === filter.value);
                            break;
                        case FilterKind.NotEqual:
                            col.and(record => record[filter.field] !== filter.value);
                            break;
                        case FilterKind.Greater:
                            col.and(record => record[filter.field] > filter.value);
                            break;
                        case FilterKind.GreaterOrEqual:
                            col.and(record => record[filter.field] >= filter.value);
                            break;
                        case FilterKind.Less:
                            col.and(record => record[filter.field] < filter.value);
                            break;
                        case FilterKind.LessOrEqual:
                            col.and(record => record[filter.field] <= filter.value);
                            break;
                        case FilterKind.Contains:
                            col.and(containsFilterFunc(filter));
                            break;
                    }
                });
            }
        } else {
            col = table.toCollection();
        }
        if (query.offset !== undefined && query.offset !== 0) {
            col.offset(query.offset);
        }
        if (query.limit !== undefined && query.limit !== 0) {
            col.limit(query.limit);
        }
        return col;
    }

    async runQueryAsync<T extends Record, Key>(
        table: Dexie.Table<T, Key>,
        col: Dexie.Collection<T, Key>,
        query?: Query<T>,
    ): Promise<T[]> {
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
                    (query.filters.length === 1 && query.filters[0].field === table.schema.primKey.keyPath)
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
            return (await col.toArray()).sort((a: T, b: T) => {
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

function isString(x: any): x is string {
    return typeof x === 'string';
}

function containsFilterFunc<T extends Record>(filter: FilterOption<T>): (record: T) => boolean {
    return (record: T): boolean => {
        const field = record[filter.field];
        if (!isString(field)) {
            throw new Error(
                `cannot use CONTAINS filter on non-string field ${filter.field} of type ${typeof record[filter.field]}`,
            );
        }
        return field.includes(filter.value);
    };
}
