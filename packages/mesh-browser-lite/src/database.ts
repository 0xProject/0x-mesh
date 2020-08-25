// tslint:disable:max-file-line-count

/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */

import Dexie from 'dexie';

import { BatchingDatastore } from './datastore';

export type Record = Order | MiniHeader | Metadata;

export interface Options {
    dataSourceName: string;
    maxOrders: number;
    maxMiniHeaders: number;
}

export interface Query<T extends Record> {
    filters?: Array<FilterOption<T>>;
    sort?: Array<SortOption<T>>;
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
    isNotPinned: number; // Used in a compound index in queries related to max expiration time.
    parsedMakerAssetData: string;
    parsedMakerFeeAssetData: string;
    lastValidatedBlockNumber: string;
    lastValidatedBlockHash: string;
}

export interface StoredOrderStatus {
    isStored: boolean;
    isMarkedRemoved: boolean;
    fillableTakerAssetAmount?: string;
}

export type OrderField = keyof Order;

export type OrderQuery = Query<Order>;

export type OrderSort = SortOption<Order>;

export type OrderFilter = FilterOption<Order>;

export interface AddOrdersResult {
    alreadyStored: string[];
    added: Order[];
    removed: Order[];
}

export interface MiniHeader {
    hash: string;
    parent: string;
    number: string;
    timestamp: string;
    logs: Uint8Array;
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
    ethRPCRequestsSentInCurrentUTCDay: number;
    startOfCurrentUTCDay: string;
}

function newNotFoundError(): Error {
    return new Error('could not find existing model or row in database');
}

function newMetadataAlreadExistsError(): Error {
    return new Error('metadata already exists in the database (use UpdateMetadata instead?)');
}

/**
 * Creates and returns a new database
 *
 * @param opts The options to use for the database
 */
export function createDatabase(opts: Options): Database {
    return new Database(opts);
}

export class Database {
    private readonly _db: Dexie;
    private readonly _maxOrders: number;
    private readonly _maxMiniHeaders: number;
    private readonly _orders: Dexie.Table<Order, string>;
    private readonly _miniHeaders: Dexie.Table<MiniHeader, string>;
    private readonly _metadata: Dexie.Table<Metadata, number>;
    private readonly _dhtstore: BatchingDatastore;
    private readonly _peerstore: BatchingDatastore;

    constructor(opts: Options) {
        this._db = new Dexie(opts.dataSourceName);
        this._maxOrders = opts.maxOrders;
        this._maxMiniHeaders = opts.maxMiniHeaders;

        this._db.version(1).stores({
            orders:
                '&hash,chainId,makerAddress,makerAssetData,makerAssetAmount,makerFee,makerFeeAssetData,takerAddress,takerAssetData,takerFeeAssetData,takerAssetAmount,takerFee,senderAddress,feeRecipientAddress,expirationTimeSeconds,salt,signature,exchangeAddress,fillableTakerAssetAmount,lastUpdated,isRemoved,isPinned,parsedMakerAssetData,parsedMakerFeeAssetData,lastValidatedBlockNumber,lastValidatedBlockHash,[isNotPinned+expirationTimeSeconds]',
            miniHeaders: '&hash,parent,number,timestamp',
            metadata: '&ethereumChainID',
            dhtstore: '&key,data',
            peerstore: '&key,data',
        });

        this._orders = this._db.table('orders');
        this._miniHeaders = this._db.table('miniHeaders');
        this._metadata = this._db.table('metadata');
        this._dhtstore = new BatchingDatastore(this._db, 'dhtstore');
        this._peerstore = new BatchingDatastore(this._db, 'peerstore');
    }

    public dhtStore(): BatchingDatastore {
        return this._dhtstore;
    }

    public peerStore(): BatchingDatastore {
        return this._peerstore;
    }

    public close(): void {
        this._db.close();
    }

    // AddOrders(orders []*types.OrderWithMetadata) (alreadyStored []common.Hash, added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
    public async addOrdersAsync(orders: Order[]): Promise<AddOrdersResult> {
        const alreadyStored: string[] = [];
        const addedMap = new Map<string, Order>();
        const removed: Order[] = [];
        await this._db.transaction('rw!', this._orders, async () => {
            for (const order of orders) {
                try {
                    await this._orders.add(order);
                } catch (e) {
                    if (e.name === 'ConstraintError') {
                        // An order with this hash already exists. Add the order hash to the
                        // array of alreadyStored
                        alreadyStored.push(order.hash);
                        continue;
                    }
                    throw e;
                }
                addedMap.set(order.hash, order);
            }

            // Remove orders with an expiration time too far in the future.
            const ordersToRemove = await this._orders
                .orderBy('[isNotPinned+expirationTimeSeconds]')
                .offset(this._maxOrders)
                .toArray();
            for (const order of ordersToRemove) {
                await this._orders.delete(order.hash);
                if (addedMap.has(order.hash)) {
                    // If the order was previously added, remove it from
                    // the added set and don't add it to the removed set.
                    addedMap.delete(order.hash);
                } else {
                    removed.push(order);
                }
            }
        });

        return {
            alreadyStored,
            added: Array.from(addedMap.values()),
            removed,
        };
    }

    // GetOrder(hash common.Hash) (*types.OrderWithMetadata, error)
    public async getOrderAsync(hash: string): Promise<Order> {
        return this._db.transaction('r!', this._orders, async () => {
            const order = await this._orders.get(hash);
            if (order === undefined) {
                throw newNotFoundError();
            }
            return order;
        });
    }

    // GetOrderStatuses(hashes []common.Hash) (statuses []StoredOrderStatus, err error)
    public async getOrderStatusesAsync(hashes: string[]): Promise<StoredOrderStatus[]> {
        let orders: Order[] = [];
        await this._db.transaction('r!', this._orders, async () => {
            orders = await this._orders.bulkGet(hashes);
        });
        const statuses: StoredOrderStatus[] = [];
        for (const order of orders) {
            if (order === undefined) {
                statuses.push({
                    isStored: false,
                    isMarkedRemoved: false,
                });
            } else if (order.isRemoved) {
                statuses.push({
                    isStored: true,
                    isMarkedRemoved: true,
                    fillableTakerAssetAmount: order.fillableTakerAssetAmount,
                });
            } else {
                statuses.push({
                    isStored: true,
                    isMarkedRemoved: false,
                    fillableTakerAssetAmount: order.fillableTakerAssetAmount,
                });
            }
        }
        return statuses;
    }

    // FindOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async findOrdersAsync(query?: OrderQuery): Promise<Order[]> {
        return this._db.transaction('r!', this._orders, async () => {
            return findRecordsAsync(this._orders, query);
        });
    }

    // CountOrders(opts *OrderQuery) (int, error)
    public async countOrdersAsync(query?: OrderQuery): Promise<number> {
        return this._db.transaction('r!', this._orders, async () => {
            if (!canUseNativeDexieIndexes(this._orders, query)) {
                // As a fallback, implement the query inefficiently (in-memory).
                // Note(albrow): If needed we can optimize specific common queries with compound indexes.
                const records = await runQueryInMemoryAsync(this._orders, query);
                return records.length;
            }
            const col = buildCollectionWithDexieIndexes(this._orders, query);
            return col.count();
        });
    }

    // DeleteOrder(hash common.Hash) error
    public async deleteOrderAsync(hash: string): Promise<void> {
        return this._db.transaction('rw!', this._orders, async () => {
            return this._orders.delete(hash);
        });
    }

    // DeleteOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async deleteOrdersAsync(query: OrderQuery | undefined): Promise<Order[]> {
        const deletedOrders: Order[] = [];
        await this._db.transaction('rw!', this._orders, async () => {
            const orders = await findRecordsAsync(this._orders, query);
            for (const order of orders) {
                await this._orders.delete(order.hash);
                deletedOrders.push(order);
            }
        });
        return deletedOrders;
    }

    // UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error
    public async updateOrderAsync(hash: string, updateFunc: (existingOrder: Order) => Order): Promise<void> {
        await this._db.transaction('rw!', this._orders, async () => {
            const existingOrder = await this._orders.get(hash);
            if (existingOrder === undefined) {
                throw newNotFoundError();
            }
            const updatedOrder = updateFunc(existingOrder);
            await this._orders.put(updatedOrder, hash);
        });
    }

    // AddMiniHeaders(miniHeaders []*types.MiniHeader) (added []*types.MiniHeader, removed []*types.MiniHeader, err error)
    public async addMiniHeadersAsync(miniHeaders: MiniHeader[]): Promise<AddMiniHeadersResult> {
        const addedMap = new Map<string, MiniHeader>();
        const removed: MiniHeader[] = [];
        await this._db.transaction('rw!', this._miniHeaders, async () => {
            for (const miniHeader of miniHeaders) {
                try {
                    await this._miniHeaders.add(miniHeader);
                } catch (e) {
                    if (e.name === 'ConstraintError') {
                        // A miniHeader with this hash already exists. This is fine based on the semantics of
                        // addMiniHeaders.
                        continue;
                    }
                    throw e;
                }
                addedMap.set(miniHeader.hash, miniHeader);
            }

            // Remove any outdated miniHeaders.
            const outdatedMiniHeaders = await this._miniHeaders
                .orderBy('number')
                .offset(this._maxMiniHeaders)
                .reverse()
                .toArray();
            for (const outdated of outdatedMiniHeaders) {
                await this._miniHeaders.delete(outdated.hash);
                if (addedMap.has(outdated.hash)) {
                    // If the order was previously added, remove it from
                    // the added set and don't add it to the removed set.
                    addedMap.delete(outdated.hash);
                } else {
                    removed.push(outdated);
                }
            }
        });

        return {
            added: Array.from(addedMap.values()),
            removed,
        };
    }

    // ResetMiniHeaders(newMiniHeaders []*types.MiniHeader) (err error)
    public async resetMiniHeadersAsync(newMiniHeaders: MiniHeader[]): Promise<void> {
        await this._db.transaction('rw!', this._miniHeaders, async () => {
            // Remove all of the existing miniheaders
            await this._miniHeaders.clear();

            for (const newMiniHeader of newMiniHeaders) {
                try {
                    await this._miniHeaders.add(newMiniHeader);
                } catch (e) {
                    if (e.name === 'ConstraintError') {
                        // A miniHeader with this hash already exists. This is
                        // fine based on the semantics of addMiniHeaders.
                        continue;
                    }
                    throw e;
                }
            }
        });
    }

    // GetMiniHeader(hash common.Hash) (*types.MiniHeader, error)
    public async getMiniHeaderAsync(hash: string): Promise<MiniHeader> {
        return this._db.transaction('r!', this._miniHeaders, async () => {
            const miniHeader = await this._miniHeaders.get(hash);
            if (miniHeader === undefined) {
                throw newNotFoundError();
            }
            return miniHeader;
        });
    }

    // FindMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async findMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        return this._db.transaction('r!', this._miniHeaders, async () => {
            return findRecordsAsync(this._miniHeaders, query);
        });
    }

    // DeleteMiniHeader(hash common.Hash) error
    public async deleteMiniHeaderAsync(hash: string): Promise<void> {
        return this._db.transaction('rw!', this._miniHeaders, async () => {
            return this._miniHeaders.delete(hash);
        });
    }

    // DeleteMiniHeaders(opts *MiniHeaderQuery) ([]*types.MiniHeader, error)
    public async deleteMiniHeadersAsync(query: MiniHeaderQuery): Promise<MiniHeader[]> {
        const deletedMiniHeaders: MiniHeader[] = [];
        await this._db.transaction('rw!', this._miniHeaders, async () => {
            const miniHeaders = await findRecordsAsync(this._miniHeaders, query);
            for (const miniHeader of miniHeaders) {
                await this._miniHeaders.delete(miniHeader.hash);
                deletedMiniHeaders.push(miniHeader);
            }
        });
        return deletedMiniHeaders;
    }

    // GetMetadata() (*types.Metadata, error)
    public async getMetadataAsync(): Promise<Metadata> {
        return this._db.transaction('r!', this._metadata, async () => {
            return this._getMetadataAsync();
        });
    }

    // SaveMetadata(metadata *types.Metadata) error
    public async saveMetadataAsync(metadata: Metadata): Promise<void> {
        await this._db.transaction('rw!', this._metadata, async () => {
            if ((await this._metadata.count()) > 0) {
                throw newMetadataAlreadExistsError();
            }
            await this._metadata.add(metadata);
        });
    }

    // UpdateMetadata(updateFunc func(oldmetadata *types.Metadata) (newMetadata *types.Metadata)) error
    public async updateMetadataAsync(updateFunc: (existingMetadata: Metadata) => Metadata): Promise<void> {
        await this._db.transaction('rw!', this._metadata, async () => {
            const existingMetadata = await this._getMetadataAsync();
            const updatedMetadata = updateFunc(existingMetadata);
            await this._metadata.put(updatedMetadata);
        });
    }

    private async _getMetadataAsync(): Promise<Metadata> {
        const count = await this._metadata.count();
        if (count === 0) {
            throw newNotFoundError();
        } else if (count > 1) {
            // This should never happen, but it's possible if a user manually messed around with
            // IndexedDB. In this case, just delete the metadata table and we should start
            // over.
            await this._metadata.clear();
            throw new Error('more than one metadata entry stored in the database');
        }
        const metadatas = await this._metadata.toArray();
        return metadatas[0];
    }
}

async function findRecordsAsync<T extends Record, Key>(table: Dexie.Table<T, Key>, query?: Query<T>): Promise<T[]> {
    if (!canUseNativeDexieIndexes(table, query)) {
        // As a fallback, implement the query inefficiently (in-memory).
        // Note(albrow): If needed we can optimize specific common queries with compound indexes.
        return runQueryInMemoryAsync(table, query);
    }
    const col = buildCollectionWithDexieIndexes(table, query);
    return col.toArray();
}

function buildCollectionWithDexieIndexes<T extends Record, Key>(
    table: Dexie.Table<T, Key>,
    query?: Query<T>,
): Dexie.Collection<T, Key> {
    if (query === null || query === undefined) {
        return table.toCollection();
    }

    // First we create the Collection based on the query fields.
    let col: Dexie.Collection<T, Key>;
    if (queryUsesFilters(query)) {
        // tslint:disable-next-line:no-non-null-assertion
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
                // Note(albrow): This iterates through all orders and is very inefficient.
                // If needed, we should try to find a way to optimize this.
                col = table.filter(containsFilterFunc(filter));
                break;
            default:
                throw new Error(`unexpected filter kind: ${filter.kind}`);
        }
        // tslint:disable-next-line:no-non-null-assertion
        if (queryUsesSortOptions(query) && query.sort![0].direction === SortDirection.Desc) {
            // Note(albrow): This is only allowed if the sort and filter are using
            // the same field. Dexie automatically returns records sorted by the filter
            // field. If the direction is Ascending, we don't need to do anything else.
            // If it the direction is Descending, we just need to call reverse().
            col.reverse();
        }
    } else if (queryUsesSortOptions(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        const sortOpt = query.sort![0];
        col = table.orderBy(sortOpt.field);
        if (sortOpt.direction === SortDirection.Desc) {
            col = col.reverse();
        }
    } else {
        // Query doesn't use filter or sort options.
        col = table.toCollection();
    }
    if (queryUsesOffset(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        col.offset(query.offset!);
    }
    if (queryUsesLimit(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        col.limit(query.limit!);
    }
    return col;
}

async function runQueryInMemoryAsync<T extends Record, Key>(
    table: Dexie.Table<T, Key>,
    query?: Query<T>,
): Promise<T[]> {
    let records = await table.toArray();
    if (query === undefined || query === null) {
        return records;
    }
    if (queryUsesFilters(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        records = filterRecords(query.filters!, records);
    }
    if (queryUsesSortOptions(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        records = sortRecords(query.sort!, records);
    }
    if (queryUsesOffset(query) && queryUsesLimit(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        records = records.slice(query.offset!, query.limit!);
    } else if (queryUsesLimit(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        records = records.slice(0, query.limit!);
    } else if (queryUsesOffset(query)) {
        // tslint:disable-next-line:no-non-null-assertion
        records = records.slice(query.offset!);
    }

    return records;
}

function filterRecords<T extends Record>(filters: Array<FilterOption<T>>, records: T[]): T[] {
    let result = records;
    // Note(albrow): As an optimization, we could use the native Dexie.js index for
    // the *first* filter when possible.
    for (const filter of filters) {
        switch (filter.kind) {
            case FilterKind.Equal:
                result = result.filter(record => record[filter.field] === filter.value);
                break;
            case FilterKind.NotEqual:
                result = result.filter(record => record[filter.field] !== filter.value);
                break;
            case FilterKind.Greater:
                result = result.filter(record => record[filter.field] > filter.value);
                break;
            case FilterKind.GreaterOrEqual:
                result = result.filter(record => record[filter.field] >= filter.value);
                break;
            case FilterKind.Less:
                result = result.filter(record => record[filter.field] < filter.value);
                break;
            case FilterKind.LessOrEqual:
                result = result.filter(record => record[filter.field] <= filter.value);
                break;
            case FilterKind.Contains:
                result = result.filter(containsFilterFunc(filter));
                break;
            default:
                throw new Error(`unexpected filter kind: ${filter.kind}`);
        }
    }

    return result;
}

function sortRecords<T extends Record>(sortOpts: Array<SortOption<T>>, records: T[]): T[] {
    // Note(albrow): As an optimization, we could use native Dexie.js ordering for
    // the *first* sort option when possible.
    const result = records;
    return result.sort((a: T, b: T) => {
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
                default:
                    throw new Error(`unexpected sort direction: ${s.direction}`);
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
    if (query === null || query === undefined) {
        return true;
    }
    // tslint:disable-next-line:no-non-null-assertion
    if (queryUsesSortOptions(query) && query.sort!.length > 1) {
        // Dexie does not support multiple sort orders.
        return false;
    }
    // tslint:disable-next-line:no-non-null-assertion
    if (queryUsesFilters(query) && query.filters!.length > 1) {
        // Dexie does not support multiple filters.
        return false;
    }
    // tslint:disable-next-line:no-non-null-assertion
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
