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
    value: Uint8Array;
    size: number;
}

enum OperationType {
    Addition,
    Deletion,
}

interface Operation {
    operationType: OperationType;
    key: string;
    value?: Uint8Array;
}

// This implements the subset of the ds.Batching interface that should be implemented
// on the Dexie side. The Go bindings for this system can be found in db/dexie_datastore.go.
// Some aspects of the ds.Batching interface make more sense to implement in Go
// for performance or dependency reasons. The most important example of this is
// that query filtering and ordering is performed on the Go side to avoid converting
// Go functions into Javascript functions.
export class BatchingDatastore {
    private readonly _db: Dexie;
    private readonly _table: Dexie.Table<{ key: string; value: Uint8Array }, string>;

    constructor(db: Dexie, tableName: string) {
        this._db = db;
        this._table = this._db._allTables[tableName] as Dexie.Table<{ key: string; value: Uint8Array }, string>;
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

    public async putAsync(key: string, value: Uint8Array): Promise<void> {
        await this._table.put({ key, value });
    }

    public async deleteAsync(key: string): Promise<void> {
        await this._table.delete(key);
    }

    public async getAsync(key: string): Promise<Uint8Array> {
        const result = await this._table.get(key);
        if (result === undefined) {
            throw new Error('datastore: key not found');
        }
        return result.value;
    }

    public async getSizeAsync(key: string): Promise<number> {
        const value = await this.getAsync(key);
        if (value === undefined) {
            throw new Error('datastore: key not found');
        }
        return value.length;
    }

    public async hasAsync(key: string): Promise<boolean> {
        const result = await this._table.get(key);
        return result !== undefined;
    }

    // NOTE(jalextowle): This function only filters the database based on prefix
    // and generates entries for each row of the database. The other query
    // operations (filtering, sorting, etc.) are implemented in
    // db/dexie_datastore.go for performance reasons. The prefixes are
    // interpreted as regular expressions to satisfy the ds.Datastore interface.
    public async queryAsync(prefix: string): Promise<Entry[]> {
        return this._db.transaction('rw!', this._table, async () => {
            const prefixRegExp = new RegExp(prefix);
            const col =
                prefix !== ''
                    ? this._table.filter((entry) => prefixRegExp.test(entry.key))
                    : this._table.toCollection();
            return (await col.toArray()).map((entry) => {
                return {
                    key: entry.key,
                    value: entry.value,
                    size: entry.value.length,
                };
            });
        });
    }
}
