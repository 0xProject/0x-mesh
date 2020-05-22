import { Dexie } from 'dexie';

export async function setupDatabase(name: string): Promise<void> {
    // TODO(albrow): Construct the database and assign it to window.
}

export interface Options {
    name: string;
    // TODO(albrow): Add other options here
}

export interface Order {
    // TODO(albrow): define this
}

export interface OrderQuery {
    // TODO(albrow): define this
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

export class Database {
    private db: Dexie;

    constructor(opts: Options) {
        this.db = new Dexie(opts.name);
        this.db.version(0).stores({
            // TODO(albrow): Add more indexes. https://dexie.org/docs/Version/Version.stores()
            orders: '&hash',
        });
    }

    // TODO(albrow): Implement all methods here.
    // AddOrders(orders []*types.OrderWithMetadata) (added []*types.OrderWithMetadata, removed []*types.OrderWithMetadata, err error)
    public async addOrdersAsync(orders: Order[]): Promise<AddOrdersResult> {
        return Promise.reject('not yet implemented');
    }
    // GetOrder(hash common.Hash) (*types.OrderWithMetadata, error)
    public async getOrderAsync(hash: string): Promise<Order> {
        return Promise.reject('not yet implemented');
    }
    // FindOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async findOrdersAsync(query: OrderQuery): Promise<Order[]> {
        return Promise.reject('not yet implemented');
    }
    // CountOrders(opts *OrderQuery) (int, error)
    public async countOrdersAsync(query: OrderQuery): Promise<number> {
        return Promise.reject('not yet implemented');
    }
    // DeleteOrder(hash common.Hash) error
    public async deleteOrderAsync(hash: string): Promise<void> {
        return Promise.reject('not yet implemented');
    }
    // DeleteOrders(opts *OrderQuery) ([]*types.OrderWithMetadata, error)
    public async deleteOrdersAsync(query: OrderQuery): Promise<Order[]> {
        return Promise.reject('not yet implemented');
    }
    // UpdateOrder(hash common.Hash, updateFunc func(existingOrder *types.OrderWithMetadata) (updatedOrder *types.OrderWithMetadata, err error)) error
    public async updateOrdersAsync(hash: string, updateFunc: (existingOrder: Order) => Order): Promise<void> {
        return Promise.reject('not yet implemented');
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
}
