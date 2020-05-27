import { createSchemaValidator } from '@0x/mesh-browser-lite/lib/schema_validator';
import { Database, Options } from './database';

(window as any).createSchemaValidator = createSchemaValidator;

(window as any).__mesh_dexie_newDatabase__ = function(opts: Options): Database {
    return new Database(opts);
};
