// tslint:disable:max-file-line-count

/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
// FIXME - Add better comments

import Dexie from 'dexie';

export enum OperationType {
    Addition,
    Removal,
}

export interface Operation {
    operationType: OperationType;
    key: string;
    value?: string;
}

export type Filter = (entry: Entry) => boolean;
export type Order = (a: Entry, b: Entry) => number;

export interface Entry {
    key: string;
    value: string;
    size: number;
}

// NOTE(jalextowle): This is missing several fields from the Query interface in
// https://github.com/ipfs/go-datastore. These fields include `returnsSizes` and
// `returnExpirations`, which are excluded because we are only satisfying the
// ds.Batching interface.
export interface Query {
    prefix: string; // namespaces the query to results whose keys have Prefix
    filters: Filter[]; // filter results. apply sequentially
    orders: Order[]; // order results. apply hierarchically
    limit: number; // maximum number of results
    offset: number; // skip given number of results
}

export class BatchingDatastore {
    private readonly _db: Dexie;
    private readonly _table: Dexie.Table;

    constructor(db: Dexie, tableName: string) {
        this._db = db;
        this._table = this._db._allTables[tableName];
        if (this._table === undefined) {
            throw new Error('BatchingDatastore: Attempting to use undefined table');
        }
    }

    /*** ds.Batching ***/

    public async commitAsync(operations: Operation[]): Promise<void> {
        await this._db.transaction('rw!', this._table, async () => {
            for (const operation of operations) {
                if (operation.operationType === OperationType.Addition) {
                    if (!operation.value) {
                        throw new Error('commitDHTAsync: no value for key');
                    }
                    await this._table.add(operation.value, operation.key);
                } else {
                    await this._table.delete(operation.key);
                }
            }
        });
    }

    /*** ds.Write ***/

    public async putAsync(key: string, value: string): Promise<void> {
        await this._table.put(value, key);
    }

    public async deleteAsync(key: string): Promise<void> {
        await this._table.delete(key);
    }

    /*** ds.Read ***/

    public async getAsync(key: string): Promise<string> {
        const value = await this._table.get(key);
        return value || '';
    }

    public async getSizeAsync(key: string): Promise<number> {
        return computeByteSize(await this.getAsync(key));
    }

    public async hasAsync(key: string): Promise<boolean> {
        return (await this.getAsync(key)) === undefined;
    }

    public async queryAsync(query: Query): Promise<Entry[]> {
        const filteredEntries: Entry[] = [];
        await this._db.transaction('rw!', this._table, async () => {
            let col =
                query.prefix === ''
                    ? this._table.toCollection()
                    : await this._table.where('key').startsWith(query.prefix);
            // FIXME - Is this the correct order for the limit and order fields?
            if (query.limit !== 0) {
                col = col.limit(query.limit);
            }
            if (query.offset !== 0) {
                col = col.offset(query.limit);
            }
            const values = await col.toArray();
            const entries = (await col.keys()).map((key, i) => {
                return {
                    key: key as string,
                    value: values[i],
                    size: computeByteSize(values[i]),
                };
            });
            for (const entry of entries) {
                let passes = true;
                for (const filter of query.filters) {
                    if (!filter(entry)) {
                        passes = false;
                        break;
                    }
                }
                if (passes) {
                    filteredEntries.push(entry);
                }
            }
            const masterComparator = createMasterComparator(query.orders);
            filteredEntries.sort(masterComparator);
        });
        return filteredEntries;
    }
}

function computeByteSize(value: string): number {
    return new TextEncoder().encode(value).length;
}

function createMasterComparator(orders: Order[]): (a: Entry, b: Entry) => number {
    return (a: Entry, b: Entry) => {
        let comparison = 0;
        for (const order of orders) {
            comparison = order(a, b);
            if (comparison !== 0) {
                return comparison;
            }
        }
        return comparison;
    };
}
