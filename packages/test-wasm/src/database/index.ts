import Dexie, { Transaction } from 'dexie';

export interface Options {
    name: string;
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
        this.db = new Dexie(opts.name);
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
        if (query.filters !== undefined && query.filters.length > 0) {
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
                throw new Error('multiple filters not yet supported');
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
        } else if (query.sort.length === 1) {
            // TODO(albrow): As an optimization we can skip calling `sortBy`
            // if the results of the query would already be sorted in
            // "natural order" based on the filters.
            if (query.sort[0].direction === SortDirection.Desc) {
                col = col.reverse();
            }
            return col.sortBy(query.sort[0].field);
        }

        // TODO(albrow): Dexie.js can't sort by more than one field. Looks like
        // we have no choice but to manually sort here. This will not be fast or
        // efficient.
        throw new Error('sorting by multiple fields is not yet implemented');
    }
}
