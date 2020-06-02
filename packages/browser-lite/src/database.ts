import Dexie from 'dexie';
import { fromTokenUnitAmount } from '@0x/utils';

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
    isRemoved: number;
    isPinned: number;
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
    ethereumChainID: number;
    maxExpirationTime: string;
    ethRPCRequestsSentInCurrentUTCDay: number;
    startOfCurrentUTCDay: string;
}

function newNotFoundError(): Error {
    return new Error('could not find existing model or row in database');
}

function newMetadataAlreadExistsError(): Error {
    return new Error('metadata already exists in the database (use UpdateMetadata instead?)');
}

export function createDatabase(opts: Options): Database {
    return new Database(opts);
}

// TODO(albrow): Implement remaining methods.
export class Database {
    private db: Dexie;
    private maxOrders: number;
    private maxMiniHeaders: number;
    private orders: Dexie.Table<Order, string>;
    private miniHeaders: Dexie.Table<MiniHeader, string>;
    private metadata: Dexie.Table<Metadata, number>;

    constructor(opts: Options) {
        this.db = new Dexie(opts.dataSourceName);
        this.maxOrders = opts.maxOrders;
        this.maxMiniHeaders = opts.maxMiniHeaders;

        this.db.version(1).stores({
            orders:
                '&hash,chainId,makerAddress,makerAssetData,makerAssetAmount,makerFee,makerFeeAssetData,takerAddress,takerAssetData,takerFeeAssetData,takerAssetAmount,takerFee,senderAddress,feeRecipientAddress,expirationTimeSeconds,salt,signature,exchangeAddress,fillableTakerAssetAmount,lastUpdated,isRemoved,isPinned,parsedMakerAssetData,parsedMakerFeeAssetData',
            miniHeaders: '&hash,parent,number,timestamp,logs',
            metadata: '&ethereumChainID',
        });

        this.orders = this.db.table('orders');
        this.miniHeaders = this.db.table('miniHeaders');
        this.metadata = this.db.table('metadata');
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
        if (!canUseNativeDexieIndexes(this.orders, query)) {
            // As a fallback, implement the query inefficiently (in-memory).
            // TODO(albrow): Find some ways to optimize specific common queries with compound indexes.
            return this.runQueryInMemoryAsync(this.orders, query!);
        }
        const col = this.buildCollectionWithDexieIndexes(this.orders, query);
        return col.toArray();
    }

    // CountOrders(opts *OrderQuery) (int, error)
    public async countOrdersAsync(query: OrderQuery): Promise<number> {
        if (!canUseNativeDexieIndexes(this.orders, query)) {
            // As a fallback, implement the query inefficiently (in-memory).
            // TODO(albrow): Find some ways to optimize specific common queries with compound indexes.
            const records = await this.runQueryInMemoryAsync(this.orders, query);
            return records.length;
        }
        const col = this.buildCollectionWithDexieIndexes(this.orders, query);
        return col.count();
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
        if (!canUseNativeDexieIndexes(this.miniHeaders, query)) {
            // As a fallback, implement the query inefficiently (in-memory).
            // TODO(albrow): Find some ways to optimize specific common queries with compound indexes.
            return this.runQueryInMemoryAsync(this.miniHeaders, query);
        }
        const col = this.buildCollectionWithDexieIndexes(this.miniHeaders, query);
        return col.toArray();
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
        const metadatas = await this.metadata.limit(1).toArray();
        if (metadatas.length == 0) {
            throw newNotFoundError();
        }
        return metadatas[0];
    }

    // SaveMetadata(metadata *types.Metadata) error
    public async saveMetadataAsync(metadata: Metadata): Promise<void> {
        await this.db.transaction('rw', this.metadata, async () => {
            if ((await this.metadata.count()) > 0) {
                throw newMetadataAlreadExistsError();
            }
            await this.metadata.add(metadata);
        });
    }

    // UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error
    public async updateMetadataAsync(updateFunc: (existingMetadata: Metadata) => Metadata): Promise<void> {
        await this.db.transaction('rw', this.metadata, async () => {
            // TODO(albrow): Pay special attention to this code and make sure it works.
            // Behavior of Dexie.js regarding transactions and function scope is a little
            // too complicated/magical.
            const existingMetadata = await this.getMetadataAsync();
            const updatedMetadata = updateFunc(existingMetadata);
            await this.metadata.put(updatedMetadata);
        });
    }

    buildCollectionWithDexieIndexes<T extends Record, Key>(
        table: Dexie.Table<T, Key>,
        query?: Query<T>,
    ): Dexie.Collection<T, Key> {
        if (query === null || query === undefined) {
            return table.toCollection();
        }

        // First we create the Collection based on the query fields.
        var col: Dexie.Collection<T, Key>;
        if (queryUsesFilters(query)) {
            const filter = query.filters![0];
            switch (filter.kind) {
                case FilterKind.Equal:
                    col = table.where(filter.field).equals(filter.value);
                    break;
                case FilterKind.NotEqual:
                    col = table.where(filter.field).notEqual(filter.value);
                    break;
                case FilterKind.Greater:
                    col = table.where(filter.field).above(filter.value);
                    break;
                case FilterKind.GreaterOrEqual:
                    col = table.where(filter.field).aboveOrEqual(filter.value);
                    break;
                case FilterKind.Less:
                    col = table.where(filter.field).below(filter.value);
                    break;
                case FilterKind.LessOrEqual:
                    col = table.where(filter.field).belowOrEqual(filter.value);
                    break;
                case FilterKind.Contains:
                    // TODO(albrow): This iterates through all orders and is very inefficient.
                    // Is there a way to optimize this?)
                    col = table.filter(containsFilterFunc(filter));
                    break;
                default:
                    throw new Error(`unexpected filter kind: ${filter.kind}`);
            }
        } else if (queryUsesSortOptions(query)) {
            const sortOpt = query.sort![0];
            col = table.orderBy(sortOpt.field);
            if (sortOpt.direction === SortDirection.Desc) {
                col = col.reverse();
            }
        } else {
            col = table.toCollection();
        }
        if (queryUsesOffset(query)) {
            col.offset(query.offset!);
        }
        if (queryUsesLimit(query)) {
            col.limit(query.limit!);
        }
        return col;
    }

    async runQueryInMemoryAsync<T extends Record, Key>(table: Dexie.Table<T, Key>, query: Query<T>): Promise<T[]> {
        let records = await table.toArray();
        if (queryUsesFilters(query)) {
            records = filterRecords(query.filters!, records);
        }
        if (queryUsesSortOptions(query)) {
            records = sortRecords(query.sort!, records);
        }
        if (queryUsesOffset(query) && queryUsesLimit(query)) {
            records = records.slice(query.offset!, query.limit!);
        } else if (queryUsesLimit(query)) {
            records = records.slice(0, query.limit!);
        } else if (queryUsesOffset(query)) {
            records = records.slice(query.offset!);
        }

        return records;
    }
}

function filterRecords<T extends Record>(filters: FilterOption<T>[], records: T[]): T[] {
    // TODO(albrow): Use the native Dexie.js index for the *first* filter when possible.
    for (let filter of filters) {
        switch (filter.kind) {
            case FilterKind.Equal:
                records = records.filter(record => record[filter.field] === filter.value);
                break;
            case FilterKind.NotEqual:
                records = records.filter(record => record[filter.field] !== filter.value);
                break;
            case FilterKind.Greater:
                records = records.filter(record => record[filter.field] > filter.value);
                break;
            case FilterKind.GreaterOrEqual:
                records = records.filter(record => record[filter.field] >= filter.value);
                break;
            case FilterKind.Less:
                records = records.filter(record => record[filter.field] < filter.value);
                break;
            case FilterKind.LessOrEqual:
                records = records.filter(record => record[filter.field] <= filter.value);
                break;
            case FilterKind.Contains:
                records = records.filter(containsFilterFunc(filter));
                break;
        }
    }

    return records;
}

function sortRecords<T extends Record>(sortOpts: SortOption<T>[], records: T[]): T[] {
    // TODO(albrow): Use native Dexie.js ordering when possible.
    return records.sort((a: T, b: T) => {
        for (const s of sortOpts) {
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

function canUseNativeDexieIndexes<T extends Record, Key>(table: Dexie.Table<T, Key>, query?: Query<T>): boolean {
    if (query === null || query === undefined) return true;
    if (queryUsesSortOptions(query) && query.sort!.length > 1) {
        // Dexie does not support multiple sort orders.
        return false;
    }
    if (queryUsesFilters(query) && query.filters!.length > 1) {
        // Dexie does not support multiple filters.
        return false;
    }
    if (queryUsesFilters(query) && queryUsesSortOptions(query) && query.filters![0].field !== query.sort![0].field) {
        // Dexie does not support sorting and filtering by two different fields.
        return false;
    }
    return true;
}

function queryUsesSortOptions<T extends Record>(query: Query<T>): boolean {
    return query.sort !== null && query.sort !== undefined && query.sort.length > 0;
}

function queryUsesFilters<T extends Record>(query: Query<T>): boolean {
    return query.filters !== null && query.filters !== undefined && query.filters.length > 0;
}

function queryUsesLimit<T extends Record>(query: Query<T>): boolean {
    return query.limit !== null && query.limit !== undefined && query.limit !== 0;
}

function queryUsesOffset<T extends Record>(query: Query<T>): boolean {
    return query.offset !== null && query.offset !== undefined && query.offset !== 0;
}

function filterUsesPrimaryKey<T extends Record, Key>(table: Dexie.Table<T, Key>, filter: FilterOption<T>): boolean {
    return filter.field === table.schema.primKey.keyPath;
}

function sortUsesPrimaryKey<T extends Record, Key>(table: Dexie.Table<T, Key>, sort: SortOption<T>): boolean {
    return sort.field === table.schema.primKey.keyPath;
}
