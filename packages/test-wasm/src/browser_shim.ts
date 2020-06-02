import { createDatabase } from '@0x/mesh-browser-lite/lib/database';
import { createSchemaValidator } from '@0x/mesh-browser-lite/lib/schema_validator';

(window as any).createSchemaValidator = createSchemaValidator;
(window as any).__mesh_dexie_newDatabase__ = createDatabase;
