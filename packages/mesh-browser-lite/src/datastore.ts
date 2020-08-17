// tslint:disable:max-file-line-count

/**
 * @hidden
 */

/**
 * NOTE(jalextowle): This comment must be here so that typedoc knows that the above
 * comment is a module comment
 */
import Dexie from 'dexie';

interface Entry {
    key: string;
    value: string;
    size: number;
}

enum OperationType {
    Addition,
    Deletion,
}

interface Operation {
    operationType: OperationType;
    key: string;
    value?: string;
}

// NOTE(jalextowle): This is missing several fields from the Query interface in
// https://github.com/ipfs/go-datastore. These fields include `returnsSizes` and
// `returnExpirations`, which are excluded because we are only satisfying the
// ds.Batching interface. Additionally, we exclude any items that require iterating
// through each key and value in a Dexie transaction. We handle that logic on the
// Go side.
interface Query {
    prefix: string; // namespaces the query to results whose keys have Prefix
    limit: number; // maximum number of results
    offset: number; // skip given number of results
}

// This implements the subset of the ds.Batching interface that should be implemented
// on the Dexie side. The Go bindings for this system can be found in db/dexie_datastore.go.
// Some aspects of the ds.Batching interface make more sense to implement in Go
// for performance or dependency reasons. The most important example of this is
// that query filtering and ordering is performed on the Go side to avoid converting
// Go functions into Javascript functions.
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

    public async commitAsync(operations: Operation[]): Promise<void> {
        await this._db.transaction('rw!', this._table, async () => {
            for (const operation of operations) {
                if (operation.operationType === OperationType.Addition) {
                    if (!operation.value) {
                        throw new Error('commitAsync: no value for key');
                    }
                    await this._table.put({ key: operation.key, value: operation.value });
                } else {
                    await this._table.delete(operation.key);
                }
            }
        });
    }

    public async putAsync(key: string, value: string): Promise<void> {
        await this._table.put({ key, value });
    }

    public async deleteAsync(key: string): Promise<void> {
        await this._table.delete(key);
    }

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
            if (query.offset !== 0) {
                col = col.offset(query.limit);
            }
            if (query.limit !== 0) {
                col = col.limit(query.limit);
            }
            const values = await col.toArray();
            const entries = (await col.keys()).map((key, i) => {
                return {
                    key: key as string,
                    value: values[i],
                    size: computeByteSize(values[i]),
                };
            });
        });
        return filteredEntries;
    }
}

function computeByteSize(value: string): number {
    return new TextEncoder().encode(value).length;
}
